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

package embedding

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

// NewPattern creates a new Pattern based on provided glob string.
//
// The resulting Pattern struct contains both original glob string and
// modified pattern suitable for matching.
//
// The modified pattern is the original one, but enclosed with the "*" wildcards,
// unless start of the line or end of the line wildcards were specified.
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
	startOfLine := strings.HasPrefix(glob, "^")
	if !startOfLine && !strings.HasPrefix(glob, "*") {
		pattern = "*" + pattern
	}
	if startOfLine {
		pattern = pattern[1:]
	}
	endOfLine := strings.HasSuffix(glob, "$")
	if !endOfLine && !strings.HasSuffix(glob, "*") {
		pattern += "*"
	}
	if endOfLine {
		pattern = pattern[:len(pattern)-1]
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

// Returns string representation of Pattern.
func (p Pattern) String() string {
	return fmt.Sprintf("Pattern %s", p.sourceGlob)
}
