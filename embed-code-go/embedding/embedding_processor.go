package embedding

import (
	"embed-code/embed-code-go/configuration"
	"fmt"
	"os"
)

type EmbeddingProcessor struct {
	DocFile string
	Config  configuration.Configuration
}

func NewEmbeddingProcessor(docFile string, config configuration.Configuration) *EmbeddingProcessor {
	return &EmbeddingProcessor{
		DocFile: docFile,
		Config:  config,
	}
}

func (ep *EmbeddingProcessor) embed() {
	context := ep.constructEmbedding()

	if context.fileContainsEmbedding && context.contentChanged() {
		err := os.WriteFile(ep.DocFile, []byte(context.result()), 0644)
		if err != nil {
			fmt.Println("Error writing to file:", err)
		}
	}
}

func (ep *EmbeddingProcessor) upToDate() bool {
	context := ep.constructEmbedding()
	return !context.contentChanged()
}

func (ep *EmbeddingProcessor) constructEmbedding() ParsingContext {
	context := NewParsingContext(ep.DocFile)

	currentState := "START"
	for currentState != "FINISH" {
		accepted := false
		for _, nextState := range TRANSITIONS[currentState] {
			transition := STATE_TO_TRANSITION[nextState]
			if transition.recognize() {
				currentState = nextState
				transition.accept(context, ep.Config)
				accepted = true
				break
			}
		}
		if !accepted {
			panic(fmt.Sprintf("Failed to parse the doc file `%s`. Context: %+v", ep.DocFile, context))
		}
	}

	return context
}

type EmbeddingContext struct {
	fileContainsEmbedding bool
	contentChanged        bool
	// Other relevant fields...
}

func (ec EmbeddingProcessor) result() string {
	// Construct the final result based on the context.
	// Return the result as a string.
}

func (ep EmbeddingProcessor) constructEmbedding() ParsingContext {

}
