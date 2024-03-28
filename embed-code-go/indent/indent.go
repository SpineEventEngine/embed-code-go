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

// Finds the maximal common indentation of given lines.
//
// If all given lines are empty, contain only whitespace, or there are no lines at all,
// returns zero.
//
// lines — a list of lines which may or may not have leading whitespaces.
//
// Returns the number of leading whitespaces in all lines except for the empty ones.
func MaxCommonIndentation(lines []string) int {
	indent := math.MaxInt32
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			trimmedLine := strings.TrimLeft(line, "\n\t ") // Check if it changes a line in-place.
			lineIndent := len(line) - len(trimmedLine)
			if lineIndent < indent {
				indent = lineIndent
			}
		}
	}

	if indent == math.MaxInt32 {
		return 0
	} else {
		return indent
	}
}

// Reduces indentation to given redundantSpaces amount.
//
// lines — a list of strings representing the lines to process.
//
// redundantSpaces — the number of leading spaces to remove from each line.
//
// Returns processed lines.
func CutIndent(lines []string, redundantSpaces int) []string {
	linesChanged := make([]string, len(lines))
	copy(linesChanged, lines)
	for i, line := range linesChanged {
		if len(line) > 0 {
			linesChanged[i] = line[redundantSpaces:]
		}
	}
	return linesChanged
}
