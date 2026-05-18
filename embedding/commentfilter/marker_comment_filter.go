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

// MarkerCommentFilter removes comments using lexical markers declared in Syntax.
type MarkerCommentFilter struct {
	Syntax Syntax
}

type blockState struct {
	active bool
	block  BlockSyntax
	keep   bool
}

type markerLineFilter struct {
	filter     MarkerCommentFilter
	line       string
	mode       Mode
	state      *blockState
	result     strings.Builder
	position   int
	hadComment bool
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
	filter := markerLineFilter{
		filter: f,
		line:   line,
		mode:   mode,
		state:  state,
	}

	return filter.filterLine()
}

// filterLine walks the current line until it reaches its end or a line comment.
func (f *markerLineFilter) filterLine() (string, bool) {
	for f.position < len(f.line) {
		if f.consumeActiveBlock() {
			continue
		}
		if f.consumeQuotedSegment() {
			continue
		}
		if consumed, stop := f.consumeComment(); consumed {
			if stop {
				break
			}
			continue
		}
		f.consumeCodeByte()
	}

	return f.result.String(), f.hadComment
}

// consumeActiveBlock consumes text while the scanner is inside a block comment.
func (f *markerLineFilter) consumeActiveBlock() bool {
	if !f.state.active {
		return false
	}
	f.hadComment = true
	end := strings.Index(f.line[f.position:], f.state.block.End)
	if end < 0 {
		if f.state.keep {
			f.result.WriteString(f.line[f.position:])
		}
		f.position = len(f.line)
		return true
	}
	endPosition := f.position + end + len(f.state.block.End)
	if f.state.keep {
		f.result.WriteString(f.line[f.position:endPosition])
	}
	f.position = endPosition
	f.state.active = false

	return true
}

// consumeQuotedSegment copies a quoted segment without scanning comment markers inside it.
func (f *markerLineFilter) consumeQuotedSegment() bool {
	quoteEnd := quotedSegmentEnd(f.line, f.position, f.filter.Syntax.QuoteChars)
	if quoteEnd <= f.position {
		return false
	}
	f.result.WriteString(f.line[f.position:quoteEnd])
	f.position = quoteEnd

	return true
}

// consumeComment consumes a comment and reports whether it consumed input and ended the line.
func (f *markerLineFilter) consumeComment() (bool, bool) {
	if _, found := documentationInlineAt(f.line, f.position, f.filter.Syntax); found {
		f.consumeInlineComment(f.mode == RetainDocumentation)
		return true, true
	}
	if block, found := documentationBlockAt(f.line, f.position, f.filter.Syntax); found {
		f.startBlockComment(block, f.mode == RetainDocumentation)
		return true, false
	}
	if _, found := inlineCommentAt(f.line, f.position, f.filter.Syntax); found {
		f.consumeInlineComment(f.mode == RetainInline || f.mode == RetainRegular)
		return true, true
	}
	if block, found := blockCommentAt(f.line, f.position, f.filter.Syntax); found {
		f.startBlockComment(block, f.mode == RetainBlock || f.mode == RetainRegular)
		return true, false
	}

	return false, false
}

// consumeInlineComment consumes the rest of the line as a line comment.
func (f *markerLineFilter) consumeInlineComment(keep bool) {
	f.hadComment = true
	if keep {
		f.result.WriteString(f.line[f.position:])
	}
	f.position = len(f.line)
}

// startBlockComment records the active block comment markers and whether to keep them.
func (f *markerLineFilter) startBlockComment(block BlockSyntax, keep bool) {
	f.hadComment = true
	f.state.active = true
	f.state.block = block
	f.state.keep = keep
}

// consumeCodeByte copies one source byte that does not belong to a recognized comment.
func (f *markerLineFilter) consumeCodeByte() {
	f.result.WriteByte(f.line[f.position])
	f.position++
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
