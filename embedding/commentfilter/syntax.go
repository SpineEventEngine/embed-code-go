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

// LineSyntax describes a single-line comment marker.
type LineSyntax struct {
	Prefix        string
	Documentation bool
}

// BlockSyntax describes a block comment marker pair.
type BlockSyntax struct {
	Start         string
	End           string
	Documentation bool
}

// Syntax describes comment markers and string delimiters for a language family.
type Syntax struct {
	Line       []LineSyntax
	Block      []BlockSyntax
	QuoteChars string
}

// SyntaxFor returns the comment syntax registered for the given file path.
func SyntaxFor(filePath string) Syntax {
	extension := normalizeExtension(filepath.Ext(filePath))
	if syntax, found := syntaxesByExtension[extension]; found {
		return syntax
	}

	return Syntax{}
}

// RegisterSyntax registers comment syntax for a source file extension.
func RegisterSyntax(extension string, syntax Syntax) {
	syntaxesByExtension[normalizeExtension(extension)] = syntax
}

// normalizeExtension returns a lowercase file extension with a leading dot.
func normalizeExtension(extension string) string {
	normalized := strings.ToLower(extension)
	if normalized == "" || strings.HasPrefix(normalized, ".") {
		return normalized
	}

	return "." + normalized
}

var cLikeSyntax = Syntax{
	Line: []LineSyntax{
		{Prefix: "///", Documentation: true},
		{Prefix: "//!", Documentation: true},
		{Prefix: "//", Documentation: false},
	},
	Block: []BlockSyntax{
		{Start: "/**", End: "*/", Documentation: true},
		{Start: "/*!", End: "*/", Documentation: true},
		{Start: "/*", End: "*/", Documentation: false},
	},
	QuoteChars: "\"'`",
}

var hashLineSyntax = Syntax{
	Line: []LineSyntax{
		{Prefix: "#", Documentation: false},
	},
	QuoteChars: "\"'",
}

var xmlSyntax = Syntax{
	Block: []BlockSyntax{
		{Start: "<!--", End: "-->", Documentation: false},
	},
	QuoteChars: "\"'",
}

var basicSyntax = Syntax{
	Line: []LineSyntax{
		{Prefix: "'", Documentation: false},
	},
	QuoteChars: "\"",
}

var syntaxesByExtension = map[string]Syntax{
	".java":       cLikeSyntax,
	".groovy":     cLikeSyntax,
	".kt":         cLikeSyntax,
	".kts":        cLikeSyntax,
	".c":          cLikeSyntax,
	".cc":         cLikeSyntax,
	".cpp":        cLikeSyntax,
	".cxx":        cLikeSyntax,
	".h":          cLikeSyntax,
	".hh":         cLikeSyntax,
	".hpp":        cLikeSyntax,
	".cs":         cLikeSyntax,
	".js":         cLikeSyntax,
	".jsx":        cLikeSyntax,
	".ts":         cLikeSyntax,
	".tsx":        cLikeSyntax,
	".go":         cLikeSyntax,
	".yml":        hashLineSyntax,
	".yaml":       hashLineSyntax,
	".xml":        xmlSyntax,
	".html":       xmlSyntax,
	".htm":        xmlSyntax,
	".vb":         basicSyntax,
	".bas":        basicSyntax,
	".vbs":        basicSyntax,
	".vbscript":   basicSyntax,
	".properties": hashLineSyntax,
}
