package parsing

import (
	"embed-code/embed-code-go/configuration"

	"embed-code/embed-code-go/embedding_instruction"
	"fmt"
	"strings"
)

// Represents an embedding instruction token of a markdown.
type EmbedInstructionToken struct{}

//
// Public methods
//

// Reports whether the current line in the parsing context starts with "<embed-code",
// and if there is no ongoing embedding and the end of the file is not reached, it returns true.
// Otherwise, it returns false.
//
// context — a context of the parsing process.
func (e EmbedInstructionToken) Recognize(context ParsingContext) bool {
	line := context.CurrentLine()
	isStatement := strings.HasPrefix(strings.TrimSpace(line), Statement)
	if context.embedding == nil && !context.ReachedEOF() && isStatement {
		return true
	}
	return false
}

// Parses the embedding instruction and extracts relevant information to update the parsing context.
//
// context — a context of the parsing process.
//
// config — a configuration of the embedding.
func (e EmbedInstructionToken) Accept(context *ParsingContext, config configuration.Configuration) {
	instructionBody := []string{}
	for !context.ReachedEOF() {
		instructionBody = append(instructionBody, context.CurrentLine())
		instruction := embedding_instruction.FromXML(strings.Join(instructionBody, ""), config)
		context.SetEmbedding(&instruction)
		context.result = append(context.result, context.CurrentLine())
		context.ToNextLine()
		if context.embedding != nil {
			break
		}
	}
	if context.embedding == nil {
		panic(fmt.Sprintf("Failed to parse an embedding instruction. Context: %v", context))
	}
}
