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

// CodeFenceStartState represents the StartState of a code fence.
type CodeFenceStartState struct{}

// Recognize reports whether the current line is not reached the end and starts with "```".
//
// context — a context of the parsing process.
func (c CodeFenceStartState) Recognize(context Context) bool {
	if !context.ReachedEOF() {
		return strings.HasPrefix(strings.TrimSpace(context.CurrentLine()), "```")
	}

	return false
}

// Accept appends the current line from the parsing context to the result, sets a flag to indicate
// that a code fence has started, calculates the indentation level of the code fence, and moves
// to the next line in the context.
//
// context — a context of the parsing process.
func (c CodeFenceStartState) Accept(context *Context, _ configuration.Configuration) error {
	line := context.CurrentLine()
	context.Result = append(context.Result, line)
	context.CodeFenceStarted = true
	leadingSpaces := len(line) - len(strings.TrimLeft(line, " "))
	context.CodeFenceIndentation = leadingSpaces
	context.ToNextLine()
	// As we accepted this state and moved to the next line, we assume that the code lines
	// start here.
	context.SetCodeStart()

	return nil
}
