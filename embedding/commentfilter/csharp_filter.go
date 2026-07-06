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
	csharpInterpolatedVerbatimStringStart = `$@"`
	csharpVerbatimInterpolatedStringStart = `@$"`
	csharpVerbatimStringStart             = `@"`
	csharpInterpolatedStringStart         = `$"`
	csharpEscapedQuote                    = `""`
	csharpEscapedOpenBrace                = `{{`
	csharpEscapedCloseBrace               = `}}`
)

// CSharpCommentFilter filters C# comments while preserving string literal text.
type CSharpCommentFilter struct{}

// csharpState tracks C# lexical state that can span source lines.
type csharpState struct {
	// blockActive reports whether scanning is inside a block comment.
	blockActive bool

	// blockKeep reports whether the active block comment should be retained.
	blockKeep bool

	// stringActive reports whether scanning is inside a string.
	stringActive bool

	// stringVerbatim reports whether the active string uses verbatim escaping.
	stringVerbatim bool

	// stringInterpolated reports whether the active string has interpolation holes.
	stringInterpolated bool

	// interpolationDepth is the active brace depth inside an interpolation expression.
	interpolationDepth int

	// interpolationFormat reports whether scanning is inside interpolation format text.
	interpolationFormat bool
}

// csharpLineFilter filters one C# source line.
type csharpLineFilter struct {
	// line is the source line being filtered.
	line string

	// mode selects which comments to retain.
	mode Mode

	// state tracks C# constructs across lines.
	state *csharpState

	// result accumulates the filtered source line.
	result strings.Builder

	// position is the current byte index in line.
	position int

	// hadComment reports whether the line contained a recognized comment.
	hadComment bool
}

// Filter removes or preserves C# comments according to mode.
//
// Parameters:
// lines - provides C# source lines.
// mode - selects comments to retain.
//
// Returns filtered source lines.
func (CSharpCommentFilter) Filter(lines []string, mode Mode) []string {
	var filtered []string
	state := csharpState{}
	for _, line := range lines {
		filteredLine, hadComment := filterCSharpLine(line, mode, &state)
		if hadComment && strings.TrimSpace(filteredLine) == "" {
			continue
		}
		filtered = append(filtered, filteredLine)
	}

	return filtered
}

// filterCSharpLine removes or preserves recognized C# comments from one line.
func filterCSharpLine(line string, mode Mode, state *csharpState) (string, bool) {
	filter := csharpLineFilter{
		line:  line,
		mode:  mode,
		state: state,
	}

	return filter.filterLine()
}

// filterLine walks the current line until it reaches its end or a line comment.
func (f *csharpLineFilter) filterLine() (string, bool) {
	for f.position < len(f.line) {
		if f.consumeActiveBlock() {
			continue
		}
		if f.consumeStringInterpolation() {
			continue
		}
		if f.consumeStringText() {
			continue
		}
		if f.consumeCharacterLiteral() {
			continue
		}
		if f.consumeStringStart() {
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

	f.closeSingleLineStringAtLineEnd()

	return f.result.String(), f.hadComment
}

// consumeActiveBlock consumes text while the scanner is inside a block comment.
func (f *csharpLineFilter) consumeActiveBlock() bool {
	if !f.state.blockActive {
		return false
	}
	f.hadComment = true
	end := strings.Index(f.line[f.position:], cStyleBlockCommentEnd)
	if end < 0 {
		if f.state.blockKeep {
			f.result.WriteString(f.line[f.position:])
		}
		f.position = len(f.line)

		return true
	}
	endPosition := f.position + end + len(cStyleBlockCommentEnd)
	if f.state.blockKeep {
		f.result.WriteString(f.line[f.position:endPosition])
	}
	f.position = endPosition
	f.state.blockActive = false

	return true
}

// consumeStringInterpolation filters code inside an active interpolation expression.
func (f *csharpLineFilter) consumeStringInterpolation() bool {
	if !f.state.stringActive || f.state.interpolationDepth == 0 {
		return false
	}
	for f.position < len(f.line) {
		if f.consumeInterpolationSegment() {
			if f.state.interpolationDepth == 0 {
				return true
			}

			continue
		}
		if f.consumeInterpolationCodeByte() {
			return true
		}
	}

	return true
}

// consumeInterpolationSegment consumes multi-byte interpolation content.
func (f *csharpLineFilter) consumeInterpolationSegment() bool {
	if f.consumeActiveBlock() {
		return true
	}
	if f.consumeInterpolationFormat() {
		return true
	}
	if f.consumeInterpolationString() {
		return true
	}
	if comment := f.consumeComment(); comment.consumed {
		return true
	}

	return false
}

// consumeInterpolationCodeByte copies expression code and reports whether interpolation closed.
func (f *csharpLineFilter) consumeInterpolationCodeByte() bool {
	switch f.line[f.position] {
	case '{':
		f.state.interpolationDepth++
		f.consumeCodeByte()

		return false
	case '}':
		f.state.interpolationDepth--
		f.consumeCodeByte()

		return f.state.interpolationDepth == 0
	default:
		f.consumeCodeByte()

		return false
	}
}

// consumeInterpolationFormat copies C# format text after a top-level interpolation colon.
func (f *csharpLineFilter) consumeInterpolationFormat() bool {
	if !f.state.interpolationFormat {
		if f.state.interpolationDepth != 1 || f.line[f.position] != ':' {
			return false
		}
		f.state.interpolationFormat = true
		f.consumeCodeByte()
	}
	for f.position < len(f.line) {
		if f.line[f.position] == '}' {
			f.consumeCodeByte()
			f.state.interpolationFormat = false
			f.state.interpolationDepth = 0

			return true
		}
		f.consumeCodeByte()
	}

	return true
}

// consumeInterpolationString copies a string literal inside interpolation code.
func (f *csharpLineFilter) consumeInterpolationString() bool {
	token, verbatim, found := csharpStringTokenAt(f.line, f.position)
	if !found {
		return false
	}
	f.result.WriteString(token)
	f.position += len(token)
	for f.position < len(f.line) {
		switch {
		case verbatim && strings.HasPrefix(f.line[f.position:], csharpEscapedQuote):
			f.result.WriteString(csharpEscapedQuote)
			f.position += len(csharpEscapedQuote)
		case !verbatim && f.line[f.position] == '\\':
			f.writeEscapedByte()
		case f.line[f.position] == '"':
			f.consumeCodeByte()

			return true
		default:
			f.consumeCodeByte()
		}
	}

	return true
}

// consumeStringText copies active string text without scanning comment markers inside it.
func (f *csharpLineFilter) consumeStringText() bool {
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

// consumeStringTextSegment consumes special syntax inside active string text.
func (f *csharpLineFilter) consumeStringTextSegment() bool {
	switch {
	case f.state.stringVerbatim && strings.HasPrefix(f.line[f.position:], csharpEscapedQuote):
		f.result.WriteString(csharpEscapedQuote)
		f.position += len(csharpEscapedQuote)
	case !f.state.stringVerbatim && f.line[f.position] == '\\':
		f.writeEscapedByte()
	case f.line[f.position] == '"':
		f.consumeCodeByte()
		f.closeString()
	case f.startsEscapedInterpolationBrace():
		f.result.WriteString(f.line[f.position : f.position+2])
		f.position += 2
	case f.state.stringInterpolated && f.line[f.position] == '{':
		f.consumeCodeByte()
		f.state.interpolationDepth = 1
	default:
		return false
	}

	return true
}

// consumeCharacterLiteral copies a C# character literal.
func (f *csharpLineFilter) consumeCharacterLiteral() bool {
	quoteEnd := quotedSegmentEnd(f.line, f.position, "'")
	if quoteEnd <= f.position {
		return false
	}
	f.result.WriteString(f.line[f.position:quoteEnd])
	f.position = quoteEnd

	return true
}

// consumeStringStart starts a C# string literal at the current position.
func (f *csharpLineFilter) consumeStringStart() bool {
	token, verbatim, found := csharpStringTokenAt(f.line, f.position)
	if !found {
		return false
	}
	interpolated := strings.HasPrefix(token, "$") || strings.HasPrefix(token, "@$")
	f.startString(token, verbatim, interpolated)

	return true
}

// csharpStringTokenAt returns a C# string opener at position.
func csharpStringTokenAt(line string, position int) (string, bool, bool) {
	switch {
	case strings.HasPrefix(line[position:], csharpInterpolatedVerbatimStringStart):
		return csharpInterpolatedVerbatimStringStart, true, true
	case strings.HasPrefix(line[position:], csharpVerbatimInterpolatedStringStart):
		return csharpVerbatimInterpolatedStringStart, true, true
	case strings.HasPrefix(line[position:], csharpVerbatimStringStart):
		return csharpVerbatimStringStart, true, true
	case strings.HasPrefix(line[position:], csharpInterpolatedStringStart):
		return csharpInterpolatedStringStart, false, true
	case line[position] == '"':
		return `"`, false, true
	default:
		return "", false, false
	}
}

// startString records the active string literal and copies its opening token.
func (f *csharpLineFilter) startString(token string, verbatim bool, interpolated bool) {
	f.result.WriteString(token)
	f.position += len(token)
	f.state.stringActive = true
	f.state.stringVerbatim = verbatim
	f.state.stringInterpolated = interpolated
}

// startsEscapedInterpolationBrace reports whether the position starts {{ or }} in string text.
func (f *csharpLineFilter) startsEscapedInterpolationBrace() bool {
	if !f.state.stringInterpolated || f.position+1 >= len(f.line) {
		return false
	}

	return strings.HasPrefix(f.line[f.position:], csharpEscapedOpenBrace) ||
		strings.HasPrefix(f.line[f.position:], csharpEscapedCloseBrace)
}

// closeSingleLineStringAtLineEnd clears invalid single-line strings at end of line.
func (f *csharpLineFilter) closeSingleLineStringAtLineEnd() {
	if f.state.stringActive && !f.state.stringVerbatim && f.state.interpolationDepth == 0 {
		f.closeString()
	}
}

// closeString clears the active C# string state.
func (f *csharpLineFilter) closeString() {
	f.state.stringActive = false
	f.state.stringVerbatim = false
	f.state.stringInterpolated = false
	f.state.interpolationDepth = 0
	f.state.interpolationFormat = false
}

// consumeComment consumes a C# comment when one starts at the scanner position.
func (f *csharpLineFilter) consumeComment() commentConsumeResult {
	if strings.HasPrefix(f.line[f.position:], cStyleDocCommentStart) {
		f.startBlockComment(f.mode == RetainDocumentation)

		return commentConsumeResult{consumed: true}
	}
	if strings.HasPrefix(f.line[f.position:], cStyleBlockCommentStart) {
		f.startBlockComment(f.mode == RetainBlock || f.mode == RetainRegular)

		return commentConsumeResult{consumed: true}
	}
	if strings.HasPrefix(f.line[f.position:], "///") {
		f.hadComment = true
		if f.mode == RetainDocumentation {
			f.result.WriteString(f.line[f.position:])
		}
		f.position = len(f.line)

		return commentConsumeResult{consumed: true, stopLine: true}
	}
	if strings.HasPrefix(f.line[f.position:], "//") {
		f.hadComment = true
		if f.mode == RetainInline || f.mode == RetainRegular {
			f.result.WriteString(f.line[f.position:])
		}
		f.position = len(f.line)

		return commentConsumeResult{consumed: true, stopLine: true}
	}

	return commentConsumeResult{}
}

// startBlockComment records the active block comment and whether to keep it.
func (f *csharpLineFilter) startBlockComment(keep bool) {
	f.hadComment = true
	f.state.blockActive = true
	f.state.blockKeep = keep
}

// writeEscapedByte copies an escaped byte pair from a regular string.
func (f *csharpLineFilter) writeEscapedByte() {
	f.result.WriteByte(f.line[f.position])
	f.position++
	if f.position < len(f.line) {
		f.result.WriteByte(f.line[f.position])
		f.position++
	}
}

// consumeCodeByte copies one source byte.
func (f *csharpLineFilter) consumeCodeByte() {
	f.result.WriteByte(f.line[f.position])
	f.position++
}
