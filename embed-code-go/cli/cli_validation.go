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
	"embed-code/embed-code-go/files"
	"errors"
	"fmt"
	"slices"
	"strings"
)

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
	isCodePathSet := isNotEmpty(userConfig.CodePath)
	isDocsPathSet := isNotEmpty(userConfig.DocsPath)
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
		return fmt.Errorf("invalid value for mode. it must be one of — %s, %s or %s",
			ModeEmbed, ModeCheck, ModeAnalyze)
	}

	return nil
}

// Validates if config is set correctly and does not have mutually exclusive params set.
func validateConfig(config Config) error {
	isCodePathSet, err := validatePathSet(config.CodePath)
	if err != nil {
		return err
	}
	isDocsPathSet, err := validatePathSet(config.DocsPath)
	if err != nil {
		return err
	}
	_, err = validatePathSet(config.FragmentsPath)
	if err != nil {
		return err
	}

	isRootsSet := isCodePathSet && isDocsPathSet
	isOneOfRootsSet := isCodePathSet || isDocsPathSet

	if isOneOfRootsSet && !isRootsSet {
		return errors.New("code-path and docs-path must both be set")
	}

	return nil
}

// Reports whether at least one of optional configs is set — code-includes, doc-includes, separator
// or fragments-path.
func validateOptionalParamsSet(config Config) bool {
	isCodeIncludesSet := isNotEmpty(config.CodeIncludes)
	isDocIncludesSet := isNotEmpty(config.DocIncludes)
	isDocExcludesSet := isNotEmpty(config.DocExcludes)
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
			return true, fmt.Errorf("the given path %s is not exist", path)
		}

		return true, nil
	}

	return false, nil
}

// Reports whether the given string is not empty.
func isNotEmpty(s string) bool {
	return !isEmpty(s)
}

// Reports whether the given string is empty.
func isEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}
