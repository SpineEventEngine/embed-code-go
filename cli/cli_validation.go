/*
 * Copyright 2026, TeamDev. All rights reserved.
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

// IllegalFolderNameChars the string with chars that are not allowed for the folder name.
const IllegalFolderNameChars = ` *?:"<>|`

// IsUsingConfigFile reports whether user configs are set with file.
func IsUsingConfigFile(config Config) bool {
	return isNotEmpty(config.ConfigPath)
}

// ValidateConfig checks the validity of provided config and returns an error if any of the
// validation rules are broken. If everything is ok, returns nil.
//
// config — a struct with user-provided args.
func ValidateConfig(config Config) error {
	err := validateMode(config.Mode)
	if err != nil {
		return err
	}

	return validateConfig(config)
}

// ValidateConfigFile performs several checks to ensure that the necessary configuration values are
// present. Also checks for the existence of the config file.
//
// userConfig — is a config given from CLI.
//
// Returns an error with a validation message. If everything is ok, returns nil.
func ValidateConfigFile(userConfig Config) error {
	// Configs should be read from file, verifying if they are not set already.
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

// Validates if mode is set to check, embed, or analyze.
func validateMode(mode string) error {
	isModeSet := isNotEmpty(mode)
	if !isModeSet {
		return errors.New("mode must be set")
	}

	validModes := []string{ModeEmbed, ModeAnalyze, ModeCheck}
	isValidMode := slices.Contains(validModes, mode)

	if !isValidMode {
		return fmt.Errorf("invalid value for mode. it must be one of — `%s`, `%s` or `%s`",
			ModeEmbed, ModeCheck, ModeAnalyze)
	}

	return nil
}

// Validates if config is set correctly and does not have mutually exclusive params set.
func validateConfig(config Config) error {
	if len(config.Embeddings) > 0 {
		return validateEmbeddingConfigs(config)
	}

	isCodePathsSet, err := validatePaths(config.BaseCodePaths)
	if err != nil {
		return err
	}
	err = findCodeSourceDuplications(config.BaseCodePaths)
	if err != nil {
		return err
	}
	isDocsPathSet, err := validatePathSet(config.BaseDocsPath)
	if err != nil {
		return err
	}
	_, err = validatePathSet(config.FragmentsPath)
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

// Validates the multi-target embedding configuration.
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

// Validates one embedding entry.
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
	if err = findCodeSourceDuplications(embedding.CodePaths); err != nil {
		return fmt.Errorf("embedding `%s`: %w", embedding.Name, err)
	}

	isDocsPathSet, err := validatePathSet(embedding.DocsPath)
	if err != nil {
		return fmt.Errorf("embedding `%s`: %w", embedding.Name, err)
	}
	_, err = validatePathSet(embedding.FragmentsPath)
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

// Reports whether at least one of optional configs is set — code-includes, doc-includes, separator
// or fragments-path.
func validateOptionalParamsSet(config Config) bool {
	isCodeIncludesSet := len(config.CodeIncludes) > 0
	isDocIncludesSet := len(config.DocIncludes) > 0
	isDocExcludesSet := len(config.DocExcludes) > 0
	isSeparatorSet := isNotEmpty(config.Separator)
	isFragmentPathSet := isNotEmpty(config.FragmentsPath)

	return isCodeIncludesSet || isDocIncludesSet || isFragmentPathSet ||
		isSeparatorSet || isDocExcludesSet
}

// Reports whether path is set or not. If it is set, checks if such path exists in a file system.
func validatePathSet(path string) (bool, error) {
	isPathSet := isNotEmpty(path)
	if isPathSet {
		exists, err := files.IsDirExist(path)
		if err != nil {
			// Since the path is set, returning true even we have an error.
			return true, err
		}
		if !exists {
			return true, fmt.Errorf("the given path `%s` does not exist", path)
		}

		return true, nil
	}

	return false, nil
}

// Reports whether all paths are valid.
//
// If paths are provided, checks whether each path exists in the file system.
//
// Returns an error if any path name is not a valid folder name.
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

// findCodeSourceDuplications checks the provided code sources for duplicate names and paths.
//
// It logs a warning for duplicate names and returns an error for duplicate paths.
func findCodeSourceDuplications(paths _type.NamedPathList) error {
	nameDuplicates := make(map[string][]string)
	pathCount := make(map[string]int)

	for _, p := range paths {
		name := p.Name
		if isEmpty(name) {
			name = "(unnamed)"
		}
		nameDuplicates[name] = append(nameDuplicates[name], p.Path)
		pathCount[p.Path]++
	}

	verifyDuplicateNames(nameDuplicates)
	return verifyDuplicatePaths(pathCount)
}

// verifyDuplicateNames logs a warning if multiple code sources share the same name.
func verifyDuplicateNames(nameDuplicates map[string][]string) {
	var warnLines []string
	for name, ps := range nameDuplicates {
		if len(ps) > 1 {
			warnLines = append(warnLines, "- "+name)
			for _, path := range ps {
				warnLines = append(warnLines, "  - "+path)
			}
		}
	}

	if len(warnLines) > 0 {
		slog.Warn(
			"Duplicate code source names detected, it may lead to " +
				"overwriting code fragments with the same relative path:\n" +
				strings.Join(warnLines, "\n"),
		)
	}
}

// verifyDuplicatePaths returns an error if multiple code sources use the same path.
func verifyDuplicatePaths(pathCount map[string]int) error {
	var errLines []string
	for path, count := range pathCount {
		if count > 1 {
			errLines = append(errLines, "- "+path)
		}
	}

	if len(errLines) > 0 {
		return fmt.Errorf(
			"duplicate code source paths detected:\n%s",
			strings.Join(errLines, "\n"),
		)
	}
	return nil
}

// Reports whether the given string is not empty.
func isNotEmpty(s string) bool {
	return !isEmpty(s)
}

// Reports whether the given string is empty.
func isEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}
