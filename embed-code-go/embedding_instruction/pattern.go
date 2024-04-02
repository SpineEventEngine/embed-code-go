package embedding_instruction

import (
	"fmt"
	"strings"

	"github.com/gobwas/glob"
)

// Pattern represents a glob-like pattern to match a line of a source file.
type Pattern struct {
	sourceGlob string
	pattern    string
}

//
// Initializers
//

// NewPattern creates a new Pattern.
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

// Match checks if the given line matches the pattern.
func (p Pattern) Match(line string) bool {
	g := glob.MustCompile(p.pattern)
	return g.Match(line)
}

func (p Pattern) String() string {
	return fmt.Sprintf("Pattern %s", p.sourceGlob)
}
