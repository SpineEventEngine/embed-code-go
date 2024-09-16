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
	"embed-code/embed-code-go/embedding"
	"fmt"
	"strings"

	"embed-code/embed-code-go/configuration"
)

// EmbedInstructionToken represents an embedding instruction token of a markdown.
type EmbedInstructionToken struct{}

// Recognize reports whether the current line in the parsing context starts with "<embed-code",
// and if there is no ongoing embedding and the end of the file is not reached, it returns true.
// Otherwise, it returns false.
//
// context — a context of the parsing process.
func (e EmbedInstructionToken) Recognize(context Context) bool {
	line := context.CurrentLine()
	isStatement := strings.HasPrefix(strings.TrimSpace(line), Statement)
	if context.Embedding == nil && !context.ReachedEOF() && isStatement {
		return true
	}

	return false
}

// Accept parses the embedding instruction and extracts relevant information to update
// the parsing context. Switches the context to the next line.
//
// context — a context of the parsing process.
//
// config — a configuration of the embedding.
//
// An error is returned if the building of the embedding instruction fails.
func (e EmbedInstructionToken) Accept(context *Context,
	config configuration.Configuration) error {
	var instructionBody []string
	for !context.ReachedEOF() {
		instructionBody = append(instructionBody, context.CurrentLine())

		instruction, err := embedding.FromXML(strings.Join(instructionBody, ""), config)
		if err == nil {
			context.SetEmbedding(&instruction)
		}

		context.Result = append(context.Result, context.CurrentLine())
		context.ToNextLine()
		if context.Embedding != nil {
			break
		}
	}
	if context.Embedding == nil {
		return fmt.Errorf("failed to parse an embedding instruction. Context: %v", context)
	}

	return nil
}
