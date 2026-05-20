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
	"fmt"
	"log/slog"
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

// filterEntry stores a comment filter and the modes that make sense for its language.
type filterEntry struct {
	filter      Filterer
	usefulModes []Mode
}

// filterFor returns the comment filter registered for the given file path and warns on odd modes.
func filterFor(
	filePath string,
	mode Mode,
	embeddingDocPath string,
	embeddingLine int,
) (Filterer, bool) {
	extension := normalizeExtension(filepath.Ext(filePath))
	entry, found := filtersByExtension[extension]
	if !found {
		warnUnsupportedCommentsMode(filePath, mode, embeddingDocPath, embeddingLine)
		return nil, false
	}
	warnUselessCommentsMode(filePath, mode, embeddingDocPath, embeddingLine, entry.usefulModes)

	return entry.filter, true
}

// normalizeExtension returns a lowercase file extension with a leading dot.
func normalizeExtension(extension string) string {
	normalized := strings.ToLower(extension)
	if normalized == "" || strings.HasPrefix(normalized, ".") {
		return normalized
	}

	return "." + normalized
}

var javaSyntax = Syntax{
	Inline: []string{"//"},
	Block: []BlockSyntax{
		{Start: "/*", End: "*/"},
	},
	Documentation: DocumentationSyntax{
		Block: []BlockSyntax{{Start: "/**", End: "*/"}},
	},
	QuoteChars: "\"'",
}

var jsSyntax = Syntax{
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

var cStyleSyntax = Syntax{
	Inline: []string{"//"},
	Block: []BlockSyntax{
		{Start: "/*", End: "*/"},
	},
	QuoteChars: "\"'",
}

var goSyntax = Syntax{
	Inline: []string{"//"},
	Block: []BlockSyntax{
		{Start: "/*", End: "*/"},
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

var allCommentModes = []Mode{
	RetainAll,
	RetainNone,
	RetainDocumentation,
	RetainRegular,
	RetainInline,
	RetainBlock,
}

var allNoneCommentModes = []Mode{RetainAll, RetainNone}

var inlineBlockCommentModes = []Mode{
	RetainAll,
	RetainNone,
	RetainInline,
	RetainBlock,
}

var regularDocCommentModes = []Mode{
	RetainAll,
	RetainNone,
	RetainDocumentation,
	RetainRegular,
}

var filtersByExtension = map[string]filterEntry{
	// Java/Kotlin
	".java":   newFilterEntry(MarkerCommentFilter{Syntax: javaSyntax}, allCommentModes),
	".kt":     newFilterEntry(MarkerCommentFilter{Syntax: javaSyntax}, allCommentModes),
	".kts":    newFilterEntry(MarkerCommentFilter{Syntax: javaSyntax}, allCommentModes),
	".groovy": newFilterEntry(MarkerCommentFilter{Syntax: javaSyntax}, allCommentModes),

	// C#
	".cs": newFilterEntry(MarkerCommentFilter{Syntax: csharpSyntax}, allCommentModes),

	// C/C++
	".c":   newFilterEntry(MarkerCommentFilter{Syntax: cStyleSyntax}, inlineBlockCommentModes),
	".h":   newFilterEntry(MarkerCommentFilter{Syntax: cStyleSyntax}, inlineBlockCommentModes),
	".cc":  newFilterEntry(MarkerCommentFilter{Syntax: cStyleSyntax}, inlineBlockCommentModes),
	".cpp": newFilterEntry(MarkerCommentFilter{Syntax: cStyleSyntax}, inlineBlockCommentModes),
	".cxx": newFilterEntry(MarkerCommentFilter{Syntax: cStyleSyntax}, inlineBlockCommentModes),
	".hh":  newFilterEntry(MarkerCommentFilter{Syntax: cStyleSyntax}, inlineBlockCommentModes),
	".hpp": newFilterEntry(MarkerCommentFilter{Syntax: cStyleSyntax}, inlineBlockCommentModes),
	".hxx": newFilterEntry(MarkerCommentFilter{Syntax: cStyleSyntax}, inlineBlockCommentModes),

	// JavaScript
	".js":  newFilterEntry(MarkerCommentFilter{Syntax: jsSyntax}, allCommentModes),
	".jsx": newFilterEntry(MarkerCommentFilter{Syntax: jsSyntax}, allCommentModes),
	".ts":  newFilterEntry(MarkerCommentFilter{Syntax: jsSyntax}, allCommentModes),
	".tsx": newFilterEntry(MarkerCommentFilter{Syntax: jsSyntax}, allCommentModes),

	// Go
	".go": newFilterEntry(MarkerCommentFilter{Syntax: goSyntax}, inlineBlockCommentModes),

	// Protobuf
	".proto": newFilterEntry(MarkerCommentFilter{Syntax: cStyleSyntax}, inlineBlockCommentModes),

	// Python
	".py":  newFilterEntry(MarkerCommentFilter{Syntax: hashLineSyntax}, allNoneCommentModes),
	".pyi": newFilterEntry(MarkerCommentFilter{Syntax: hashLineSyntax}, allNoneCommentModes),
	".pyw": newFilterEntry(MarkerCommentFilter{Syntax: hashLineSyntax}, allNoneCommentModes),

	// YAML
	".yml":  newFilterEntry(MarkerCommentFilter{Syntax: hashLineSyntax}, allNoneCommentModes),
	".yaml": newFilterEntry(MarkerCommentFilter{Syntax: hashLineSyntax}, allNoneCommentModes),

	// XML
	".xml": newFilterEntry(MarkerCommentFilter{Syntax: xmlSyntax}, allNoneCommentModes),

	// HTML
	".html": newFilterEntry(MarkerCommentFilter{Syntax: xmlSyntax}, allNoneCommentModes),
	".htm":  newFilterEntry(MarkerCommentFilter{Syntax: xmlSyntax}, allNoneCommentModes),

	// Visual Basic
	".vb":       newFilterEntry(VisualBasicCommentFilter{}, regularDocCommentModes),
	".bas":      newFilterEntry(VisualBasicCommentFilter{}, regularDocCommentModes),
	".vbs":      newFilterEntry(VisualBasicCommentFilter{}, regularDocCommentModes),
	".vbscript": newFilterEntry(VisualBasicCommentFilter{}, regularDocCommentModes),
}

// newFilterEntry creates a filter registry entry.
func newFilterEntry(filter Filterer, usefulModes []Mode) filterEntry {
	return filterEntry{
		filter:      filter,
		usefulModes: usefulModes,
	}
}

// warnUnsupportedCommentsMode logs when comments filtering is requested for an unsupported file.
func warnUnsupportedCommentsMode(
	filePath string,
	mode Mode,
	embeddingDocPath string,
	embeddingLine int,
) {
	if mode == RetainAll {
		return
	}
	slog.Warn(
		fmt.Sprintf(
			"`comments=\"%s\"` was requested in `%s` for `%s`, "+
				"but comment filtering is not supported for this file extension.",
			mode,
			fileURL(embeddingDocPath, embeddingLine),
			filePath,
		),
	)
}

// warnUselessCommentsMode logs when the selected mode has no distinct meaning for a file.
func warnUselessCommentsMode(
	filePath string,
	mode Mode,
	embeddingDocPath string,
	embeddingLine int,
	usefulModes []Mode,
) {
	if containsMode(usefulModes, mode) {
		return
	}
	var wrappedModes []string
	for _, mode := range usefulModes {
		wrappedModes = append(wrappedModes, fmt.Sprintf("`%s`", mode))
	}

	slog.Warn(
		fmt.Sprintf(
			"`comments=\"%s\"` was requested in `%s` for `%s`, but this mode does not have "+
				"a distinct meaning for this file type. Supported modes are: %s.",
			mode,
			fileURL(embeddingDocPath, embeddingLine),
			filePath,
			strings.Join(wrappedModes, ", "),
		),
	)
}

// fileURL returns an absolute file URL for a local path and line.
func fileURL(path string, line int) string {
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return "file://" + path
	}

	url := "file://" + absolutePath
	if line > 0 {
		url = fmt.Sprintf("%s:%d", url, line)
	}

	return url
}

// containsMode reports whether the list includes the given mode.
func containsMode(modes []Mode, mode Mode) bool {
	for _, usefulMode := range modes {
		if usefulMode == mode {
			return true
		}
	}

	return false
}
