package parsing

import (
	"bufio"
	"embed-code/embed-code-go/embedding_instruction"
	"fmt"
	"os"
)

// Represents the context for parsing a file containing code embeddings.
//
// embedding - a pointer to the embedding instruction.
//
// source - a list of strings representing the original markdown file.
//
// markdownFile - a path to the markdown file.
//
// lineIndex - an index of the current line in the markdown file.
//
// result - a list of strings representing the markdown file updated with embedding.
//
// codeFenceStarted - a flag indicating whether a code fence has been started.
//
// codeFenceIndentation - an indentation of the markdown's code fences.
//
// file_contains_embedding - a flag indicating whether the file contains an embedding instruction.
type ParsingContext struct {
	embedding               *embedding_instruction.EmbeddingInstruction
	source                  []string
	markdownFile            string
	lineIndex               int
	result                  []string
	codeFenceStarted        bool
	codeFenceIndentation    int
	file_contains_embedding bool
}

//
// Initializers
//

// Creates and returns a new ParsingContext struct
// with initial values for markdownFile, source, lineIndex, and result.
func NewParsingContext(markdownFile string) ParsingContext {
	return ParsingContext{
		markdownFile: markdownFile,
		source:       readLines(markdownFile),
		lineIndex:    0,
		result:       make([]string, 0),
	}
}

//
// Public methods
//

// Returns the line of source code at the current ParsingContext.lineIndex.
func (pc ParsingContext) CurrentLine() string {
	return pc.source[pc.lineIndex]
}

// Increments ParsingContext.lineIndex field by 1.
func (pc *ParsingContext) ToNextLine() {
	pc.lineIndex++
}

// Reports whether the end of the source code file has been reached.
func (pc ParsingContext) ReachedEOF() bool {
	return pc.lineIndex >= len(pc.source)
}

// Reports whether the content of the code file has changed
// compared to the embedding of the markdown file.
func (pc ParsingContext) IsContentChanged() bool {
	for i := 0; i < pc.lineIndex; i++ {
		if pc.source[i] != pc.result[i] {
			return true
		}
	}
	return false
}

// Reports whether the doc file contains an embedding.
func (pc ParsingContext) IsContainsEmbedding() bool {
	return pc.file_contains_embedding
}

// Sets an embedding to ParsingContext.
//
// Also sets file_contains_embedding flag.
func (pc *ParsingContext) SetEmbedding(embedding *embedding_instruction.EmbeddingInstruction) {
	if embedding != nil {
		pc.file_contains_embedding = true
	}
	pc.embedding = embedding
}

// Returns the result lines of the ParsingContext.
func (pc ParsingContext) GetResult() []string {
	return pc.result
}

// Returns a string representation of ParsingContext.
func (pc ParsingContext) String() string {
	return fmt.Sprintf("ParsingContext[embedding=`%s`, file=`%s`, line=`%d`]",
		pc.embedding, pc.markdownFile, pc.lineIndex)
}

//
// Static functions
//

// Returns the content of a file placed at filepath as a list of strings.
func readLines(filepath string) []string {
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return nil
	}
	return lines
}
