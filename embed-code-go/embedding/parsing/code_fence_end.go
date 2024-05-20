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
	"embed-code/embed-code-go/configuration"
	"strings"
)

//
// Public methods
//

// Represents the end of a code fence.
type CodeFenceEnd struct{}

// Reports whether the current line is the end of a code fence.
//
// The line is a code fence end if:
//   - the end is not reached;
//   - the code fence has started;
//   - the current line starts with the appropriate indentation and "```"
//
// context — a context of the parsing process.
func (c CodeFenceEnd) Recognize(context ParsingContext) bool {
	if !context.ReachedEOF() {
		indentation := strings.Repeat(" ", context.CodeFenceIndentation)
		return context.CodeFenceStarted && strings.HasPrefix(context.CurrentLine(), indentation+"```")
	}
	return false
}

// Processes the end of a code fence by adding the current line to the result,
// resetting certain context variables, and moving to the next line.
//
// context — a context of the parsing process.
//
// config — a configuration of the embedding.
func (c CodeFenceEnd) Accept(context *ParsingContext, config configuration.Configuration) {
	line := context.CurrentLine()
	renderSample(context)
	context.Result = append(context.Result, line)
	context.SetEmbedding(nil)
	context.CodeFenceStarted = false
	context.CodeFenceIndentation = 0
	context.ToNextLine()
}

//
// Private methods
//

func renderSample(context *ParsingContext) {
	for _, line := range context.Embedding.Content() {
		indentation := strings.Repeat(" ", context.CodeFenceIndentation)
		context.Result = append(context.Result, indentation+line)
	}
}
