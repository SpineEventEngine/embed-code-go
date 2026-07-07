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
	var wrongTagErr WrongInstructionTagError
	if errors.As(parseErr, &wrongTagErr) {
		return wrongTagErr.Error()
	}
	if !openingTagClosed(instructionBody) {
		return fmt.Sprintf("the opening `<%s>` tag is not closed; add `>` or `/>` before the code fence",
			EmbeddingTag,
		)
	}
	if !instructionClosed(instructionBody) {
		return fmt.Sprintf("the `<%s>` instruction is not closed; add `</%s>` or use `/>`",
			EmbeddingTag,
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

// openingTagClosed reports whether the first embed-code opening tag reaches `>`.
func openingTagClosed(instructionBody []string) bool {
	return scanInstructionBody(instructionBody, func(remaining string) (bool, bool) {
		switch {
		case strings.HasPrefix(remaining, ">"):
			return true, true
		case strings.HasPrefix(remaining, "</"+EmbeddingTag):
			return false, true
		case strings.HasPrefix(remaining, "<"+EmbeddingTag):
			return false, true
		default:
			return false, false
		}
	})
}

// instructionClosed reports whether the instruction has a self-closing or closing tag.
func instructionClosed(instructionBody []string) bool {
	return scanInstructionBody(instructionBody, func(remaining string) (bool, bool) {
		switch {
		case strings.HasPrefix(remaining, "/>"):
			return true, true
		case strings.HasPrefix(remaining, "</"+EmbeddingTag+">"):
			return true, true
		default:
			return false, false
		}
	})
}

// scanInstructionBody scans instruction text outside quoted attribute values.
func scanInstructionBody(
	instructionBody []string,
	check func(remaining string) (bool, bool),
) bool {
	inTag := false
	inQuote := false
	var quote byte
	for _, line := range instructionBody {
		if inTag && !inQuote && markdownCodeFenceStart(line) {
			return false
		}
		matched, done := scanInstructionLine(line, check, &inTag, &inQuote, &quote)
		if done {
			return matched
		}
	}

	return false
}

// scanInstructionLine scans one instruction line outside quoted attribute values.
func scanInstructionLine(
	line string,
	check func(remaining string) (bool, bool),
	inTag *bool,
	inQuote *bool,
	quote *byte,
) (bool, bool) {
	for lineCursor := 0; lineCursor < len(line); lineCursor++ {
		if !*inTag {
			if !startInstructionScan(line, &lineCursor, inTag) {
				break
			}

			continue
		}
		if *inQuote {
			skipQuotedCharacter(line, &lineCursor, inQuote, quote)

			continue
		}
		if line[lineCursor] == '"' || line[lineCursor] == '\'' {
			*inQuote = true
			*quote = line[lineCursor]

			continue
		}
		matched, done := check(line[lineCursor:])
		if done {
			return matched, true
		}
	}

	return false, false
}

// startInstructionScan moves lineCursor to the first embed-code tag.
func startInstructionScan(line string, lineCursor *int, inTag *bool) bool {
	tagStart := strings.Index(line[*lineCursor:], "<"+EmbeddingTag)
	if tagStart < 0 {
		return false
	}
	*lineCursor += tagStart + len("<"+EmbeddingTag) - 1
	*inTag = true

	return true
}

// skipQuotedCharacter updates quote state while scanning an attribute value.
func skipQuotedCharacter(line string, lineCursor *int, inQuote *bool, quote *byte) {
	if line[*lineCursor] == '\\' {
		*lineCursor++

		return
	}
	if line[*lineCursor] == *quote {
		*inQuote = false
	}
}

// markdownCodeFenceStart reports whether a line starts a Markdown code fence.
func markdownCodeFenceStart(line string) bool {
	trimmedLine := strings.TrimSpace(line)
	if !strings.HasPrefix(trimmedLine, "```") {
		return false
	}

	return len(codeFenceMarker(trimmedLine)) >= minCodeFenceMarkerLength
}
