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
	"encoding/xml"
	"errors"
	"fmt"
	"strings"

	"embed-code/embed-code-go/configuration"
)

// EmbedInstructionTokenState represents an embedding instruction token of a markdown.
type EmbedInstructionTokenState struct{}

// InstructionParseError reports a failed embedding instruction parse and its source line.
type InstructionParseError struct {
	Line   int
	Reason string
}

// Error returns a user-facing description of an embedding instruction parse failure.
func (e InstructionParseError) Error() string {
	return fmt.Sprintf(
		"failed to parse an embedding instruction: %s",
		e.Reason,
	)
}

// Recognize reports whether the current line in the parsing context starts with "<embed-code",
// and if there is no ongoing embedding and the end of the file is not reached, it returns true.
// Otherwise, it returns false.
//
// context — a context of the parsing process.
func (e EmbedInstructionTokenState) Recognize(context Context) bool {
	line := context.CurrentLine()
	isStatement := strings.HasPrefix(strings.TrimSpace(line), "<"+EmbeddingTag)
	if context.EmbeddingInstruction == nil && !context.ReachedEOF() && isStatement {
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
// Returns an error if the building of the embedding instruction fails.
func (e EmbedInstructionTokenState) Accept(context *Context,
	config configuration.Configuration) error {
	var instructionBody []string
	startLine := context.CurrentIndex()
	var parseErr error
	for !context.ReachedEOF() && context.EmbeddingInstruction == nil {
		line := context.CurrentLine()
		instructionBody = append(instructionBody, line)

		instruction, err := FromXML(strings.Join(instructionBody, " "), config)
		if err == nil {
			instruction.DocumentationFile = context.MarkdownFilePath
			instruction.DocumentationLine = startLine
			context.SetEmbedding(&instruction)
		} else {
			parseErr = err
		}

		context.Result = append(context.Result, line)
		context.ToNextLine()
	}
	if context.EmbeddingInstruction == nil {
		return InstructionParseError{
			Line:   startLine,
			Reason: parseFailureReason(instructionBody, parseErr),
		}
	}

	return nil
}

// parseFailureReason explains why an embedding instruction could not be parsed.
func parseFailureReason(instructionBody []string, parseErr error) string {
	instruction := strings.TrimSpace(strings.Join(instructionBody, " "))
	if !strings.Contains(instruction, "/>") &&
		!strings.Contains(instruction, "</"+EmbeddingTag+">") {
		return fmt.Sprintf("the `<%s>` tag is not closed",
			EmbeddingTag,
		)
	}
	if parseErr != nil {
		var syntaxErr *xml.SyntaxError
		if errors.As(parseErr, &syntaxErr) {
			return syntaxErr.Msg
		}
		return parseErr.Error()
	}

	return "invalid embedding instruction"
}
