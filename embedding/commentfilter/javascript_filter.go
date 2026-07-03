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

const jsTemplateInterpolationStart = "${"

// JavaScriptCommentFilter filters JavaScript and TypeScript comments while preserving literal text.
type JavaScriptCommentFilter struct{}

// javascriptState tracks JavaScript lexical state that can span source lines.
type javascriptState struct {
	// blockActive reports whether scanning is inside a block comment.
	blockActive bool

	// blockKeep reports whether the active block comment should be retained.
	blockKeep bool

	// blockEnd contains the closing marker for the active block comment.
	blockEnd string

	// template reports whether scanning is inside template literal text.
	template bool

	// templateInterpolationDepth is the active brace depth of a template interpolation.
	templateInterpolationDepth int
}

// javascriptLineFilter filters one JavaScript or TypeScript source line.
type javascriptLineFilter struct {
	// line is the source line being filtered.
	line string

	// mode selects which comments to retain.
	mode Mode

	// state tracks JavaScript constructs across lines.
	state *javascriptState

	// result accumulates the filtered source line.
	result strings.Builder

	// position is the current byte index in line.
	position int

	// hadComment reports whether the line contained a recognized comment.
	hadComment bool
}

// Filter removes or preserves JavaScript and TypeScript comments according to mode.
//
// Parameters:
// lines - provides JavaScript or TypeScript source lines.
// mode - selects comments to retain.
//
// Returns filtered source lines.
func (JavaScriptCommentFilter) Filter(lines []string, mode Mode) []string {
	var filtered []string
	state := javascriptState{}
	for _, line := range lines {
		filteredLine, hadComment := filterJavaScriptLine(line, mode, &state)
		if hadComment && strings.TrimSpace(filteredLine) == "" {
			continue
		}
		filtered = append(filtered, filteredLine)
	}

	return filtered
}

// filterJavaScriptLine removes or preserves recognized JavaScript comments from one line.
func filterJavaScriptLine(line string, mode Mode, state *javascriptState) (string, bool) {
	filter := javascriptLineFilter{
		line:  line,
		mode:  mode,
		state: state,
	}

	return filter.filterLine()
}

// filterLine walks the current line until it reaches its end or a line comment.
func (f *javascriptLineFilter) filterLine() (string, bool) {
	for f.position < len(f.line) {
		if f.consumeActiveBlock() {
			continue
		}
		if f.consumeTemplateInterpolation() {
			continue
		}
		if f.consumeTemplateText() {
			continue
		}
		if f.consumeString() {
			continue
		}
		if f.consumeRegexLiteral() {
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
func (f *javascriptLineFilter) consumeActiveBlock() bool {
	if !f.state.blockActive {
		return false
	}
	f.hadComment = true
	end := strings.Index(f.line[f.position:], f.state.blockEnd)
	if end < 0 {
		if f.state.blockKeep {
			f.result.WriteString(f.line[f.position:])
		}
		f.position = len(f.line)

		return true
	}
	endPosition := f.position + end + len(f.state.blockEnd)
	if f.state.blockKeep {
		f.result.WriteString(f.line[f.position:endPosition])
	}
	f.position = endPosition
	f.state.blockActive = false
	f.state.blockEnd = ""

	return true
}

// consumeTemplateInterpolation resumes JavaScript expression scanning inside `${...}`.
func (f *javascriptLineFilter) consumeTemplateInterpolation() bool {
	if f.state.templateInterpolationDepth == 0 {
		return false
	}
	f.consumeInterpolationDepth(&f.state.templateInterpolationDepth)
	if f.state.templateInterpolationDepth == 0 {
		f.state.template = true
	}

	return true
}

// consumeTemplateText copies template text and filters comments inside `${...}` code.
func (f *javascriptLineFilter) consumeTemplateText() bool {
	if !f.state.template && f.line[f.position] != '`' {
		return false
	}
	if !f.state.template {
		f.state.template = true
		f.consumeCodeByte()
	}
	for f.position < len(f.line) {
		switch {
		case f.line[f.position] == '\\':
			f.writeEscapedByte()
		case f.line[f.position] == '`':
			f.consumeCodeByte()
			f.state.template = false

			return true
		case strings.HasPrefix(f.line[f.position:], jsTemplateInterpolationStart):
			f.result.WriteString(jsTemplateInterpolationStart)
			f.position += len(jsTemplateInterpolationStart)
			f.state.template = false
			f.state.templateInterpolationDepth = 1
			f.consumeTemplateInterpolation()
			if f.state.templateInterpolationDepth > 0 {
				return true
			}
		default:
			f.consumeCodeByte()
		}
	}

	return true
}

// consumeString copies a quoted string without scanning comment markers inside it.
func (f *javascriptLineFilter) consumeString() bool {
	if f.position >= len(f.line) {
		return false
	}
	switch f.line[f.position] {
	case '"', '\'':
		quoteEnd := quotedSegmentEnd(f.line, f.position, "\"'")
		f.result.WriteString(f.line[f.position:quoteEnd])
		f.position = quoteEnd

		return true
	default:
		return false
	}
}

// consumeRegexLiteral copies a regular-expression literal without treating its content as comments.
func (f *javascriptLineFilter) consumeRegexLiteral() bool {
	if !f.regexStartsHere() {
		return false
	}
	f.consumeCodeByte()
	inClass := false
	for f.position < len(f.line) {
		switch f.line[f.position] {
		case '\\':
			f.writeEscapedByte()
		case '[':
			inClass = true
			f.consumeCodeByte()
		case ']':
			inClass = false
			f.consumeCodeByte()
		case '/':
			if inClass {
				f.consumeCodeByte()

				continue
			}
			f.consumeCodeByte()
			f.consumeRegexFlags()

			return true
		default:
			f.consumeCodeByte()
		}
	}

	return true
}

// regexStartsHere reports whether the slash at the current position can start a regex literal.
func (f *javascriptLineFilter) regexStartsHere() bool {
	if f.line[f.position] != '/' {
		return false
	}
	if strings.HasPrefix(f.line[f.position:], "//") ||
		strings.HasPrefix(f.line[f.position:], cStyleBlockCommentStart) {
		return false
	}
	previous := previousSignificantByte(f.line[:f.position])
	if previous == 0 {
		return true
	}

	return strings.ContainsRune("([{=,:;!&|?+-*~^<>%", rune(previous))
}

// consumeRegexFlags copies identifier characters after a regex literal closing slash.
func (f *javascriptLineFilter) consumeRegexFlags() {
	for f.position < len(f.line) {
		char := f.line[f.position]
		if !isASCIIIdentifierByte(char) {
			return
		}
		f.consumeCodeByte()
	}
}

// consumeInterpolationDepth filters comments inside interpolation code until depth closes or line ends.
func (f *javascriptLineFilter) consumeInterpolationDepth(depth *int) {
	for f.position < len(f.line) {
		if f.consumeActiveBlock() {
			continue
		}
		if f.consumeString() {
			continue
		}
		if f.consumeRegexLiteral() {
			continue
		}
		if consumed, stop := f.consumeComment(); consumed {
			if stop {
				return
			}

			continue
		}
		var done bool
		*depth, done = f.consumeInterpolationCode(*depth)
		if done {
			*depth = 0

			return
		}
	}
}

// consumeInterpolationCode copies expression code and updates interpolation brace depth.
func (f *javascriptLineFilter) consumeInterpolationCode(depth int) (int, bool) {
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

// consumeComment consumes a JavaScript comment and reports whether it ended the line.
func (f *javascriptLineFilter) consumeComment() (bool, bool) {
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

// startBlockComment records the active block comment markers and whether to keep them.
func (f *javascriptLineFilter) startBlockComment(keep bool) {
	f.hadComment = true
	f.state.blockActive = true
	f.state.blockKeep = keep
	f.state.blockEnd = cStyleBlockCommentEnd
}

// writeEscapedByte copies an escaped byte pair from a literal.
func (f *javascriptLineFilter) writeEscapedByte() {
	f.result.WriteByte(f.line[f.position])
	f.position++
	if f.position < len(f.line) {
		f.result.WriteByte(f.line[f.position])
		f.position++
	}
}

// consumeCodeByte copies one source byte.
func (f *javascriptLineFilter) consumeCodeByte() {
	f.result.WriteByte(f.line[f.position])
	f.position++
}

// previousSignificantByte returns the last non-space byte in text.
func previousSignificantByte(text string) byte {
	for position := len(text) - 1; position >= 0; position-- {
		if !isASCIISpace(text[position]) {
			return text[position]
		}
	}

	return 0
}

// isASCIISpace reports whether char is an ASCII whitespace byte.
func isASCIISpace(char byte) bool {
	return char == ' ' || char == '\t' || char == '\n' || char == '\r'
}

// isASCIIIdentifierByte reports whether char can appear in a JavaScript identifier or regex flag.
func isASCIIIdentifierByte(char byte) bool {
	return char == '_' ||
		char == '$' ||
		(char >= 'a' && char <= 'z') ||
		(char >= 'A' && char <= 'Z') ||
		(char >= '0' && char <= '9')
}
