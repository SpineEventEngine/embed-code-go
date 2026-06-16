//go:build showcase

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

package showcase_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// TestShowcase runs the showcase example suite.
func TestShowcase(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Showcase Suite")
}

var _ = Describe("Showcase", func() {
	var repoRoot string

	BeforeEach(func() {
		repoRoot = findRepoRoot()
	})

	Describe("embedding examples", func() {
		It("should check, detect staleness, embed, and recheck positive examples", func() {
			docsRoot := copyShowcaseDocs(repoRoot, filepath.Join("embedding", "positive"))
			configPath := writeShowcaseConfig(repoRoot, docsRoot)

			checkOutput, err := runEmbedCode(repoRoot, "check", configPath)
			Expect(err).ShouldNot(HaveOccurred(), "expected positive showcase check to pass:\n%s", checkOutput)

			staleDoc := filepath.Join(docsRoot, "whole-file-source.md")
			replaceInFile(staleDoc, "package org.showcase;", "package stale.showcase;")

			staleOutput, err := runEmbedCode(repoRoot, "check", configPath)
			Expect(err).Should(HaveOccurred(), "expected stale showcase check to fail:\n%s", staleOutput)
			Expect(staleOutput).Should(ContainSubstring("File to update:"))
			Expect(staleOutput).Should(ContainSubstring("whole-file-source.md"))

			embedOutput, err := runEmbedCode(repoRoot, "embed", configPath)
			Expect(err).ShouldNot(HaveOccurred(), "expected positive showcase embed to repair stale doc:\n%s", embedOutput)
			Expect(embedOutput).Should(ContainSubstring("Embedding process finished."))

			finalOutput, err := runEmbedCode(repoRoot, "check", configPath)
			Expect(err).ShouldNot(HaveOccurred(), "expected positive showcase check to pass after embed:\n%s", finalOutput)
		})

		Describe("negative examples", func() {
			for _, tc := range negativeShowcaseCases() {
				tc := tc

				It("should report "+tc.name, func() {
					docsRoot := copyShowcaseDocs(repoRoot, filepath.Join("embedding", "negative", "docs"))
					configPath := writeSingleDocConfig(
						docsRoot,
						[]namedSource{javaSource(repoRoot)},
						tc.doc,
					)

					output, err := runEmbedCode(repoRoot, "check", configPath)
					Expect(err).Should(HaveOccurred(), "expected negative scenario to fail:\n%s", output)
					for _, expected := range tc.expected {
						Expect(output).Should(ContainSubstring(expected))
					}
				})
			}
		})
	})

	Describe("configuration examples", func() {
		for _, config := range []string{
			"single-source.yml",
			"named-sources.yml",
			"include-exclude.yml",
			"multiple-embeddings.yml",
		} {
			config := config

			It("should check "+config, func() {
				configPath := filepath.Join("showcase", "configuration", config)

				output, err := runEmbedCode(repoRoot, "check", configPath)
				Expect(err).ShouldNot(HaveOccurred(), "expected configuration example to pass:\n%s", output)
			})
		}
	})
})

// negativeShowcaseCase describes one intentionally broken showcase document.
type negativeShowcaseCase struct {
	name     string
	doc      string
	expected []string
}

// namedSource is the named code source path.
type namedSource struct {
	name string
	path string
}

// negativeShowcaseCases returns the expected failures for the broken embedding examples.
func negativeShowcaseCases() []negativeShowcaseCase {
	return []negativeShowcaseCase{
		{
			name: "missing source",
			doc:  "missing-source.md",
			expected: []string{
				"code file `$java/org/showcase/DoesNotExist.java",
				"not found",
			},
		},
		{
			name: "missing fragment",
			doc:  "missing-fragment.md",
			expected: []string{
				"fragment `does not exist`",
				"not found",
			},
		},
		{
			name: "missing pattern",
			doc:  "missing-pattern.md",
			expected: []string{
				"matches the line pattern",
				"doesNotExistPattern",
			},
		},
		{
			name: "invalid attributes",
			doc:  "invalid-attributes.md",
			expected: []string{
				"must NOT specify both a fragment name and start/end/line patterns",
			},
		},
		{
			name: "missing code fence",
			doc:  "missing-code-fence.md",
			expected: []string{
				"expected a markdown code fence after the embedding instruction",
			},
		},
		{
			name: "unclosed code fence",
			doc:  "unclosed-code-fence.md",
			expected: []string{
				"the markdown code fence after the embedding instruction is not closed",
			},
		},
		{
			name: "stale snippet",
			doc:  "stale-snippet.md",
			expected: []string{
				"File to update:",
				"stale-snippet.md",
				"the documentation files are not up-to-date with code files",
			},
		},
	}
}

// javaSource returns the Java showcase source root.
func javaSource(repoRoot string) namedSource {
	return namedSource{
		name: "java",
		path: filepath.Join(repoRoot, "showcase", "code", "java"),
	}
}

// findRepoRoot returns the repository root by walking up from this test file.
func findRepoRoot() string {
	GinkgoHelper()

	_, filePath, _, ok := runtime.Caller(0)
	Expect(ok).Should(BeTrue(), "could not locate showcase test file")

	return filepath.Clean(filepath.Join(filepath.Dir(filePath), ".."))
}

// copyShowcaseDocs copies one showcase documentation folder to a temporary test directory.
func copyShowcaseDocs(repoRoot string, relativeSource string) string {
	GinkgoHelper()

	sourceRoot := filepath.Join(repoRoot, "showcase", relativeSource)
	tempRoot, err := os.MkdirTemp("", "embed-code-showcase-docs-*")
	Expect(err).ShouldNot(HaveOccurred())
	DeferCleanup(os.RemoveAll, tempRoot)

	targetRoot := filepath.Join(tempRoot, "docs")
	copyDir(sourceRoot, targetRoot)

	return targetRoot
}

// copyDir recursively copies a directory tree while preserving regular file permissions.
func copyDir(sourceRoot string, targetRoot string) {
	GinkgoHelper()

	err := filepath.WalkDir(sourceRoot, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		relativePath, err := filepath.Rel(sourceRoot, path)
		if err != nil {
			return err
		}
		targetPath := filepath.Join(targetRoot, relativePath)
		if entry.IsDir() {
			return os.MkdirAll(targetPath, 0o755)
		}
		if !entry.Type().IsRegular() {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}

		return os.WriteFile(targetPath, data, info.Mode())
	})
	Expect(err).ShouldNot(HaveOccurred(), "failed to copy showcase docs")
}

// writeShowcaseConfig creates a temp config that points at copied positive docs.
func writeShowcaseConfig(repoRoot string, docsRoot string) string {
	GinkgoHelper()

	return writeConfig(docsRoot, []namedSource{
		javaSource(repoRoot),
		{name: "kotlin", path: filepath.Join(repoRoot, "showcase", "code", "kotlin")},
		{name: "text", path: filepath.Join(repoRoot, "showcase", "code", "text")},
	}, []string{"**/*.md", "**/*.html"})
}

// writeSingleDocConfig creates a temp config for one negative showcase document.
func writeSingleDocConfig(
	docsRoot string,
	sources []namedSource,
	docInclude string,
) string {
	GinkgoHelper()

	return writeConfig(docsRoot, sources, []string{docInclude})
}

// writeConfig writes a YAML config with absolute source and documentation paths.
func writeConfig(
	docsRoot string,
	sources []namedSource,
	includes []string,
) string {
	GinkgoHelper()

	var builder strings.Builder
	builder.WriteString("code-path:\n")
	for _, source := range sources {
		builder.WriteString(fmt.Sprintf("  - name: %s\n", source.name))
		builder.WriteString(fmt.Sprintf("    path: %s\n", filepath.ToSlash(source.path)))
	}
	builder.WriteString(fmt.Sprintf("docs-path: %s\n", filepath.ToSlash(docsRoot)))
	builder.WriteString("doc-includes:\n")
	for _, include := range includes {
		builder.WriteString(fmt.Sprintf("  - %q\n", include))
	}
	builder.WriteString("separator: \"// ...\"\n")

	tempRoot, err := os.MkdirTemp("", "embed-code-showcase-config-*")
	Expect(err).ShouldNot(HaveOccurred())
	DeferCleanup(os.RemoveAll, tempRoot)

	configPath := filepath.Join(tempRoot, "embed-code.yml")
	Expect(os.WriteFile(configPath, []byte(builder.String()), 0o644)).
		Should(Succeed(), "failed to write temp config")

	return configPath
}

// runEmbedCode executes the CLI through `go run` and returns combined output.
func runEmbedCode(repoRoot string, mode string, configPath string) (string, error) {
	GinkgoHelper()

	cmd := exec.Command("go", "run", "./main.go", "-mode="+mode, "-config-path="+configPath)
	cmd.Dir = repoRoot
	output, err := cmd.CombinedOutput()

	return string(output), err
}

// replaceInFile replaces one expected substring in a copied documentation file.
func replaceInFile(path string, oldText string, newText string) {
	GinkgoHelper()

	data, err := os.ReadFile(path)
	Expect(err).ShouldNot(HaveOccurred(), "failed to read %s", path)
	content := string(data)
	Expect(content).Should(ContainSubstring(oldText), "expected %s to contain %q", path, oldText)
	content = strings.Replace(content, oldText, newText, 1)
	Expect(os.WriteFile(path, []byte(content), 0o644)).
		Should(Succeed(), "failed to write %s", path)
}
