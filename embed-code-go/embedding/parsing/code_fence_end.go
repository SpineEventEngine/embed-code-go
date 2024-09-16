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

// CodeFenceEnd represents the end of a code fence.
type CodeFenceEnd struct{}

// Recognize reports whether the current line is the end of a code fence.
//
// The line is a code fence end if:
//   - the end is not reached;
//   - the code fence has started;
//   - the current line starts with the appropriate indentation and "```"
//
// context — a context of the parsing process.
func (c CodeFenceEnd) Recognize(context Context) bool {
	if !context.ReachedEOF() {
		indentation := strings.Repeat(" ", context.CodeFenceIndentation)

		return context.CodeFenceStarted && strings.HasPrefix(context.CurrentLine(), indentation+"```")
	}

	return false
}

// Accept processes the end of a code fence by adding the current line to the result,
// resetting certain context variables, and moving to the next line.
//
// context — a context of the parsing process.
//
// config — a configuration of the embedding.
//
// Returns an error if the rendering was not successful.
func (c CodeFenceEnd) Accept(context *Context, _ configuration.Configuration) error {
	line := context.CurrentLine()
	err := renderSample(context)
	context.SetEmbedding(nil)
	if err == nil {
		context.Result = append(context.Result, line)
	} else {
		context.ResolveEmbeddingNotFound()
	}
	context.CodeFenceStarted = false
	context.CodeFenceIndentation = 0
	context.ToNextLine()

	return err
}

// Renders the sample content of the embedding.
//
// context — a context of the parsing process.
//
// Returns an error if the reading of the embedding's content was not successful.
func renderSample(context *Context) error {
	content, err := context.Embedding.Content()
	if err != nil {
		return err
	}
	for _, line := range content {
		indentation := strings.Repeat(" ", context.CodeFenceIndentation)
		context.Result = append(context.Result, indentation+line)
	}

	return nil
}
