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
	"strings"
	"unicode"
)

const (
	commentPrefix = '\''
	docPrefix     = "'''"
	rem           = "rem"
)

// VisualBasicCommentFilter filters the Visual Basic comment forms:
//   - documentation comments starting with `”'`;
//   - apostrophe comments starting with `'`;
//   - REM comments starting with `REM`.
type VisualBasicCommentFilter struct{}

// Filter removes or preserves Visual Basic comments according to mode.
func (VisualBasicCommentFilter) Filter(lines []string, mode CommentFilterMode) []string {
	var filtered []string
	for _, line := range lines {
		filteredLine, hadComment := filterVisualBasicLine(line, mode)
		if hadComment && strings.TrimSpace(filteredLine) == "" {
			continue
		}
		filtered = append(filtered, filteredLine)
	}

	return filtered
}

// filterVisualBasicLine removes or preserves one Visual Basic comment.
func filterVisualBasicLine(line string, mode CommentFilterMode) (string, bool) {
	var result strings.Builder
	position := 0
	for position < len(line) {
		if quoteEnd := quotedSegmentEnd(line, position, "\""); quoteEnd > position {
			result.WriteString(line[position:quoteEnd])
			position = quoteEnd
			continue
		}
		if strings.HasPrefix(line[position:], docPrefix) {
			if mode == RetainDocumentation {
				result.WriteString(line[position:])
			}
			return result.String(), true
		}
		if line[position] == commentPrefix || remCommentAt(line, position) {
			if mode == RetainInline || mode == RetainRegular {
				result.WriteString(line[position:])
			}
			return result.String(), true
		}
		result.WriteByte(line[position])
		position++
	}

	return result.String(), false
}

// remCommentAt reports whether a Visual Basic REM comment starts at position.
func remCommentAt(line string, position int) bool {
	if len(line[position:]) < len(rem) ||
		!strings.EqualFold(
			line[position:position+len(rem)],
			rem,
		) {
		return false
	}
	return remPrefixBoundary(line, position) &&
		remSuffixBoundary(line, position+len(rem))
}

// remPrefixBoundary reports whether REM appears where a statement can start.
func remPrefixBoundary(line string, position int) bool {
	for cursor := position - 1; cursor >= 0; cursor-- {
		if unicode.IsSpace(rune(line[cursor])) {
			continue
		}
		return line[cursor] == ':'
	}

	return true
}

// remSuffixBoundary reports whether REM is followed by whitespace or the end of line.
func remSuffixBoundary(line string, position int) bool {
	return position >= len(line) || unicode.IsSpace(rune(line[position]))
}
