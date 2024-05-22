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

package embedding

import (
	"embed-code/embed-code-go/configuration"
	"embed-code/embed-code-go/embedding/parsing"
	"fmt"
	"os"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
)

// Represents read and write permissions for the owner of the file, and read-only permissions for group and others.
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

// Constructs embedding and modifies the doc file if embedding is needed.
func (ep EmbeddingProcessor) Embed() error {
	context, err := ep.constructEmbedding()
	if err != nil {
		return EmbeddingError{Context: context, OriginalError: err}
	}

	if context.IsContainsEmbedding() && context.IsContentChanged() {
		err := os.WriteFile(ep.DocFile, []byte(strings.Join(context.GetResult(), "\n")), filePermission)
		if err != nil {
			return EmbeddingError{Context: context, OriginalError: err}
		}
	}
	return nil
}

// Reports whether the embedding of the target markdown is up-to-date with the code file.
func (ep EmbeddingProcessor) IsUpToDate() bool {
	context, err := ep.constructEmbedding()
	if err != nil {
		panic(err)
	}
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
func (ep EmbeddingProcessor) constructEmbedding() (parsing.ParsingContext, error) {
	context := parsing.NewParsingContext(ep.DocFile)

	currentState := "START"
	for currentState != "FINISH" {
		accepted := false
		for _, nextState := range parsing.Transitions[currentState] {
			transition := parsing.StateToTransition[nextState]
			if transition.Recognize(context) {
				currentState = nextState
				err := transition.Accept(&context, ep.Config)
				if err != nil {
					return context, err
				}
				accepted = true
				break
			}
		}
		if !accepted {
			return context, fmt.Errorf(fmt.Sprintf("failed to parse the doc file `%s`. Context: %+v", ep.DocFile, context))
		}
	}

	return context, nil
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
			err := processor.Embed()
			if err != nil {
				panic(err)
			}
		}
	}
}

// Raises an error if the documentation files are not up-to-date with code files.
//
// config — a configuration for embedding.
func CheckUpToDate(config configuration.Configuration) {
	changedFiles := findChangedFiles(config)
	if len(changedFiles) > 0 {
		panic(UnexpectedDiffError{changedFiles})
	}
}

// Returns a list of documentation files that are not up-to-date with their code files.
//
// config — a configuration for embedding.
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
