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
	syntax := SyntaxFor(filePath)
	if len(syntax.Line) == 0 && len(syntax.Block) == 0 {
		return lines
	}

	return filterLines(lines, syntax, mode)
}

type blockState struct {
	active bool
	syntax BlockSyntax
	keep   bool
}

// filterLines removes or preserves recognized comments across all lines.
func filterLines(lines []string, syntax Syntax, mode Mode) []string {
	var filtered []string
	state := blockState{}
	for _, line := range lines {
		filteredLine, hadComment := filterLine(line, syntax, mode, &state)
		if hadComment && strings.TrimSpace(filteredLine) == "" {
			continue
		}
		filtered = append(filtered, filteredLine)
	}

	return filtered
}

// filterLine removes or preserves recognized comments from a single source line.
func filterLine(line string, syntax Syntax, mode Mode, state *blockState) (string, bool) {
	var result strings.Builder
	position := 0
	hadComment := false

	for position < len(line) {
		if state.active {
			hadComment = true
			end := strings.Index(line[position:], state.syntax.End)
			if end < 0 {
				if state.keep {
					result.WriteString(line[position:])
				}
				return result.String(), hadComment
			}
			endPosition := position + end + len(state.syntax.End)
			if state.keep {
				result.WriteString(line[position:endPosition])
			}
			position = endPosition
			state.active = false
			continue
		}

		if quoteEnd := quotedSegmentEnd(line, position, syntax.QuoteChars); quoteEnd > position {
			result.WriteString(line[position:quoteEnd])
			position = quoteEnd
			continue
		}
		if lineSyntax, found := lineCommentAt(line, position, syntax); found {
			hadComment = true
			if keepLineComment(lineSyntax, mode) {
				result.WriteString(line[position:])
			}
			break
		}
		if blockSyntax, found := blockCommentAt(line, position, syntax); found {
			hadComment = true
			state.active = true
			state.syntax = blockSyntax
			state.keep = keepBlockComment(blockSyntax, mode)
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

// lineCommentAt reports whether a line comment starts at the given position.
func lineCommentAt(line string, position int, syntax Syntax) (LineSyntax, bool) {
	for _, lineSyntax := range syntax.Line {
		if strings.HasPrefix(line[position:], lineSyntax.Prefix) {
			return lineSyntax, true
		}
	}

	return LineSyntax{}, false
}

// blockCommentAt reports whether a block comment starts at the given position.
func blockCommentAt(line string, position int, syntax Syntax) (BlockSyntax, bool) {
	for _, blockSyntax := range syntax.Block {
		if strings.HasPrefix(line[position:], blockSyntax.Start) {
			return blockSyntax, true
		}
	}

	return BlockSyntax{}, false
}

// keepLineComment reports whether the mode retains the given line comment kind.
func keepLineComment(lineSyntax LineSyntax, mode Mode) bool {
	return mode == RetainInline || mode == RetainDocumentation && lineSyntax.Documentation
}

// keepBlockComment reports whether the mode retains the given block comment kind.
func keepBlockComment(blockSyntax BlockSyntax, mode Mode) bool {
	return mode == RetainBlock || mode == RetainDocumentation && blockSyntax.Documentation
}
