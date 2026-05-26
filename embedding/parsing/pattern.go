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

	"github.com/gobwas/glob"
)

// Pattern represents a glob-like pattern to match a line of a source file.
//
// Contains both original glob string and modified pattern suitable for matching.
//
// sourceGlob — a glob-like string, e.g. "*main*" or "^main".
//
// pattern — a pattern to search for.
type Pattern struct {
	sourceGlob string
	pattern    string
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
//	p := NewPattern("*.txt")
//	fmt.Println("Original glob:", p.sourceGlob) // "*.txt"
//	fmt.Println("Modified pattern:", p.pattern) // "*.txt*"
//
//	p := NewPattern("^.txt")
//	fmt.Println("Original glob:", p.sourceGlob) // "*.txt"
//	fmt.Println("Modified pattern:", p.pattern) // ".txt*"
func NewPattern(glob string) Pattern {
	pattern := glob

	startOfLine := strings.HasPrefix(glob, lineStart)
	if !startOfLine && !strings.HasPrefix(glob, anyCharacterSequence) {
		pattern = anyCharacterSequence + pattern
	}
	if startOfLine {
		pattern = pattern[1:]
	}

	endOfLine := strings.HasSuffix(glob, lineEnd)
	if !endOfLine && !strings.HasSuffix(glob, anyCharacterSequence) {
		pattern += anyCharacterSequence
	}
	if endOfLine {
		lastIndex := len(pattern) - 1
		pattern = pattern[:lastIndex]
	}

	return Pattern{
		sourceGlob: glob,
		pattern:    pattern,
	}
}

// Match reports whether given line matches the pattern.
//
// line — a line to check the match for.
func (p Pattern) Match(line string) bool {
	g := glob.MustCompile(p.pattern)

	return g.Match(line)
}

// HasLineSeparator reports whether the pattern contains an escaped line separator.
func (p Pattern) HasLineSeparator() bool {
	_, hasSeparator := p.linePatterns()

	return hasSeparator
}

// MatchLineSequence reports whether source lines match the escaped-line-separated pattern.
func (p Pattern) MatchLineSequence(lines []string) bool {
	patternLines, _ := p.linePatterns()
	if len(patternLines) != len(lines) {
		return false
	}
	for i, patternLine := range patternLines {
		pattern := NewPattern(patternLine)
		if !pattern.Match(lines[i]) {
			return false
		}
	}

	return true
}

// linePatterns returns trimmed pattern lines separated by an escaped newline.
func (p Pattern) linePatterns() ([]string, bool) {
	var patternLines []string
	var line strings.Builder
	hasSeparator := false
	for i := 0; i < len(p.sourceGlob); {
		remaining := p.sourceGlob[i:]
		switch {
		case strings.HasPrefix(remaining, escapedLineSeparator):
			line.WriteString(escapedLineSeparator)
			i += len(escapedLineSeparator)
		case strings.HasPrefix(remaining, lineSeparator):
			patternLines = append(patternLines, strings.TrimSpace(line.String()))
			line.Reset()
			hasSeparator = true
			i += len(lineSeparator)
		default:
			line.WriteByte(p.sourceGlob[i])
			i++
		}
	}
	patternLines = append(patternLines, strings.TrimSpace(line.String()))

	return patternLines, hasSeparator
}

// Returns string representation of Pattern.
func (p Pattern) String() string {
	return fmt.Sprintf("Pattern %s", p.sourceGlob)
}
