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
	"embed-code/embed-code-go/files"
	"embed-code/embed-code-go/instruction"
	"errors"
	"fmt"
	"os"
	"strings"

	"embed-code/embed-code-go/configuration"
	"embed-code/embed-code-go/embedding/parsing"

	"github.com/bmatcuk/doublestar/v4"
)

// Processor entity processes a single documentation file and embeds code snippets
// into it based on the provided configuration.
//
// DocFilePath — the path to the documentation file.
//
// Config — a configuration for embedding.
type Processor struct {
	DocFilePath    string
	Config         configuration.Configuration
	TransitionsMap parsing.TransitionMap
}

// NewProcessor creates and returns new Processor with given docFile and config.
func NewProcessor(docFile string, config configuration.Configuration) Processor {
	return Processor{
		DocFilePath:    docFile,
		Config:         config,
		TransitionsMap: parsing.Transitions,
	}
}

// NewProcessorWithTransitions Creates and returns new Processor with given docFile, config
// and transitions.
func NewProcessorWithTransitions(docFile string, config configuration.Configuration,
	transitions parsing.TransitionMap) Processor {
	return Processor{
		DocFilePath:    docFile,
		Config:         config,
		TransitionsMap: transitions,
	}
}

// Embed Constructs embedding and modifies the doc file if embedding is needed.
//
// If any problems faced, an error is returned.
func (p Processor) Embed() error {
	context, err := p.constructEmbedding()
	if err != nil {
		return EmbeddingError{Context: context}
	}

	if context.IsContainsEmbedding() && context.IsContentChanged() {
		err = os.WriteFile(p.DocFilePath, []byte(strings.Join(context.GetResult(), "\n")),
			os.FileMode(files.ReadWriteExecPermission))
		if err != nil {
			return EmbeddingError{Context: context}
		}
	}

	return nil
}

// FindChangedEmbeddings Returns the list of EmbeddingInstruction that are changed in the
// markdown file.
//
// If any problems during the embedding construction faced, an error is returned.
func (p Processor) FindChangedEmbeddings() ([]instruction.Instruction, error) {
	context, err := p.constructEmbedding()
	changedEmbeddings := context.FindChangedEmbeddings()
	if err != nil {
		return changedEmbeddings, EmbeddingError{Context: context}
	}

	return changedEmbeddings, nil
}

// IsUpToDate reports whether the embedding of the target markdown is up-to-date with the code file.
func (p Processor) IsUpToDate() bool {
	context, err := p.constructEmbedding()
	if err != nil {
		panic(err)
	}

	return !context.IsContentChanged()
}

// Creates and returns new ParsingContext based on Processor.DocFilePath and Processor.Config.
//
// If any problems faced, an error is returned.
//
// Processes an embedding by iterating through different states based on transitions until it
// reaches the finish state. If a transition is recognized, it updates the current state and
// accepts the transition. If no transition is accepted, the error indicating the failure to parse
// the document file is returned.
func (p Processor) constructEmbedding() (parsing.Context, error) {
	context := parsing.NewContext(p.DocFilePath)
	isErrorFaced := false
	errorStr := fmt.Sprintf(
		"an error was occurred during embedding construction for doc file `%s`", p.DocFilePath)
	var constructEmbeddingError = errors.New(errorStr)

	var currentState parsing.Transition
	currentState = parsing.Start{StateName: "START"}
	finishState := parsing.Finish{StateName: "FINISH"}

	for currentState.State() != finishState.State() {
		accepted := false
		for _, nextState := range parsing.Transitions[currentState] {
			if nextState.Recognize(context) {
				currentState = nextState
				err := nextState.Accept(&context, p.Config)
				if err != nil {
					isErrorFaced = true
				}
				accepted = true

				break
			}
		}
		if !accepted {
			currentState = parsing.RegularLine{StateName: "REGULAR_LINE"}
			context.ResolveUnacceptedEmbedding()
			isErrorFaced = true
		}
	}

	var err error
	if isErrorFaced {
		err = constructEmbeddingError
	}

	return context, err
}

// EmbedAll processes embedding for multiple documentation files based on provided config.
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
			processor := NewProcessor(documentationFile, config)
			err := processor.Embed()
			if err != nil {
				panic(err)
			}
		}
	}
}

// CheckUpToDate raises an error if the documentation files are not up-to-date with code files.
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
			upToDate := NewProcessor(documentationFile, config).IsUpToDate()
			if !upToDate {
				changedFiles = append(changedFiles, documentationFile)
			}
		}
	}

	return changedFiles
}
