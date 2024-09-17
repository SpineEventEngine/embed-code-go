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

// BlankLine represents a blank line of a markdown.
type BlankLine struct{}

// Recognize reports whether the current line is a blank line.
//
// Checks if the current line is empty and not part of a code fence,
// and if there is an embedding. If these conditions are met, it returns true.
// Otherwise, it returns false.
func (b BlankLine) Recognize(context Context) bool {
	isEmptyString := strings.TrimSpace(context.CurrentLine()) == ""
	if !context.ReachedEOF() && isEmptyString {
		return !context.CodeFenceStarted && context.Embedding != nil
	}

	return false
}

// Accept processes a blank line of a markdown.
//
// Appends the current line of the context to the result, and moves to the next line.
func (b BlankLine) Accept(context *Context, _ configuration.Configuration) error {
	line := context.CurrentLine()
	context.Result = append(context.Result, line)
	context.ToNextLine()

	return nil
}
