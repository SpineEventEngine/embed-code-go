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

// EmbeddingCommentFilter filters comments for one embed-code instruction.
type EmbeddingCommentFilter struct {
	filePath         string
	embeddingDocPath string
	embeddingLine    int
}

// CommentFilter strips source comments according to the requested mode.
type CommentFilter interface {
	Filter(lines []string, mode Mode) []string
}

// Filter returns source lines with comments stripped according to the requested mode.
func Filter(
	lines []string,
	filePath string,
	mode Mode,
	embeddingDocPath string,
	embeddingLine int,
) []string {
	filter := EmbeddingCommentFilter{
		filePath:         filePath,
		embeddingDocPath: embeddingDocPath,
		embeddingLine:    embeddingLine,
	}

	return filter.Filter(lines, mode)
}

// Filter strips comments using the filter registered in the filtersByExtension.
func (f EmbeddingCommentFilter) Filter(lines []string, mode Mode) []string {
	if mode == RetainAll {
		return lines
	}
	filter, found := filterFor(f.filePath, mode, f.embeddingDocPath, f.embeddingLine)
	if !found {
		return lines
	}

	return filter.Filter(lines, mode)
}

// filterFor returns the comment filter registered for the given file path and warns on odd modes.
func filterFor(
	filePath string,
	mode Mode,
	embeddingDocPath string,
	embeddingLine int,
) (CommentFilter, bool) {
	extension := normalizeExtension(filepath.Ext(filePath))
	entry, found := filtersByExtension[extension]
	if !found {
		warnUnsupportedFileType(filePath, mode, embeddingDocPath, embeddingLine)
		return nil, false
	}
	warnUnsupportedCommentsMode(filePath, mode, embeddingDocPath, embeddingLine, entry.supportedModes)

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

// warnUnsupportedFileType logs when comments filtering is requested for an unsupported file.
func warnUnsupportedFileType(
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

// warnUnsupportedCommentsMode logs when the selected mode is not supported for a file.
func warnUnsupportedCommentsMode(
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
