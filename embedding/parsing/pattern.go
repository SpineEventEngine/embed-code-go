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
// Pattern contains the original glob string and compiled matchers for each source-line pattern.
type Pattern struct {
	// sourceGlob is the original glob-like pattern.
	sourceGlob string

	// matchers contains one compiled matcher per source-line pattern.
	matchers []lineMatcher
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
	// compiled is the compiled glob matcher.
	compiled glob.Glob
}

// NewPattern creates a new Pattern from an embed-code source-line pattern.
//
// NewPattern encloses the original pattern with "*" wildcards,
// unless start of the line or end of the line wildcards were specified.
//
// A multi-line pattern uses "\n" as a separator between consecutive source-line
// patterns. For example, "Test \n adds two values" matches a line matching "Test"
// followed by a line matching "adds two values". Each part separated by "\n" is
// compiled separately and follows the same wildcard rules.
// Use "\\n" to match literal "\n" text instead of starting the next pattern line.
//
// Supported wildcards:
//   - "*" - matches any sequence of characters;
//   - "^" - matches the start of the line;
//   - "$" - matches the end of the line.
//
// Parameters:
// globString - provides source-line pattern text.
//
// Returns:
// Pattern - compiled source-line pattern.
// error - when any pattern line cannot be compiled.
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

	var matcher lineMatcher
	var compileErr error
	func() {
		defer func() {
			recovered := recover()
			if recovered != nil {
				compileErr = fmt.Errorf("glob pattern compiler panicked: %v", recovered)
			}
		}()

		compiledGlob, err := glob.Compile(pattern)
		if err != nil {
			compileErr = err

			return
		}
		matcher = lineMatcher{compiled: compiledGlob}
	}()
	if compileErr != nil {
		return lineMatcher{}, compileErr
	}

	return matcher, nil
}

// matches reports whether the source line matches the compiled pattern.
func (m lineMatcher) matches(line string) bool {
	if m.compiled == nil {
		return false
	}

	return matchGlob(m.compiled, line)
}

// matchGlob reports whether a glob matches and treats dependency panics as misses.
//
// A miss lets the instruction layer report PatternNotFoundError with source
// context instead of exposing a third-party matcher panic to users.
// Recovery stays at single-line matcher granularity so one panicking candidate
// does not abort the rest of the source scan.
func matchGlob(compiledGlob glob.Glob, line string) bool {
	var matched bool
	func() {
		defer func() {
			recovered := recover()
			if recovered != nil {
				matched = false
			}
		}()

		matched = compiledGlob.Match(line)
	}()

	return matched
}

// FindIn returns the first source-line range matching the pattern.
//
// Parameters:
// lines - provides source lines to scan.
// startFrom - provides the first index to scan.
//
// Returns:
// int - inclusive start index.
// int - inclusive end index.
// bool - whether a match was found.
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
				line.WriteString(sourceGlob[cursor : cursor+size])
				cursor += size

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
//
// Returns diagnostic pattern text.
func (p Pattern) String() string {
	return fmt.Sprintf("Pattern %s", p.sourceGlob)
}
