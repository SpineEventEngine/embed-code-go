package embedding_instruction

import (
	"fmt"
	"strings"

	"github.com/gobwas/glob"
)

// Represents a glob-like pattern to match a line of a source file.
//
// Contains both the original glob string and a modified pattern suitable for matching.
//
// sourceGlob — a glob-like string, e.g. "*main*" or "^main".
//
// pattern — a pattern to search for.
type Pattern struct {
	sourceGlob string
	pattern    string
}

//
// Initializers
//

// Creates a new Pattern based on the provided glob string.
//
// The resulting Pattern struct contains both the original glob string and a
// modified pattern suitable for matching.
//
// The modified pattern is the original one, but enclosed with the "*" wildcards,
// unless start of the line or end of the line wildcards weren't specified.
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

//
// Public methods
//

// Reports whether given line matches the pattern.
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
