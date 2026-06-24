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
	"fmt"
	"os"
	"regexp"

	"embed-code/embed-code-go/fragmentation"
)

// Context represents the context for parsing a file containing code embeddings.
//
// EmbeddingInstruction - a pointer to the embedding instruction.
//
// MarkdownFilePath - a path to the markdown file.
//
// Result - a list of strings representing the markdown file updated with embedding.
//
// CodeFenceStarted - a flag indicating whether a code fence has been started.
//
// CodeFenceIndentation - an indentation of the markdown's code fences.
//
// EmbeddingsNotFound - a list of embedding instructions that are not found in the code.
//
// UnacceptedEmbeddings - a list of embedding instructions that are not accepted by the parser.
type Context struct {
	EmbeddingInstruction     *Instruction
	MarkdownFilePath         string
	Result                   []string
	CodeFenceStarted         bool
	CodeFenceMarker          string
	CodeFenceIndentation     int
	MarkdownFenceStarted     bool
	MarkdownFenceMarker      string
	MarkdownFenceIndentation int
	EmbeddingsNotFound       []Instruction
	UnacceptedEmbeddings     []Instruction
	// source - a list of strings representing the original markdown file.
	source []string
	// lineIndex - an index of the current line in the markdown file.
	lineIndex int
	// fileContainsEmbedding - a flag indicating whether the file contains an embedding instruction.
	fileContainsEmbedding bool
	// embeddings - a list of embedding instructions found in the markdown file.
	embeddings []EmbeddingContext
	// resolver owns source fragmentation cache state for this processing operation.
	resolver *fragmentation.Resolver
}

// EmbeddingsCount returns number of found embeddings.
func (c *Context) EmbeddingsCount() int {
	return len(c.embeddings)
}

// EmbeddingContext contains an instruction and its position in the source Markdown file.
//
// embeddingInstruction - an Instruction, containing all the needed embedding information.
//
// SourceStartIndex - an index of the StartState line in the original markdown file.
//
// SourceEndIndex - an index of the end line in the original markdown file.
type EmbeddingContext struct {
	embeddingInstruction Instruction
	SourceStartIndex     int
	SourceEndIndex       int
}

// NewContext Creates and returns a new Context struct with initial values for markdownFile, source,
// lineIndex, and result.
func NewContext(markdownFile string) (Context, error) {
	return NewContextWithResolver(markdownFile, fragmentation.NewResolver())
}

// NewContextWithResolver creates a parsing context using the provided source resolver.
func NewContextWithResolver(
	markdownFile string,
	resolver *fragmentation.Resolver,
) (Context, error) {
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
func NewEmptyContext(markdownFile string) Context {
	return Context{
		MarkdownFilePath: markdownFile,
		Result:           make([]string, 0),
		lineIndex:        0,
	}
}

// CurrentLine returns the line of source code at the current Context line index.
func (c *Context) CurrentLine() string {
	return c.source[c.lineIndex]
}

// CurrentIndex returns the current one-based source line number.
func (c *Context) CurrentIndex() int {
	return c.lineIndex + 1
}

// ToNextLine advances the parser to the next source line.
func (c *Context) ToNextLine() {
	c.lineIndex++
}

// ReachedEOF reports whether the end of the source code file has been reached.
func (c *Context) ReachedEOF() bool {
	return c.lineIndex >= len(c.source)
}

// IsContentChanged Reports whether the content of the code file has changed compared to the
// embedding of the markdown file.
func (c *Context) IsContentChanged() bool {
	for i := 0; i < c.lineIndex; i++ {
		if c.source[i] != c.Result[i] {
			return true
		}
	}

	return false
}

// IsContainsEmbedding reports whether the doc file contains an embedding.
func (c *Context) IsContainsEmbedding() bool {
	return c.fileContainsEmbedding
}

// ResolveEmbeddingNotFound writes the source content of the markdown file if embedding
// is not found.
func (c *Context) ResolveEmbeddingNotFound() {
	currentEmbedding := *c.CurrentEmbedding()
	source := c.readEmbeddingSource(currentEmbedding)
	c.Result = append(c.Result, source...)
	c.EmbeddingsNotFound = append(c.EmbeddingsNotFound, currentEmbedding.embeddingInstruction)
}

// ResolveUnacceptedEmbedding deletes embedding from the list of embeddings if it is not accepted.
//
// Also appends it to the list of such embeddings for logging.
func (c *Context) ResolveUnacceptedEmbedding() {
	currentEmbeddingInstruction := c.CurrentEmbedding().embeddingInstruction
	c.UnacceptedEmbeddings = append(c.UnacceptedEmbeddings, currentEmbeddingInstruction)
	c.embeddings = c.embeddings[:c.currentEmbeddingIndex()]
	c.EmbeddingInstruction = nil
}

// StartEmbedding records an instruction as the current embedding.
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

// SetCodeStart sets the current line as a start of a code lines in the result. It's needed to not
// include instructions in the embedding.
func (c *Context) SetCodeStart() {
	if c.fileContainsEmbedding {
		lastEmbedding := c.CurrentEmbedding()
		lastEmbedding.SourceStartIndex = c.lineIndex
	}
}

// GetResult returns the result lines of the Context.
func (c *Context) GetResult() []string {
	return c.Result
}

// Returns a string representation of Context.
func (c *Context) String() string {
	return fmt.Sprintf("Context[embedding=`%s`, file=`%s`, line=`%d`]",
		c.EmbeddingInstruction, c.MarkdownFilePath, c.lineIndex)
}

// CurrentEmbedding returns the embedding currently being parsed.
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

// readLines returns the content of a file placed at filepath as a list of strings.
func readLines(filepath string) ([]string, error) {
	bytes, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	str := string(bytes)
	lines := regexp.MustCompile("\r?\n").Split(str, -1)

	return lines, nil
}
