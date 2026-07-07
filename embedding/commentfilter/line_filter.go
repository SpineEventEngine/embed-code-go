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

// interpolationForm describes how a language opens and closes interpolation expressions.
//
// start and openBrace play different roles: start is the marker that opens an
// interpolation in literal text, while openBrace is the plain brace that nests
// inside expression code. In `${items.map(i => { return i })}` the interpolation
// opens with `${`, but the lambda `{` — which no longer follows a `$` — must
// still deepen the nesting so that its matching closeBrace does not end the
// interpolation early.
type interpolationForm struct {
	// start opens an interpolation expression in literal text, e.g. `${`.
	start string

	// openBrace nests the expression code one level deeper.
	openBrace byte

	// closeBrace closes one nesting level.
	//
	// The expression ends at the closeBrace that returns the nesting depth to zero.
	closeBrace byte
}

// interpolatedLiteral describes a string literal form that embeds interpolation expressions.
type interpolatedLiteral struct {
	// delimiter opens and closes the literal text.
	delimiter string

	// escapes reports whether backslashes escape literal bytes.
	escapes bool

	// interpolation opens and closes interpolation expressions inside the literal.
	interpolation interpolationForm
}

// lineFilter carries the scanning state shared by the per-language line filters.
type lineFilter struct {
	// line is the source line being filtered.
	line string

	// mode selects which comments to retain.
	mode Mode

	// result accumulates the filtered source line.
	result strings.Builder

	// position is the current byte index in line.
	position int

	// hadComment reports whether the line contained a recognized comment.
	hadComment bool
}

// commentConsumeResult describes a consumed source comment.
type commentConsumeResult struct {
	// consumed reports whether a recognized comment marker was consumed.
	consumed bool

	// stopLine reports whether the consumed comment reaches the end of the source line.
	stopLine bool
}

// blockCommentState tracks a non-nested block comment across source lines.
type blockCommentState struct {
	// active reports whether scanning is inside a block comment.
	active bool

	// keep reports whether the active block comment should be retained.
	keep bool
}

// hasPrefix reports whether text starts at the scanner position.
func (f *lineFilter) hasPrefix(text string) bool {
	return strings.HasPrefix(f.line[f.position:], text)
}

// consumeCodeByte copies one source byte.
func (f *lineFilter) consumeCodeByte() {
	f.result.WriteByte(f.line[f.position])
	f.position++
}

// consumeMarker copies marker text and advances past it.
func (f *lineFilter) consumeMarker(marker string) {
	f.result.WriteString(marker)
	f.position += len(marker)
}

// writeEscapedByte copies an escaped byte pair from a literal.
func (f *lineFilter) writeEscapedByte() {
	f.result.WriteByte(f.line[f.position])
	f.position++
	if f.position < len(f.line) {
		f.result.WriteByte(f.line[f.position])
		f.position++
	}
}

// consumeQuotedSegment copies a quoted literal when one starts at the scanner position.
func (f *lineFilter) consumeQuotedSegment(quoteChars string) bool {
	quoteEnd := quotedSegmentEnd(f.line, f.position, quoteChars)
	if quoteEnd <= f.position {
		return false
	}
	f.result.WriteString(f.line[f.position:quoteEnd])
	f.position = quoteEnd

	return true
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

// consumeLineComment consumes the rest of the line as a line comment.
func (f *lineFilter) consumeLineComment(keep bool) {
	f.hadComment = true
	if keep {
		f.result.WriteString(f.line[f.position:])
	}
	f.position = len(f.line)
}

// consumeCStyleComment consumes a C-style comment when one starts at the scanner position.
//
// Parameters:
// docLineMarker - optional documentation line-comment marker such as `///`; empty when absent.
// startBlock - language hook that records an opened block comment.
//
// Returns comment consume result.
func (f *lineFilter) consumeCStyleComment(
	docLineMarker string,
	startBlock func(keep bool),
) commentConsumeResult {
	if f.hasPrefix(cStyleDocCommentStart) {
		startBlock(f.mode == RetainDocumentation)

		return commentConsumeResult{consumed: true}
	}
	if f.hasPrefix(cStyleBlockCommentStart) {
		startBlock(f.mode == RetainBlock || f.mode == RetainRegular)

		return commentConsumeResult{consumed: true}
	}
	if docLineMarker != "" && f.hasPrefix(docLineMarker) {
		f.consumeLineComment(f.mode == RetainDocumentation)

		return commentConsumeResult{consumed: true, stopLine: true}
	}
	if f.hasPrefix("//") {
		f.consumeLineComment(f.mode == RetainInline || f.mode == RetainRegular)

		return commentConsumeResult{consumed: true, stopLine: true}
	}

	return commentConsumeResult{}
}

// startBlockComment records and consumes an opened non-nested C-style block comment.
func (f *lineFilter) startBlockComment(state *blockCommentState, keep bool) {
	f.hadComment = true
	state.active = true
	state.keep = keep
	if keep {
		f.result.WriteString(cStyleBlockCommentStart)
	}
	f.position += len(cStyleBlockCommentStart)
}

// consumeActiveBlock consumes text while the scanner is inside a non-nested block comment.
func (f *lineFilter) consumeActiveBlock(state *blockCommentState) bool {
	if !state.active {
		return false
	}
	f.hadComment = true
	end := strings.Index(f.line[f.position:], cStyleBlockCommentEnd)
	if end < 0 {
		if state.keep {
			f.result.WriteString(f.line[f.position:])
		}
		f.position = len(f.line)

		return true
	}
	endPosition := f.position + end + len(cStyleBlockCommentEnd)
	if state.keep {
		f.result.WriteString(f.line[f.position:endPosition])
	}
	f.position = endPosition
	state.active = false

	return true
}

// consumeInterpolationCodeByte copies one expression byte and updates interpolation brace depth.
//
// Returns true when the consumed byte closed the interpolation expression.
func (f *lineFilter) consumeInterpolationCodeByte(form interpolationForm, depth *int) bool {
	char := f.line[f.position]
	f.consumeCodeByte()
	switch char {
	case form.openBrace:
		*depth++
	case form.closeBrace:
		*depth--

		return *depth == 0
	}

	return false
}

// consumeInterpolatedText copies interpolated literal text and enters interpolation code.
//
// Parameters:
// literal - describes the literal delimiter, escape rule, and interpolation opener.
// active - tracks whether the scanner is inside the literal text across lines.
// depth - tracks the interpolation brace depth across lines.
// resumeInterpolation - language hook that scans interpolation code until depth closes.
//
// Returns true when literal text was consumed.
func (f *lineFilter) consumeInterpolatedText(
	literal interpolatedLiteral,
	active *bool,
	depth *int,
	resumeInterpolation func() bool,
) bool {
	if !*active && !f.hasPrefix(literal.delimiter) {
		return false
	}
	if !*active {
		*active = true
		f.consumeMarker(literal.delimiter)
	}
	for f.position < len(f.line) {
		if f.consumeInterpolatedTextSegment(literal, active, depth, resumeInterpolation) {
			return true
		}
	}

	return true
}

// consumeInterpolatedTextSegment consumes one piece of interpolated literal text.
//
// Returns true when the literal closed or its interpolation stayed open at the line end.
func (f *lineFilter) consumeInterpolatedTextSegment(
	literal interpolatedLiteral,
	active *bool,
	depth *int,
	resumeInterpolation func() bool,
) bool {
	switch {
	case literal.escapes && f.line[f.position] == '\\':
		f.writeEscapedByte()
	case f.hasPrefix(literal.delimiter):
		f.consumeMarker(literal.delimiter)
		*active = false

		return true
	case f.hasPrefix(literal.interpolation.start):
		f.consumeMarker(literal.interpolation.start)
		*active = false
		*depth = 1
		resumeInterpolation()

		return *depth > 0
	default:
		f.consumeCodeByte()
	}

	return false
}
