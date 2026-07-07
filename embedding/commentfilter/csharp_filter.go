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
	csharpInterpolatedRawStringStart      = `$"""`
	csharpRawStringDelimiter              = `"""`
	csharpEscapedQuote                    = `""`
	csharpEscapedOpenBrace                = `{{`
	csharpEscapedCloseBrace               = `}}`
	csharpDocLineComment                  = `///`
)

// csharpInterpolation describes the C# `{...}` string interpolation form.
// Holes open with a bare `{` inside interpolated string text.
var csharpInterpolation = interpolationForm{
	start:      "{",
	openBrace:  '{',
	closeBrace: '}',
}

// csharpStringStart describes a C# string opening token.
type csharpStringStart struct {
	// token is the source text that opens the string.
	token string

	// verbatim reports whether the string uses verbatim escaping.
	verbatim bool

	// interpolated reports whether the string contains interpolation holes.
	interpolated bool

	// rawDelimiter closes a raw string; empty for regular and verbatim strings.
	rawDelimiter string
}

// CSharpCommentFilter filters C# comments while preserving string literal text.
type CSharpCommentFilter struct{}

// csharpState tracks C# lexical state that can span source lines.
type csharpState struct {
	// block tracks a block comment across source lines.
	block blockCommentState

	// stringActive reports whether scanning is inside a string.
	stringActive bool

	// stringVerbatim reports whether the active string uses verbatim escaping.
	stringVerbatim bool

	// stringInterpolated reports whether the active string has interpolation holes.
	stringInterpolated bool

	// stringRawDelimiter closes an active raw string.
	stringRawDelimiter string

	// interpolationDepth is the active brace depth inside an interpolation expression.
	interpolationDepth int

	// interpolationFormat reports whether scanning is inside interpolation format text.
	interpolationFormat bool
}

// csharpLineFilter filters one C# source line.
type csharpLineFilter struct {
	lineFilter

	// state tracks C# constructs across lines.
	state *csharpState
}

// Filter removes or preserves C# comments according to mode.
//
// Parameters:
// lines - provides C# source lines.
// mode - selects comments to retain.
//
// Returns filtered source lines.
func (CSharpCommentFilter) Filter(lines []string, mode Mode) []string {
	state := csharpState{}

	return filterLines(lines, func(line string) (string, bool) {
		filter := csharpLineFilter{
			lineFilter: lineFilter{line: line, mode: mode},
			state:      &state,
		}

		return filter.filterLine()
	})
}

// filterLine walks the current line until it reaches its end or a line comment.
func (f *csharpLineFilter) filterLine() (string, bool) {
	for f.position < len(f.line) {
		if f.consumeActiveBlock(&f.state.block) {
			continue
		}
		if f.consumeStringInterpolation() {
			continue
		}
		if f.consumeStringText() {
			continue
		}
		if f.consumeQuotedSegment("'") {
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
		if f.consumeInterpolationCodeByte(csharpInterpolation, &f.state.interpolationDepth) {
			return true
		}
	}

	return true
}

// consumeInterpolationSegment consumes multi-byte interpolation content.
func (f *csharpLineFilter) consumeInterpolationSegment() bool {
	if f.consumeActiveBlock(&f.state.block) {
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

// consumeInterpolationFormat copies C# format text after a top-level interpolation colon.
// This lightweight scanner does not track parentheses inside interpolation expressions,
// so a ternary colon in parentheses is treated as format text.
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

// consumeInterpolationString copies a line-local string literal inside interpolation code.
// Nested multi-line verbatim strings intentionally resume as interpolation code on the next line.
func (f *csharpLineFilter) consumeInterpolationString() bool {
	start, found := csharpStringStartAt(f.line, f.position)
	if !found {
		return false
	}
	f.consumeMarker(start.token)
	if start.rawDelimiter != "" {
		return f.consumeLineLocalRawString(start.rawDelimiter)
	}
	for f.position < len(f.line) {
		switch {
		case start.verbatim && f.hasPrefix(csharpEscapedQuote):
			f.consumeMarker(csharpEscapedQuote)
		case !start.verbatim && f.line[f.position] == '\\':
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

// consumeLineLocalRawString copies raw string text until the closing delimiter or line end.
func (f *csharpLineFilter) consumeLineLocalRawString(delimiter string) bool {
	for f.position < len(f.line) {
		if f.hasPrefix(delimiter) {
			f.consumeMarker(delimiter)

			return true
		}
		f.consumeCodeByte()
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
	if f.state.stringRawDelimiter != "" {
		return f.consumeRawStringTextSegment()
	}

	return f.consumeRegularStringTextSegment()
}

// consumeRawStringTextSegment consumes special syntax inside active raw string text.
func (f *csharpLineFilter) consumeRawStringTextSegment() bool {
	switch {
	case f.hasPrefix(f.state.stringRawDelimiter):
		f.consumeMarker(f.state.stringRawDelimiter)
		f.closeString()
	case f.startsEscapedInterpolationBrace():
		f.consumeMarker(f.line[f.position : f.position+2])
	case f.state.stringInterpolated && f.line[f.position] == '{':
		f.consumeCodeByte()
		f.state.interpolationDepth = 1
	default:
		return false
	}

	return true
}

// consumeRegularStringTextSegment consumes special syntax inside regular string text.
func (f *csharpLineFilter) consumeRegularStringTextSegment() bool {
	switch {
	case f.state.stringVerbatim && f.hasPrefix(csharpEscapedQuote):
		f.consumeMarker(csharpEscapedQuote)
	case !f.state.stringVerbatim && f.line[f.position] == '\\':
		f.writeEscapedByte()
	case f.line[f.position] == '"':
		f.consumeCodeByte()
		f.closeString()
	case f.startsEscapedInterpolationBrace():
		f.consumeMarker(f.line[f.position : f.position+2])
	case f.state.stringInterpolated && f.line[f.position] == '{':
		f.consumeCodeByte()
		f.state.interpolationDepth = 1
	default:
		return false
	}

	return true
}

// consumeStringStart starts a C# string literal at the current position.
func (f *csharpLineFilter) consumeStringStart() bool {
	start, found := csharpStringStartAt(f.line, f.position)
	if !found {
		return false
	}
	f.consumeMarker(start.token)
	f.state.stringActive = true
	f.state.stringVerbatim = start.verbatim
	f.state.stringInterpolated = start.interpolated
	f.state.stringRawDelimiter = start.rawDelimiter

	return true
}

// csharpStringStartAt returns a C# string opener at position.
func csharpStringStartAt(line string, position int) (csharpStringStart, bool) {
	switch {
	case strings.HasPrefix(line[position:], csharpInterpolatedRawStringStart):
		return csharpStringStart{
			token:        csharpInterpolatedRawStringStart,
			interpolated: true,
			rawDelimiter: csharpRawStringDelimiter,
		}, true
	case strings.HasPrefix(line[position:], csharpInterpolatedVerbatimStringStart):
		return csharpStringStart{
			token:        csharpInterpolatedVerbatimStringStart,
			verbatim:     true,
			interpolated: true,
		}, true
	case strings.HasPrefix(line[position:], csharpVerbatimInterpolatedStringStart):
		return csharpStringStart{
			token:        csharpVerbatimInterpolatedStringStart,
			verbatim:     true,
			interpolated: true,
		}, true
	case strings.HasPrefix(line[position:], csharpVerbatimStringStart):
		return csharpStringStart{
			token:    csharpVerbatimStringStart,
			verbatim: true,
		}, true
	case strings.HasPrefix(line[position:], csharpInterpolatedStringStart):
		return csharpStringStart{
			token:        csharpInterpolatedStringStart,
			interpolated: true,
		}, true
	case strings.HasPrefix(line[position:], csharpRawStringDelimiter):
		return csharpStringStart{
			token:        csharpRawStringDelimiter,
			rawDelimiter: csharpRawStringDelimiter,
		}, true
	case line[position] == '"':
		return csharpStringStart{token: `"`}, true
	default:
		return csharpStringStart{}, false
	}
}

// startsEscapedInterpolationBrace reports whether the position starts {{ or }} in string text.
func (f *csharpLineFilter) startsEscapedInterpolationBrace() bool {
	if !f.state.stringInterpolated || f.position+1 >= len(f.line) {
		return false
	}

	return f.hasPrefix(csharpEscapedOpenBrace) || f.hasPrefix(csharpEscapedCloseBrace)
}

// closeSingleLineStringAtLineEnd clears invalid single-line strings at end of line.
func (f *csharpLineFilter) closeSingleLineStringAtLineEnd() {
	if f.state.stringActive &&
		!f.state.stringVerbatim &&
		f.state.stringRawDelimiter == "" &&
		f.state.interpolationDepth == 0 {
		f.closeString()
	}
}

// closeString clears the active C# string state.
func (f *csharpLineFilter) closeString() {
	f.state.stringActive = false
	f.state.stringVerbatim = false
	f.state.stringInterpolated = false
	f.state.stringRawDelimiter = ""
	f.state.interpolationDepth = 0
	f.state.interpolationFormat = false
}

// consumeComment consumes a C# comment when one starts at the scanner position.
func (f *csharpLineFilter) consumeComment() commentConsumeResult {
	return f.consumeCStyleComment(csharpDocLineComment, f.startBlockComment)
}

// startBlockComment records the active block comment and whether to keep it.
func (f *csharpLineFilter) startBlockComment(keep bool) {
	f.lineFilter.startBlockComment(&f.state.block, keep)
}
