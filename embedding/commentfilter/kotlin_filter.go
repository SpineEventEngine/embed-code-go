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

const kotlinRawStringDelimiter = "\"\"\""

// kotlinInterpolation describes the Kotlin `${...}` string interpolation form.
var kotlinInterpolation = interpolationForm{
	start:      "${",
	openBrace:  '{',
	closeBrace: '}',
}

// kotlinRawStringLiteral describes the Kotlin raw triple-quoted string form.
// Raw strings have no backslash escapes.
var kotlinRawStringLiteral = interpolatedLiteral{
	delimiter:     kotlinRawStringDelimiter,
	interpolation: kotlinInterpolation,
}

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

	// rawInterpolationDepth is the active brace depth of a raw-string interpolation.
	rawInterpolationDepth int
}

// kotlinLineFilter filters one Kotlin source line.
type kotlinLineFilter struct {
	// lineFilter provides shared line scanning state and helpers.
	lineFilter

	// state tracks Kotlin constructs across lines.
	state *kotlinState
}

// Filter removes or preserves Kotlin comments according to mode.
//
// Parameters:
// lines - provides Kotlin source lines.
// mode - selects comments to retain.
//
// Returns filtered source lines.
func (KotlinCommentFilter) Filter(lines []string, mode Mode) []string {
	state := kotlinState{}

	return filterLines(lines, func(line string) (string, bool) {
		filter := kotlinLineFilter{
			lineFilter: lineFilter{line: line, mode: mode},
			state:      &state,
		}

		return filter.filterLine()
	})
}

// filterLine walks the current line until it reaches its end or a line comment.
func (f *kotlinLineFilter) filterLine() (string, bool) {
	for f.position < len(f.line) {
		if f.consumeActiveBlock() {
			continue
		}
		if f.consumeRawInterpolation() {
			continue
		}
		if f.consumeRawString() {
			continue
		}
		if f.consumeString() {
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

// consumeActiveBlock consumes a possibly nested Kotlin block comment.
func (f *kotlinLineFilter) consumeActiveBlock() bool {
	if f.state.blockDepth == 0 {
		return false
	}
	f.hadComment = true
	for f.position < len(f.line) {
		switch {
		case f.hasPrefix(cStyleBlockCommentStart):
			f.writeBlockText(cStyleBlockCommentStart)
			f.state.blockDepth++
			f.position += len(cStyleBlockCommentStart)
		case f.hasPrefix(cStyleBlockCommentEnd):
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
	return f.consumeInterpolatedText(
		kotlinRawStringLiteral,
		&f.state.rawString,
		&f.state.rawInterpolationDepth,
		f.consumeRawInterpolation,
	)
}

// consumeRawInterpolation resumes Kotlin expression scanning inside a raw-string interpolation.
func (f *kotlinLineFilter) consumeRawInterpolation() bool {
	if f.state.rawInterpolationDepth == 0 {
		return false
	}
	f.consumeInterpolationDepth(&f.state.rawInterpolationDepth)
	if f.state.rawInterpolationDepth == 0 {
		f.state.rawString = true
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
		return f.consumeQuotedSegment("'")
	default:
		return false
	}
}

// consumeQuotedString copies a Kotlin quoted string and filters comments inside `${...}`.
func (f *kotlinLineFilter) consumeQuotedString() {
	f.consumeCodeByte()
	for f.position < len(f.line) {
		switch {
		case f.line[f.position] == '\\':
			f.writeEscapedByte()
		case f.line[f.position] == '"':
			f.consumeCodeByte()

			return
		case f.hasPrefix(kotlinInterpolation.start):
			f.consumeMarker(kotlinInterpolation.start)
			f.consumeInterpolation()
		default:
			f.consumeCodeByte()
		}
	}
}

// consumeInterpolation filters comments inside a Kotlin string interpolation expression.
func (f *kotlinLineFilter) consumeInterpolation() {
	depth := 1
	f.consumeInterpolationDepth(&depth)
}

// consumeInterpolationDepth filters comments inside interpolation code
// until depth closes or line ends.
func (f *kotlinLineFilter) consumeInterpolationDepth(depth *int) {
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
		if comment := f.consumeComment(); comment.consumed {
			if comment.stopLine {
				return
			}

			continue
		}
		if f.consumeInterpolationCodeByte(kotlinInterpolation, depth) {
			return
		}
	}
}

// consumeComment consumes a Kotlin comment when one starts at the scanner position.
func (f *kotlinLineFilter) consumeComment() commentConsumeResult {
	return f.consumeCStyleComment("", f.startBlockComment)
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
