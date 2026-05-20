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

package commentfilter

import (
	"bytes"
	"log/slog"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// TestCommentFilter runs the comment filter test suite.
func TestCommentFilter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Comment Filter Suite")
}

var _ = Describe("Comment filter", func() {
	Describe("YAML", func() {
		It("should strip all comments", func() {
			lines := []string{
				"name: test # inline",
				"# standalone",
				"value: \"# literal\"",
			}

			expected := []string{
				"name: test ",
				"value: \"# literal\"",
			}

			assertFiltered("config.yml", RetainNone, lines, expected)
		})
	})

	Describe("XML", func() {
		It("should strip all comments", func() {
			lines := []string{
				"<root>",
				"  <!-- hidden -->",
				"  <item name=\"<!-- literal -->\"/>",
				"</root>",
			}

			expected := []string{
				"<root>",
				"  <item name=\"<!-- literal -->\"/>",
				"</root>",
			}

			assertFiltered("layout.xml", RetainNone, lines, expected)
		})
	})

	Describe("Java-style languages", func() {
		It("should keep documentation comments", func() {
			lines := []string{
				"/** API docs. */",
				"// implementation note",
				"fun call() = \"// literal\"",
			}

			expected := []string{
				"/** API docs. */",
				"fun call() = \"// literal\"",
			}

			assertFiltered("api.kt", RetainDocumentation, lines, expected)
		})

		It("should keep block comments", func() {
			lines := []string{
				"/** API docs. */",
				"/* implementation note */",
				"String create();",
			}

			expected := []string{
				"/* implementation note */",
				"String create();",
			}

			assertFiltered("Api.java", RetainBlock, lines, expected)
		})

		It("should keep regular comments", func() {
			lines := []string{
				"/** API docs. */",
				"/* implementation note */",
				"String create(); // inline note",
			}

			expected := []string{
				"/* implementation note */",
				"String create(); // inline note",
			}

			assertFiltered("Api.java", RetainRegular, lines, expected)
		})
	})

	Describe("JavaScript and TypeScript", func() {
		It("should strip comments without treating template literals as comments", func() {
			lines := []string{
				"// module comment",
				"const url = `http://example.org/*not-comment*/`;",
				"const value = 42; // inline comment",
			}

			expected := []string{
				"const url = `http://example.org/*not-comment*/`;",
				"const value = 42; ",
			}

			assertFiltered("sample.ts", RetainNone, lines, expected)
		})
	})

	Describe("C#", func() {
		It("should keep XML documentation comments", func() {
			lines := []string{
				"/// <summary>Creates a value.</summary>",
				"// implementation note",
				"public string Create() => \"// literal\";",
			}

			expected := []string{
				"/// <summary>Creates a value.</summary>",
				"public string Create() => \"// literal\";",
			}

			assertFiltered("Api.cs", RetainDocumentation, lines, expected)
		})

		It("should keep inline comments", func() {
			lines := []string{
				"/// <summary>Creates a value.</summary>",
				"// implementation note",
				"public string Create() => \"// literal\";",
			}

			expected := []string{
				"// implementation note",
				"public string Create() => \"// literal\";",
			}

			assertFiltered("Api.cs", RetainInline, lines, expected)
		})
	})

	Describe("C and C++", func() {
		It("should strip all comments without treating literals as comments", func() {
			lines := []string{
				"// header comment",
				"#include <stdio.h>",
				"",
				"/* block comment */",
				"const char slash = '/';",
				"const char* url = \"http://example.org\";",
				"int create() { return 1; } // inline comment",
			}

			expected := []string{
				"#include <stdio.h>",
				"",
				"const char slash = '/';",
				"const char* url = \"http://example.org\";",
				"int create() { return 1; } ",
			}

			assertFiltered("sample.cpp", RetainNone, lines, expected)
		})

		It("should keep inline comments", func() {
			lines := []string{
				"// header comment",
				"int create();",
				"/* block comment */",
				"int count(); // inline comment",
			}

			expected := []string{
				"// header comment",
				"int create();",
				"int count(); // inline comment",
			}

			assertFiltered("sample.cpp", RetainInline, lines, expected)
		})

		It("should keep block comments", func() {
			lines := []string{
				"// header comment",
				"int create();",
				"/* block comment */",
				"int count(); // inline comment",
			}

			expected := []string{
				"int create();",
				"/* block comment */",
				"int count(); ",
			}

			assertFiltered("sample.hpp", RetainBlock, lines, expected)
		})
	})

	Describe("Go", func() {
		It("should strip all comments without treating literals as comments", func() {
			lines := []string{
				"// package comment",
				"package sample",
				"",
				"/* block comment */",
				"const slash = '/'",
				"const url = \"http://example.org\"",
				"const raw = `/* not a comment */`",
				"func create() {} // inline comment",
			}

			expected := []string{
				"package sample",
				"",
				"const slash = '/'",
				"const url = \"http://example.org\"",
				"const raw = `/* not a comment */`",
				"func create() {} ",
			}

			assertFiltered("sample.go", RetainNone, lines, expected)
		})

		It("should keep inline comments", func() {
			lines := []string{
				"// package comment",
				"package sample",
				"/* block comment */",
				"func create() {} // inline comment",
			}

			expected := []string{
				"// package comment",
				"package sample",
				"func create() {} // inline comment",
			}

			assertFiltered("sample.go", RetainInline, lines, expected)
		})

		It("should keep block comments", func() {
			lines := []string{
				"// package comment",
				"package sample",
				"/* block comment */",
				"func create() {} // inline comment",
			}

			expected := []string{
				"package sample",
				"/* block comment */",
				"func create() {} ",
			}

			assertFiltered("sample.go", RetainBlock, lines, expected)
		})
	})

	Describe("Protobuf", func() {
		It("should strip all comments without treating literals as comments", func() {
			lines := []string{
				"// file comment",
				"syntax = \"proto3\";",
				"",
				"/* message comment */",
				"message Sample {",
				"  string url = 1 [default = 'http://example.org'];",
				"  int32 count = 2; // inline comment",
				"}",
			}

			expected := []string{
				"syntax = \"proto3\";",
				"",
				"message Sample {",
				"  string url = 1 [default = 'http://example.org'];",
				"  int32 count = 2; ",
				"}",
			}

			assertFiltered("sample.proto", RetainNone, lines, expected)
		})

		It("should keep inline comments", func() {
			lines := []string{
				"// file comment",
				"syntax = \"proto3\";",
				"/* message comment */",
				"message Sample {} // inline comment",
			}

			expected := []string{
				"// file comment",
				"syntax = \"proto3\";",
				"message Sample {} // inline comment",
			}

			assertFiltered("sample.proto", RetainInline, lines, expected)
		})

		It("should keep block comments", func() {
			lines := []string{
				"// file comment",
				"syntax = \"proto3\";",
				"/* message comment */",
				"message Sample {} // inline comment",
			}

			expected := []string{
				"syntax = \"proto3\";",
				"/* message comment */",
				"message Sample {} ",
			}

			assertFiltered("sample.proto", RetainBlock, lines, expected)
		})
	})

	Describe("Python", func() {
		It("should strip all comments", func() {
			lines := []string{
				"# module comment",
				"name = 'hash # literal'",
				"value = 1 # inline comment",
			}

			expected := []string{
				"name = 'hash # literal'",
				"value = 1 ",
			}

			assertFiltered("module.py", RetainNone, lines, expected)
		})
	})

	Describe("Visual Basic", func() {
		It("should strip all comments", func() {
			lines := []string{
				"' file comment",
				"REM module comment",
				"Dim text = \"REM not a comment\"",
				"Dim value = 1 ' inline",
				"Dim ready = True : Rem after statement separator",
				"Dim reminder = 1",
			}

			expected := []string{
				"Dim text = \"REM not a comment\"",
				"Dim value = 1 ",
				"Dim ready = True : ",
				"Dim reminder = 1",
			}

			assertFiltered("Module.vb", RetainNone, lines, expected)
		})

		It("should keep regular comments", func() {
			lines := []string{
				"''' <summary>Creates a value.</summary>",
				"' file comment",
				"REM module comment",
				"Dim value = 1 ' inline",
			}

			expected := []string{
				"' file comment",
				"REM module comment",
				"Dim value = 1 ' inline",
			}

			assertFiltered("Module.vb", RetainRegular, lines, expected)
		})

		It("should keep documentation comments", func() {
			lines := []string{
				"''' <summary>Creates a value.</summary>",
				"' implementation note",
				"REM module comment",
				"Public Function Create() As String",
			}

			expected := []string{
				"''' <summary>Creates a value.</summary>",
				"Public Function Create() As String",
			}

			assertFiltered("Module.vb", RetainDocumentation, lines, expected)
		})
	})

	Describe("unsupported extensions", func() {
		It("should return unsupported files unchanged", func() {
			lines := []string{
				"# docs",
				"sub call { } # inline",
			}

			assertFiltered("service.pl", RetainAll, lines, lines)
		})

		It("should warn about unsupported comment modes", func() {
			output := captureWarnings(func() {
				Filter([]string{"# comment"}, "service.pl", RetainNone, "docs/guide.md", 12)
			})

			Expect(output).Should(ContainSubstring(
				"comment filtering is not supported for this file extension",
			))
			Expect(output).Should(ContainSubstring("file://"))
			Expect(output).Should(ContainSubstring("guide.md:12"))
		})
	})

	Describe("warnings", func() {
		It("should warn about modes without language-specific meaning", func() {
			output := captureWarnings(func() {
				Filter([]string{"<!-- comment -->"}, "layout.xml", RetainDocumentation, "docs/guide.md", 12)
			})

			Expect(output).Should(ContainSubstring("documentation"))
			Expect(output).Should(ContainSubstring("layout.xml"))
			Expect(output).Should(ContainSubstring("file://"))
			Expect(output).Should(ContainSubstring("guide.md:12"))
			Expect(output).Should(ContainSubstring("does not have a distinct meaning"))
		})
	})
})

// assertFiltered verifies filtering output for one file path and mode.
func assertFiltered(
	filePath string,
	mode Mode,
	lines []string,
	expected []string,
) {
	got := Filter(lines, filePath, mode, "docs/guide.md", 12)

	Expect(got).Should(Equal(expected))
}

// captureWarnings runs action and returns slog warning output.
func captureWarnings(action func()) string {
	var output bytes.Buffer
	previous := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&output, &slog.HandlerOptions{
		Level: slog.LevelWarn,
	})))
	defer slog.SetDefault(previous)

	action()

	return output.String()
}
