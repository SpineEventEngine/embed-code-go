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
	"fmt"
	"os"
	"regexp"

	"embed-code/embed-code-go/fragmentation"
)

// Context represents the state of parsing a documentation file containing code embeddings.
type Context struct {
	// EmbeddingInstruction is the instruction currently being parsed.
	EmbeddingInstruction *Instruction

	// MarkdownFilePath is the path to the documentation file.
	MarkdownFilePath string

	// Result contains the documentation lines produced by parsing.
	Result []string

	// CodeFenceStarted reports whether an embedding code fence is open.
	CodeFenceStarted bool

	// CodeFenceMarker is the marker used by the open embedding code fence.
	CodeFenceMarker string

	// CodeFenceIndentation is the indentation of the embedding code fence.
	CodeFenceIndentation int

	// MarkdownFenceStarted reports whether an ordinary Markdown code fence is open.
	MarkdownFenceStarted bool

	// MarkdownFenceMarker is the marker used by the open ordinary Markdown code fence.
	MarkdownFenceMarker string

	// MarkdownFenceIndentation is the indentation of the ordinary Markdown code fence.
	MarkdownFenceIndentation int

	// EmbeddingsNotFound contains instructions whose source fragments were not found.
	EmbeddingsNotFound []Instruction

	// UnacceptedEmbeddings contains instructions rejected by the parser.
	UnacceptedEmbeddings []Instruction

	// source contains the original documentation lines.
	source []string

	// lineIndex is the zero-based index of the current documentation line.
	lineIndex int

	// fileContainsEmbedding reports whether the file contains an embedding instruction.
	fileContainsEmbedding bool

	// embeddings contains accepted embedding instructions and their source positions.
	embeddings []EmbeddingContext

	// resolver owns source fragmentation cache state for this processing operation.
	resolver *fragmentation.Resolver
}

// EmbeddingsCount returns the number of found embeddings.
//
// Returns accepted embedding count.
func (c *Context) EmbeddingsCount() int {
	return len(c.embeddings)
}

// EmbeddingContext contains an instruction and its position in the source Markdown file.
type EmbeddingContext struct {
	// embeddingInstruction contains the embedding parameters.
	embeddingInstruction Instruction

	// SourceStartIndex is the first source-line index belonging to the embedding.
	SourceStartIndex int

	// SourceEndIndex is the first source-line index after the embedding.
	SourceEndIndex int
}

// NewContext creates a parsing context for a documentation file.
//
// It initializes MarkdownFilePath, source lines, line index, and result buffer.
//
// Parameters:
// markdownFile - identifies the documentation file to parse.
//
// Returns:
// Context - initialized parsing context.
// error - when the documentation file cannot be read.
func NewContext(markdownFile string) (Context, error) {
	resolver, err := fragmentation.NewResolver(fragmentation.DefaultResolverCacheLimit)
	if err != nil {
		return Context{}, err
	}

	return NewContextWithResolver(markdownFile, resolver)
}

// NewContextWithResolver creates a parsing context using the provided source resolver.
//
// If resolver is nil, it creates a resolver with the default cache limit.
func NewContextWithResolver(
	markdownFile string,
	resolver *fragmentation.Resolver,
) (Context, error) {
	if resolver == nil {
		var err error
		resolver, err = fragmentation.NewResolver(fragmentation.DefaultResolverCacheLimit)
		if err != nil {
			return Context{}, err
		}
	}

	source, err := readLines(markdownFile)
	if err != nil {
		return Context{}, err
	}

	return Context{
		MarkdownFilePath: markdownFile,
		Result:           make([]string, 0),
		source:           source,
		lineIndex:        0,
		resolver:         resolver,
	}, nil
}

// NewEmptyContext creates a Context for a documentation file that was not parsed.
//
// Parameters:
// markdownFile - identifies the skipped documentation file.
//
// Returns empty parsing context.
func NewEmptyContext(markdownFile string) Context {
	return Context{
		MarkdownFilePath: markdownFile,
		Result:           make([]string, 0),
		lineIndex:        0,
	}
}

// CurrentLine returns the documentation line at the current parser index.
//
// Returns the current documentation source line.
func (c *Context) CurrentLine() string {
	return c.source[c.lineIndex]
}

// CurrentIndex returns the current one-based documentation source line number.
//
// Returns one-based line number.
func (c *Context) CurrentIndex() int {
	return c.lineIndex + 1
}

// ToNextLine advances the parser to the next source line.
func (c *Context) ToNextLine() {
	c.lineIndex++
}

// ReachedEOF reports whether the parser reached the end of the documentation source file.
//
// Returns true when the parser index is at or beyond the source length.
func (c *Context) ReachedEOF() bool {
	return c.lineIndex >= len(c.source)
}

// IsContentChanged reports whether generated documentation differs from the source content.
//
// It compares generated result lines with original documentation source lines.
//
// Returns true when generated lines differ from original source lines.
func (c *Context) IsContentChanged() bool {
	if len(c.Result) != c.lineIndex {
		return true
	}

	for i := 0; i < c.lineIndex; i++ {
		if c.source[i] != c.Result[i] {
			return true
		}
	}

	return false
}

// IsContainsEmbedding reports whether the doc file contains an embedding.
//
// Returns true after at least one embedding instruction is recognized.
func (c *Context) IsContainsEmbedding() bool {
	return c.fileContainsEmbedding
}

// ResolveEmbeddingNotFound preserves the original Markdown when source content is missing.
//
// It also records the instruction for logging.
func (c *Context) ResolveEmbeddingNotFound() {
	currentEmbedding := *c.CurrentEmbedding()
	source := c.readEmbeddingSource(currentEmbedding)
	c.Result = append(c.Result, source...)
	c.EmbeddingsNotFound = append(c.EmbeddingsNotFound, currentEmbedding.embeddingInstruction)
}

// ResolveUnacceptedEmbedding records and removes an instruction rejected by the parser.
//
// It also records the instruction for logging.
func (c *Context) ResolveUnacceptedEmbedding() {
	currentEmbeddingInstruction := c.CurrentEmbedding().embeddingInstruction
	c.UnacceptedEmbeddings = append(c.UnacceptedEmbeddings, currentEmbeddingInstruction)
	c.embeddings = c.embeddings[:c.currentEmbeddingIndex()]
	c.EmbeddingInstruction = nil
}

// StartEmbedding records an instruction as the current embedding.
//
// Parameters:
// instruction - provides parsed embedding instruction data.
func (c *Context) StartEmbedding(instruction Instruction) {
	c.fileContainsEmbedding = true
	instruction.resolver = c.resolver
	embeddingContext := EmbeddingContext{
		embeddingInstruction: instruction,
	}

	c.embeddings = append(c.embeddings, embeddingContext)
	c.EmbeddingInstruction = &c.CurrentEmbedding().embeddingInstruction
}

// FinishEmbedding closes the current embedding at the current parser position.
func (c *Context) FinishEmbedding() {
	currentEmbedding := c.CurrentEmbedding()
	currentEmbedding.SourceEndIndex = c.lineIndex
	c.EmbeddingInstruction = nil
}

// SetCodeStart records the first source line belonging to the current embedding fence.
//
// It excludes the instruction and opening fence from the original embedded source range.
func (c *Context) SetCodeStart() {
	if c.fileContainsEmbedding {
		lastEmbedding := c.CurrentEmbedding()
		lastEmbedding.SourceStartIndex = c.lineIndex
	}
}

// GetResult returns the generated documentation lines.
//
// Returns generated documentation lines.
func (c *Context) GetResult() []string {
	return c.Result
}

// String returns a string representation of Context.
//
// Returns diagnostic context text.
func (c *Context) String() string {
	return fmt.Sprintf("Context[embedding=`%s`, file=`%s`, line=`%d`]",
		c.EmbeddingInstruction, c.MarkdownFilePath, c.lineIndex)
}

// CurrentEmbedding returns the embedding currently being parsed.
//
// Returns current embedding context.
func (c *Context) CurrentEmbedding() *EmbeddingContext {
	return &c.embeddings[c.currentEmbeddingIndex()]
}

// currentEmbeddingIndex returns the index of the latest embedding.
func (c *Context) currentEmbeddingIndex() int {
	return len(c.embeddings) - 1
}

// readEmbeddingSource returns original Markdown lines for one embedding.
func (c *Context) readEmbeddingSource(context EmbeddingContext) []string {
	return c.source[context.SourceStartIndex:context.SourceEndIndex]
}

// readLines returns file content as lines split on Unix or Windows line endings.
//
// Parameters:
// filepath - provides the file to read.
//
// Returns:
// []string - file content lines.
// error - when the file cannot be read.
func readLines(filepath string) ([]string, error) {
	bytes, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	str := string(bytes)
	lines := regexp.MustCompile("\r?\n").Split(str, -1)

	return lines, nil
}
