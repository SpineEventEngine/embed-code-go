package embedding

import (
	"bufio"
	"embed-code/embed-code-go/embedding_instruction"
	"fmt"
	"os"
)

type ParsingContext struct {
	embedding               *embedding_instruction.EmbeddingInstruction
	source                  []string
	markdownFile            string
	lineIndex               int
	result                  []string
	codeFenceStarted        bool
	codeFenceIndentation    int
	fragmentsDir            string
	file_contains_embedding bool
}

func NewParsingContext(markdownFile string) ParsingContext {
	return ParsingContext{
		markdownFile: markdownFile,
		source:       readLines(markdownFile),
		lineIndex:    0,
		result:       make([]string, 0),
	}
}

func (pc ParsingContext) currentLine() string {
	return pc.source[pc.lineIndex]
}

func (pc *ParsingContext) toNextLine() {
	pc.lineIndex++
}

func (pc ParsingContext) reachedEOF() bool {
	return pc.lineIndex >= len(pc.source)
}

func (pc ParsingContext) checkContentChanged() bool {
	for i := 0; i < pc.lineIndex; i++ {
		if pc.source[i] != pc.result[i] {
			return true
		}
	}
	return false
}

func (pc ParsingContext) checkContainsEmbedding() bool {
	return pc.file_contains_embedding
}

func (pc *ParsingContext) setEmbedding(embedding *embedding_instruction.EmbeddingInstruction) {
	pc.file_contains_embedding = true
	pc.embedding = embedding
}

func (pc ParsingContext) String() string {
	return fmt.Sprintf("ParsingContext[embedding=`%s`, file=`%s`, line=`%d`]",
		pc.embedding, pc.markdownFile, pc.lineIndex)
}

func readLines(filename string) []string {
	file, err := os.Open(filename)
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
