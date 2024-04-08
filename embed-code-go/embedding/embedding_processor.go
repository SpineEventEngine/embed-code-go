package embedding

import (
	"embed-code/embed-code-go/configuration"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type EmbeddingProcessor struct {
	DocFile string
	Config  configuration.Configuration
}

func NewEmbeddingProcessor(docFile string, config configuration.Configuration) EmbeddingProcessor {
	return EmbeddingProcessor{
		DocFile: docFile,
		Config:  config,
	}
}

func (ep EmbeddingProcessor) Embed() {
	context := ep.constructEmbedding()

	if context.checkContainsEmbedding() && context.checkContentChanged() {
		err := os.WriteFile(ep.DocFile, []byte(strings.Join(context.result, "\n")), 0644)
		if err != nil {
			panic(err)
		}
	}
}

func (ep EmbeddingProcessor) CheckUpToDate() bool {
	context := ep.constructEmbedding()
	return !context.checkContentChanged()
}

func (ep EmbeddingProcessor) constructEmbedding() ParsingContext {
	context := NewParsingContext(ep.DocFile)

	currentState := "START"
	for currentState != "FINISH" {
		accepted := false
		for _, nextState := range Transitions[currentState] {
			transition := StateToTransition[nextState]
			if transition.recognize(context) {
				currentState = nextState
				transition.accept(&context, ep.Config)
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

func EmbedAll(configuration configuration.Configuration) {
	documentationRoot := configuration.DocumentationRoot
	docPatterns := configuration.DocIncludes
	for _, pattern := range docPatterns {
		documentationFiles, _ := filepath.Glob(filepath.Join(documentationRoot, pattern))
		for _, documentationFile := range documentationFiles {
			processor := NewEmbeddingProcessor(documentationFile, configuration)
			processor.Embed()
		}
	}
}
