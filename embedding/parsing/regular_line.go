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

const minCodeFenceMarkerLength = 3

// RegularLineState represents ordinary Markdown content.
type RegularLineState struct{}

// Recognize accepts any current line as ordinary Markdown content.
//
// Parameters:
// context - provides current parser state.
//
// Returns true for every parser context.
func (r RegularLineState) Recognize(_ Context) bool {
	return true
}

// Accept appends the current line and advances to the next documentation line.
//
// Parameters:
// context - provides mutable parser state.
//
// Returns nil.
func (r RegularLineState) Accept(context *Context, _ configuration.Configuration) error {
	line := context.CurrentLine()
	updateMarkdownFenceContext(context, line)
	context.Result = append(context.Result, line)
	context.ToNextLine()

	return nil
}

// updateMarkdownFenceContext tracks ordinary Markdown fences outside embeddings.
func updateMarkdownFenceContext(context *Context, line string) {
	if context.EmbeddingInstruction != nil {
		return
	}
	leadingSpaces := len(line) - len(strings.TrimLeft(line, " "))
	trimmedLine := strings.TrimSpace(line)
	if !strings.HasPrefix(trimmedLine, "```") {
		return
	}
	marker := codeFenceMarker(trimmedLine)
	if len(marker) < minCodeFenceMarkerLength {
		return
	}
	if !context.MarkdownFenceStarted {
		context.MarkdownFenceStarted = true
		context.MarkdownFenceMarker = marker
		context.MarkdownFenceIndentation = leadingSpaces

		return
	}
	if context.MarkdownFenceIndentation != leadingSpaces {
		return
	}
	if marker[0] != context.MarkdownFenceMarker[0] || len(marker) < len(context.MarkdownFenceMarker) {
		return
	}
	if strings.TrimSpace(trimmedLine[len(marker):]) != "" {
		return
	}
	context.MarkdownFenceStarted = false
	context.MarkdownFenceMarker = ""
	context.MarkdownFenceIndentation = 0
}
