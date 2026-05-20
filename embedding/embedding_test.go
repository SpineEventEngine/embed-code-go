// Copyright 2026, TeamDev. All rights reserved.
//
// Redistribution and use in source and/or binary forms, with or without
// modification, must retain the above copyright notice and the following
// disclaimer.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package embedding_test

import (
	"embed-code/embed-code-go/files"
	_type "embed-code/embed-code-go/type"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"embed-code/embed-code-go/configuration"
	"embed-code/embed-code-go/embedding"
	"embed-code/embed-code-go/embedding/parsing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const temporaryTestDir = "../test/docs"

func TestEmbedding(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Data Suite")
}

var _ = Describe("Embedding", func() {
	var config configuration.Configuration

	BeforeEach(func() {
		currentDir, err := os.Getwd()
		if err != nil {
			Fail("unexpected error during the test setup: " + err.Error())
		}
		err = os.Chdir(currentDir)
		if err != nil {
			Fail("unexpected error during the test setup: " + err.Error())
		}
		config = buildConfigWithSourceFiles()

		// Copying files not to edit them directly during the test run.
		copyDirRecursive("../test/resources/docs", config.DocumentationRoot)
	})

	AfterEach(func() {
		if err := os.RemoveAll(temporaryTestDir); err != nil {
			Fail(err.Error())
		}
	})

	It("should be up to date", func() {
		docPath := fmt.Sprintf("%s/whole-file-fragment.md", config.DocumentationRoot)
		processor := embedding.NewProcessor(docPath, config)

		Expect(processor.Embed()).Error().ShouldNot(HaveOccurred())
		Expect(processor.IsUpToDate()).Should(BeTrue())
	})

	It("should be up to date as there is nothing to update", func() {
		docPath := fmt.Sprintf("%s/no-embedding-doc.md", config.DocumentationRoot)
		processor := embedding.NewProcessor(docPath, config)

		Expect(processor.Embed()).Error().ShouldNot(HaveOccurred())
		Expect(processor.IsUpToDate()).Should(BeTrue())
	})

	// TODO:olena-zmiiova:https://github.com/SpineEventEngine/embed-code/issues/59
	It("should have error as it has invalid transition map", func() {
		Skip(
			"Temporarily disabled, see " +
				"[issue #59](https://github.com/SpineEventEngine/embed-code/issues/59).",
		)
		docPath := fmt.Sprintf("%s/split-lines.md", config.DocumentationRoot)

		falseTransitions := parsing.TransitionMap{
			parsing.Start:       {parsing.Finish, parsing.EmbedInstruction, parsing.RegularLine},
			parsing.RegularLine: {parsing.CodeFenceEnd},
		}

		falseProcessor := embedding.NewProcessorWithTransitions(docPath, config, falseTransitions)
		Expect(falseProcessor.Embed()).Error().Should(HaveOccurred())
	})

	It("should successfully embed with multi lined tag", func() {
		docPath := fmt.Sprintf("%s/multi-lined-tag.md", config.DocumentationRoot)
		processor := embedding.NewProcessor(docPath, config)
		Expect(processor.Embed()).Error().ShouldNot(HaveOccurred())

		Expect(processor.IsUpToDate()).Should(BeTrue())
	})

	It("should embed directly from source", func() {
		docPath := fmt.Sprintf("%s/doc.md", config.DocumentationRoot)
		processor := embedding.NewProcessor(docPath, config)

		Expect(processor.Embed()).Error().ShouldNot(HaveOccurred())

		Expect(processor.IsUpToDate()).Should(BeTrue())
	})

	It("should report files that are not up to date", func() {
		config.DocIncludes = []string{"doc.md"}
		docPath := fmt.Sprintf("%s/doc.md", config.DocumentationRoot)

		Expect(embedding.CheckUpToDate(config)).Should(ContainElement(docPath))
	})

	It("should ignore embed-code samples inside markdown code fences", func() {
		docPath := fmt.Sprintf("%s/embed-code-sample-in-fence.md", config.DocumentationRoot)
		processor := embedding.NewProcessor(docPath, config)

		Expect(processor.Embed()).Error().ShouldNot(HaveOccurred())
		Expect(processor.IsUpToDate()).Should(BeTrue())
	})

	It("should detect markdown fences by triple-or-more backticks only", func() {
		docPath := fmt.Sprintf("%s/triple-backticks-only-fence.md", config.DocumentationRoot)
		processor := embedding.NewProcessor(docPath, config)

		Expect(processor.Embed()).Error().ShouldNot(HaveOccurred())

		docContent, err := os.ReadFile(docPath)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(strings.Count(string(docContent), "System.out.println(\"Hello world\");")).
			Should(Equal(2))
		Expect(processor.IsUpToDate()).Should(BeTrue())
	})

	It("should report all check errors", func() {
		config.DocIncludes = []string{"missing-closing-tag.md", "unclosed-nested-tag.md"}

		var recovered any
		func() {
			defer func() {
				recovered = recover()
			}()
			embedding.CheckUpToDate(config)
		}()

		Expect(recovered).ShouldNot(BeNil())
		Expect(fmt.Sprint(recovered)).Should(And(
			ContainSubstring("missing-closing-tag.md"),
			ContainSubstring("the `<embed-code>` tag is not closed"),
			ContainSubstring("unclosed-nested-tag.md"),
			ContainSubstring("element <unexpected> closed by </embed-code>"),
		))
	})

	It("should embed with multi lined tag attributes", func() {
		docPath := fmt.Sprintf("%s/multi-lined-valid-tag-attributes.md", config.DocumentationRoot)
		processor := embedding.NewProcessor(docPath, config)
		Expect(processor.Embed()).Error().ShouldNot(HaveOccurred())

		Expect(processor.IsUpToDate()).Should(BeTrue())
	})

	It("should report a missing closing tag", func() {
		docPath := fmt.Sprintf("%s/missing-closing-tag.md", config.DocumentationRoot)
		processor := embedding.NewProcessor(docPath, config)

		_, err := processor.Embed()

		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(ContainSubstring(
			"missing-closing-tag.md:3`: " +
				"failed to parse an embedding instruction: " +
				"the `<embed-code>` tag is not closed",
		))
	})

	It("should report the XML parser error", func() {
		docPath := fmt.Sprintf("%s/unclosed-nested-tag.md", config.DocumentationRoot)
		processor := embedding.NewProcessor(docPath, config)

		_, err := processor.Embed()

		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(ContainSubstring(
			"unclosed-nested-tag.md:3`: " +
				"failed to parse an embedding instruction: " +
				"element <unexpected> closed by </embed-code>",
		))
	})

	// TODO:olena-zmiiova:https://github.com/SpineEventEngine/embed-code/issues/65
	It("should successfully embed to a file in a nested dir", func() {
		Skip(
			"Temporarily disabled, see " +
				"[issue #65](https://github.com/SpineEventEngine/embed-code/issues/65).",
		)
		docPath := fmt.Sprintf("%s/nested-dir-1/nested-dir-2/nested-dir-doc.md",
			config.DocumentationRoot)
		processor := embedding.NewProcessor(docPath, config)

		Expect(func() {
			embedding.EmbedAll(config)
		}).NotTo(Panic())

		Expect(processor.IsUpToDate()).Should(BeTrue())
	})

	It("should not embed to a file matched the `doc-excludes` pattern", func() {
		config.DocExcludes = []string{"**/excluded-doc.*"}

		docPath := fmt.Sprintf("%s/excluded-doc.md", config.DocumentationRoot)
		processor := embedding.NewProcessor(docPath, config)

		Expect(processor.Embed()).Error().ShouldNot(HaveOccurred())
		Expect(processor.IsUpToDate()).Should(BeTrue())
	})
})

func buildConfigWithSourceFiles() configuration.Configuration {
	var config = configuration.NewConfiguration()
	config.DocumentationRoot = temporaryTestDir
	config.CodeRoots = _type.NamedPathList{_type.NamedPath{Path: "../test/resources/code/java"}}

	return config
}

func copyDirRecursive(sourceDirPath string, targetDirPath string) {
	info, err := os.Stat(sourceDirPath)
	if err != nil {
		panic(err)
	}

	err = os.MkdirAll(targetDirPath, info.Mode())
	if err != nil {
		panic(err)
	}

	entries, err := os.ReadDir(sourceDirPath)
	if err != nil {
		panic(err)
	}

	for _, entry := range entries {
		sourcePath := filepath.Join(sourceDirPath, entry.Name())
		targetPath := filepath.Join(targetDirPath, entry.Name())

		if entry.IsDir() {
			copyDirRecursive(sourcePath, targetPath)
		} else {
			err = copyFile(sourcePath, targetPath)
			if err != nil {
				panic(err)
			}
		}
	}
}

func copyFile(sourceFilePath string, targetFilePath string) (err error) {
	sourceFile, err := os.Open(sourceFilePath)
	if err != nil {
		Fail(err.Error())
	}

	defer func(sourceFile *os.File) {
		err = sourceFile.Close()
		if err != nil {
			Fail(err.Error())
		}
	}(sourceFile)

	targetFile, err := os.Create(targetFilePath)
	if err != nil {
		return
	}
	defer func() {
		err = targetFile.Close()
		if err != nil {
			Fail(err.Error())
		}
	}()

	if _, err = io.Copy(targetFile, sourceFile); err != nil {
		return
	}

	err = os.Chmod(targetFilePath, os.FileMode(files.WritePermission))

	return
}
