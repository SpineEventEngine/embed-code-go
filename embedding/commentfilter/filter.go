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

import "strings"

// Filter returns source lines with comments retained according to the requested mode.
func Filter(lines []string, filePath string, mode Mode) []string {
	if mode == RetainAll {
		return lines
	}
	filter, found := filterFor(filePath)
	if !found {
		return lines
	}

	return filter.Filter(lines, mode)
}

// MarkerCommentFilter removes comments using lexical markers declared in Syntax.
type MarkerCommentFilter struct {
	Syntax Syntax
}

type blockState struct {
	active bool
	block  BlockSyntax
	keep   bool
}

// Filter removes or preserves recognized comments across all lines.
func (f MarkerCommentFilter) Filter(lines []string, mode Mode) []string {
	var filtered []string
	state := blockState{}
	for _, line := range lines {
		filteredLine, hadComment := f.filterLine(line, mode, &state)
		if hadComment && strings.TrimSpace(filteredLine) == "" {
			continue
		}
		filtered = append(filtered, filteredLine)
	}

	return filtered
}

// filterLine removes or preserves recognized comments from a single source line.
func (f MarkerCommentFilter) filterLine(line string, mode Mode, state *blockState) (string, bool) {
	var result strings.Builder
	position := 0
	hadComment := false

	for position < len(line) {
		if state.active {
			hadComment = true
			end := strings.Index(line[position:], state.block.End)
			if end < 0 {
				if state.keep {
					result.WriteString(line[position:])
				}
				return result.String(), hadComment
			}
			endPosition := position + end + len(state.block.End)
			if state.keep {
				result.WriteString(line[position:endPosition])
			}
			position = endPosition
			state.active = false
			continue
		}

		if quoteEnd := quotedSegmentEnd(line, position, f.Syntax.QuoteChars); quoteEnd > position {
			result.WriteString(line[position:quoteEnd])
			position = quoteEnd
			continue
		}
		if _, found := documentationInlineAt(line, position, f.Syntax); found {
			hadComment = true
			if mode == RetainDocumentation {
				result.WriteString(line[position:])
			}
			break
		}
		if block, found := documentationBlockAt(line, position, f.Syntax); found {
			hadComment = true
			state.active = true
			state.block = block
			state.keep = mode == RetainDocumentation
			continue
		}
		if _, found := inlineCommentAt(line, position, f.Syntax); found {
			hadComment = true
			if mode == RetainInline || mode == RetainRegular {
				result.WriteString(line[position:])
			}
			break
		}
		if block, found := blockCommentAt(line, position, f.Syntax); found {
			hadComment = true
			state.active = true
			state.block = block
			state.keep = mode == RetainBlock || mode == RetainRegular
			continue
		}

		result.WriteByte(line[position])
		position++
	}

	return result.String(), hadComment
}

// quotedSegmentEnd returns the end offset of a quoted string starting at position.
func quotedSegmentEnd(line string, position int, quoteChars string) int {
	if position >= len(line) || !strings.ContainsRune(quoteChars, rune(line[position])) {
		return position
	}
	quote := line[position]
	cursor := position + 1
	for cursor < len(line) {
		if line[cursor] == '\\' {
			cursor += 2
			continue
		}
		if line[cursor] == quote {
			return cursor + 1
		}
		cursor++
	}

	return len(line)
}

// documentationInlineAt reports whether a documentation line comment starts at the position.
func documentationInlineAt(line string, position int, syntax Syntax) (string, bool) {
	return prefixAt(line, position, syntax.Documentation.Inline)
}

// documentationBlockAt reports whether a documentation block comment starts at the position.
func documentationBlockAt(line string, position int, syntax Syntax) (BlockSyntax, bool) {
	return blockAt(line, position, syntax.Documentation.Block)
}

// inlineCommentAt reports whether an inline comment starts at the given position.
func inlineCommentAt(line string, position int, syntax Syntax) (string, bool) {
	return prefixAt(line, position, syntax.Inline)
}

// blockCommentAt reports whether a block comment starts at the given position.
func blockCommentAt(line string, position int, syntax Syntax) (BlockSyntax, bool) {
	return blockAt(line, position, syntax.Block)
}

// prefixAt reports whether one of the given prefixes starts at the position.
func prefixAt(line string, position int, prefixes []string) (string, bool) {
	for _, prefix := range prefixes {
		if strings.HasPrefix(line[position:], prefix) {
			return prefix, true
		}
	}

	return "", false
}

// blockAt reports whether one of the given block markers starts at the position.
func blockAt(line string, position int, blocks []BlockSyntax) (BlockSyntax, bool) {
	for _, block := range blocks {
		if strings.HasPrefix(line[position:], block.Start) {
			return block, true
		}
	}

	return BlockSyntax{}, false
}
