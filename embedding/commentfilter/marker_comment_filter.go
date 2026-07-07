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

// activeSegment describes a multi-line construct currently being scanned.
type activeSegment struct {
	// end closes the active construct.
	end string

	// keep reports whether the active construct is copied to output.
	keep bool

	// escapes reports whether backslashes escape bytes while searching for end.
	escapes bool

	// comment reports whether the active construct is a source comment.
	comment bool
}

// markerState tracks active multi-line lexical constructs across source lines.
type markerState struct {
	// segment contains the active construct configuration, if one is open.
	segment *activeSegment
}

// markerLineFilter tracks lexical comment filtering state for one source line.
type markerLineFilter struct {
	lineFilter

	// filter contains the language syntax configuration.
	filter MarkerCommentFilter

	// state tracks multi-line lexical constructs across lines.
	state *markerState
}

// Filter removes or preserves recognized comments across all lines.
//
// Parameters:
// lines - provides source lines.
// mode - selects comments to retain.
//
// Returns filtered source lines.
func (f MarkerCommentFilter) Filter(lines []string, mode Mode) []string {
	state := markerState{}

	return filterLines(lines, func(line string) (string, bool) {
		filter := markerLineFilter{
			lineFilter: lineFilter{line: line, mode: mode},
			filter:     f,
			state:      &state,
		}

		return filter.filterLine()
	})
}

// filterLine walks the current line until it reaches its end or a line comment.
func (f *markerLineFilter) filterLine() (string, bool) {
	for f.position < len(f.line) {
		if f.consumeActiveSegment() {
			continue
		}
		if f.consumeTextBlockStart() {
			continue
		}
		if f.consumeQuotedSegment(f.filter.Syntax.QuoteChars) {
			continue
		}
		if comment := f.consumeComment(); comment.consumed {
			if comment.stopLine {
				break
			}

			continue
		}
		f.consumeCodeByte()
	}

	return f.result.String(), f.hadComment
}

// consumeActiveSegment consumes text while the scanner is inside a multi-line construct.
func (f *markerLineFilter) consumeActiveSegment() bool {
	segment := f.state.segment
	if segment == nil {
		return false
	}
	if segment.comment {
		f.hadComment = true
	}
	endPosition, found := segmentEnd(f.line, f.position, *segment)
	if !found {
		if segment.keep {
			f.result.WriteString(f.line[f.position:])
		}
		f.position = len(f.line)

		return true
	}
	if segment.keep {
		f.result.WriteString(f.line[f.position:endPosition])
	}
	f.position = endPosition
	f.state.segment = nil

	return true
}

// segmentEnd returns the end offset of an active segment close delimiter.
func segmentEnd(line string, position int, segment activeSegment) (int, bool) {
	for cursor := position; cursor < len(line); {
		if segment.escapes && line[cursor] == '\\' {
			cursor += 2

			continue
		}
		if strings.HasPrefix(line[cursor:], segment.end) {
			return cursor + len(segment.end), true
		}
		cursor++
	}

	return len(line), false
}

// consumeTextBlockStart starts a configured text block literal.
func (f *markerLineFilter) consumeTextBlockStart() bool {
	marker, found := textBlockAt(f.line, f.position, f.filter.Syntax.TextBlocks)
	if !found {
		return false
	}
	f.consumeMarker(marker.Delimiter)
	f.state.segment = &activeSegment{
		end:     marker.Delimiter,
		keep:    true,
		escapes: marker.Escapes,
	}

	return true
}

// consumeComment consumes a comment when one starts at the scanner position.
func (f *markerLineFilter) consumeComment() commentConsumeResult {
	if prefixAt(f.line, f.position, f.filter.Syntax.Documentation.Inline) {
		f.consumeLineComment(f.mode == RetainDocumentation)

		return commentConsumeResult{consumed: true, stopLine: true}
	}
	if block, found := blockAt(f.line, f.position, f.filter.Syntax.Documentation.Block); found {
		f.startBlockComment(block, f.mode == RetainDocumentation)

		return commentConsumeResult{consumed: true}
	}
	if prefixAt(f.line, f.position, f.filter.Syntax.Inline) {
		f.consumeLineComment(f.mode == RetainInline || f.mode == RetainRegular)

		return commentConsumeResult{consumed: true, stopLine: true}
	}
	if block, found := blockAt(f.line, f.position, f.filter.Syntax.Block); found {
		f.startBlockComment(block, f.mode == RetainBlock || f.mode == RetainRegular)

		return commentConsumeResult{consumed: true}
	}

	return commentConsumeResult{}
}

// startBlockComment records the active block comment markers and whether to keep them.
func (f *markerLineFilter) startBlockComment(block BlockMarker, keep bool) {
	f.hadComment = true
	start := block.Start
	if block.End == cStyleBlockCommentEnd && strings.HasPrefix(start, cStyleDocCommentStart) {
		start = cStyleBlockCommentStart
	}
	if keep {
		f.result.WriteString(start)
	}
	f.position += len(start)
	f.state.segment = &activeSegment{
		end:     block.End,
		keep:    keep,
		comment: true,
	}
}

// prefixAt reports whether one of the given prefixes starts at the position.
func prefixAt(line string, position int, prefixes []string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(line[position:], prefix) {
			return true
		}
	}

	return false
}

// textBlockAt reports whether one of the given text block markers starts at the position.
func textBlockAt(line string, position int, markers []TextBlockMarker) (TextBlockMarker, bool) {
	for _, marker := range markers {
		if strings.HasPrefix(line[position:], marker.Delimiter) {
			return marker, true
		}
	}

	return TextBlockMarker{}, false
}

// blockAt reports whether one of the given block markers starts at the position.
func blockAt(line string, position int, blocks []BlockMarker) (BlockMarker, bool) {
	for _, block := range blocks {
		if strings.HasPrefix(line[position:], block.Start) {
			return block, true
		}
	}

	return BlockMarker{}, false
}
