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

// filtersByExtension is a mapping of the file extension to its comment filter.
var filtersByExtension = map[string]filterEntry{
	// Java/Kotlin
	".java":   filterConfig(MarkerCommentFilter{Syntax: javaSyntax}, allModes),
	".kt":     filterConfig(MarkerCommentFilter{Syntax: javaSyntax}, allModes),
	".kts":    filterConfig(MarkerCommentFilter{Syntax: javaSyntax}, allModes),
	".groovy": filterConfig(MarkerCommentFilter{Syntax: javaSyntax}, allModes),

	// C#
	".cs": filterConfig(MarkerCommentFilter{Syntax: csharpSyntax}, allModes),

	// C/C++
	".c":   filterConfig(MarkerCommentFilter{Syntax: cStyleSyntax}, regularModes),
	".h":   filterConfig(MarkerCommentFilter{Syntax: cStyleSyntax}, regularModes),
	".cc":  filterConfig(MarkerCommentFilter{Syntax: cStyleSyntax}, regularModes),
	".cpp": filterConfig(MarkerCommentFilter{Syntax: cStyleSyntax}, regularModes),
	".cxx": filterConfig(MarkerCommentFilter{Syntax: cStyleSyntax}, regularModes),
	".hh":  filterConfig(MarkerCommentFilter{Syntax: cStyleSyntax}, regularModes),
	".hpp": filterConfig(MarkerCommentFilter{Syntax: cStyleSyntax}, regularModes),
	".hxx": filterConfig(MarkerCommentFilter{Syntax: cStyleSyntax}, regularModes),

	// JavaScript
	".js":  filterConfig(MarkerCommentFilter{Syntax: jsSyntax}, allModes),
	".jsx": filterConfig(MarkerCommentFilter{Syntax: jsSyntax}, allModes),
	".ts":  filterConfig(MarkerCommentFilter{Syntax: jsSyntax}, allModes),
	".tsx": filterConfig(MarkerCommentFilter{Syntax: jsSyntax}, allModes),

	// Go
	".go": filterConfig(MarkerCommentFilter{Syntax: goSyntax}, regularModes),

	// Protobuf
	".proto": filterConfig(MarkerCommentFilter{Syntax: cStyleSyntax}, regularModes),

	// Python
	".py":  filterConfig(MarkerCommentFilter{Syntax: hashLineSyntax}, noneMode),
	".pyi": filterConfig(MarkerCommentFilter{Syntax: hashLineSyntax}, noneMode),
	".pyw": filterConfig(MarkerCommentFilter{Syntax: hashLineSyntax}, noneMode),

	// YAML
	".yml":  filterConfig(MarkerCommentFilter{Syntax: hashLineSyntax}, noneMode),
	".yaml": filterConfig(MarkerCommentFilter{Syntax: hashLineSyntax}, noneMode),

	// XML
	".xml": filterConfig(MarkerCommentFilter{Syntax: xmlSyntax}, noneMode),

	// HTML
	".html": filterConfig(MarkerCommentFilter{Syntax: xmlSyntax}, noneMode),
	".htm":  filterConfig(MarkerCommentFilter{Syntax: xmlSyntax}, noneMode),

	// Visual Basic
	".vb":       filterConfig(VisualBasicCommentFilter{}, documentationModes),
	".bas":      filterConfig(VisualBasicCommentFilter{}, documentationModes),
	".vbs":      filterConfig(VisualBasicCommentFilter{}, documentationModes),
	".vbscript": filterConfig(VisualBasicCommentFilter{}, documentationModes),
}

// filterEntry stores a comment filter and supported modes for its language.
type filterEntry struct {
	filter         CommentFilter
	supportedModes []Mode
}

var javaSyntax = CommentMarker{
	Inline: []string{"//"},
	Block: []BlockMarker{
		{Start: "/*", End: "*/"},
	},
	Documentation: DocumentationMarker{
		Block: []BlockMarker{{Start: "/**", End: "*/"}},
	},
	QuoteChars: "\"'",
}

var jsSyntax = CommentMarker{
	Inline: []string{"//"},
	Block: []BlockMarker{
		{Start: "/*", End: "*/"},
	},
	Documentation: DocumentationMarker{
		Block: []BlockMarker{{Start: "/**", End: "*/"}},
	},
	QuoteChars: "\"'`",
}

var csharpSyntax = CommentMarker{
	Inline: []string{"//"},
	Block: []BlockMarker{
		{Start: "/*", End: "*/"},
	},
	Documentation: DocumentationMarker{
		Inline: []string{"///"},
		Block:  []BlockMarker{{Start: "/**", End: "*/"}},
	},
	QuoteChars: "\"'`",
}

var cStyleSyntax = CommentMarker{
	Inline: []string{"//"},
	Block: []BlockMarker{
		{Start: "/*", End: "*/"},
	},
	QuoteChars: "\"'",
}

var goSyntax = CommentMarker{
	Inline: []string{"//"},
	Block: []BlockMarker{
		{Start: "/*", End: "*/"},
	},
	QuoteChars: "\"'`",
}

var hashLineSyntax = CommentMarker{
	Inline:     []string{"#"},
	QuoteChars: "\"'",
}

var xmlSyntax = CommentMarker{
	Block: []BlockMarker{
		{Start: "<!--", End: "-->"},
	},
	QuoteChars: "\"'",
}

// allModes lists all comment filtering modes.
var allModes = []Mode{
	RetainAll,
	RetainNone,
	RetainDocumentation,
	RetainRegular,
	RetainInline,
	RetainBlock,
}

// noneMode lists modes for languages whose comments are not separated into supported subtypes.
var noneMode = []Mode{RetainAll, RetainNone}

// regularModes lists modes for languages that distinguish inline and block comments,
// but do not expose documentation comments as a separate supported type.
var regularModes = []Mode{
	RetainAll,
	RetainNone,
	RetainInline,
	RetainBlock,
}

// documentationModes lists modes for languages that distinguish documentation and regular comments,
// but do not expose inline and block comments as separate supported types.
var documentationModes = []Mode{
	RetainAll,
	RetainNone,
	RetainDocumentation,
	RetainRegular,
}

// filterConfig creates a filter registry entry.
func filterConfig(filter CommentFilter, supportedModes []Mode) filterEntry {
	return filterEntry{
		filter:         filter,
		supportedModes: supportedModes,
	}
}
