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

	// template reports whether scanning is inside template literal text.
	template bool

	// templateInterpolationDepth is the active brace depth of a template interpolation.
	templateInterpolationDepth int

	// nestedTemplate reports whether interpolation scanning is inside nested template text.
	nestedTemplate bool
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

// commentConsumeResult describes a consumed JavaScript comment.
type commentConsumeResult struct {
	// consumed reports whether a recognized comment marker was consumed.
	consumed bool

	// stopLine reports whether the consumed comment reaches the end of the source line.
	stopLine bool
}

// interpolationCodeResult describes the effect of one consumed interpolation byte.
type interpolationCodeResult struct {
	// depth is the brace depth after consuming the byte at the scanner position.
	depth int

	// closed reports whether the consumed byte closed the current interpolation expression.
	closed bool
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

// consumeActiveBlock consumes text while the scanner is inside a block comment.
func (f *javascriptLineFilter) consumeActiveBlock() bool {
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
	previous := previousSignificantToken(f.line[:f.position])
	if previous == "" {
		return true
	}
	if previous == "++" || previous == "--" {
		return false
	}
	if regexPrecedingKeyword(previous) {
		return true
	}
	if len(previous) != 1 {
		return false
	}

	return strings.ContainsRune("([{=,:;!&|?+-*~^<>%", rune(previous[0]))
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

// consumeInterpolationDepth filters interpolation code until depth closes or line ends.
//
// Parameters:
// depth - current brace depth of the interpolation expression; updated in place.
func (f *javascriptLineFilter) consumeInterpolationDepth(depth *int) {
	for f.position < len(f.line) {
		if f.consumeActiveBlock() {
			continue
		}
		if f.consumeNestedTemplateLiteral() {
			continue
		}
		if f.consumeString() {
			continue
		}
		if f.consumeRegexLiteral() {
			continue
		}
		if comment := f.consumeComment(); comment.consumed {
			if comment.stopLine {
				return
			}

			continue
		}
		code := f.consumeInterpolationCode(*depth)
		*depth = code.depth
		if code.closed {
			*depth = 0

			return
		}
	}
}

// consumeNestedTemplateLiteral copies a template literal found inside interpolation code.
func (f *javascriptLineFilter) consumeNestedTemplateLiteral() bool {
	if !f.startOrResumeNestedTemplateLiteral() {
		return false
	}
	for f.position < len(f.line) {
		switch {
		case f.line[f.position] == '\\':
			f.writeEscapedByte()
		case f.line[f.position] == '`':
			f.consumeCodeByte()
			f.state.nestedTemplate = false

			return true
		case strings.HasPrefix(f.line[f.position:], jsTemplateInterpolationStart):
			f.result.WriteString(jsTemplateInterpolationStart)
			f.position += len(jsTemplateInterpolationStart)
			depth := 1
			f.consumeInterpolationDepth(&depth)
			if depth > 0 {
				return true
			}
		default:
			f.consumeCodeByte()
		}
	}

	return true
}

// startOrResumeNestedTemplateLiteral enters or resumes nested template scanning.
func (f *javascriptLineFilter) startOrResumeNestedTemplateLiteral() bool {
	if f.state.nestedTemplate {
		return true
	}
	if f.position >= len(f.line) || f.line[f.position] != '`' {
		return false
	}
	f.state.nestedTemplate = true
	f.consumeCodeByte()

	return true
}

// consumeInterpolationCode copies expression code and updates interpolation brace depth.
//
// Parameters:
// depth - current brace depth before consuming the byte at the scanner position.
//
// Returns interpolation code result.
func (f *javascriptLineFilter) consumeInterpolationCode(depth int) interpolationCodeResult {
	switch f.line[f.position] {
	case '{':
		depth++
		f.consumeCodeByte()

		return interpolationCodeResult{depth: depth}
	case '}':
		depth--
		f.consumeCodeByte()

		return interpolationCodeResult{depth: depth, closed: depth == 0}
	default:
		f.consumeCodeByte()

		return interpolationCodeResult{depth: depth}
	}
}

// consumeComment consumes a JavaScript comment when one starts at the scanner position.
//
// Returns comment consume result.
func (f *javascriptLineFilter) consumeComment() commentConsumeResult {
	if strings.HasPrefix(f.line[f.position:], cStyleDocCommentStart) {
		f.startBlockComment(f.mode == RetainDocumentation)

		return commentConsumeResult{consumed: true}
	}
	if strings.HasPrefix(f.line[f.position:], cStyleBlockCommentStart) {
		f.startBlockComment(f.mode == RetainBlock || f.mode == RetainRegular)

		return commentConsumeResult{consumed: true}
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

// startBlockComment records the active block comment markers and whether to keep them.
func (f *javascriptLineFilter) startBlockComment(keep bool) {
	f.hadComment = true
	f.state.blockActive = true
	f.state.blockKeep = keep
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

// previousSignificantToken returns the last non-space token in text.
func previousSignificantToken(text string) string {
	for position := len(text) - 1; position >= 0; position-- {
		if isASCIISpace(text[position]) {
			continue
		}
		if isASCIIIdentifierByte(text[position]) {
			end := position + 1
			for position >= 0 && isASCIIIdentifierByte(text[position]) {
				position--
			}

			return text[position+1 : end]
		}
		if position > 0 {
			token := text[position-1 : position+1]
			if token == "++" || token == "--" {
				return token
			}
		}

		return text[position : position+1]
	}

	return ""
}

// regexPrecedingKeyword reports whether keyword can precede a regex literal.
func regexPrecedingKeyword(keyword string) bool {
	switch keyword {
	case "await", "case", "delete", "else", "in", "instanceof", "return",
		"throw", "typeof", "void", "yield":
		return true
	default:
		return false
	}
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
