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

// Pattern represents a glob-like pattern to match a line of a source file.
//
// Contains both original glob string, modified pattern suitable for matching,
// and a compiled matcher for the modified pattern.
//
// sourceGlob — a glob-like string, e.g. "*main*" or "^main".
//
// pattern — a pattern to search for.
type Pattern struct {
	sourceGlob string
	pattern    string
	matcher    glob.Glob
}

const (
	anyCharacterSequence = "*"
	lineStart            = "^"
	lineEnd              = "$"
	lineSeparator        = `\n`
	escapedLineSeparator = `\\n`
)

// NewPattern creates a new Pattern based on provided glob string.
//
// The resulting Pattern struct contains both original glob string and
// modified pattern suitable for matching.
//
// The modified pattern is the original one, but enclosed with the "*" wildcards,
// unless start of the line or end of the line wildcards were specified.
//
// A multi-line pattern uses "\n" as a separator between consecutive source-line
// patterns. For example, "Test \n adds two values" matches a line matching "Test"
// followed by a line matching "adds two values". Each part separated by "\n" is
// converted to Pattern separately and follows the same wildcard rules.
// Use "\\n" to match literal "\n" text instead of starting the next pattern line.
//
// glob — a string that represents a pattern that can include such wildcards:
//   - "*" — matches any sequence of characters;
//   - "^" — matches the start of the line;
//   - "$" — matches the end of the line.
//
// Example usage:
//
//	p, err := NewPattern("*.txt")
//	fmt.Println("Original glob:", p.sourceGlob) // "*.txt"
//	fmt.Println("Modified pattern:", p.pattern) // "*.txt*"
//
//	p, err = NewPattern("^.txt")
//	fmt.Println("Original glob:", p.sourceGlob) // "^.txt"
//	fmt.Println("Modified pattern:", p.pattern) // ".txt*"
//
// Returns an error if the modified glob pattern cannot be compiled.
func NewPattern(globString string) (Pattern, error) {
	pattern := globString

	startOfLine := strings.HasPrefix(globString, lineStart)
	if !startOfLine && !strings.HasPrefix(globString, anyCharacterSequence) {
		pattern = anyCharacterSequence + pattern
	}
	if startOfLine {
		pattern = pattern[1:]
	}

	endOfLine := strings.HasSuffix(globString, lineEnd)
	if !endOfLine && !strings.HasSuffix(globString, anyCharacterSequence) {
		pattern += anyCharacterSequence
	}
	if endOfLine {
		lastIndex := len(pattern) - 1
		pattern = pattern[:lastIndex]
	}

	matcher, err := glob.Compile(pattern)
	if err != nil {
		return Pattern{}, err
	}

	return Pattern{
		sourceGlob: globString,
		pattern:    pattern,
		matcher:    matcher,
	}, nil
}

// Match reports whether given line matches the pattern.
//
// line — a line to check the match for.
func (p Pattern) Match(line string) bool {
	if p.matcher == nil {
		return false
	}

	return p.matcher.Match(line)
}

// HasLineSeparator reports whether the pattern contains an escaped line separator.
func (p Pattern) HasLineSeparator() bool {
	_, hasSeparator := p.linePatterns()

	return hasSeparator
}

// MatchLineSequence reports whether source lines match the escaped-line-separated pattern.
func (p Pattern) MatchLineSequence(lines []string) bool {
	patterns, err := p.lineSequencePatterns()
	if err != nil {
		return false
	}

	return matchLineSequencePatterns(patterns, lines)
}

// matchLineSequencePatterns reports whether compiled Patterns match source lines in order.
func matchLineSequencePatterns(patterns []Pattern, lines []string) bool {
	if len(patterns) != len(lines) {
		return false
	}
	for i, pattern := range patterns {
		if !pattern.Match(lines[i]) {
			return false
		}
	}

	return true
}

// lineSequencePatterns returns the Patterns for each part of a multi-line pattern.
func (p Pattern) lineSequencePatterns() ([]Pattern, error) {
	patternLines, _ := p.linePatterns()
	patterns := make([]Pattern, 0, len(patternLines))
	for _, patternLine := range patternLines {
		pattern, err := NewPattern(patternLine)
		if err != nil {
			return nil, err
		}
		patterns = append(patterns, pattern)
	}

	return patterns, nil
}

// linePatterns returns trimmed pattern lines separated by an escaped newline.
func (p Pattern) linePatterns() ([]string, bool) {
	var patternLines []string
	var line strings.Builder
	hasSeparator := false
	trimLeft := false
	for cursor := 0; cursor < len(p.sourceGlob); {
		remaining := p.sourceGlob[cursor:]
		switch {
		case strings.HasPrefix(remaining, escapedLineSeparator):
			line.WriteString(escapedLineSeparator)
			cursor += len(escapedLineSeparator)
		case strings.HasPrefix(remaining, lineSeparator):
			patternLines = append(patternLines, strings.TrimRightFunc(line.String(), unicode.IsSpace))
			line.Reset()
			hasSeparator = true
			trimLeft = true
			cursor += len(lineSeparator)
		case trimLeft:
			r, size := utf8.DecodeRuneInString(remaining)
			if !unicode.IsSpace(r) {
				trimLeft = false
				line.WriteByte(p.sourceGlob[cursor])
				cursor++

				continue
			}
			cursor += size
		default:
			trimLeft = false
			line.WriteByte(p.sourceGlob[cursor])
			cursor++
		}
	}
	patternLines = append(patternLines, line.String())

	return patternLines, hasSeparator
}

// String returns a string representation of Pattern.
func (p Pattern) String() string {
	return fmt.Sprintf("Pattern %s", p.sourceGlob)
}
