/*
 * Copyright 2026, TeamDev. All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 * Redistribution and use in source and/or binary forms, with or without
 * modification, must retain the above copyright notice and the following
 * disclaimer.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
 * "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
 * LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
 * A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
 * OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
 * SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
 * LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
 * DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
 * THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 * OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

package embedding_test

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"embed-code/embed-code-go/configuration"
	"embed-code/embed-code-go/embedding"
	"embed-code/embed-code-go/embedding/parsing"
	_type "embed-code/embed-code-go/type"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// TestEmbedding runs the embedding test suite.
func TestEmbedding(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Data Suite")
}

var _ = Describe("Embedding", func() {
	var config configuration.Configuration

	BeforeEach(func() {
		config = buildConfigWithSourceFiles(GinkgoT().TempDir())
		Expect(os.CopyFS(
			config.DocumentationRoot,
			os.DirFS("../test/resources/docs"),
		)).To(Succeed())
	})

	It("should be up to date", func() {
		docPath := testDocPath(config, "whole-file-fragment.md")
		processor := newProcessor(docPath, config)

		Expect(processor.Embed()).Error().ShouldNot(HaveOccurred())
		Expect(processor.IsUpToDate()).Should(BeTrue())
	})

	It("should be up to date as there is nothing to update", func() {
		docPath := testDocPath(config, "no-embedding-doc.md")
		processor := newProcessor(docPath, config)

		Expect(processor.Embed()).Error().ShouldNot(HaveOccurred())
		Expect(processor.IsUpToDate()).Should(BeTrue())
	})

	It("should successfully embed with multi lined tag", func() {
		docPath := testDocPath(config, "multi-lined-tag.md")
		processor := newProcessor(docPath, config)
		Expect(processor.Embed()).Error().ShouldNot(HaveOccurred())

		Expect(processor.IsUpToDate()).Should(BeTrue())
	})

	It("should embed directly from source", func() {
		docPath := testDocPath(config, "doc.md")
		processor := newProcessor(docPath, config)

		Expect(processor.Embed()).Error().ShouldNot(HaveOccurred())

		Expect(processor.IsUpToDate()).Should(BeTrue())
	})

	It("should report files that are not up to date", func() {
		config.DocIncludes = []string{"doc.md"}
		docPath := testDocPath(config, "doc.md")

		outdatedFiles, err := embedding.CheckUpToDate(config)

		Expect(err).ShouldNot(HaveOccurred())
		Expect(outdatedFiles).Should(ContainElement(docPath))
	})

	It("should ignore embed-code samples inside markdown code fences", func() {
		docPath := testDocPath(config, "embed-code-sample-in-fence.md")
		processor := newProcessor(docPath, config)

		Expect(processor.Embed()).Error().ShouldNot(HaveOccurred())
		Expect(processor.IsUpToDate()).Should(BeTrue())
	})

	It("should detect markdown fences by triple-or-more backticks only", func() {
		docPath := testDocPath(config, "triple-backticks-only-fence.md")
		processor := newProcessor(docPath, config)

		Expect(processor.Embed()).Error().ShouldNot(HaveOccurred())

		docContent, err := os.ReadFile(docPath)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(strings.Count(string(docContent), "System.out.println(\"Hello world\");")).
			Should(Equal(2))
		Expect(processor.IsUpToDate()).Should(BeTrue())
	})

	It("should report all check errors", func() {
		config.DocIncludes = []string{"missing-closing-tag.md", "unclosed-nested-tag.md"}

		_, err := embedding.CheckUpToDate(config)

		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(And(
			ContainSubstring("missing-closing-tag.md"),
			ContainSubstring(
				"the `<embed-code>` instruction is not closed; add `</embed-code>` or use `/>`",
			),
			ContainSubstring("unclosed-nested-tag.md"),
			ContainSubstring("element <unexpected> closed by </embed-code>"),
		))
		var processingErr embedding.ProcessingError
		Expect(errors.As(err, &processingErr)).Should(BeTrue())

		var parseErr parsing.InstructionParseError
		Expect(errors.As(err, &parseErr)).Should(BeTrue())
	})

	It("should report all pattern matching errors", func() {
		config.DocIncludes = []string{"missing-start-pattern.md", "missing-end-pattern.md"}

		_, err := embedding.CheckUpToDate(config)

		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(And(
			ContainSubstring("missing-start-pattern.md:3"),
			ContainSubstring(
				"no line in code file `file://",
			),
			ContainSubstring(
				"` matches the start pattern "+
					"`*doesNotExistStart*`",
			),
			ContainSubstring("missing-end-pattern.md:3"),
			ContainSubstring(
				"` matches the end pattern "+
					"`*doesNotExistEnd*`",
			),
		))
		var patternErr parsing.PatternNotFoundError
		Expect(errors.As(err, &patternErr)).Should(BeTrue())
		Expect(patternErr.Line).Should(Equal(3))
	})

	It("should embed with multi lined tag attributes", func() {
		docPath := testDocPath(config, "multi-lined-valid-tag-attributes.md")
		processor := newProcessor(docPath, config)
		Expect(processor.Embed()).Error().ShouldNot(HaveOccurred())

		Expect(processor.IsUpToDate()).Should(BeTrue())
	})

	It("should embed a method with escaped newline patterns", func() {
		config.DocIncludes = []string{"escaped-newline-pattern.md"}
		docPath := testDocPath(config, "escaped-newline-pattern.md")
		processor := newProcessor(docPath, config)

		Expect(processor.Embed()).Error().ShouldNot(HaveOccurred())

		docContent, err := os.ReadFile(docPath)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(string(docContent)).Should(ContainSubstring("@Test\n" +
			"@DisplayName(\"adds two values\")"))
		Expect(string(docContent)).Should(ContainSubstring("assertEquals(2, value);\n}"))
		Expect(string(docContent)).ShouldNot(ContainSubstring("subtractsTwoValues"))
	})

	It("should embed a method with exact escaped newline patterns", func() {
		config.DocIncludes = []string{"escaped-newline-exact-pattern.md"}
		docPath := testDocPath(config, "escaped-newline-exact-pattern.md")
		processor := newProcessor(docPath, config)

		Expect(processor.Embed()).Error().ShouldNot(HaveOccurred())

		docContent, err := os.ReadFile(docPath)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(string(docContent)).Should(ContainSubstring("@Test\n" +
			"@DisplayName(\"adds two values\")"))
		Expect(string(docContent)).Should(ContainSubstring("assertEquals(2, value);\n}"))
		Expect(string(docContent)).ShouldNot(ContainSubstring("subtractsTwoValues"))
	})

	It("should embed matching lines with an escaped newline line pattern", func() {
		config.DocIncludes = []string{"escaped-newline-line-pattern.md"}
		docPath := testDocPath(config, "escaped-newline-line-pattern.md")
		processor := newProcessor(docPath, config)

		Expect(processor.Embed()).Error().ShouldNot(HaveOccurred())

		docContent, err := os.ReadFile(docPath)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(string(docContent)).Should(ContainSubstring("@Test\n" +
			"@DisplayName(\"adds two values\")"))
		Expect(string(docContent)).ShouldNot(ContainSubstring("void addsTwoValues"))
		Expect(string(docContent)).ShouldNot(ContainSubstring("subtractsTwoValues"))
	})

	It("should embed a line with an escaped newline literal pattern", func() {
		config.DocIncludes = []string{"escaped-newline-literal-pattern.md"}
		docPath := testDocPath(config, "escaped-newline-literal-pattern.md")
		processor := newProcessor(docPath, config)

		Expect(processor.Embed()).Error().ShouldNot(HaveOccurred())

		docContent, err := os.ReadFile(docPath)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(string(docContent)).Should(ContainSubstring(
			"private static final String MY_STRING = \"\\n\";",
		))
	})

	It("should report a missing closing tag", func() {
		docPath := testDocPath(config, "missing-closing-tag.md")
		processor := newProcessor(docPath, config)

		_, err := processor.Embed()

		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(ContainSubstring(
			"missing-closing-tag.md:3`: " +
				"failed to parse an embedding instruction: " +
				"the `<embed-code>` instruction is not closed; add `</embed-code>` or use `/>`",
		))
	})

	It("should preserve typed parser errors after adding document context", func() {
		docPath := testDocPath(config, "missing-closing-tag.md")
		processor := newProcessor(docPath, config)

		_, err := processor.Embed()

		Expect(err).Should(HaveOccurred())
		var processingErr embedding.ProcessingError
		Expect(errors.As(err, &processingErr)).Should(BeTrue())
		Expect(processingErr.DocFilePath).Should(Equal(docPath))
		Expect(processingErr.Line).Should(Equal(3))

		var parseErr parsing.InstructionParseError
		Expect(errors.As(err, &parseErr)).Should(BeTrue())
		Expect(parseErr.Line).Should(Equal(3))
		Expect(parseErr.Reason).Should(Equal(
			"the `<embed-code>` instruction is not closed; add `</embed-code>` or use `/>`",
		))
	})

	It("should report the XML parser error", func() {
		docPath := testDocPath(config, "unclosed-nested-tag.md")
		processor := newProcessor(docPath, config)

		_, err := processor.Embed()

		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(ContainSubstring(
			"unclosed-nested-tag.md:3`: " +
				"failed to parse an embedding instruction: " +
				"element <unexpected> closed by </embed-code>",
		))
	})

	It("should report a missing code fence after the instruction", func() {
		docPath := testDocPath(config, "missing-code-fence.md")
		processor := newProcessor(docPath, config)

		_, err := processor.Embed()

		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(ContainSubstring(
			"missing-code-fence.md:3`: " +
				"expected a markdown code fence after the embedding instruction",
		))
	})

	It("should report a missing file attribute with documentation context", func() {
		docPath := testDocPath(config, "missing-file-attribute.md")
		Expect(os.WriteFile(
			docPath,
			[]byte("# Missing file attribute\n\n"+
				"<embed-code fragment=\"main()\"/>\n"+
				"```java\n"+
				"```\n"),
			0600,
		)).To(Succeed())
		processor := newProcessor(docPath, config)

		_, err := processor.Embed()

		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(ContainSubstring(
			"missing-file-attribute.md:3`: " +
				"failed to parse an embedding instruction: " +
				"<embed-code> must specify a non-empty `file` attribute",
		))
	})

	It("should report an unclosed code fence after the instruction", func() {
		docPath := testDocPath(config, "unclosed-code-fence.md")
		processor := newProcessor(docPath, config)

		_, err := processor.Embed()

		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(ContainSubstring(
			"unclosed-code-fence.md:3`: " +
				"the markdown code fence after the embedding instruction is not closed",
		))
	})

	It("should successfully embed to a file in a nested dir", func() {
		config.CodeRoots = _type.NamedPathList{_type.NamedPath{Path: "../test/resources/code/kotlin"}}
		config.DocIncludes = []string{"nested-dir-1/nested-dir-2/nested-dir-doc.md"}
		docPath := testDocPath(config, "nested-dir-1/nested-dir-2/nested-dir-doc.md")
		processor := newProcessor(docPath, config)

		result, err := embedding.EmbedAll(config)

		Expect(err).ShouldNot(HaveOccurred())
		Expect(result.UpdatedTargetFiles).Should(ContainElement(docPath))
		Expect(processor.IsUpToDate()).Should(BeTrue())
	})

	It("should not embed to a file matched the `doc-excludes` pattern", func() {
		config.DocExcludes = []string{"**/excluded-doc.*"}

		docPath := testDocPath(config, "excluded-doc.md")
		processor := newProcessor(docPath, config)

		context, err := processor.Embed()

		Expect(err).ShouldNot(HaveOccurred())
		Expect(context).ShouldNot(BeNil())
		Expect(context.EmbeddingsCount()).Should(Equal(0))
		Expect(context.IsContainsEmbedding()).Should(BeFalse())
		Expect(processor.IsUpToDate()).Should(BeTrue())
	})
})

// buildConfigWithSourceFiles builds an embedding config with an isolated documentation root.
func buildConfigWithSourceFiles(documentationRoot string) configuration.Configuration {
	var config = configuration.NewConfiguration()
	config.DocumentationRoot = documentationRoot
	config.CodeRoots = _type.NamedPathList{_type.NamedPath{Path: "../test/resources/code/java"}}

	return config
}

// testDocPath returns the normalized path to a copied documentation fixture.
func testDocPath(config configuration.Configuration, name string) string {
	return filepath.ToSlash(filepath.Join(config.DocumentationRoot, name))
}

func newProcessor(
	docPath string,
	config configuration.Configuration,
) embedding.Processor {
	processor, err := embedding.NewProcessor(docPath, config)

	Expect(err).ShouldNot(HaveOccurred())

	return processor
}
