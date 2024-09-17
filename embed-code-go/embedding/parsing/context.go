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

// EmbeddingInParsingContext contains the information about the position in the source and the
// resulting Markdown files.
//
// Embedding - an Instruction, containing all the needed embedding information.
//
// SourceStartIndex - an index of the start line in the original markdown file.
//
// SourceEndIndex - an index of the end line in the original markdown file.
//
// ResultStartIndex - an index of the start line in the result markdown file.
//
// ResultEndIndex - an index of the end line in the result markdown file.
type EmbeddingInParsingContext struct {
	Embedding        Instruction
	SourceStartIndex int
	SourceEndIndex   int
	ResultStartIndex int
	ResultEndIndex   int
}

// Context represents the context for parsing a file containing code embeddings.
//
// Embedding - a pointer to the embedding instruction.
//
// Source - a list of strings representing the original markdown file.
//
// MarkdownFilePath - a path to the markdown file.
//
// LineIndex - an index of the current line in the markdown file.
//
// Result - a list of strings representing the markdown file updated with embedding.
//
// CodeFenceStarted - a flag indicating whether a code fence has been started.
//
// CodeFenceIndentation - an indentation of the markdown's code fences.
//
// FileContainsEmbedding - a flag indicating whether the file contains an embedding instruction.
//
// Embeddings - a list of embedding instructions found in the markdown file.
//
// EmbeddingsNotFound - a list of embedding instructions that are not found in the code.
//
// UnacceptedEmbeddings - a list of embedding instructions that are not accepted by the parser.
type Context struct {
	Embedding             *Instruction
	Source                []string
	MarkdownFilePath      string
	LineIndex             int
	Result                []string
	CodeFenceStarted      bool
	CodeFenceIndentation  int
	FileContainsEmbedding bool
	Embeddings            []EmbeddingInParsingContext
	EmbeddingsNotFound    []Instruction
	UnacceptedEmbeddings  []Instruction
}

// NewContext Creates and returns a new Context struct with initial values for markdownFile, source,
// lineIndex, and result.
func NewContext(markdownFile string) Context {
	return Context{
		MarkdownFilePath: markdownFile,
		Source:           readLines(markdownFile),
		LineIndex:        0,
		Result:           make([]string, 0),
	}
}

// CurrentLine returns the line of source code at the current ParsingContext.lineIndex.
func (c *Context) CurrentLine() string {
	return c.Source[c.LineIndex]
}

// ToNextLine increments ParsingContext.lineIndex field by 1.
func (c *Context) ToNextLine() {
	c.LineIndex++
}

// ReachedEOF reports whether the end of the source code file has been reached.
func (c *Context) ReachedEOF() bool {
	return c.LineIndex >= len(c.Source)
}

// IsContentChanged Reports whether the content of the code file has changed compared to the
// embedding of the markdown file.
func (c *Context) IsContentChanged() bool {
	for i := 0; i < c.LineIndex; i++ {
		if c.Source[i] != c.Result[i] {
			return true
		}
	}

	return false
}

// FindChangedEmbeddings returns a list of changed embeddings.
func (c *Context) FindChangedEmbeddings() []Instruction {
	var changedEmbeddings []Instruction
	for _, embedding := range c.Embeddings {
		sourceContent := c.readEmbeddingSource(embedding)
		resultContent := c.readEmbeddingResult(embedding)
		if !isStringSlicesEqual(sourceContent, resultContent) {
			changedEmbeddings = append(changedEmbeddings, embedding.Embedding)
		}
	}

	return changedEmbeddings
}

// IsContainsEmbedding reports whether the doc file contains an embedding.
func (c *Context) IsContainsEmbedding() bool {
	return c.FileContainsEmbedding
}

// ResolveEmbeddingNotFound writes the source content of the markdown file if embedding
// is not found.
func (c *Context) ResolveEmbeddingNotFound() {
	currentEmbedding := c.Embeddings[len(c.Embeddings)-1]
	source := c.readEmbeddingSource(currentEmbedding)
	c.Result = append(c.Result, source...)
	c.EmbeddingsNotFound = append(c.EmbeddingsNotFound, currentEmbedding.Embedding)
}

// ResolveUnacceptedEmbedding deletes embedding from the list of embeddings if it is not accepted.
//
// Also appends it to the list of such embeddings for logging.
func (c *Context) ResolveUnacceptedEmbedding() {
	currentEmbedding := c.Embeddings[len(c.Embeddings)-1]
	c.UnacceptedEmbeddings = append(c.UnacceptedEmbeddings, currentEmbedding.Embedding)
	c.Embeddings = c.Embeddings[:len(c.Embeddings)-1]
	c.SetEmbedding(nil)
}

// SetEmbedding sets an embedding to Context. Also sets FileContainsEmbedding flag.
func (c *Context) SetEmbedding(embedding *Instruction) {
	// TODO:2024-09-05:olena-zmiiova: https://github.com/SpineEventEngine/embed-code/issues/48
	indexIncrease := 2 // +2 for instruction and code fence.
	if embedding != nil {
		c.FileContainsEmbedding = true
		c.Embeddings = append(c.Embeddings, EmbeddingInParsingContext{
			Embedding:        *embedding,
			SourceStartIndex: c.LineIndex + indexIncrease,
			ResultStartIndex: len(c.Result) + indexIncrease,
		})
	} else {
		c.Embeddings[len(c.Embeddings)-1].SourceEndIndex = c.LineIndex
		c.Embeddings[len(c.Embeddings)-1].ResultEndIndex = len(c.Result)
	}
	c.Embedding = embedding
}

// GetResult returns the result lines of the Context.
func (c *Context) GetResult() []string {
	return c.Result
}

// Returns a string representation of Context.
func (c *Context) String() string {
	return fmt.Sprintf("ParsingContext[embedding=`%s`, file=`%s`, line=`%d`]",
		c.Embedding, c.MarkdownFilePath, c.LineIndex)
}

func (c *Context) readEmbeddingSource(context EmbeddingInParsingContext) []string {
	return c.Source[context.SourceStartIndex : context.SourceEndIndex+1]
}

func (c *Context) readEmbeddingResult(context EmbeddingInParsingContext) []string {
	return c.Result[context.ResultStartIndex : context.ResultEndIndex+1]
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
