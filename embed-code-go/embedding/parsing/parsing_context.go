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
	"embed-code/embed-code-go/embedding_instruction"
	"fmt"
	"os"
	"strings"
)

// Represents an embedding in the parsing context.
//
// Contains the information about the position of it in the source and the resulting Markdown files.
//
// Embedding - an EmbeddingInstruction, containing all the needed embedding information.
//
// SourceStartLineIndex - an index of the start line in the original markdown file.
//
// SourceEndLineIndex - an index of the end line in the original markdown file.
//
// ResultStartLineIndex - an index of the start line in the result markdown file.
//
// ResultEndLineIndex - an index of the end line in the result markdown file.
type EmbeddingInParsingContext struct {
	Embedding            embedding_instruction.EmbeddingInstruction
	SourceStartLineIndex int
	SourceEndLineIndex   int
	ResultStartLineIndex int
	ResultEndLineIndex   int
}

// Represents the context for parsing a file containing code embeddings.
//
// Embedding - a pointer to the embedding instruction.
//
// Source - a list of strings representing the original markdown file.
//
// MarkdownFile - a path to the markdown file.
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
type ParsingContext struct {
	Embedding             *embedding_instruction.EmbeddingInstruction
	Source                []string
	MarkdownFile          string
	LineIndex             int
	Result                []string
	CodeFenceStarted      bool
	CodeFenceIndentation  int
	FileContainsEmbedding bool
	Embeddings            []EmbeddingInParsingContext
	EmbeddingsNotFound    []embedding_instruction.EmbeddingInstruction
	EmbeddingsNotAccepted []embedding_instruction.EmbeddingInstruction
}

//
// Initializers
//

// Creates and returns a new ParsingContext struct
// with initial values for markdownFile, source, lineIndex, and result.
func NewParsingContext(markdownFile string) ParsingContext {
	return ParsingContext{
		MarkdownFile: markdownFile,
		Source:       readLines(markdownFile),
		LineIndex:    0,
		Result:       make([]string, 0),
	}
}

//
// Public methods
//

// Returns the line of source code at the current ParsingContext.lineIndex.
func (pc ParsingContext) CurrentLine() string {
	return pc.Source[pc.LineIndex]
}

// Increments ParsingContext.lineIndex field by 1.
func (pc *ParsingContext) ToNextLine() {
	pc.LineIndex++
}

// Reports whether the end of the source code file has been reached.
func (pc ParsingContext) ReachedEOF() bool {
	return pc.LineIndex >= len(pc.Source)
}

// Reports whether the content of the code file has changed
// compared to the embedding of the markdown file.
func (pc ParsingContext) IsContentChanged() bool {
	for i := 0; i < pc.LineIndex; i++ {
		if pc.Source[i] != pc.Result[i] {
			return true
		}
	}
	return false
}

// Returns a list of changed embeddings.
func (pc ParsingContext) FindChangedEmbeddings() []embedding_instruction.EmbeddingInstruction {
	changedEmbeddings := make([]embedding_instruction.EmbeddingInstruction, 0)
	for _, embedding := range pc.Embeddings {
		sourceContent := pc.readEmbeddingSource(embedding)
		resultContent := pc.readEmbeddingResult(embedding)
		if !isStringSlicesEqual(sourceContent, resultContent) {
			changedEmbeddings = append(changedEmbeddings, embedding.Embedding)
		}
	}
	return changedEmbeddings
}

// Reports whether the doc file contains an embedding.
func (pc ParsingContext) IsContainsEmbedding() bool {
	return pc.FileContainsEmbedding
}

// Writes the source content of the markdown file if embedding is not found.
func (pc *ParsingContext) ResolveEmbeddingNotFound() {
	currentEmbedding := pc.Embeddings[len(pc.Embeddings)-1]
	source := pc.readEmbeddingSource(currentEmbedding)
	pc.Result = append(pc.Result, source...)
	pc.EmbeddingsNotFound = append(pc.EmbeddingsNotFound, currentEmbedding.Embedding)
}

// Deletes embedding from the list of embeddings if it is not accepted. 
//
// Also appends it to the list of such embeddings for logging.
func (pc *ParsingContext) ResolveEmbeddingNotAccepted() {
	currentEmbedding := pc.Embeddings[len(pc.Embeddings)-1]
	pc.EmbeddingsNotAccepted = append(pc.EmbeddingsNotAccepted, currentEmbedding.Embedding)
	pc.Embeddings = pc.Embeddings[:len(pc.Embeddings)-1]
	pc.SetEmbedding(nil)
}

// Sets an embedding to ParsingContext.
//
// Also sets FileContainsEmbedding flag.
func (pc *ParsingContext) SetEmbedding(embedding *embedding_instruction.EmbeddingInstruction) {
	if embedding != nil {
		pc.FileContainsEmbedding = true
		pc.Embeddings = append(pc.Embeddings, EmbeddingInParsingContext{
			Embedding:            *embedding,
			SourceStartLineIndex: pc.LineIndex + 2,   // +2 for instruction and code fence.
			ResultStartLineIndex: len(pc.Result) + 2, // +2 for instruction and code fence.
		})
	} else {
		pc.Embeddings[len(pc.Embeddings)-1].SourceEndLineIndex = pc.LineIndex
		pc.Embeddings[len(pc.Embeddings)-1].ResultEndLineIndex = len(pc.Result)
	}
	pc.Embedding = embedding
}

// Returns the result lines of the ParsingContext.
func (pc ParsingContext) GetResult() []string {
	return pc.Result
}

// Returns a string representation of ParsingContext.
func (pc ParsingContext) String() string {
	return fmt.Sprintf("ParsingContext[embedding=`%s`, file=`%s`, line=`%d`]",
		pc.Embedding, pc.MarkdownFile, pc.LineIndex)
}

//
// Private methods
//

func (pc ParsingContext) readEmbeddingSource(
	embeddingInParsingContext EmbeddingInParsingContext) []string {

	return pc.Source[embeddingInParsingContext.SourceStartLineIndex : embeddingInParsingContext.SourceEndLineIndex+1]
}

func (pc ParsingContext) readEmbeddingResult(
	embeddingInParsingContext EmbeddingInParsingContext) []string {

	return pc.Result[embeddingInParsingContext.ResultStartLineIndex : embeddingInParsingContext.ResultEndLineIndex+1]
}

//
// Static functions
//

// Returns the content of a file placed at filepath as a list of strings.
func readLines(filepath string) []string {
	bytes, err := os.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	str := string(bytes)
	str = strings.ReplaceAll(str, "\r\n", "\n")
	lines := strings.Split(str, "\n")
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
