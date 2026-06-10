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
)

// TestShowcasePositiveFlow verifies the showcase docs can be checked, detected as stale,
// repaired with embed mode, and checked again without changing repository files.
func TestShowcasePositiveFlow(t *testing.T) {
	repoRoot := findRepoRoot(t)
	docsRoot := copyShowcaseDocs(t, repoRoot, "docs")
	configPath := writeShowcaseConfig(t, repoRoot, docsRoot)

	checkOutput, err := runEmbedCode(t, repoRoot, "check", configPath)
	if err != nil {
		t.Fatalf("expected positive showcase check to pass:\n%s", checkOutput)
	}

	staleDoc := filepath.Join(docsRoot, "01-whole-file-source.md")
	replaceInFile(t, staleDoc, "package org.showcase;", "package stale.showcase;")

	staleOutput, err := runEmbedCode(t, repoRoot, "check", configPath)
	if err == nil {
		t.Fatalf("expected stale showcase check to fail:\n%s", staleOutput)
	}
	assertOutputContains(t, staleOutput, "File to update:")
	assertOutputContains(t, staleOutput, "01-whole-file-source.md")

	embedOutput, err := runEmbedCode(t, repoRoot, "embed", configPath)
	if err != nil {
		t.Fatalf("expected positive showcase embed to repair stale doc:\n%s", embedOutput)
	}
	assertOutputContains(t, embedOutput, "Embedding process finished.")

	finalOutput, err := runEmbedCode(t, repoRoot, "check", configPath)
	if err != nil {
		t.Fatalf("expected positive showcase check to pass after embed:\n%s", finalOutput)
	}
}

// TestShowcaseNegativeScenarios verifies each negative document fails with its expected reason.
func TestShowcaseNegativeScenarios(t *testing.T) {
	repoRoot := findRepoRoot(t)

	cases := []struct {
		name     string
		doc      string
		sources  []namedSource
		expected []string
	}{
		{
			name: "missing source",
			doc:  "missing-source.md",
			sources: []namedSource{
				{name: "java", path: filepath.Join(repoRoot, "examples", "showcase", "code", "java")},
			},
			expected: []string{
				"code file `$java/org/showcase/DoesNotExist.java",
				"not found",
			},
		},
		{
			name: "missing fragment",
			doc:  "missing-fragment.md",
			sources: []namedSource{
				{name: "java", path: filepath.Join(repoRoot, "examples", "showcase", "code", "java")},
			},
			expected: []string{
				"fragment `does not exist`",
				"not found",
			},
		},
		{
			name: "missing pattern",
			doc:  "missing-pattern.md",
			sources: []namedSource{
				{name: "java", path: filepath.Join(repoRoot, "examples", "showcase", "code", "java")},
			},
			expected: []string{
				"matches the line pattern",
				"doesNotExistPattern",
			},
		},
		{
			name: "invalid attributes",
			doc:  "invalid-attributes.md",
			sources: []namedSource{
				{name: "java", path: filepath.Join(repoRoot, "examples", "showcase", "code", "java")},
			},
			expected: []string{
				"must NOT specify both a fragment name and start/end/line patterns",
			},
		},
		{
			name: "missing code fence",
			doc:  "missing-code-fence.md",
			sources: []namedSource{
				{name: "java", path: filepath.Join(repoRoot, "examples", "showcase", "code", "java")},
			},
			expected: []string{
				"expected a markdown code fence after the embedding instruction",
			},
		},
		{
			name: "unclosed code fence",
			doc:  "unclosed-code-fence.md",
			sources: []namedSource{
				{name: "java", path: filepath.Join(repoRoot, "examples", "showcase", "code", "java")},
			},
			expected: []string{
				"the markdown code fence after the embedding instruction is not closed",
			},
		},
		{
			name: "stale snippet",
			doc:  "stale-snippet.md",
			sources: []namedSource{
				{name: "java", path: filepath.Join(repoRoot, "examples", "showcase", "code", "java")},
			},
			expected: []string{
				"File to update:",
				"stale-snippet.md",
				"the documentation files are not up-to-date with code files",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			docsRoot := copyShowcaseDocs(t, repoRoot, filepath.Join("negative", "docs"))
			configPath := writeSingleDocConfig(t, docsRoot, tc.sources, tc.doc)

			output, err := runEmbedCode(t, repoRoot, "check", configPath)
			if err == nil {
				t.Fatalf("expected negative scenario to fail:\n%s", output)
			}
			for _, expected := range tc.expected {
				assertOutputContains(t, output, expected)
			}
		})
	}
}

// TestShowcaseConfigurationExamples verifies the runnable configuration examples.
func TestShowcaseConfigurationExamples(t *testing.T) {
	repoRoot := findRepoRoot(t)

	configs := []string{
		"root-source.yml",
		"single-source.yml",
		"named-sources.yml",
		"include-exclude.yml",
		"multiple-embeddings.yml",
	}

	for _, config := range configs {
		t.Run(config, func(t *testing.T) {
			configPath := filepath.Join("examples", "showcase", "configuration", config)
			output, err := runEmbedCode(t, repoRoot, "check", configPath)
			if err != nil {
				t.Fatalf("expected configuration example to pass:\n%s", output)
			}
		})
	}
}

type namedSource struct {
	name string
	path string
}

// findRepoRoot returns the repository root by walking up from this test file.
func findRepoRoot(t *testing.T) string {
	t.Helper()

	_, filePath, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("could not locate showcase test file")
	}

	return filepath.Clean(filepath.Join(filepath.Dir(filePath), "..", ".."))
}

// copyShowcaseDocs copies one showcase documentation folder to a temporary test directory.
func copyShowcaseDocs(t *testing.T, repoRoot string, relativeSource string) string {
	t.Helper()

	sourceRoot := filepath.Join(repoRoot, "examples", "showcase", relativeSource)
	targetRoot := filepath.Join(t.TempDir(), "docs")
	copyDir(t, sourceRoot, targetRoot)

	return targetRoot
}

// copyDir recursively copies a directory tree while preserving regular file permissions.
func copyDir(t *testing.T, sourceRoot string, targetRoot string) {
	t.Helper()

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
	if err != nil {
		t.Fatalf("failed to copy showcase docs: %v", err)
	}
}

// writeShowcaseConfig creates a temp config that points at copied positive docs.
func writeShowcaseConfig(t *testing.T, repoRoot string, docsRoot string) string {
	t.Helper()

	return writeConfig(t, docsRoot, []namedSource{
		{name: "java", path: filepath.Join(repoRoot, "examples", "showcase", "code", "java")},
		{name: "kotlin", path: filepath.Join(repoRoot, "examples", "showcase", "code", "kotlin")},
		{name: "text", path: filepath.Join(repoRoot, "examples", "showcase", "code", "text")},
	}, []string{"**/*.md", "**/*.html"}, []string{"ignored-by-exclude.md"})
}

// writeSingleDocConfig creates a temp config for one negative showcase document.
func writeSingleDocConfig(
	t *testing.T,
	docsRoot string,
	sources []namedSource,
	docInclude string,
) string {
	t.Helper()

	return writeConfig(t, docsRoot, sources, []string{docInclude}, nil)
}

// writeConfig writes a YAML config with absolute source and documentation paths.
func writeConfig(
	t *testing.T,
	docsRoot string,
	sources []namedSource,
	includes []string,
	excludes []string,
) string {
	t.Helper()

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
	if len(excludes) > 0 {
		builder.WriteString("doc-excludes:\n")
		for _, exclude := range excludes {
			builder.WriteString(fmt.Sprintf("  - %q\n", exclude))
		}
	}
	builder.WriteString("separator: \"// ...\"\n")

	configPath := filepath.Join(t.TempDir(), "embed-code.yml")
	if err := os.WriteFile(configPath, []byte(builder.String()), 0o644); err != nil {
		t.Fatalf("failed to write temp config: %v", err)
	}

	return configPath
}

// runEmbedCode executes the CLI through `go run` and returns combined output.
func runEmbedCode(t *testing.T, repoRoot string, mode string, configPath string) (string, error) {
	t.Helper()

	cmd := exec.Command("go", "run", "./main.go", "-mode", mode, "-config-path", configPath)
	cmd.Dir = repoRoot
	output, err := cmd.CombinedOutput()

	return string(output), err
}

// replaceInFile replaces one expected substring in a copied documentation file.
func replaceInFile(t *testing.T, path string, oldText string, newText string) {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read %s: %v", path, err)
	}
	content := string(data)
	if !strings.Contains(content, oldText) {
		t.Fatalf("expected %s to contain %q", path, oldText)
	}
	content = strings.Replace(content, oldText, newText, 1)
	if err = os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write %s: %v", path, err)
	}
}

// assertOutputContains fails the test when a command output does not include a substring.
func assertOutputContains(t *testing.T, output string, expected string) {
	t.Helper()

	if !strings.Contains(output, expected) {
		t.Fatalf("expected output to contain %q:\n%s", expected, output)
	}
}
