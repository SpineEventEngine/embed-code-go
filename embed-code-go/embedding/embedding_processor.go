package embedding

import (
	"embed-code/embed-code-go/configuration"
	"embed-code/embed-code-go/embedding/parsing"
	"fmt"
	"os"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
)

const filePermission = 0644

// The EmbeddingProcessor entity processes a single documentation file and embeds code snippets
// into it based on the provided configuration.
//
// DocFile — the path to the documentation file.
//
// Config — a configuration for embedding.
type EmbeddingProcessor struct {
	DocFile        string
	Config         configuration.Configuration
	TransitionsMap map[string][]string
}

//
// Initializers
//

// Creates and returns new EmbeddingProcessor with given docFile and config.
func NewEmbeddingProcessor(docFile string, config configuration.Configuration) EmbeddingProcessor {
	return EmbeddingProcessor{
		DocFile:        docFile,
		Config:         config,
		TransitionsMap: parsing.Transitions,
	}
}

// Creates and returns new EmbeddingProcessor with given docFile, config and transitions.
func NewEmbeddingProcessorWithTransitions(docFile string,
	config configuration.Configuration,
	transitions map[string][]string,
) EmbeddingProcessor {
	return EmbeddingProcessor{
		DocFile:        docFile,
		Config:         config,
		TransitionsMap: transitions,
	}
}

//
// Public methods
//

// Constructs embedding and modifys the doc file if embedding is needed.
func (ep EmbeddingProcessor) Embed() {
	context := ep.constructEmbedding()

	if context.IsContainsEmbedding() && context.IsContentChanged() {
		err := os.WriteFile(ep.DocFile, []byte(strings.Join(context.GetResult(), "\n")), filePermission)
		if err != nil {
			panic(err)
		}
	}
}

// Reports whether the embedding of the target markdown is up-to-date with the code file.
func (ep EmbeddingProcessor) IsUpToDate() bool {
	context := ep.constructEmbedding()
	return !context.IsContentChanged()
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
		globString := strings.Join([]string{documentationRoot, pattern}, "/")
		documentationFiles, _ := doublestar.FilepathGlob(globString)
		for _, documentationFile := range documentationFiles {
			processor := NewEmbeddingProcessor(documentationFile, config)
			processor.Embed()
		}
	}
}

func CheckUpToDate(config configuration.Configuration) {
	changedFiles := findChangedFiles(config)
	if len(changedFiles) > 0 {
		panic(UnexpectedDiffError{changedFiles})
	}
}

func findChangedFiles(config configuration.Configuration) []string {
	documentationRoot := config.DocumentationRoot
	docPatterns := config.DocIncludes
	var changedFiles []string

	for _, pattern := range docPatterns {
		globString := strings.Join([]string{documentationRoot, pattern}, "/")
		matches, err := doublestar.FilepathGlob(globString)
		if err != nil {
			panic(err)
		}

		for _, documentationFile := range matches {
			upToDate := NewEmbeddingProcessor(documentationFile, config).IsUpToDate()
			if !upToDate {
				changedFiles = append(changedFiles, documentationFile)
			}
		}
	}

	return changedFiles
}
