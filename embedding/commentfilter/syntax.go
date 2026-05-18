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
	"path/filepath"
	"strings"
)

// BlockSyntax describes a block comment marker pair.
type BlockSyntax struct {
	Start string
	End   string
}

// DocumentationSyntax describes API documentation comment markers.
type DocumentationSyntax struct {
	Inline []string
	Block  []BlockSyntax
}

// Syntax describes lexical comment markers and string delimiters for a language family.
type Syntax struct {
	Inline        []string
	Block         []BlockSyntax
	Documentation DocumentationSyntax
	QuoteChars    string
}

// Filterer removes or preserves source comments according to the requested mode.
type Filterer interface {
	Filter(lines []string, mode Mode) []string
}

// filterFor returns the comment filter registered for the given file path.
func filterFor(filePath string) (Filterer, bool) {
	extension := normalizeExtension(filepath.Ext(filePath))
	filter, found := filtersByExtension[extension]
	return filter, found
}

// normalizeExtension returns a lowercase file extension with a leading dot.
func normalizeExtension(extension string) string {
	normalized := strings.ToLower(extension)
	if normalized == "" || strings.HasPrefix(normalized, ".") {
		return normalized
	}

	return "." + normalized
}

var javaStyleSyntax = Syntax{
	Inline: []string{"//"},
	Block: []BlockSyntax{
		{Start: "/*", End: "*/"},
	},
	Documentation: DocumentationSyntax{
		Block: []BlockSyntax{{Start: "/**", End: "*/"}},
	},
	QuoteChars: "\"'`",
}

var csharpSyntax = Syntax{
	Inline: []string{"//"},
	Block: []BlockSyntax{
		{Start: "/*", End: "*/"},
	},
	Documentation: DocumentationSyntax{
		Inline: []string{"///"},
		Block:  []BlockSyntax{{Start: "/**", End: "*/"}},
	},
	QuoteChars: "\"'`",
}

var hashLineSyntax = Syntax{
	Inline:     []string{"#"},
	QuoteChars: "\"'",
}

var xmlSyntax = Syntax{
	Block: []BlockSyntax{
		{Start: "<!--", End: "-->"},
	},
	QuoteChars: "\"'",
}

var filtersByExtension = map[string]Filterer{
	// Java/Kotlin
	".java":   MarkerCommentFilter{Syntax: javaStyleSyntax},
	".kt":     MarkerCommentFilter{Syntax: javaStyleSyntax},
	".kts":    MarkerCommentFilter{Syntax: javaStyleSyntax},
	".groovy": MarkerCommentFilter{Syntax: javaStyleSyntax},

	// C#
	".cs": MarkerCommentFilter{Syntax: csharpSyntax},

	// JavaScript
	".js":  MarkerCommentFilter{Syntax: javaStyleSyntax},
	".jsx": MarkerCommentFilter{Syntax: javaStyleSyntax},
	".ts":  MarkerCommentFilter{Syntax: javaStyleSyntax},
	".tsx": MarkerCommentFilter{Syntax: javaStyleSyntax},

	// YAML
	".yml":  MarkerCommentFilter{Syntax: hashLineSyntax},
	".yaml": MarkerCommentFilter{Syntax: hashLineSyntax},

	// XML
	".xml": MarkerCommentFilter{Syntax: xmlSyntax},

	// HTML
	".html": MarkerCommentFilter{Syntax: xmlSyntax},
	".htm":  MarkerCommentFilter{Syntax: xmlSyntax},

	// Visual Basic
	".vb":       VisualBasicCommentFilter{},
	".bas":      VisualBasicCommentFilter{},
	".vbs":      VisualBasicCommentFilter{},
	".vbscript": VisualBasicCommentFilter{},
}
