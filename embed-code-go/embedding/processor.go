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
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"embed-code/embed-code-go/configuration"
	"embed-code/embed-code-go/embedding/parsing"
	"embed-code/embed-code-go/files"

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
	context, err := p.fillEmbeddingContext()
	if err != nil {
		return UnexpectedProcessingError{Context: context}
	}

	if context.IsContainsEmbedding() && context.IsContentChanged() {
		data := []byte(strings.Join(context.GetResult(), "\n"))
		err = os.WriteFile(p.DocFilePath, data, os.FileMode(files.ReadWriteExecPermission))
		if err != nil {
			return UnexpectedProcessingError{context}
		}
	}

	return nil
}

// FindChangedEmbeddings Returns the list of EmbeddingInstruction that are changed in the
// markdown file.
//
// If any problems during the embedding construction faced, an error is returned.
func (p Processor) FindChangedEmbeddings() ([]parsing.Instruction, error) {
	context, err := p.fillEmbeddingContext()
	changedEmbeddings := context.FindChangedEmbeddings()
	if err != nil {
		return changedEmbeddings, UnexpectedProcessingError{context}
	}

	return changedEmbeddings, nil
}

// IsUpToDate reports whether the embedding of the target markdown is up-to-date with the code file.
func (p Processor) IsUpToDate() bool {
	context, err := p.fillEmbeddingContext()
	if err != nil {
		panic(err)
	}

	return !context.IsContentChanged()
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
		matches, _ := doublestar.FilepathGlob(globString)
		for _, match := range matches {

			err := filepath.Walk(match, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				processor := NewProcessor(match, config)
				err = processor.Embed()
				if err != nil {
					return err
				}
				return nil
			})

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

// Iterates through the doc file line by line considering them as a states of an embedding.
// Such way, transits from the state to the next possible one until it reaches the end of a file.
// By the transition process, fills the parsing.Context accordingly, so it is ready to retrieve
// the result. Returns a parsing.Context and an error if any occurs.
func (p Processor) fillEmbeddingContext() (parsing.Context, error) {
	context := parsing.NewContext(p.DocFilePath)
	errorStr := "an error was occurred during embedding construction for doc file `%s` at line %v"

	var currentState parsing.State
	currentState = parsing.Start
	finishState := parsing.Finish

	for currentState != finishState {
		accepted, newState, err := p.moveToNextState(&currentState, &context)
		if err != nil {
			return parsing.Context{}, fmt.Errorf(errorStr, p.DocFilePath, context.CurrentIndex())
		}
		if !accepted {
			currentState = &parsing.RegularLineState{}
			context.ResolveUnacceptedEmbedding()
			return context, fmt.Errorf(errorStr, p.DocFilePath, context.CurrentIndex())
		}
		currentState = *newState
	}

	return context, nil
}

// Moves to the next state accordingly to a transition map from the current state. Reports whether
// it successfully moved to the next state and returns the new state.
func (p Processor) moveToNextState(state *parsing.State, context *parsing.Context) (
	bool, *parsing.State, error) {
	for _, nextState := range parsing.Transitions[*state] {
		if nextState.Recognize(*context) {
			err := nextState.Accept(context, p.Config)
			if err != nil {
				return false, &nextState, err
			}

			return true, &nextState, nil
		}
	}

	return false, state, nil
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
		for _, match := range matches {

			err := filepath.Walk(match, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				upToDate := NewProcessor(match, config).IsUpToDate()
				if !upToDate {
					changedFiles = append(changedFiles, match)
				}
				return nil
			})

			if err != nil {
				fmt.Println("error walking through directory:", err)
			}

		}
	}

	return changedFiles
}
