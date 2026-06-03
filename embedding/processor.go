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
	"path/filepath"
	"slices"
	"strings"

	"embed-code/embed-code-go/configuration"
	"embed-code/embed-code-go/embedding/parsing"
	"embed-code/embed-code-go/files"
	"embed-code/embed-code-go/logging"

	"github.com/bmatcuk/doublestar/v4"
)

// Processor entity processes a single documentation file and embeds code snippets
// into it based on the provided configuration.
//
// DocFilePath — the path to the documentation file.
//
// Config — a configuration for embedding.
type Processor struct {
	DocFilePath      string
	Config           configuration.Configuration
	TransitionsMap   parsing.TransitionMap
	requiredDocPaths []string
}

// EmbedAllResult is result of the EmbedAll method.
//
// TargetFiles is the list of target documentation files.
//
// TotalEmbeddings is the total number of embeddings found in the target documentation files.
//
// UpdatedTargetFiles is the list of updated target documentation files.
type EmbedAllResult struct {
	TargetFiles        []string
	TotalEmbeddings    int
	UpdatedTargetFiles []string
}

// NewProcessor creates and returns new Processor with given docFile and config.
func NewProcessor(docFile string, config configuration.Configuration) (Processor, error) {
	requiredDocPaths, err := requiredDocs(config)
	if err != nil {
		return Processor{}, err
	}

	return newProcessor(docFile, config, parsing.Transitions, requiredDocPaths), nil
}

// NewProcessorWithTransitions Creates and returns new Processor with given docFile, config
// and transitions.
func NewProcessorWithTransitions(docFile string, config configuration.Configuration,
	transitions parsing.TransitionMap) (Processor, error) {
	requiredDocPaths, err := requiredDocs(config)
	if err != nil {
		return Processor{}, err
	}

	return newProcessor(docFile, config, transitions, requiredDocPaths), nil
}

// newProcessor creates a Processor with a precomputed documentation file list.
func newProcessor(
	docFile string,
	config configuration.Configuration,
	transitions parsing.TransitionMap,
	requiredDocPaths []string,
) Processor {
	return Processor{
		DocFilePath:      docFile,
		Config:           config,
		TransitionsMap:   transitions,
		requiredDocPaths: requiredDocPaths,
	}
}

// Embed constructs embedding and modifies the doc file if embedding is needed.
//
// Returns an empty context without parsing the file when it is excluded by configuration.
// If any problems faced, an error is returned.
func (p Processor) Embed() (*parsing.Context, error) {
	if !slices.Contains(p.requiredDocPaths, p.DocFilePath) {
		slog.Info(fmt.Sprintf("Skipping `%s`; it is excluded by the configuration.",
			logging.FileReference(p.DocFilePath)))
		context := parsing.NewEmptyContext(p.DocFilePath)

		return &context, nil
	}

	slog.Info(fmt.Sprintf("Started processing doc file `%s`.", logging.FileReference(p.DocFilePath)))
	context, err := p.fillEmbeddingContext()
	if err != nil {
		return nil, err
	}
	if context.IsContainsEmbedding() && context.IsContentChanged() {
		data := []byte(strings.Join(context.GetResult(), "\n"))
		err = os.WriteFile(p.DocFilePath, data, os.FileMode(files.ReadWriteExecPermission))
		if err != nil {
			return &context, err
		}
		slog.Info(fmt.Sprintf("Updated `%s` after processing %d embedding(s).",
			logging.FileReference(p.DocFilePath), context.EmbeddingsCount()))
	} else {
		slog.Info(fmt.Sprintf(
			"Documentation is up-to-date in `%s`.",
			logging.FileReference(p.DocFilePath),
		))
	}

	return &context, nil
}

// FindChangedEmbeddings Returns the list of EmbeddingInstruction that are changed in the
// markdown file.
//
// If any problems during the embedding construction faced, an error is returned.
func (p Processor) FindChangedEmbeddings() ([]parsing.Instruction, error) {
	if !slices.Contains(p.requiredDocPaths, p.DocFilePath) {
		return nil, nil
	}
	context, err := p.fillEmbeddingContext()
	changedEmbeddings := context.FindChangedEmbeddings()
	if err != nil {
		return changedEmbeddings, err
	}

	return changedEmbeddings, nil
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
	if !slices.Contains(p.requiredDocPaths, p.DocFilePath) {
		slog.Info(fmt.Sprintf("Skipping `%s`; it is excluded by the configuration.",
			logging.FileReference(p.DocFilePath)))

		return true, nil
	}
	slog.Info(fmt.Sprintf("Checking `%s`.", logging.FileReference(p.DocFilePath)))
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
		logging.FileReference(p.DocFilePath), context.EmbeddingsCount(), status))

	return upToDate, nil
}

// EmbedAll processes embedding for multiple documentation files based on provided config.
//
// Iterates over patterns in the configuration, finds documentation files matching those patterns,
// creates an EmbeddingProcessor for each file, and embeds code fragments in them.
//
// config — a configuration for embedding.
func EmbedAll(config configuration.Configuration) (EmbedAllResult, error) {
	requiredDocPaths, err := requiredDocs(config)
	if err != nil {
		return EmbedAllResult{}, err
	}
	totalEmbeddings := 0
	var updatedTargetFiles []string
	var embeddingErrors []error
	for _, doc := range requiredDocPaths {
		processor := newProcessor(doc, config, parsing.Transitions, requiredDocPaths)
		context, err := processor.Embed()
		if err != nil {
			embeddingErrors = append(embeddingErrors, err)

			continue
		}
		totalEmbeddings += context.EmbeddingsCount()
		if context.IsContentChanged() {
			updatedTargetFiles = append(updatedTargetFiles, doc)
		}
	}
	if len(embeddingErrors) > 0 {
		return EmbedAllResult{}, errors.Join(embeddingErrors...)
	}
	if totalEmbeddings > 0 {
		slog.Info(
			fmt.Sprintf(
				"Processed %d documentation file(s) with %d embedding(s) in `%s`%s.",
				len(requiredDocPaths), totalEmbeddings,
				logging.FileReference(config.DocumentationRoot),
				configNameLabel(config),
			),
		)
	} else {
		slog.Warn(
			fmt.Sprintf("No embedding instructions were found in documentation folder `%s`%s.",
				logging.FileReference(config.DocumentationRoot), configNameLabel(config)),
		)
	}

	return EmbedAllResult{
		TargetFiles:        requiredDocPaths,
		TotalEmbeddings:    totalEmbeddings,
		UpdatedTargetFiles: updatedTargetFiles,
	}, nil
}

// configNameLabel formats a configuration name for summary log messages.
func configNameLabel(config configuration.Configuration) string {
	if config.Name == "" {
		return ""
	}

	return fmt.Sprintf(" for `%s` embedding setup", config.Name)
}

// CheckUpToDate returns documentation files that are not up-to-date with code files.
//
// config — a configuration for embedding.
func CheckUpToDate(config configuration.Configuration) ([]string, error) {
	changedFiles, checkErrors := findChangedFiles(config)
	if len(checkErrors) > 0 {
		return nil, errors.Join(checkErrors...)
	}

	return changedFiles, nil
}

// Iterates through the doc file line by line considering them as a states of an embedding.
// Such way, transits from the state to the next possible one until it reaches the end of a file.
// By the transition process, fills the parsing.Context accordingly, so it is ready to retrieve
// the result.
//
// Returns a parsing.Context and an error if any occurs.
func (p Processor) fillEmbeddingContext() (parsing.Context, error) {
	context, err := parsing.NewContext(p.DocFilePath)
	if err != nil {
		return context, err
	}
	absDocPath, _ := filepath.Abs(p.DocFilePath)
	errorStr := "failed to embed code fragment into doc file `file://%s:%d`: %s"

	var currentState parsing.State
	currentState = parsing.Start
	finishState := parsing.Finish

	for currentState != finishState {
		accepted, newState, err := p.moveToNextState(&currentState, &context)
		if err != nil {
			return context, fmt.Errorf(errorStr, absDocPath, errorLine(context, err), err)
		}
		if !accepted {
			err = unacceptedTransitionError(context)
			currentState = &parsing.RegularLineState{}
			if context.EmbeddingInstruction != nil {
				context.ResolveUnacceptedEmbedding()
			}

			return context, fmt.Errorf(errorStr, absDocPath, errorLine(context, err), err)
		}
		currentState = *newState
	}

	return context, nil
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
	for _, nextState := range p.TransitionsMap[*state] {
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

// findChangedFiles returns documentation files that are not up-to-date with their code files.
//
// config — a configuration for embedding.
func findChangedFiles(config configuration.Configuration) ([]string, []error) {
	requiredDocPaths, err := requiredDocs(config)
	if err != nil {
		return nil, []error{err}
	}
	var changedFiles []string
	var checkErrors []error
	for _, doc := range requiredDocPaths {
		upToDate, err := newProcessor(
			doc, config, parsing.Transitions, requiredDocPaths,
		).isUpToDate()
		if err != nil {
			checkErrors = append(checkErrors, err)

			continue
		}
		if !upToDate {
			changedFiles = append(changedFiles, doc)
		}
	}

	return changedFiles, checkErrors
}

// requiredDocs returns documentation files matched by includes minus excludes.
func requiredDocs(config configuration.Configuration) ([]string, error) {
	documentationRoot := config.DocumentationRoot
	includedPatterns := config.DocIncludes
	excludedPatterns := config.DocExcludes

	includedDocs, err := getFilesByPatterns(documentationRoot, includedPatterns)
	if err != nil {
		return nil, err
	}

	excludedDocs, err := getFilesByPatterns(documentationRoot, excludedPatterns)
	if err != nil {
		return nil, err
	}
	if len(excludedDocs) == 0 {
		slog.Info(fmt.Sprintf(
			"Found %d documentation file(s) from `%s` matching include pattern(s) %s.",
			len(includedDocs), logging.FileReference(documentationRoot),
			patternsLabel(includedPatterns),
		))

		return includedDocs, nil
	}

	result := removeElements(includedDocs, excludedDocs)
	slog.Info(fmt.Sprintf(
		"Found %d documentation file(s) from `%s` matching include pattern(s) %s "+
			"and exclude pattern(s) %s.",
		len(result), logging.FileReference(documentationRoot), patternsLabel(includedPatterns),
		patternsLabel(excludedPatterns),
	))

	return result, nil
}

// patternsLabel formats glob patterns for human-readable log messages.
func patternsLabel(patterns []string) string {
	if len(patterns) == 0 {
		return "nothing"
	}

	return "`" + strings.Join(patterns, "`, `") + "`"
}

// getFilesByPatterns expands documentation glob patterns relative to the given root.
func getFilesByPatterns(root string, patterns []string) ([]string, error) {
	var result []string
	for _, pattern := range patterns {
		globString := filepath.Join(root, filepath.FromSlash(pattern))
		matches, err := doublestar.FilepathGlob(globString)
		if err != nil {
			return nil, err
		}
		for _, match := range matches {
			result = append(result, filepath.ToSlash(match))
		}
	}

	return result, nil
}

// Returns the elements of the first array excluding those present in the second array.
func removeElements(first, second []string) []string {
	secondMap := make(map[string]struct{})
	for _, value := range second {
		secondMap[value] = struct{}{}
	}

	var result []string
	for _, value := range first {
		if _, exists := secondMap[value]; !exists {
			result = append(result, value)
		}
	}

	return result
}
