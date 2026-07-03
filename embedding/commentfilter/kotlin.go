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

const kotlinRawStringDelimiter = "\"\"\""

// KotlinCommentFilter filters Kotlin comments while preserving Kotlin string forms.
type KotlinCommentFilter struct{}

// kotlinState tracks Kotlin lexical state that can span source lines.
type kotlinState struct {
	// blockDepth is the current nested block comment depth.
	blockDepth int

	// blockKeep reports whether the active block comment should be retained.
	blockKeep bool

	// rawString reports whether scanning is inside a raw triple-quoted string.
	rawString bool
}

// kotlinLineFilter filters one Kotlin source line.
type kotlinLineFilter struct {
	// line is the source line being filtered.
	line string

	// mode selects which comments to retain.
	mode Mode

	// state tracks Kotlin constructs across lines.
	state *kotlinState

	// result accumulates the filtered source line.
	result strings.Builder

	// position is the current byte index in line.
	position int

	// hadComment reports whether the line contained a recognized comment.
	hadComment bool
}

// Filter removes or preserves Kotlin comments according to mode.
//
// Parameters:
// lines - provides Kotlin source lines.
// mode - selects comments to retain.
//
// Returns filtered source lines.
func (KotlinCommentFilter) Filter(lines []string, mode Mode) []string {
	var filtered []string
	state := kotlinState{}
	for _, line := range lines {
		filteredLine, hadComment := filterKotlinLine(line, mode, &state)
		if hadComment && strings.TrimSpace(filteredLine) == "" {
			continue
		}
		filtered = append(filtered, filteredLine)
	}

	return filtered
}

// filterKotlinLine removes or preserves recognized Kotlin comments from one line.
func filterKotlinLine(line string, mode Mode, state *kotlinState) (string, bool) {
	filter := kotlinLineFilter{
		line:  line,
		mode:  mode,
		state: state,
	}

	return filter.filterLine()
}

// filterLine walks the current line until it reaches its end or a line comment.
func (f *kotlinLineFilter) filterLine() (string, bool) {
	for f.position < len(f.line) {
		if f.consumeActiveBlock() {
			continue
		}
		if f.consumeRawString() {
			continue
		}
		if f.consumeString() {
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

// consumeActiveBlock consumes a possibly nested Kotlin block comment.
func (f *kotlinLineFilter) consumeActiveBlock() bool {
	if f.state.blockDepth == 0 {
		return false
	}
	f.hadComment = true
	for f.position < len(f.line) {
		switch {
		case strings.HasPrefix(f.line[f.position:], cStyleBlockCommentStart):
			f.writeBlockText(cStyleBlockCommentStart)
			f.state.blockDepth++
			f.position += len(cStyleBlockCommentStart)
		case strings.HasPrefix(f.line[f.position:], cStyleBlockCommentEnd):
			f.writeBlockText(cStyleBlockCommentEnd)
			f.state.blockDepth--
			f.position += len(cStyleBlockCommentEnd)
			if f.state.blockDepth == 0 {
				return true
			}
		default:
			if f.state.blockKeep {
				f.result.WriteByte(f.line[f.position])
			}
			f.position++
		}
	}

	return true
}

// consumeRawString copies Kotlin raw-string text and filters `${...}` interpolation code.
//
// It treats the first three quotes in a run of four or more quotes as the raw-string delimiter.
func (f *kotlinLineFilter) consumeRawString() bool {
	if !f.state.rawString && !strings.HasPrefix(f.line[f.position:], kotlinRawStringDelimiter) {
		return false
	}
	if !f.state.rawString {
		f.state.rawString = true
		f.result.WriteString(kotlinRawStringDelimiter)
		f.position += len(kotlinRawStringDelimiter)
	}
	for f.position < len(f.line) {
		switch {
		case strings.HasPrefix(f.line[f.position:], kotlinRawStringDelimiter):
			f.result.WriteString(kotlinRawStringDelimiter)
			f.position += len(kotlinRawStringDelimiter)
			f.state.rawString = false

			return true
		case strings.HasPrefix(f.line[f.position:], "${"):
			f.result.WriteString("${")
			f.position += len("${")
			f.state.rawString = false
			f.consumeInterpolation()
			f.state.rawString = true
		default:
			f.consumeCodeByte()
		}
	}

	return true
}

// consumeString copies Kotlin string and character literals, filtering interpolated expressions.
func (f *kotlinLineFilter) consumeString() bool {
	if f.position >= len(f.line) {
		return false
	}
	switch f.line[f.position] {
	case '"':
		f.consumeQuotedString()

		return true
	case '\'':
		quoteEnd := quotedSegmentEnd(f.line, f.position, "'")
		f.result.WriteString(f.line[f.position:quoteEnd])
		f.position = quoteEnd

		return true
	default:
		return false
	}
}

// consumeQuotedString copies a Kotlin quoted string and filters comments inside `${...}`.
func (f *kotlinLineFilter) consumeQuotedString() {
	f.result.WriteByte(f.line[f.position])
	f.position++
	for f.position < len(f.line) {
		switch {
		case f.line[f.position] == '\\':
			f.writeEscapedByte()
		case f.line[f.position] == '"':
			f.result.WriteByte(f.line[f.position])
			f.position++

			return
		case strings.HasPrefix(f.line[f.position:], "${"):
			f.result.WriteString("${")
			f.position += len("${")
			f.consumeInterpolation()
		default:
			f.consumeCodeByte()
		}
	}
}

// consumeInterpolation filters comments inside a Kotlin string interpolation expression.
func (f *kotlinLineFilter) consumeInterpolation() {
	depth := 1
	for f.position < len(f.line) {
		if f.consumeActiveBlock() {
			continue
		}
		if f.consumeRawString() {
			continue
		}
		if f.consumeString() {
			continue
		}
		if consumed, stop := f.consumeComment(); consumed {
			if stop {
				return
			}

			continue
		}
		var done bool
		depth, done = f.consumeInterpolationCode(depth)
		if done {
			return
		}
	}
}

// consumeInterpolationCode copies expression code and updates interpolation brace depth.
func (f *kotlinLineFilter) consumeInterpolationCode(depth int) (int, bool) {
	switch f.line[f.position] {
	case '{':
		depth++
		f.consumeCodeByte()

		return depth, false
	case '}':
		depth--
		f.consumeCodeByte()

		return depth, depth == 0
	default:
		f.consumeCodeByte()

		return depth, false
	}
}

// consumeComment consumes a Kotlin comment and reports whether it ended the line.
func (f *kotlinLineFilter) consumeComment() (bool, bool) {
	if strings.HasPrefix(f.line[f.position:], cStyleDocCommentStart) {
		f.startBlockComment(f.mode == RetainDocumentation)

		return true, false
	}
	if strings.HasPrefix(f.line[f.position:], cStyleBlockCommentStart) {
		f.startBlockComment(f.mode == RetainBlock || f.mode == RetainRegular)

		return true, false
	}
	if strings.HasPrefix(f.line[f.position:], "//") {
		f.hadComment = true
		if f.mode == RetainInline || f.mode == RetainRegular {
			f.result.WriteString(f.line[f.position:])
		}
		f.position = len(f.line)

		return true, true
	}

	return false, false
}

// startBlockComment starts a Kotlin block comment with nesting depth one.
func (f *kotlinLineFilter) startBlockComment(keep bool) {
	f.hadComment = true
	f.state.blockDepth = 1
	f.state.blockKeep = keep
	if keep {
		f.result.WriteString(cStyleBlockCommentStart)
	}
	f.position += len(cStyleBlockCommentStart)
}

// writeBlockText appends block comment text when the active mode retains it.
func (f *kotlinLineFilter) writeBlockText(text string) {
	if f.state.blockKeep {
		f.result.WriteString(text)
	}
}

// writeEscapedByte copies an escaped byte pair from a quoted string.
func (f *kotlinLineFilter) writeEscapedByte() {
	f.result.WriteByte(f.line[f.position])
	f.position++
	if f.position < len(f.line) {
		f.result.WriteByte(f.line[f.position])
		f.position++
	}
}

// consumeCodeByte copies one source byte.
func (f *kotlinLineFilter) consumeCodeByte() {
	f.result.WriteByte(f.line[f.position])
	f.position++
}
