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
	usefulModes map[Mode]struct{}
}

// filterFor returns the comment filter registered for the given file path and warns on odd modes.
func filterFor(filePath string, mode Mode, embeddingDocPath string) (Filterer, bool) {
	extension := normalizeExtension(filepath.Ext(filePath))
	entry, found := filtersByExtension[extension]
	if !found {
		warnUnsupportedCommentsMode(filePath, mode, embeddingDocPath)
		return nil, false
	}
	warnUselessCommentsMode(filePath, mode, embeddingDocPath, entry.usefulModes)

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

var allCommentModes = usefulModes(
	RetainAll,
	RetainNone,
	RetainDocumentation,
	RetainRegular,
	RetainInline,
	RetainBlock,
)

var allOrNoneCommentModes = usefulModes(RetainAll, RetainNone)

var regularAndDocCommentModes = usefulModes(
	RetainAll,
	RetainNone,
	RetainDocumentation,
	RetainRegular,
)

var filtersByExtension = map[string]filterEntry{
	// Java/Kotlin
	".java":   newFilterEntry(MarkerCommentFilter{Syntax: javaStyleSyntax}, allCommentModes),
	".kt":     newFilterEntry(MarkerCommentFilter{Syntax: javaStyleSyntax}, allCommentModes),
	".kts":    newFilterEntry(MarkerCommentFilter{Syntax: javaStyleSyntax}, allCommentModes),
	".groovy": newFilterEntry(MarkerCommentFilter{Syntax: javaStyleSyntax}, allCommentModes),

	// C#
	".cs": newFilterEntry(MarkerCommentFilter{Syntax: csharpSyntax}, allCommentModes),

	// JavaScript
	".js":  newFilterEntry(MarkerCommentFilter{Syntax: javaStyleSyntax}, allCommentModes),
	".jsx": newFilterEntry(MarkerCommentFilter{Syntax: javaStyleSyntax}, allCommentModes),
	".ts":  newFilterEntry(MarkerCommentFilter{Syntax: javaStyleSyntax}, allCommentModes),
	".tsx": newFilterEntry(MarkerCommentFilter{Syntax: javaStyleSyntax}, allCommentModes),

	// YAML
	".yml":  newFilterEntry(MarkerCommentFilter{Syntax: hashLineSyntax}, allOrNoneCommentModes),
	".yaml": newFilterEntry(MarkerCommentFilter{Syntax: hashLineSyntax}, allOrNoneCommentModes),

	// XML
	".xml": newFilterEntry(MarkerCommentFilter{Syntax: xmlSyntax}, allOrNoneCommentModes),

	// HTML
	".html": newFilterEntry(MarkerCommentFilter{Syntax: xmlSyntax}, allOrNoneCommentModes),
	".htm":  newFilterEntry(MarkerCommentFilter{Syntax: xmlSyntax}, allOrNoneCommentModes),

	// Visual Basic
	".vb":       newFilterEntry(VisualBasicCommentFilter{}, regularAndDocCommentModes),
	".bas":      newFilterEntry(VisualBasicCommentFilter{}, regularAndDocCommentModes),
	".vbs":      newFilterEntry(VisualBasicCommentFilter{}, regularAndDocCommentModes),
	".vbscript": newFilterEntry(VisualBasicCommentFilter{}, regularAndDocCommentModes),
}

// usefulModes creates a lookup set for comment modes that make sense for a language.
func usefulModes(modes ...Mode) map[Mode]struct{} {
	result := make(map[Mode]struct{}, len(modes))
	for _, mode := range modes {
		result[mode] = struct{}{}
	}

	return result
}

// newFilterEntry creates a filter registry entry.
func newFilterEntry(filter Filterer, usefulModes map[Mode]struct{}) filterEntry {
	return filterEntry{
		filter:      filter,
		usefulModes: usefulModes,
	}
}

// warnUnsupportedCommentsMode logs when comments filtering is requested for an unsupported file.
func warnUnsupportedCommentsMode(filePath string, mode Mode, embeddingDocPath string) {
	if mode == RetainAll {
		return
	}
	slog.Warn(
		fmt.Sprintf(
			"`comments=\"%s\"` was requested in `%s` for `%s`, "+
				"but comment filtering is not supported for this file extension.",
			mode,
			fileURL(embeddingDocPath),
			filePath,
		),
	)
}

// warnUselessCommentsMode logs when the selected mode has no distinct meaning for a file.
func warnUselessCommentsMode(
	filePath string,
	mode Mode,
	embeddingDocPath string,
	usefulModes map[Mode]struct{},
) {
	if _, found := usefulModes[mode]; found {
		return
	}
	slog.Warn(
		fmt.Sprintf(
			"`comments=\"%s\"` was requested in `%s` for `%s`, but this mode does not have "+
				"a distinct meaning for this file type. Useful modes are: %s.",
			mode,
			fileURL(embeddingDocPath),
			filePath,
			formatModes(usefulModes),
		),
	)
}

// fileURL returns an absolute file URL for a local path.
func fileURL(path string) string {
	if path == "" {
		return "file://<unknown>"
	}
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return "file://" + path
	}

	return "file://" + absolutePath
}

// formatModes formats modes for a warning message.
func formatModes(modes map[Mode]struct{}) string {
	order := []Mode{
		RetainAll,
		RetainNone,
		RetainDocumentation,
		RetainRegular,
		RetainInline,
		RetainBlock,
	}
	var result []string
	for _, mode := range order {
		if _, found := modes[mode]; found {
			result = append(result, fmt.Sprintf("`%s`", mode))
		}
	}

	return strings.Join(result, ", ")
}
