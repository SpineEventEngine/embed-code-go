/*
 * Copyright 2026, TeamDev. All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 * Redistribution and use in source and/or binary forms, with or without
 * modification, must retain the above copyright notice and the following
 * disclaimer.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
 * "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
 * LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
 * A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
 * OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
 * SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
 * LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
 * DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
 * THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 * OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

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
	// Line is the line containing the start of the malformed instruction.
	Line int

	// Reason describes why parsing failed.
	Reason string
}

// MissingCodeFenceError reports that an embedding instruction is not followed by a code fence.
type MissingCodeFenceError struct {
	// Line is the line containing the embedding instruction.
	Line int
}

// UnclosedCodeFenceError reports that an embedding code fence is not closed.
type UnclosedCodeFenceError struct {
	// Line is the line containing the embedding instruction.
	Line int
}

// Error returns a user-facing description of an embedding instruction parse failure.
//
// Returns formatted parse failure text.
func (e InstructionParseError) Error() string {
	return fmt.Sprintf(
		"failed to parse an embedding instruction: %s",
		e.Reason,
	)
}

// Error returns a user-facing description of a missing code fence after an instruction.
//
// Returns formatted missing-fence text.
func (e MissingCodeFenceError) Error() string {
	return "expected a markdown code fence after the embedding instruction"
}

// Error returns a user-facing description of an unclosed embedding code fence.
//
// Returns formatted unclosed-fence text.
func (e UnclosedCodeFenceError) Error() string {
	return "the markdown code fence after the embedding instruction is not closed"
}

// Recognize reports whether the current line starts an embedding instruction.
//
// Parameters:
// context - provides current parser state.
//
// Returns true when the current line starts an embed-code instruction.
func (e EmbedInstructionTokenState) Recognize(context Context) bool {
	line := context.CurrentLine()
	isStatement := strings.HasPrefix(strings.TrimSpace(line), "<"+EmbeddingTag)
	if context.EmbeddingInstruction == nil &&
		!context.ReachedEOF() &&
		!context.MarkdownFenceStarted &&
		isStatement {
		return true
	}

	return false
}

// Accept parses the embedding instruction and records it in the parsing context.
//
// Parameters:
// context - provides mutable parser state.
// config - provides embedding configuration.
//
// Returns an error when the instruction cannot be parsed.
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
			context.StartEmbedding(instruction)
		} else {
			parseErr = err
		}

		context.Result = append(context.Result, line)
		context.ToNextLine()

		// Once the tag is syntactically closed, the following lines are not
		// part of the instruction. Stop here instead of consuming the rest of
		// the document trying to parse a complete but invalid instruction.
		if context.EmbeddingInstruction == nil && instructionClosed(instructionBody) {
			break
		}
	}
	if context.EmbeddingInstruction == nil {
		return InstructionParseError{
			Line:   startLine,
			Reason: parseFailureReason(instructionBody, parseErr),
		}
	}

	return nil
}

// instructionClosed reports whether the accumulated instruction body contains a
// closing tag, meaning any following lines are not part of the instruction.
//
// The terminator is only recognized outside quoted attribute values, so a
// value such as `line="<br/>"` does not end the instruction early.
func instructionClosed(instructionBody []string) bool {
	instruction := strings.Join(instructionBody, " ")
	closingTag := "</" + EmbeddingTag + ">"
	insideValue := false
	for i := 0; i < len(instruction); i++ {
		char := instruction[i]
		if char == '\\' && i+1 < len(instruction) && instruction[i+1] == '"' {
			i++

			continue
		}
		if char == '"' {
			insideValue = !insideValue

			continue
		}
		if insideValue {
			continue
		}
		if strings.HasPrefix(instruction[i:], "/>") ||
			strings.HasPrefix(instruction[i:], closingTag) {
			return true
		}
	}

	return false
}

// parseFailureReason explains why an embedding instruction could not be parsed.
func parseFailureReason(instructionBody []string, parseErr error) string {
	if !instructionClosed(instructionBody) {
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
