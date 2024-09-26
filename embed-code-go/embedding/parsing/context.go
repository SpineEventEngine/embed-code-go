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
	EmbeddingInstruction *Instruction
	MarkdownFilePath     string
	Result               []string
	CodeFenceStarted     bool
	CodeFenceIndentation int
	EmbeddingsNotFound   []Instruction
	UnacceptedEmbeddings []Instruction
	// source - a list of strings representing the original markdown file.
	source []string
	// lineIndex - an index of the current line in the markdown file.
	lineIndex int
	// fileContainsEmbedding - a flag indicating whether the file contains an embedding instruction.
	fileContainsEmbedding bool
	// embeddings - a list of embedding instructions found in the markdown file.
	embeddings []parsingContext
}

// parsingContext contains the information about the position in the source and the
// resulting Markdown files.
//
// embeddingInstruction - an Instruction, containing all the needed embedding information.
//
// sourceStartIndex - an index of the StartState line in the original markdown file.
//
// sourceEndIndex - an index of the end line in the original markdown file.
//
// resultStartIndex - an index of the StartState line in the result markdown file.
//
// resultEndIndex - an index of the end line in the result markdown file.
type parsingContext struct {
	embeddingInstruction Instruction
	sourceStartIndex     int
	sourceEndIndex       int
	resultStartIndex     int
	resultEndIndex       int
}

// NewContext Creates and returns a new Context struct with initial values for markdownFile, source,
// lineIndex, and result.
func NewContext(markdownFile string) Context {
	return Context{
		MarkdownFilePath: markdownFile,
		Result:           make([]string, 0),
		source:           readLines(markdownFile),
		lineIndex:        0,
	}
}

// CurrentLine returns the line of source code at the current ParsingContext.lineIndex.
func (c *Context) CurrentLine() string {
	return c.source[c.lineIndex]
}

// ToNextLine increments ParsingContext.lineIndex field by 1.
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

// FindChangedEmbeddings returns a list of changed embeddings.
func (c *Context) FindChangedEmbeddings() []Instruction {
	var changedEmbeddings []Instruction
	for _, embedding := range c.embeddings {
		sourceContent := c.readEmbeddingSource(embedding)
		resultContent := c.readEmbeddingResult(embedding)
		if !isStringSlicesEqual(sourceContent, resultContent) {
			changedEmbeddings = append(changedEmbeddings, embedding.embeddingInstruction)
		}
	}

	return changedEmbeddings
}

// IsContainsEmbedding reports whether the doc file contains an embedding.
func (c *Context) IsContainsEmbedding() bool {
	return c.fileContainsEmbedding
}

// ResolveEmbeddingNotFound writes the source content of the markdown file if embedding
// is not found.
func (c *Context) ResolveEmbeddingNotFound() {
	currentEmbedding := *c.currentEmbedding()
	source := c.readEmbeddingSource(currentEmbedding)
	c.Result = append(c.Result, source...)
	c.EmbeddingsNotFound = append(c.EmbeddingsNotFound, currentEmbedding.embeddingInstruction)
}

// ResolveUnacceptedEmbedding deletes embedding from the list of embeddings if it is not accepted.
//
// Also appends it to the list of such embeddings for logging.
func (c *Context) ResolveUnacceptedEmbedding() {
	currentEmbeddingInstruction := c.currentEmbedding().embeddingInstruction
	c.UnacceptedEmbeddings = append(c.UnacceptedEmbeddings, currentEmbeddingInstruction)
	c.embeddings = c.embeddings[:c.currentEmbeddingIndex()]
	c.SetEmbedding(nil)
}

// SetEmbedding sets an embedding to Context. Also sets fileContainsEmbedding flag.
func (c *Context) SetEmbedding(embedding *Instruction) {
	sourceIndex := c.lineIndex
	resultIndex := len(c.Result) + 1

	if embedding == nil {
		c.currentEmbedding().sourceEndIndex = sourceIndex
		c.currentEmbedding().resultEndIndex = resultIndex
	} else {
		c.fileContainsEmbedding = true
		context := parsingContext{
			embeddingInstruction: *embedding,
		}

		c.embeddings = append(c.embeddings, context)
	}
	c.EmbeddingInstruction = embedding
}

// SetCodeStart sets the current line as a start of a code lines in the result. It's needed to not
// include instructions in the embedding.
func (c *Context) SetCodeStart() {
	if c.fileContainsEmbedding {
		lastEmbedding := c.currentEmbedding()
		lastEmbedding.sourceStartIndex = c.lineIndex
		lastEmbedding.resultStartIndex = len(c.Result) + 1
	}
}

// GetResult returns the result lines of the Context.
func (c *Context) GetResult() []string {
	return c.Result
}

// Returns a string representation of Context.
func (c *Context) String() string {
	return fmt.Sprintf("ParsingContext[embedding=`%s`, file=`%s`, line=`%d`]",
		c.EmbeddingInstruction, c.MarkdownFilePath, c.lineIndex)
}

func (c *Context) currentEmbedding() *parsingContext {
	return &c.embeddings[c.currentEmbeddingIndex()]
}

func (c *Context) currentEmbeddingIndex() int {
	return len(c.embeddings) - 1
}

func (c *Context) readEmbeddingSource(context parsingContext) []string {
	return c.source[context.sourceStartIndex:context.sourceEndIndex]
}

func (c *Context) readEmbeddingResult(context parsingContext) []string {
	return c.Result[context.resultStartIndex:context.resultEndIndex]
}

// Returns the content of a file placed at filepath as a list of strings.
func readLines(filepath string) []string {
	bytes, err := os.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	str := string(bytes)
	lines := regexp.MustCompile("\r?\n").Split(str, -1)

	return lines
}

func isStringSlicesEqual(first, second []string) bool {
	if len(first) != len(second) {
		return false
	}
	for i := range first {
		if first[i] != second[i] {
			return false
		}
	}

	return true
}
