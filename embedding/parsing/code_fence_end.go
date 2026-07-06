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
	"strings"

	"embed-code/embed-code-go/configuration"
)

// CodeFenceEndState represents the end of a code fence.
type CodeFenceEndState struct{}

// Recognize reports whether the current line closes the active embedding fence.
//
// It requires EOF not reached, an active code fence, matching fence indentation,
// and a closing fence marker compatible with the opening marker.
//
// Parameters:
// context - provides current parser state.
//
// Returns true when the current line closes the active embedding code fence.
func (c CodeFenceEndState) Recognize(context Context) bool {
	if context.ReachedEOF() {
		return false
	}
	if !context.CodeFenceStarted {
		return false
	}
	indentation := strings.Repeat(" ", context.CodeFenceIndentation)
	line := strings.TrimPrefix(context.CurrentLine(), indentation)
	if line == context.CurrentLine() && context.CodeFenceIndentation > 0 {
		return false
	}
	if context.CodeFenceMarker == "" {
		return false
	}

	return isClosingCodeFence(line, context.CodeFenceMarker)
}

// Accept renders the embedding content and closes the active embedding fence.
//
// It appends the closing fence when rendering succeeds. When rendering fails,
// it restores the original Markdown for the embedding before advancing.
//
// Parameters:
// context - provides mutable parser state.
//
// Returns an error when embedded content cannot be produced.
func (c CodeFenceEndState) Accept(context *Context, _ configuration.Configuration) error {
	line := context.CurrentLine()
	err := renderSample(context)
	context.FinishEmbedding()
	if err == nil {
		context.Result = append(context.Result, line)
	} else {
		context.ResolveEmbeddingNotFound()
	}
	context.CodeFenceStarted = false
	context.CodeFenceMarker = ""
	context.CodeFenceIndentation = 0
	context.ToNextLine()

	return err
}

// renderSample appends rendered embedding source lines to the parse result.
//
// Parameters:
// context - provides mutable parser state and the current embedding instruction.
//
// Returns an error when reading the embedding content fails.
func renderSample(context *Context) error {
	content, err := context.EmbeddingInstruction.Content()
	if err != nil {
		return err
	}
	for _, line := range content {
		indentation := strings.Repeat(" ", context.CodeFenceIndentation)
		context.Result = append(context.Result, indentation+line)
	}

	return nil
}

// isClosingCodeFence reports whether line closes a fence opened with marker.
func isClosingCodeFence(line string, marker string) bool {
	if line == "" {
		return false
	}
	markerChar := marker[0]
	index := 0
	for index < len(line) && line[index] == markerChar {
		index++
	}
	if index < len(marker) {
		return false
	}

	return strings.TrimSpace(line[index:]) == ""
}
