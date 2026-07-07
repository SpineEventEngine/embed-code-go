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

const (
	pythonEscapedOpenBrace  = "{{"
	pythonEscapedCloseBrace = "}}"
)

// pythonInterpolation describes the Python f-string `{...}` interpolation form.
var pythonInterpolation = interpolationForm{
	start:      "{",
	openBrace:  '{',
	closeBrace: '}',
}

// pythonStringStart describes a Python string opening token.
type pythonStringStart struct {
	// token is the source text that opens the string.
	token string

	// delimiter closes the string.
	delimiter string

	// multiline reports whether the string can continue across source lines.
	multiline bool

	// interpolated reports whether the string contains f-string expression holes.
	interpolated bool
}

// PythonCommentFilter filters Python comments while preserving string literal text.
type PythonCommentFilter struct{}

// pythonState tracks Python lexical state that can span source lines.
type pythonState struct {
	// stringActive reports whether scanning is inside string text.
	stringActive bool

	// stringDelimiter closes the active string.
	stringDelimiter string

	// stringMultiline reports whether the active string can continue across source lines.
	stringMultiline bool

	// stringInterpolated reports whether the active string has interpolation holes.
	stringInterpolated bool

	// interpolationDepth is the active brace depth inside an f-string expression.
	interpolationDepth int

	// interpolationFormat reports whether scanning is inside f-string format text.
	interpolationFormat bool

	// interpolationParenDepth is the parenthesis depth inside an f-string expression.
	interpolationParenDepth int

	// interpolationBracketDepth is the bracket depth inside an f-string expression.
	interpolationBracketDepth int

	// expressionStringActive reports whether expression scanning is inside a string.
	expressionStringActive bool

	// expressionStringDelimiter closes the active expression string.
	expressionStringDelimiter string

	// expressionStringMultiline reports whether the expression string can span lines.
	expressionStringMultiline bool
}

// pythonLineFilter filters one Python source line.
type pythonLineFilter struct {
	// lineFilter provides shared line scanning state and helpers.
	lineFilter

	// state tracks Python constructs across lines.
	state *pythonState
}

// Filter removes or preserves Python comments according to mode.
//
// Parameters:
// lines - provides Python source lines.
// mode - selects comments to retain.
//
// Returns filtered source lines.
func (PythonCommentFilter) Filter(lines []string, mode Mode) []string {
	state := pythonState{}

	return filterLines(lines, func(line string) (string, bool) {
		filter := pythonLineFilter{
			lineFilter: lineFilter{line: line, mode: mode},
			state:      &state,
		}

		return filter.filterLine()
	})
}

// filterLine walks the current line until it reaches its end or a line comment.
func (f *pythonLineFilter) filterLine() (string, bool) {
	for f.position < len(f.line) {
		if f.consumeStringInterpolation() {
			continue
		}
		if f.consumeStringText() {
			continue
		}
		if f.consumeStringStart() {
			continue
		}
		if f.consumeComment().consumed {
			break
		}
		f.consumeCodeByte()
	}

	f.closeSingleLineStateAtLineEnd()

	return f.result.String(), f.hadComment
}

// consumeStringInterpolation filters code inside an active f-string expression.
func (f *pythonLineFilter) consumeStringInterpolation() bool {
	if f.state.interpolationDepth == 0 {
		return false
	}
	for f.position < len(f.line) {
		if f.consumeInterpolationSegment() {
			if f.state.interpolationDepth == 0 {
				f.state.stringActive = true

				return true
			}

			continue
		}
		if f.consumeInterpolationCodeByte() {
			f.state.stringActive = true

			return true
		}
	}

	return true
}

// consumeInterpolationSegment consumes multi-byte f-string expression content.
func (f *pythonLineFilter) consumeInterpolationSegment() bool {
	if f.consumeExpressionString() {
		return true
	}
	if f.consumeInterpolationFormat() {
		return true
	}
	if comment := f.consumeComment(); comment.consumed {
		return true
	}

	return false
}

// consumeInterpolationFormat copies f-string format text after a top-level colon.
func (f *pythonLineFilter) consumeInterpolationFormat() bool {
	if !f.state.interpolationFormat {
		if !f.startsFormatSpec() {
			return false
		}
		f.state.interpolationFormat = true
		f.consumeCodeByte()
	}
	for f.position < len(f.line) {
		if f.line[f.position] == '}' {
			f.consumeCodeByte()
			f.closeInterpolation()

			return true
		}
		f.consumeCodeByte()
	}

	return true
}

// startsFormatSpec reports whether the current colon starts f-string format text.
func (f *pythonLineFilter) startsFormatSpec() bool {
	return f.state.interpolationDepth == 1 &&
		f.state.interpolationParenDepth == 0 &&
		f.state.interpolationBracketDepth == 0 &&
		f.line[f.position] == ':'
}

// consumeInterpolationCodeByte copies one expression byte and updates Python nesting state.
func (f *pythonLineFilter) consumeInterpolationCodeByte() bool {
	char := f.line[f.position]
	f.consumeCodeByte()
	switch char {
	case '(':
		f.state.interpolationParenDepth++
	case ')':
		if f.state.interpolationParenDepth > 0 {
			f.state.interpolationParenDepth--
		}
	case '[':
		f.state.interpolationBracketDepth++
	case ']':
		if f.state.interpolationBracketDepth > 0 {
			f.state.interpolationBracketDepth--
		}
	case pythonInterpolation.openBrace:
		f.state.interpolationDepth++
	case pythonInterpolation.closeBrace:
		f.state.interpolationDepth--
		if f.state.interpolationDepth == 0 {
			f.closeInterpolation()

			return true
		}
	}

	return false
}

// consumeExpressionString copies strings nested inside an f-string expression.
func (f *pythonLineFilter) consumeExpressionString() bool {
	if f.state.expressionStringActive {
		return f.consumeActiveExpressionString()
	}
	start, found := pythonStringStartAt(f.line, f.position)
	if !found {
		return false
	}
	f.consumeMarker(start.token)
	f.state.expressionStringActive = true
	f.state.expressionStringDelimiter = start.delimiter
	f.state.expressionStringMultiline = start.multiline

	return true
}

// consumeActiveExpressionString copies active expression string text.
func (f *pythonLineFilter) consumeActiveExpressionString() bool {
	for f.position < len(f.line) {
		switch {
		case f.line[f.position] == '\\':
			f.writeEscapedByte()
		case f.hasPrefix(f.state.expressionStringDelimiter):
			f.consumeMarker(f.state.expressionStringDelimiter)
			f.closeExpressionString()

			return true
		default:
			f.consumeCodeByte()
		}
	}

	return true
}

// consumeStringText copies active Python string text without scanning comment markers inside it.
func (f *pythonLineFilter) consumeStringText() bool {
	if !f.state.stringActive || f.state.interpolationDepth > 0 {
		return false
	}
	for f.position < len(f.line) {
		if f.consumeStringTextSegment() {
			return true
		}
		f.consumeCodeByte()
	}

	return true
}

// consumeStringTextSegment consumes special syntax inside active Python string text.
func (f *pythonLineFilter) consumeStringTextSegment() bool {
	switch {
	case f.line[f.position] == '\\':
		f.writeEscapedByte()
	case f.hasPrefix(f.state.stringDelimiter):
		f.consumeMarker(f.state.stringDelimiter)
		f.closeString()
	case f.startsEscapedInterpolationBrace():
		f.consumeMarker(f.line[f.position : f.position+2])
	case f.state.stringInterpolated && f.line[f.position] == '{':
		f.consumeCodeByte()
		f.state.stringActive = false
		f.state.interpolationDepth = 1
	default:
		return false
	}

	return true
}

// consumeStringStart starts a Python string literal at the current position.
func (f *pythonLineFilter) consumeStringStart() bool {
	start, found := pythonStringStartAt(f.line, f.position)
	if !found {
		return false
	}
	f.consumeMarker(start.token)
	f.state.stringActive = true
	f.state.stringDelimiter = start.delimiter
	f.state.stringMultiline = start.multiline
	f.state.stringInterpolated = start.interpolated

	return true
}

// startsEscapedInterpolationBrace reports whether the position starts {{ or }} in f-string text.
func (f *pythonLineFilter) startsEscapedInterpolationBrace() bool {
	if !f.state.stringInterpolated || f.position+1 >= len(f.line) {
		return false
	}

	return f.hasPrefix(pythonEscapedOpenBrace) || f.hasPrefix(pythonEscapedCloseBrace)
}

// closeSingleLineStateAtLineEnd clears line-local Python string state at end of line.
func (f *pythonLineFilter) closeSingleLineStateAtLineEnd() {
	if f.state.stringActive && !f.state.stringMultiline && f.state.interpolationDepth == 0 {
		f.closeString()
	}
	if f.state.expressionStringActive && !f.state.expressionStringMultiline {
		f.closeExpressionString()
	}
}

// closeString clears active Python string state.
func (f *pythonLineFilter) closeString() {
	f.state.stringActive = false
	f.state.stringDelimiter = ""
	f.state.stringMultiline = false
	f.state.stringInterpolated = false
	f.closeInterpolation()
}

// closeExpressionString clears active f-string expression string state.
func (f *pythonLineFilter) closeExpressionString() {
	f.state.expressionStringActive = false
	f.state.expressionStringDelimiter = ""
	f.state.expressionStringMultiline = false
}

// closeInterpolation clears active Python f-string expression state.
func (f *pythonLineFilter) closeInterpolation() {
	f.state.interpolationDepth = 0
	f.state.interpolationFormat = false
	f.state.interpolationParenDepth = 0
	f.state.interpolationBracketDepth = 0
}

// consumeComment consumes a Python line comment when one starts at the scanner position.
func (f *pythonLineFilter) consumeComment() commentConsumeResult {
	if !f.hasPrefix("#") {
		return commentConsumeResult{}
	}
	f.consumeLineComment(f.mode == RetainAll)

	return commentConsumeResult{consumed: true, stopLine: true}
}

// pythonStringStartAt returns a Python string opener at position.
func pythonStringStartAt(line string, position int) (pythonStringStart, bool) {
	for _, prefixLength := range []int{2, 1, 0} {
		if position+prefixLength >= len(line) {
			continue
		}
		prefix := line[position : position+prefixLength]
		interpolated, valid := pythonStringPrefix(prefix)
		if !valid {
			continue
		}
		if delimiter, multiline, found := pythonStringDelimiterAt(line, position+prefixLength); found {
			return pythonStringStart{
				token:        prefix + delimiter,
				delimiter:    delimiter,
				multiline:    multiline,
				interpolated: interpolated,
			}, true
		}
	}

	return pythonStringStart{}, false
}

// pythonStringDelimiterAt returns a Python string delimiter at position.
func pythonStringDelimiterAt(line string, position int) (string, bool, bool) {
	switch {
	case strings.HasPrefix(line[position:], pythonDoubleQuoteBlock):
		return pythonDoubleQuoteBlock, true, true
	case strings.HasPrefix(line[position:], pythonSingleQuoteBlock):
		return pythonSingleQuoteBlock, true, true
	case line[position] == '"':
		return `"`, false, true
	case line[position] == '\'':
		return `'`, false, true
	default:
		return "", false, false
	}
}

// pythonStringPrefix reports whether prefix can precede a Python string literal.
func pythonStringPrefix(prefix string) (bool, bool) {
	switch strings.ToLower(prefix) {
	case "":
		return false, true
	case "f":
		return true, true
	case "r", "u", "b", "br", "rb":
		return false, true
	case "fr", "rf":
		return true, true
	default:
		return false, false
	}
}
