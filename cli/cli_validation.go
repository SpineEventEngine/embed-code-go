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

package cli

import (
	"embed-code/embed-code-go/files"
	_type "embed-code/embed-code-go/type"
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"
	"slices"
	"strings"
)

// IllegalFolderNameChars contains characters that are not allowed in folder names.
const IllegalFolderNameChars = `/\ *?:"<>|`

// IsUsingConfigFile reports whether user configs are set with file.
//
// Parameters:
// config - provides user CLI settings.
//
// Returns true when ConfigPath is not empty.
func IsUsingConfigFile(config Config) bool {
	return isNotEmpty(config.ConfigPath)
}

// ValidateConfig checks user args and returns the first validation error.
//
// Parameters:
// config - provides user CLI or YAML settings.
//
// Returns an error when mode or path settings are invalid.
func ValidateConfig(config Config) error {
	err := validateMode(config.Mode)
	if err != nil {
		return err
	}

	return validateConfig(config)
}

// ValidateConfigFile checks that config-file mode is used correctly.
//
// Parameters:
// userConfig - provides command-line settings before YAML loading.
//
// Returns an error when config-file mode is invalid or the file is missing.
func ValidateConfigFile(userConfig Config) error {
	// Config values should be read from file, so other root or optional params
	// must not be set already.
	isCodePathSet := len(userConfig.BaseCodePaths) > 0 &&
		isNotEmpty(userConfig.BaseCodePaths[0].Path)
	isDocsPathSet := isNotEmpty(userConfig.BaseDocsPath)
	areOptionalParamsSet := validateOptionalParamsSet(userConfig)
	isOneOfRootsSet := isCodePathSet || isDocsPathSet

	if isOneOfRootsSet || areOptionalParamsSet {
		return errors.New(
			"config path cannot be set when code-path, docs-path or optional params are set")
	}

	exists, err := files.IsFileExist(userConfig.ConfigPath)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	return errors.New("expected to use config file, but it does not exist")
}

// validateMode checks if mode is set to check or embed.
func validateMode(mode string) error {
	isModeSet := isNotEmpty(mode)
	if !isModeSet {
		return errors.New("mode must be set")
	}

	validModes := []string{ModeEmbed, ModeCheck}
	isValidMode := slices.Contains(validModes, mode)

	if !isValidMode {
		return fmt.Errorf("invalid value for mode. it must be one of — `%s` or `%s`",
			ModeEmbed, ModeCheck)
	}

	return nil
}

// validateConfig checks if config is set correctly and has no mutually exclusive params.
func validateConfig(config Config) error {
	if len(config.Embeddings) > 0 {
		return validateEmbeddingConfigs(config)
	}

	isCodePathsSet, err := validatePaths(config.BaseCodePaths)
	if err != nil {
		return err
	}
	err = validateCodeSources(config.BaseCodePaths)
	if err != nil {
		return err
	}
	isDocsPathSet, err := validatePathSet(config.BaseDocsPath)
	if err != nil {
		return err
	}
	isRootsSet := isCodePathsSet && isDocsPathSet
	isOneOfRootsSet := isCodePathsSet || isDocsPathSet

	if isOneOfRootsSet && !isRootsSet {
		return errors.New("`code-path` and `docs-path` must both be set")
	}

	return nil
}

// validateEmbeddingConfigs checks the multi-target embedding configuration.
func validateEmbeddingConfigs(config Config) error {
	isCodePathsSet, err := validatePaths(config.BaseCodePaths)
	if err != nil {
		return err
	}
	isDocsPathSet, err := validatePathSet(config.BaseDocsPath)
	if err != nil {
		return err
	}
	if isCodePathsSet || isDocsPathSet {
		return errors.New("`code-path` and `docs-path` cannot be set when `embeddings` are set")
	}
	if validateOptionalParamsSet(config) {
		return errors.New("root optional embedding options cannot be set when `embeddings` are set")
	}

	for i, embedding := range config.Embeddings {
		if err = validateEmbeddingConfig(embedding, i); err != nil {
			return err
		}
	}

	if err = findEmbeddingNameDuplications(config.Embeddings); err != nil {
		return err
	}
	verifyDuplicateEmbeddingDocsPaths(config.Embeddings)

	return nil
}

// validateEmbeddingConfig checks one embedding entry.
func validateEmbeddingConfig(embedding EmbeddingConfig, index int) error {
	if isEmpty(embedding.Name) {
		return fmt.Errorf("embedding #%d: `name` must be set", index+1)
	}
	if strings.ContainsAny(embedding.Name, IllegalFolderNameChars) {
		return fmt.Errorf("embedding `%s`: `name` `%s` is not valid, "+
			"those characters are not allowed `%s`",
			embedding.Name, embedding.Name, IllegalFolderNameChars)
	}

	isCodePathsSet, err := validatePaths(embedding.CodePaths)
	if err != nil {
		return fmt.Errorf("embedding `%s`: %w", embedding.Name, err)
	}
	if err = validateCodeSources(embedding.CodePaths); err != nil {
		return fmt.Errorf("embedding `%s`: %w", embedding.Name, err)
	}

	isDocsPathSet, err := validatePathSet(embedding.DocsPath)
	if err != nil {
		return fmt.Errorf("embedding `%s`: %w", embedding.Name, err)
	}
	isRootsSet := isCodePathsSet && isDocsPathSet
	if !isRootsSet {
		return fmt.Errorf("embedding `%s`: `code-path` and `docs-path` must both be set",
			embedding.Name)
	}

	return nil
}

// findEmbeddingNameDuplications returns an error if multiple embeddings use the same name.
func findEmbeddingNameDuplications(embeddings []EmbeddingConfig) error {
	nameCount := make(map[string]int)
	for _, embedding := range embeddings {
		nameCount[embedding.Name]++
	}

	var errLines []string
	for name, count := range nameCount {
		if count > 1 {
			errLines = append(errLines, "- "+name)
		}
	}

	if len(errLines) > 0 {
		slices.Sort(errLines)

		return fmt.Errorf(
			"duplicate embedding names detected:\n%s",
			strings.Join(errLines, "\n"),
		)
	}

	return nil
}

// verifyDuplicateEmbeddingDocsPaths logs a warning if multiple embeddings use the same docs path.
func verifyDuplicateEmbeddingDocsPaths(embeddings []EmbeddingConfig) {
	docsPathEmbeddings := make(map[string][]string)
	for _, embedding := range embeddings {
		docsPath := filepath.Clean(embedding.DocsPath)
		docsPathEmbeddings[docsPath] = append(
			docsPathEmbeddings[docsPath],
			embedding.Name,
		)
	}

	var warnLines []string
	for docsPath, names := range docsPathEmbeddings {
		if len(names) > 1 {
			slices.Sort(names)
			warnLines = append(warnLines, fmt.Sprintf("- `%s`: %s", docsPath, strings.Join(names, ", ")))
		}
	}

	if len(warnLines) > 0 {
		slices.Sort(warnLines)
		slog.Warn(
			"Multiple `embeddings` use the same `docs-path`. " +
				"Make sure they are intended to process the same documentation root:\n" +
				strings.Join(warnLines, "\n"),
		)
	}
}

// validateOptionalParamsSet reports whether at least one optional config is set.
func validateOptionalParamsSet(config Config) bool {
	isDocIncludesSet := len(config.DocIncludes) > 0
	isDocExcludesSet := len(config.DocExcludes) > 0
	isSeparatorSet := isNotEmpty(config.Separator)

	return isDocIncludesSet || isSeparatorSet || isDocExcludesSet
}

// validatePathSet reports whether path is set and checks if it exists.
func validatePathSet(path string) (bool, error) {
	isPathSet := isNotEmpty(path)
	if isPathSet {
		exists, err := files.IsDirExist(path)
		if err != nil {
			// Since the path is set, return true even when the path check fails.
			return true, err
		}
		if !exists {
			return true, fmt.Errorf("the given path `%s` does not exist", path)
		}

		return true, nil
	}

	return false, nil
}

// validatePaths reports whether all paths are set and valid.
//
// It checks whether each provided path exists in the file system.
//
// Parameters:
// paths - provides source paths to validate.
//
// Returns:
// bool - whether all paths are set.
// error - when any path does not exist or any path name is invalid.
func validatePaths(paths _type.NamedPathList) (bool, error) {
	allPathsSet := true
	if len(paths) == 0 {
		return false, nil
	}
	for _, path := range paths {
		isPathSet, err := validatePathSet(path.Path)
		if err != nil {
			return true, fmt.Errorf("the given path `%s` does not exist", path)
		}
		if strings.ContainsAny(path.Name, IllegalFolderNameChars) {
			return true, fmt.Errorf("the given code path name `%s` "+
				"is not a valid name for the folder, those characters are not allowed `%s`",
				path.Name, IllegalFolderNameChars)
		}
		if !isPathSet {
			allPathsSet = false
		}
	}

	return allPathsSet, nil
}

// validateCodeSources checks that code sources can be resolved unambiguously.
func validateCodeSources(paths _type.NamedPathList) error {
	nameCount := make(map[string]int)
	pathCount := make(map[string]int)
	pathNames := make(map[string][]string)
	unnamedCount := 0
	hasNamed := false

	for _, pathEntry := range paths {
		if isEmpty(pathEntry.Path) {
			continue
		}
		if isEmpty(pathEntry.Name) {
			unnamedCount++
		} else {
			hasNamed = true
			nameCount[pathEntry.Name]++
		}
		pathCount[pathEntry.Path]++
		pathNames[pathEntry.Path] = append(pathNames[pathEntry.Path], pathEntry.Name)
	}

	if err := verifyCodeSourceNames(nameCount); err != nil {
		return err
	}
	if unnamedCount > 1 {
		return errors.New("only one unnamed source code path is allowed")
	}
	if hasNamed && unnamedCount > 0 {
		return errors.New("named and unnamed source code paths cannot be mixed")
	}

	warnDuplicatePaths(pathCount, pathNames)

	return nil
}

// verifyCodeSourceNames returns an error if multiple code sources share the same name.
func verifyCodeSourceNames(nameCount map[string]int) error {
	var errLines []string
	for name, count := range nameCount {
		if count > 1 {
			errLines = append(errLines, "- "+name)
		}
	}

	if len(errLines) > 0 {
		slices.Sort(errLines)

		return fmt.Errorf(
			"duplicate source code path names detected:\n%s",
			strings.Join(errLines, "\n"),
		)
	}

	return nil
}

// warnDuplicatePaths logs a warning if multiple code sources use the same path.
func warnDuplicatePaths(pathCount map[string]int, pathNames map[string][]string) {
	var warnLines []string
	for path, count := range pathCount {
		if count > 1 {
			names := pathNames[path]
			slices.Sort(names)
			warnLines = append(warnLines, fmt.Sprintf("- %s: %s", path, strings.Join(names, ", ")))
		}
	}

	if len(warnLines) > 0 {
		slices.Sort(warnLines)
		slog.Warn(
			"Duplicate source code paths detected:\n" + strings.Join(warnLines, "\n"),
		)
	}
}

// isNotEmpty reports whether the given string is not empty.
func isNotEmpty(s string) bool {
	return !isEmpty(s)
}

// isEmpty reports whether the given string is empty.
func isEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}
