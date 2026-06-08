// Copyright 2024, TeamDev. All rights reserved.
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

package parsing

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/gobwas/glob"
)

// Pattern represents a glob-like pattern to match consecutive source lines.
//
// Contains the original glob string and compiled matchers for each source-line pattern.
//
// sourceGlob — a glob-like string, e.g. "*main*" or "^main".
type Pattern struct {
	sourceGlob string
	matchers   []lineMatcher
}

const (
	anyCharacterSequence = "*"
	lineStart            = "^"
	lineEnd              = "$"
	lineSeparator        = `\n`
	escapedLineSeparator = `\\n`
)

// lineMatcher matches a single source line using the compiled glob pattern.
type lineMatcher struct {
	compiled glob.Glob
}

// NewPattern creates a new Pattern based on provided glob string.
//
// The modified pattern is the original one, but enclosed with the "*" wildcards,
// unless start of the line or end of the line wildcards were specified.
//
// A multi-line pattern uses "\n" as a separator between consecutive source-line
// patterns. For example, "Test \n adds two values" matches a line matching "Test"
// followed by a line matching "adds two values". Each part separated by "\n" is
// compiled separately and follows the same wildcard rules.
// Use "\\n" to match literal "\n" text instead of starting the next pattern line.
//
// glob — a string that represents a pattern that can include such wildcards:
//   - "*" — matches any sequence of characters;
//   - "^" — matches the start of the line;
//   - "$" — matches the end of the line.
//
// Returns an error if any modified glob pattern cannot be compiled.
func NewPattern(globString string) (Pattern, error) {
	patternLines := splitPatternLines(globString)
	matchers := make([]lineMatcher, 0, len(patternLines))
	for _, patternLine := range patternLines {
		matcher, err := compileLineMatcher(patternLine)
		if err != nil {
			return Pattern{}, err
		}
		matchers = append(matchers, matcher)
	}

	return Pattern{
		sourceGlob: globString,
		matchers:   matchers,
	}, nil
}

// compileLineMatcher compiles one source-line pattern into a glob matcher.
func compileLineMatcher(patternLine string) (lineMatcher, error) {
	pattern := patternLine

	startOfLine := strings.HasPrefix(patternLine, lineStart)
	if !startOfLine && !strings.HasPrefix(patternLine, anyCharacterSequence) {
		pattern = anyCharacterSequence + pattern
	}
	if startOfLine {
		pattern = pattern[1:]
	}

	endOfLine := strings.HasSuffix(patternLine, lineEnd)
	if !endOfLine && !strings.HasSuffix(patternLine, anyCharacterSequence) {
		pattern += anyCharacterSequence
	}
	if endOfLine {
		lastIndex := len(pattern) - 1
		pattern = pattern[:lastIndex]
	}

	compiledGlob, err := glob.Compile(pattern)
	if err != nil {
		return lineMatcher{}, err
	}

	return lineMatcher{compiled: compiledGlob}, nil
}

// matches reports whether the source line matches the compiled pattern.
func (m lineMatcher) matches(line string) bool {
	return m.compiled != nil && m.compiled.Match(line)
}

// FindIn returns the first source-line range matching the pattern.
func (p Pattern) FindIn(lines []string, startFrom int) (int, int, bool) {
	if len(p.matchers) == 0 || startFrom < 0 {
		return 0, 0, false
	}
	lastStart := len(lines) - len(p.matchers)
	for start := startFrom; start <= lastStart; start++ {
		if p.matchesAt(lines, start) {
			return start, start + len(p.matchers) - 1, true
		}
	}

	return 0, 0, false
}

// matchesAt reports whether the compiled matchers match source lines at start.
func (p Pattern) matchesAt(lines []string, start int) bool {
	if len(p.matchers) == 0 || start < 0 || start+len(p.matchers) > len(lines) {
		return false
	}
	for i, matcher := range p.matchers {
		if !matcher.matches(lines[start+i]) {
			return false
		}
	}

	return true
}

// splitPatternLines returns trimmed pattern lines separated by an escaped newline.
func splitPatternLines(sourceGlob string) []string {
	var patternLines []string
	var line strings.Builder
	trimLeft := false
	for cursor := 0; cursor < len(sourceGlob); {
		remaining := sourceGlob[cursor:]
		switch {
		case strings.HasPrefix(remaining, escapedLineSeparator):
			line.WriteString(escapedLineSeparator)
			cursor += len(escapedLineSeparator)
		case strings.HasPrefix(remaining, lineSeparator):
			patternLines = append(patternLines, strings.TrimRightFunc(line.String(), unicode.IsSpace))
			line.Reset()
			trimLeft = true
			cursor += len(lineSeparator)
		case trimLeft:
			r, size := utf8.DecodeRuneInString(remaining)
			if !unicode.IsSpace(r) {
				trimLeft = false
				line.WriteByte(sourceGlob[cursor])
				cursor++

				continue
			}
			cursor += size
		default:
			trimLeft = false
			line.WriteByte(sourceGlob[cursor])
			cursor++
		}
	}
	patternLines = append(patternLines, line.String())

	return patternLines
}

// String returns a string representation of Pattern.
func (p Pattern) String() string {
	return fmt.Sprintf("Pattern %s", p.sourceGlob)
}
