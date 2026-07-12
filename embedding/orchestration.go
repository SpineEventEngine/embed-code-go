/*
 * Copyright 2026, TeamDev. All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 * Redistribution and use in source and/or binary forms, with or without
 * modification, must retain the above copyright notice and the following
 * disclaimer.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
 * "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
 * LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
 * A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
 * OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
 * SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
 * LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
 * DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
 * THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 * OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

package embedding

import (
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"

	"embed-code/embed-code-go/configuration"
	"embed-code/embed-code-go/embedding/parsing"
	"embed-code/embed-code-go/logging"

	"github.com/bmatcuk/doublestar/v4"
)

// EmbedAllResult contains the result of an EmbedAll operation.
type EmbedAllResult struct {
	// TotalEmbeddings is the total number of embeddings found in the target documentation files.
	TotalEmbeddings int

	// UpdatedTargetFiles contains documentation files changed by embedding.
	UpdatedTargetFiles []string
}

// processorHandler applies one processing mode to a discovered documentation file.
type processorHandler func(docFilePath string, processor Processor) error

// EmbedAll embeds code fragments into all documentation files selected by config.
//
// It resolves documentation files from configured patterns, creates a Processor
// for each file, and embeds code fragments into those documents.
//
// Parameters:
// config - provides embedding configuration.
//
// Returns:
// EmbedAllResult - embedding result.
// error - when selected documents fail to process.
func EmbedAll(config configuration.Configuration) (EmbedAllResult, error) {
	totalEmbeddings := 0
	var updatedTargetFiles []string
	requiredDocPaths, embeddingErrors := processRequiredDocs(config, func(
		_ string,
		processor Processor,
	) error {
		context, err := processor.Embed()
		if err != nil {
			return err
		}
		totalEmbeddings += context.EmbeddingsCount()
		if context.IsContentChanged() {
			updatedTargetFiles = append(updatedTargetFiles, context.MarkdownFilePath)
		}

		return nil
	})
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
		TotalEmbeddings:    totalEmbeddings,
		UpdatedTargetFiles: updatedTargetFiles,
	}, nil
}

// configNameLabel formats a configuration name for summary log messages.
//
// A non-empty label starts with a space so callers can append it directly.
func configNameLabel(config configuration.Configuration) string {
	if config.Name == "" {
		return ""
	}

	return fmt.Sprintf(" for `%s` embedding setup", config.Name)
}

// CheckUpToDate returns documentation files that are not up-to-date with code files.
//
// Parameters:
// config - provides embedding configuration.
//
// Returns:
// []string - stale documentation file paths.
// error - when selected documents fail to process.
func CheckUpToDate(config configuration.Configuration) ([]string, error) {
	changedFiles, checkErrors := findChangedFiles(config)
	if len(checkErrors) > 0 {
		return nil, errors.Join(checkErrors...)
	}

	return changedFiles, nil
}

// findChangedFiles returns documentation files that are not up-to-date with their code files.
func findChangedFiles(config configuration.Configuration) ([]string, []error) {
	var changedFiles []string
	_, checkErrors := processRequiredDocs(config, func(
		docFilePath string,
		processor Processor,
	) error {
		upToDate, err := processor.isUpToDate()
		if err != nil {
			return err
		}
		if !upToDate {
			changedFiles = append(changedFiles, docFilePath)
		}

		return nil
	})

	return changedFiles, checkErrors
}

// processRequiredDocs applies a processing handler to every documentation file in config.
func processRequiredDocs(
	config configuration.Configuration,
	handle processorHandler,
) ([]string, []error) {
	requiredDocPaths, err := requiredDocs(config)
	if err != nil {
		return nil, []error{err}
	}

	var processingErrors []error
	resolver := newDefaultResolver()
	for _, doc := range requiredDocPaths {
		processor := newProcessor(doc, config, parsing.Transitions, requiredDocPaths, resolver)
		if err := handle(doc, processor); err != nil {
			processingErrors = append(processingErrors, err)
		}
	}

	return requiredDocPaths, processingErrors
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

// removeElements returns values from first that are not present in second.
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
