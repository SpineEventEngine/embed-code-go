package embedding

import (
	"embed-code/embed-code-go/configuration"
	"embed-code/embed-code-go/embedding/parsing"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// The EmbeddingProcessor entity processes a single documentation file and embeds code snippets
// into it based on the provided configuration.
//
// DocFile — the path to the documentation file.
//
// Config — a configuration for embedding.
type EmbeddingProcessor struct {
	DocFile string
	Config  configuration.Configuration
}

//
// Initializers
//

// Creates and returns new EmbeddingProcessor with the given docFile and config.
func NewEmbeddingProcessor(docFile string, config configuration.Configuration) EmbeddingProcessor {
	return EmbeddingProcessor{
		DocFile: docFile,
		Config:  config,
	}
}

//
// Public methods
//

// Constructs embedding and modifys the doc file if embedding is needed.
func (ep EmbeddingProcessor) Embed() {
	context := ep.constructEmbedding()

	if context.СheckContainsEmbedding() && context.СheckContentChanged() {
		err := os.WriteFile(ep.DocFile, []byte(strings.Join(context.GetResult(), "\n")), 0644)
		if err != nil {
			panic(err)
		}
	}
}

// Reports whether the embedding of the target markdown is up-to-date with the code file.
func (ep EmbeddingProcessor) CheckUpToDate() bool {
	context := ep.constructEmbedding()
	return !context.СheckContentChanged()
}

//
//	Private methods
//

// Creates and returns new ParsingContext based on
// EmbeddingProcessor.DocFile and EmbeddingProcessor.Config.
//
// Processes an embedding by iterating through different states based on transitions
// until it reaches the finish state. If a transition is recognized,
// it updates the current state and accepts the transition.
// If no transition is accepted, it panics with a message
// indicating the failure to parse the document file.
func (ep EmbeddingProcessor) constructEmbedding() parsing.ParsingContext {
	context := parsing.NewParsingContext(ep.DocFile)

	currentState := "START"
	for currentState != "FINISH" {
		accepted := false
		for _, nextState := range parsing.Transitions[currentState] {
			transition := parsing.StateToTransition[nextState]
			if transition.Recognize(context) {
				currentState = nextState
				transition.Accept(&context, ep.Config)
				accepted = true
				break
			}
		}
		if !accepted {
			panic(fmt.Sprintf("failed to parse the doc file `%s`. Context: %+v", ep.DocFile, context))
		}
	}

	return context
}

//
// Static functions
//

// Processes embedding for multiple documentation files based on provided config.
//
// Iterates over patterns in the configuration, finds documentation files matching those patterns,
// creates an EmbeddingProcessor for each file, and embeds code fragments in them.
//
// config — a configuration for embedding.
func EmbedAll(config configuration.Configuration) {
	documentationRoot := config.DocumentationRoot
	docPatterns := config.DocIncludes
	for _, pattern := range docPatterns {
		documentationFiles, _ := filepath.Glob(filepath.Join(documentationRoot, pattern))
		for _, documentationFile := range documentationFiles {
			processor := NewEmbeddingProcessor(documentationFile, config)
			processor.Embed()
		}
	}
}
