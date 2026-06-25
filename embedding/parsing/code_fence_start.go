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
	"strings"

	"embed-code/embed-code-go/configuration"
)

// CodeFenceStartState represents the start state of an embedding code fence.
type CodeFenceStartState struct{}

// Recognize reports whether the current line starts a code fence.
//
// Parameters:
// context - provides current parser state.
//
// Returns true when EOF is not reached and the current line starts a Markdown code fence.
func (c CodeFenceStartState) Recognize(context Context) bool {
	if !context.ReachedEOF() {
		return strings.HasPrefix(strings.TrimSpace(context.CurrentLine()), "```")
	}

	return false
}

// Accept records code fence state and advances to the first embedded source line.
//
// It appends the current line to the result, records that the code fence has started,
// records its indentation, and advances to the next line.
//
// Parameters:
// context - provides mutable parser state.
//
// Returns nil.
func (c CodeFenceStartState) Accept(context *Context, _ configuration.Configuration) error {
	line := context.CurrentLine()
	trimmedLine := strings.TrimSpace(line)
	context.Result = append(context.Result, line)
	context.CodeFenceStarted = true
	context.CodeFenceMarker = codeFenceMarker(trimmedLine)
	leadingSpaces := len(line) - len(strings.TrimLeft(line, " "))
	context.CodeFenceIndentation = leadingSpaces
	context.ToNextLine()
	// After accepting the opening fence and moving to the next line,
	// embedded source lines start at the current context position.
	context.SetCodeStart()

	return nil
}

// codeFenceMarker returns the repeated fence marker characters at the start of line.
func codeFenceMarker(line string) string {
	if line == "" {
		return ""
	}
	markerChar := line[0]
	index := 0
	for index < len(line) && line[index] == markerChar {
		index++
	}

	return line[:index]
}
