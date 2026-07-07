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

package commentfilter

// BlockMarker describes a block comment marker pair.
type BlockMarker struct {
	// Start is the block comment opening marker.
	Start string

	// End is the block comment closing marker.
	End string
}

// DocumentationMarker describes API documentation comment markers.
type DocumentationMarker struct {
	// Inline contains documentation line-comment markers.
	Inline []string

	// Block contains documentation block-comment marker pairs.
	Block []BlockMarker
}

// TextBlockMarker describes a multi-line text literal delimiter.
type TextBlockMarker struct {
	// Delimiter opens and closes the text literal.
	Delimiter string

	// Escapes reports whether backslashes escape delimiter bytes.
	Escapes bool
}

// CommentMarker describes lexical comment markers and string delimiters for a language family.
type CommentMarker struct {
	// Inline contains line-comment markers.
	Inline []string

	// Block contains block-comment marker pairs.
	Block []BlockMarker

	// Documentation contains API documentation comment markers.
	Documentation DocumentationMarker

	// TextBlocks contains markers that open and close multi-line text literals.
	TextBlocks []TextBlockMarker

	// QuoteChars contains characters that open and close quoted strings.
	QuoteChars string
}

// MarkerCommentFilter removes comments using lexical markers declared in CommentMarker.
type MarkerCommentFilter struct {
	// Syntax contains the comment markers and string delimiters to recognize.
	Syntax CommentMarker
}
