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
	"os"
	"path/filepath"
	"strings"

	"embed-code/embed-code-go/configuration"
	"embed-code/embed-code-go/embedding"
	_type "embed-code/embed-code-go/type"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Orchestration", func() {
	It("should share resolver cache across documentation files in one operation", func() {
		documentationRoot := GinkgoT().TempDir()
		config := configuration.NewConfiguration()
		config.DocumentationRoot = documentationRoot
		config.CodeRoots = _type.NamedPathList{_type.NamedPath{Path: documentationRoot}}
		config.DocIncludes = []string{"source.md", "second.md"}
		sourceDoc := filepath.ToSlash(filepath.Join(documentationRoot, "source.md"))
		secondDoc := filepath.ToSlash(filepath.Join(documentationRoot, "second.md"))
		writeSourceEmbeddingDoc(sourceDoc)
		writeEmbeddingDoc(secondDoc)

		_, err := embedding.EmbedAll(config)

		Expect(err).ShouldNot(HaveOccurred())
		secondDocContent, err := os.ReadFile(secondDoc)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(strings.Count(string(secondDocContent), "original source line")).
			Should(Equal(1))
	})
})

// writeSourceEmbeddingDoc writes a source file that also acts as a target document.
func writeSourceEmbeddingDoc(path string) {
	Expect(os.WriteFile(
		path,
		[]byte("# Source\n\noriginal source line\n\n<embed-code file=\"source.md\"/>\n```md\n```\n"),
		0600,
	)).To(Succeed())
}

// writeEmbeddingDoc writes a target documentation file with one whole-file embedding.
func writeEmbeddingDoc(path string) {
	Expect(os.WriteFile(
		path,
		[]byte("# Second\n\n<embed-code file=\"source.md\"/>\n```md\n```\n"),
		0600,
	)).To(Succeed())
}
