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

package indent

import (
	"math"
	"strings"
)

// MaxCommonIndentation finds the maximal common indentation of given lines.
//
// Parameters:
// lines - provides source lines to inspect.
//
// Returns maximal common indentation, or zero when no non-blank lines exist.
func MaxCommonIndentation(lines []string) int {
	indent := math.MaxInt32
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			trimmedLine := strings.TrimLeft(line, "\n\t ")
			lineIndent := len(line) - len(trimmedLine)
			if lineIndent < indent {
				indent = lineIndent
			}
		}
	}

	if indent == math.MaxInt32 {
		return 0
	}

	return indent
}

// CutIndent reduces indentation to given redundantSpaces amount.
//
// It copies lines before trimming, so the input slice is not modified.
// If a line is shorter than redundantSpaces, the whole line is removed.
//
// Parameters:
// lines - provides source lines to trim.
// redundantSpaces - provides the maximum indentation to remove.
//
// Returns source lines with indentation removed.
func CutIndent(lines []string, redundantSpaces int) []string {
	linesChanged := make([]string, len(lines))
	copy(linesChanged, lines)
	for i, line := range linesChanged {
		if len(line) > 0 {
			cutLength := min(redundantSpaces, len(line))
			linesChanged[i] = line[cutLength:]
		}
	}

	return linesChanged
}
