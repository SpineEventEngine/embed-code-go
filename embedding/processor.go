// Copyright 2026, TeamDev. All rights reserved.
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
	"errors"
	"fmt"
	"log/slog"
	"os"
	"slices"
	"strings"

	"embed-code/embed-code-go/configuration"
	"embed-code/embed-code-go/embedding/parsing"
	"embed-code/embed-code-go/files"
	"embed-code/embed-code-go/logging"
)

// Processor processes a single documentation file using the provided embedding configuration.
type Processor struct {
	// docFilePath is the path to the documentation file.
	docFilePath string

	// config contains the embedding settings.
	config configuration.Configuration

	// transitionsMap defines valid parser state transitions.
	transitionsMap parsing.TransitionMap

	// requiredDocPaths contains documentation files included by the configuration.
	requiredDocPaths []string
}

// NewProcessor creates and returns new Processor with given docFile and config.
func NewProcessor(docFile string, config configuration.Configuration) (Processor, error) {
	requiredDocPaths, err := requiredDocs(config)
	if err != nil {
		return Processor{}, err
	}

	return newProcessor(docFile, config, parsing.Transitions, requiredDocPaths), nil
}

// newProcessor creates a Processor with a precomputed documentation file list.
func newProcessor(
	docFile string,
	config configuration.Configuration,
	transitions parsing.TransitionMap,
	requiredDocPaths []string,
) Processor {
	return Processor{
		docFilePath:      docFile,
		config:           config,
		transitionsMap:   transitions,
		requiredDocPaths: requiredDocPaths,
	}
}

// Embed constructs embedding and modifies the doc file if embedding is needed.
//
// Returns an empty context without parsing the file when it is excluded by configuration.
// If any problems faced, an error is returned.
func (p Processor) Embed() (*parsing.Context, error) {
	if !slices.Contains(p.requiredDocPaths, p.docFilePath) {
		slog.Info(fmt.Sprintf("Skipping `%s`; it is excluded by the configuration.",
			logging.FileReference(p.docFilePath)))
		context := parsing.NewEmptyContext(p.docFilePath)

		return &context, nil
	}

	slog.Info(fmt.Sprintf("Started processing doc file `%s`.", logging.FileReference(p.docFilePath)))
	context, err := p.fillEmbeddingContext()
	if err != nil {
		return nil, err
	}
	if context.IsContainsEmbedding() && context.IsContentChanged() {
		data := []byte(strings.Join(context.GetResult(), "\n"))
		err = os.WriteFile(p.docFilePath, data, os.FileMode(files.DocumentationFilePermission))
		if err != nil {
			return &context, err
		}
		slog.Info(fmt.Sprintf("Updated `%s` after processing %d embedding(s).",
			logging.FileReference(p.docFilePath), context.EmbeddingsCount()))
	} else {
		slog.Info(fmt.Sprintf(
			"Documentation is up-to-date in `%s`.",
			logging.FileReference(p.docFilePath),
		))
	}

	return &context, nil
}

// IsUpToDate reports whether the embedding of the target markdown is up-to-date with the code file.
func (p Processor) IsUpToDate() bool {
	upToDate, err := p.isUpToDate()
	if err != nil {
		return false
	}

	return upToDate
}

// isUpToDate reports whether the target markdown is up-to-date and returns processing errors.
func (p Processor) isUpToDate() (bool, error) {
	if !slices.Contains(p.requiredDocPaths, p.docFilePath) {
		slog.Info(fmt.Sprintf("Skipping `%s`; it is excluded by the configuration.",
			logging.FileReference(p.docFilePath)))

		return true, nil
	}
	slog.Info(fmt.Sprintf("Checking `%s`.", logging.FileReference(p.docFilePath)))
	context, err := p.fillEmbeddingContext()
	if err != nil {
		return false, err
	}

	upToDate := !context.IsContentChanged()
	status := "up to date"
	if !upToDate {
		status = "needs an update"
	}
	slog.Info(fmt.Sprintf("Checked `%s`: %d embedding(s), %s.",
		logging.FileReference(p.docFilePath), context.EmbeddingsCount(), status))

	return upToDate, nil
}

// Iterates through the doc file line by line considering them as a states of an embedding.
// Such way, transits from the state to the next possible one until it reaches the end of a file.
// By the transition process, fills the parsing.Context accordingly, so it is ready to retrieve
// the result.
//
// Returns a parsing.Context and an error if any occurs.
func (p Processor) fillEmbeddingContext() (parsing.Context, error) {
	context, err := parsing.NewContext(p.docFilePath)
	if err != nil {
		return context, err
	}

	var currentState parsing.State
	currentState = parsing.Start
	finishState := parsing.Finish

	for currentState != finishState {
		accepted, newState, err := p.moveToNextState(&currentState, &context)
		if err != nil {
			return context, p.processingError(context, err)
		}
		if !accepted {
			err = unacceptedTransitionError(context)
			currentState = &parsing.RegularLineState{}
			if context.EmbeddingInstruction != nil {
				context.ResolveUnacceptedEmbedding()
			}

			return context, p.processingError(context, err)
		}
		currentState = *newState
	}

	return context, nil
}

// processingError wraps a parsing error with the current documentation location.
func (p Processor) processingError(context parsing.Context, err error) ProcessingError {
	return ProcessingError{
		DocFilePath: p.docFilePath,
		Line:        errorLine(context, err),
		Err:         err,
	}
}

// errorLine returns the source line that should be used in the embedding error location.
func errorLine(context parsing.Context, err error) int {
	var parseErr parsing.InstructionParseError
	if errors.As(err, &parseErr) {
		return parseErr.Line
	}
	var missingFenceErr parsing.MissingCodeFenceError
	if errors.As(err, &missingFenceErr) {
		return missingFenceErr.Line
	}
	var unclosedFenceErr parsing.UnclosedCodeFenceError
	if errors.As(err, &unclosedFenceErr) {
		return unclosedFenceErr.Line
	}
	var patternErr parsing.PatternNotFoundError
	if errors.As(err, &patternErr) {
		return patternErr.Line
	}
	if context.EmbeddingsCount() > 0 {
		return context.CurrentEmbedding().SourceStartIndex - 1
	}

	return context.CurrentIndex()
}

// unacceptedTransitionError explains why the parser could not accept the current state.
func unacceptedTransitionError(context parsing.Context) error {
	if context.EmbeddingInstruction != nil && context.CodeFenceStarted {
		return parsing.UnclosedCodeFenceError{
			Line: context.EmbeddingInstruction.DocumentationLine,
		}
	}
	if context.EmbeddingInstruction != nil && !context.CodeFenceStarted {
		return parsing.MissingCodeFenceError{
			Line: context.EmbeddingInstruction.DocumentationLine,
		}
	}

	return fmt.Errorf("unexpected parser state at line %d", context.CurrentIndex())
}

// Moves to the next state accordingly to a transition map from the current state. Reports whether
// it successfully moved to the next state and returns the new state.
func (p Processor) moveToNextState(state *parsing.State, context *parsing.Context) (
	bool, *parsing.State, error) {
	for _, nextState := range p.transitionsMap[*state] {
		if nextState.Recognize(*context) {
			err := nextState.Accept(context, p.config)
			if err != nil {
				return false, &nextState, err
			}

			return true, &nextState, nil
		}
	}

	return false, state, nil
}
