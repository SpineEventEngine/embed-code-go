/*
 * Copyright 2024, TeamDev. All rights reserved.
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
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"
)

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
// path — a path to a yaml configuration file.
//
// Returns an error with a validation message. If everything is ok, returns nil.
func ValidateConfigFile(path string) error {
	exists, err := isFileExist(path)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("config file is not exist")
	}

	configFields := readConfigFields(path)
	if isEmpty(configFields.CodePath) || isEmpty(configFields.DocsPath) {
		return errors.New("config must include both code-path and docs-path fields")
	}

	return nil
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
		return fmt.Errorf("invalid value for mode. it must be one of — %s, %s or %s",
			ModeEmbed, ModeCheck, ModeAnalyze)
	}

	return nil
}

// Validates if config is set correctly and does not have mutually exclusive params set.
func validateConfig(config Config) error {
	isConfigSet := isNotEmpty(config.ConfigPath)
	isCodePathSet, err := validatePathIfSet(config.CodePath)
	if err != nil {
		return err
	}
	isDocsPathSet, err := validatePathIfSet(config.DocsPath)
	if err != nil {
		return err
	}
	isOptionalParamsSet, err := validateIfOptionalParamsAreSet(config)
	if err != nil {
		return err
	}

	err = validateParamsCombinations(isConfigSet, isCodePathSet, isDocsPathSet, isOptionalParamsSet)
	if err != nil {
		return err
	}

	return nil
}

// Validates if config path and arguments are not conflicting.
func validateParamsCombinations(
	isConfigSet bool, isCodePathSet bool, isDocsPathSet bool, isOptionalParamsSet bool) error {
	isRootsSet := isCodePathSet && isDocsPathSet
	isOneOfRootsSet := isCodePathSet || isDocsPathSet

	if isConfigSet && (isOneOfRootsSet || isOptionalParamsSet) {
		return errors.New(
			"config path cannot be set when code-path, docs-path or optional params are set")
	}
	if isOneOfRootsSet && !isRootsSet {
		return errors.New(
			"if one of code-path and docs-path is set, the another one must be set as well")
	}
	if !(isRootsSet || isConfigSet) {
		return errors.New("embed code should be used with either config-path or both " +
			"code-path and docs-path being set")
	}

	return nil
}

// Reports whether at least one of optional configs is set — code-includes, doc-includes, separator
// or fragments-path.
func validateIfOptionalParamsAreSet(config Config) (bool, error) {
	isCodeIncludesSet := isNotEmpty(config.CodeIncludes)
	isDocIncludesSet := isNotEmpty(config.DocIncludes)
	isSeparatorSet := isNotEmpty(config.Separator)
	isFragmentPathSet, err := validatePathIfSet(config.FragmentsPath)
	if err != nil {
		return false, err
	}

	return isCodeIncludesSet || isDocIncludesSet || isFragmentPathSet || isSeparatorSet, nil
}

// Reports whether path is set or not. If it is set, checks if such path exists in a file system.
func validatePathIfSet(path string) (bool, error) {
	isPathSet := isNotEmpty(path)
	if isPathSet {
		exists, err := isDirExist(path)
		if err != nil {
			// Since the path is set, returning true even we have an error.
			return true, err
		}
		if !exists {
			return true, errors.New("config file is not exist")
		}

		return true, nil
	}

	return false, nil
}

// Reports whether the given path to a file exists in the file system.
func isFileExist(filePath string) (bool, error) {
	exists, info, err := validateIfPathExists(filePath)
	if err != nil {
		return false, err
	}
	if exists {
		if info.IsDir() {
			return false, fmt.Errorf("%s is a directory, the file was expected", filePath)
		}

		return true, nil
	}

	return false, nil
}

// Reports whether the given directory exists in the file system.
func isDirExist(path string) (bool, error) {
	exists, info, err := validateIfPathExists(path)
	if err != nil {
		return false, err
	}
	if exists {
		if info.IsDir() {
			return true, nil
		}

		return false, fmt.Errorf("%s is a file, the directory was expected", path)
	}

	return false, nil
}

// Reports whether the given path is valid and exist in the file system. Also returns a FileInfo if
// the path exists.
func validateIfPathExists(path string) (bool, os.FileInfo, error) {
	info, err := os.Stat(path)

	if err != nil {
		if os.IsNotExist(err) {
			return false, nil, fmt.Errorf("the path %s is not exist", path)
		}

		return false, nil, err
	}

	return true, info, nil
}

// Reports whether the given string is not empty.
func isNotEmpty(s string) bool {
	return !isEmpty(s)
}

// Reports whether the given string is empty.
func isEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}
