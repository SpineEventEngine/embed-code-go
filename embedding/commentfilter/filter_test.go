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
	"reflect"
	"strings"
	"testing"
)

// TestFilterYaml verifies YAML line comment filtering.
func TestFilterYaml(t *testing.T) {
	lines := []string{
		"name: test # inline",
		"# standalone",
		"value: \"# literal\"",
	}

	expected := []string{
		"name: test ",
		"value: \"# literal\"",
	}

	assertFiltered(t, "config.yml", RetainNone, lines, expected)
}

// TestFilterXml verifies XML block comment filtering.
func TestFilterXml(t *testing.T) {
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

	assertFiltered(t, "layout.xml", RetainNone, lines, expected)
}

// TestFilterJavaStyle verifies Java-family marker-based filtering.
func TestFilterJavaStyle(t *testing.T) {
	t.Run("documentation", func(t *testing.T) {
		lines := []string{
			"/** API docs. */",
			"// implementation note",
			"fun call() = \"// literal\"",
		}

		expected := []string{
			"/** API docs. */",
			"fun call() = \"// literal\"",
		}

		assertFiltered(t, "api.kt", RetainDocumentation, lines, expected)
	})

	t.Run("block", func(t *testing.T) {
		lines := []string{
			"/** API docs. */",
			"/* implementation note */",
			"String create();",
		}

		expected := []string{
			"/* implementation note */",
			"String create();",
		}

		assertFiltered(t, "Api.java", RetainBlock, lines, expected)
	})

	t.Run("regular", func(t *testing.T) {
		lines := []string{
			"/** API docs. */",
			"/* implementation note */",
			"String create(); // inline note",
		}

		expected := []string{
			"/* implementation note */",
			"String create(); // inline note",
		}

		assertFiltered(t, "Api.java", RetainRegular, lines, expected)
	})
}

// TestFilterCSharp verifies C# XML documentation comment filtering.
func TestFilterCSharp(t *testing.T) {
	t.Run("documentation", func(t *testing.T) {
		lines := []string{
			"/// <summary>Creates a value.</summary>",
			"// implementation note",
			"public string Create() => \"// literal\";",
		}

		expected := []string{
			"/// <summary>Creates a value.</summary>",
			"public string Create() => \"// literal\";",
		}

		assertFiltered(t, "Api.cs", RetainDocumentation, lines, expected)
	})

	t.Run("inline", func(t *testing.T) {
		lines := []string{
			"/// <summary>Creates a value.</summary>",
			"// implementation note",
			"public string Create() => \"// literal\";",
		}

		expected := []string{
			"// implementation note",
			"public string Create() => \"// literal\";",
		}

		assertFiltered(t, "Api.cs", RetainInline, lines, expected)
	})
}

// TestFilterVisualBasic verifies Visual Basic comment filtering.
func TestFilterVisualBasic(t *testing.T) {
	t.Run("none", func(t *testing.T) {
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

		assertFiltered(t, "Module.vb", RetainNone, lines, expected)
	})

	t.Run("regular", func(t *testing.T) {
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

		assertFiltered(t, "Module.vb", RetainRegular, lines, expected)
	})

	t.Run("documentation", func(t *testing.T) {
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

		assertFiltered(t, "Module.vb", RetainDocumentation, lines, expected)
	})
}

// TestFilterUnsupportedExtension verifies unsupported files are returned unchanged.
func TestFilterUnsupportedExtension(t *testing.T) {
	lines := []string{
		"# docs",
		"sub call { } # inline",
	}

	assertFiltered(t, "service.pl", RetainAll, lines, lines)
}

// TestFilterWarnsAboutUselessMode verifies warnings for modes without language-specific meaning.
func TestFilterWarnsAboutUselessMode(t *testing.T) {
	output := captureWarnings(func() {
		Filter([]string{"<!-- comment -->"}, "layout.xml", RetainDocumentation, "docs/guide.md", 12)
	})

	if !strings.Contains(output, "documentation") ||
		!strings.Contains(output, "layout.xml") ||
		!strings.Contains(output, "file://") ||
		!strings.Contains(output, "docs/guide.md:12") ||
		!strings.Contains(output, "does not have a distinct meaning") {
		t.Fatalf("warning output = %q", output)
	}
}

// TestFilterWarnsAboutUnsupportedExtension verifies warnings for unsupported file extensions.
func TestFilterWarnsAboutUnsupportedExtension(t *testing.T) {
	output := captureWarnings(func() {
		Filter([]string{"# comment"}, "service.pl", RetainNone, "docs/guide.md", 12)
	})

	if !strings.Contains(output, "comment filtering is not supported for this file extension") ||
		!strings.Contains(output, "file://") ||
		!strings.Contains(output, "docs/guide.md:12") {
		t.Fatalf("warning output = %q", output)
	}
}

// assertFiltered verifies filtering output for one file path and mode.
func assertFiltered(t *testing.T, filePath string, mode Mode, lines []string, expected []string) {
	t.Helper()

	got := Filter(lines, filePath, mode, "docs/guide.md", 12)
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("Filter() = %#v, expected %#v", got, expected)
	}
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
