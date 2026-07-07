/*
 * Copyright 2026, TeamDev. All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 * Redistribution and use in source and/or binary forms, with or without
 * modification, must retain the above copyright notice and the following
 * disclaimer.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
 * "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
 * LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
 * A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
 * OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
 * SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
 * LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
 * DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
 * THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 * OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

package commentfilter

import "strings"

const jsTemplateDelimiter = "`"

// jsInterpolation describes the JavaScript `${...}` template interpolation form.
var jsInterpolation = interpolationForm{
	start:      "${",
	openBrace:  '{',
	closeBrace: '}',
}

// jsTemplateLiteral describes the JavaScript template literal form.
var jsTemplateLiteral = interpolatedLiteral{
	delimiter:     jsTemplateDelimiter,
	escapes:       true,
	interpolation: jsInterpolation,
}

// JavaScriptCommentFilter filters JavaScript and TypeScript comments while preserving literal text.
type JavaScriptCommentFilter struct{}

// javascriptState tracks JavaScript lexical state that can span source lines.
type javascriptState struct {
	// block tracks a block comment across source lines.
	block blockCommentState

	// template reports whether scanning is inside template literal text.
	template bool

	// templateInterpolationDepth is the active brace depth of a template interpolation.
	templateInterpolationDepth int

	// nestedTemplate reports whether interpolation scanning is inside nested template text.
	nestedTemplate bool

	// nestedTemplateInterpolationDepth is the brace depth inside nested template code.
	nestedTemplateInterpolationDepth int
}

// javascriptLineFilter filters one JavaScript or TypeScript source line.
type javascriptLineFilter struct {
	// lineFilter provides shared line scanning state and helpers.
	lineFilter

	// state tracks JavaScript constructs across lines.
	state *javascriptState
}

// Filter removes or preserves JavaScript and TypeScript comments according to mode.
//
// Parameters:
// lines - provides JavaScript or TypeScript source lines.
// mode - selects comments to retain.
//
// Returns filtered source lines.
func (JavaScriptCommentFilter) Filter(lines []string, mode Mode) []string {
	state := javascriptState{}

	return filterLines(lines, func(line string) (string, bool) {
		filter := javascriptLineFilter{
			lineFilter: lineFilter{line: line, mode: mode},
			state:      &state,
		}

		return filter.filterLine()
	})
}

// filterLine walks the current line until it reaches its end or a line comment.
func (f *javascriptLineFilter) filterLine() (string, bool) {
	for f.position < len(f.line) {
		if f.consumeActiveBlock(&f.state.block) {
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

// consumeTemplateInterpolation resumes JavaScript expression scanning inside `${...}`.
func (f *javascriptLineFilter) consumeTemplateInterpolation() bool {
	if f.state.templateInterpolationDepth == 0 {
		return false
	}
	if f.consumeNestedTemplateInterpolation() {
		return true
	}
	f.consumeInterpolationDepth(&f.state.templateInterpolationDepth)
	if f.state.templateInterpolationDepth == 0 {
		f.state.template = true
	}

	return true
}

// consumeTemplateText copies template text and filters comments inside `${...}` code.
func (f *javascriptLineFilter) consumeTemplateText() bool {
	return f.consumeInterpolatedText(
		jsTemplateLiteral,
		&f.state.template,
		&f.state.templateInterpolationDepth,
		f.consumeTemplateInterpolation,
	)
}

// consumeString copies a quoted string without scanning comment markers inside it.
func (f *javascriptLineFilter) consumeString() bool {
	return f.consumeQuotedSegment("\"'")
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
	if f.hasPrefix("//") || f.hasPrefix(cStyleBlockCommentStart) {
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
		if f.consumeActiveBlock(&f.state.block) {
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
		if f.consumeInterpolationCodeByte(jsInterpolation, depth) {
			return
		}
	}
}

// consumeNestedTemplateLiteral copies a template literal found inside interpolation code.
func (f *javascriptLineFilter) consumeNestedTemplateLiteral() bool {
	return f.consumeInterpolatedText(
		jsTemplateLiteral,
		&f.state.nestedTemplate,
		&f.state.nestedTemplateInterpolationDepth,
		f.consumeNestedTemplateInterpolation,
	)
}

// consumeNestedTemplateInterpolation resumes code inside nested template `${...}`.
func (f *javascriptLineFilter) consumeNestedTemplateInterpolation() bool {
	if f.state.nestedTemplateInterpolationDepth == 0 {
		return false
	}
	f.consumeInterpolationDepth(&f.state.nestedTemplateInterpolationDepth)
	if f.state.nestedTemplateInterpolationDepth == 0 {
		f.state.nestedTemplate = true
	}

	return true
}

// consumeComment consumes a JavaScript comment when one starts at the scanner position.
func (f *javascriptLineFilter) consumeComment() commentConsumeResult {
	return f.consumeCStyleComment("", f.startBlockComment)
}

// startBlockComment records the active block comment and whether to keep it.
func (f *javascriptLineFilter) startBlockComment(keep bool) {
	f.lineFilter.startBlockComment(&f.state.block, keep)
}

// previousSignificantToken returns the last non-space token in text.
func previousSignificantToken(text string) string {
	for position := len(text) - 1; position >= 0; position-- {
		if isASCIISpace(text[position]) {
			continue
		}
		if previousPosition, skipped := skipTrailingBlockComment(text, position); skipped {
			position = previousPosition + 1

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

// skipTrailingBlockComment skips a block comment ending at position.
func skipTrailingBlockComment(text string, position int) (int, bool) {
	if position < len(cStyleBlockCommentEnd)-1 ||
		text[position-len(cStyleBlockCommentEnd)+1:position+1] != cStyleBlockCommentEnd {
		return position, false
	}
	start := strings.LastIndex(text[:position-len(cStyleBlockCommentEnd)+1], cStyleBlockCommentStart)
	if start < 0 {
		return position, false
	}

	return start - 1, true
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
