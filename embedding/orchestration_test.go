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

package embedding

import (
	"os"
	"path/filepath"

	"embed-code/embed-code-go/configuration"
	_type "embed-code/embed-code-go/type"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Orchestration", func() {
	It("should share resolver cache across documentation files in one operation", func() {
		documentationRoot := GinkgoT().TempDir()
		sourceRoot := GinkgoT().TempDir()
		config := configuration.NewConfiguration()
		config.DocumentationRoot = documentationRoot
		config.CodeRoots = _type.NamedPathList{_type.NamedPath{Path: sourceRoot}}
		config.DocIncludes = []string{"first.md", "second.md"}
		sourcePath := filepath.Join(sourceRoot, "Example.java")
		firstDoc := filepath.ToSlash(filepath.Join(documentationRoot, "first.md"))
		secondDoc := filepath.ToSlash(filepath.Join(documentationRoot, "second.md"))
		writeEmbeddingDoc(firstDoc)
		writeEmbeddingDoc(secondDoc)
		writeSource(sourcePath, "class Example { String version = \"first\"; }")

		_, processingErrors := processRequiredDocs(config, func(
			docFilePath string,
			processor Processor,
		) error {
			_, err := processor.Embed()
			if err != nil {
				return err
			}
			if docFilePath == firstDoc {
				writeSource(sourcePath, "class Example { String version = \"second\"; }")
			}

			return nil
		})

		Expect(processingErrors).Should(BeEmpty())
		secondDocContent, err := os.ReadFile(secondDoc)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(string(secondDocContent)).Should(ContainSubstring(
			"class Example { String version = \"first\"; }",
		))
		Expect(string(secondDocContent)).ShouldNot(ContainSubstring(
			"class Example { String version = \"second\"; }",
		))
	})
})

// writeEmbeddingDoc writes a target documentation file with one whole-file embedding.
func writeEmbeddingDoc(path string) {
	Expect(os.WriteFile(
		path,
		[]byte("<embed-code file=\"Example.java\"/>\n```java\n```\n"),
		0600,
	)).To(Succeed())
}

// writeSource writes source content used by the embedding resolver.
func writeSource(path string, content string) {
	Expect(os.WriteFile(path, []byte(content), 0600)).To(Succeed())
}
