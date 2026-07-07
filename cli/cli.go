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
	_type "embed-code/embed-code-go/type"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"embed-code/embed-code-go/configuration"
	"embed-code/embed-code-go/embedding"
	"embed-code/embed-code-go/logging"

	"gopkg.in/yaml.v3"
)

// Config contains user-specified embed-code settings.
type Config struct {
	// BaseCodePaths contains directories with source code files.
	BaseCodePaths _type.NamedPathList `yaml:"code-path"`

	// BaseDocsPath is the root directory containing documentation files.
	BaseDocsPath string `yaml:"docs-path"`

	// DocIncludes contains patterns selecting documentation files to process.
	// Patterns are resolved relative to the documentation root.
	// For example, "docs/**/*.md,guides/*.html". The default is "**/*.md,**/*.html".
	DocIncludes _type.StringList `yaml:"doc-includes"`

	// DocExcludes contains patterns selecting documentation files to skip.
	DocExcludes _type.StringList `yaml:"doc-excludes"`

	// Separator is inserted between multiple partitions of one fragment.
	// The default is "...".
	Separator string `yaml:"separator"`

	// Embeddings contains independent embedding target configurations.
	Embeddings []EmbeddingConfig `yaml:"embeddings"`

	// Info reports whether info-level logs should be shown.
	Info bool `yaml:"info"`

	// Stacktrace reports whether panic stack traces should be shown.
	Stacktrace bool `yaml:"stacktrace"`

	// ConfigPath is the path to the YAML configuration file.
	ConfigPath string

	// Mode selects check or embed execution.
	Mode string
}

// EmbeddingConfig contains a complete configuration for one embedding target.
type EmbeddingConfig struct {
	// Name identifies the embedding target.
	Name string `yaml:"name"`

	// CodePaths contains directories with source code files.
	CodePaths _type.NamedPathList `yaml:"code-path"`

	// DocsPath is the root directory containing documentation files.
	DocsPath string `yaml:"docs-path"`

	// DocIncludes contains patterns selecting documentation files to process.
	DocIncludes _type.StringList `yaml:"doc-includes"`

	// DocExcludes contains patterns selecting documentation files to skip.
	DocExcludes _type.StringList `yaml:"doc-excludes"`

	// Separator is inserted between multiple partitions of one fragment.
	Separator string `yaml:"separator"`
}

// EmbedCodeSamplesResult contains the result of an EmbedCodeSamples operation.
type EmbedCodeSamplesResult struct {
	// EmbedAllResult contains the underlying embedding result.
	embedding.EmbedAllResult
}

const (
	// ModeCheck checks whether documentation snippets are up-to-date.
	ModeCheck = "check"

	// ModeEmbed rewrites documentation snippets from source code.
	ModeEmbed = "embed"
)

// CheckCodeSamples returns documentation files that are not up-to-date with code files.
//
// Parameters:
// config - provides embedding configuration.
//
// Returns:
// []string - stale documentation file paths.
// error - when selected documents fail to process.
func CheckCodeSamples(config configuration.Configuration) ([]string, error) {
	return embedding.CheckUpToDate(config)
}

// EmbedCodeSamples embeds code fragments in documentation files.
//
// Parameters:
// config - provides embedding configuration.
//
// Returns:
// EmbedCodeSamplesResult - embedding result.
// error - when selected documents fail to process or write.
func EmbedCodeSamples(config configuration.Configuration) (EmbedCodeSamplesResult, error) {
	embeddingResult, err := embedding.EmbedAll(config)
	if err != nil {
		return EmbedCodeSamplesResult{}, err
	}

	return EmbedCodeSamplesResult{
		embeddingResult,
	}, nil
}

// ReadArgs reads user-specified command-line args.
//
// Returns command-line configuration.
func ReadArgs() Config {
	codePath := flag.String("code-path", "", "a path to a root directory with code files")
	docsPath := flag.String("docs-path", "", "a path to a root directory with docs files")
	docIncludes := flag.String("doc-includes", "",
		"a comma-separated string of glob patterns for docs files to include")
	docExcludes := flag.String("doc-excludes", "",
		"a comma-separated string of glob patterns for docs files to exclude")
	separator := flag.String("separator", "",
		"a string that's inserted between multiple partitions of a single fragment")
	configPath := flag.String("config-path", "", "a path to a yaml configuration file")
	mode := flag.String("mode", "",
		"a mode of embed-code execution, which can be 'check' or 'embed'")
	info := flag.Bool("info", false,
		"an info-level logging setter that enables info logs when set to 'true'")
	stacktrace := flag.Bool("stacktrace", false,
		"a stack trace setter that enables stack traces in error logs when set to 'true'")

	flag.Parse()

	return Config{
		BaseCodePaths: _type.NamedPathList{_type.NamedPath{Path: *codePath}},
		BaseDocsPath:  *docsPath,
		DocIncludes:   parseListArgument(*docIncludes),
		DocExcludes:   parseListArgument(*docExcludes),
		Separator:     *separator,
		ConfigPath:    *configPath,
		Mode:          *mode,
		Info:          *info,
		Stacktrace:    *stacktrace,
	}
}

// FillArgsFromConfigFile fills args with values read from the configured YAML file.
//
// Parameters:
// args - provides the config file path and command-line defaults.
//
// Returns:
// Config - merged configuration.
// error - when the YAML file cannot be read or decoded.
func FillArgsFromConfigFile(args Config) (Config, error) {
	configFields, err := readConfigFields(args.ConfigPath)
	if err != nil {
		return args, err
	}
	slog.Info(fmt.Sprintf(
		"Loaded config file `%s`. Found %d embedding setup(s).",
		logging.FileReference(args.ConfigPath), len(configFields.Embeddings),
	))
	args.BaseDocsPath = configFields.BaseDocsPath
	args.BaseCodePaths = configFields.BaseCodePaths

	if len(configFields.Embeddings) > 0 {
		args.Embeddings = configFields.Embeddings
	}
	if len(configFields.DocIncludes) > 0 {
		args.DocIncludes = configFields.DocIncludes
	}
	if len(configFields.DocExcludes) > 0 {
		args.DocExcludes = configFields.DocExcludes
	}
	if isNotEmpty(configFields.Separator) {
		args.Separator = configFields.Separator
	}
	args.Info = configFields.Info
	args.Stacktrace = configFields.Stacktrace

	return args, nil
}

// BuildEmbedCodeConfiguration builds normalized embedding configurations from user args.
//
// Parameters:
// userArgs - provides command-line and YAML configuration values.
//
// Returns normalized embedding configurations.
func BuildEmbedCodeConfiguration(userArgs Config) []configuration.Configuration {
	embedCodeConfigs := make([]configuration.Configuration, 0)

	if len(userArgs.Embeddings) > 0 {
		for _, embedding := range userArgs.Embeddings {
			slog.Info(fmt.Sprintf(
				"Processing the `%s` embedding setup. Documentation folder: `%s`. %s.",
				embedding.Name, logging.FileReference(embedding.DocsPath),
				sourceFoldersLabel(embedding.CodePaths),
			))
			embedCodeConfigs = append(embedCodeConfigs, configFromEmbedding(embedding))
		}

		return embedCodeConfigs
	}

	embedCodeConfig := configWithOptionalParams(userArgs)
	embedCodeConfig.CodeRoots = userArgs.BaseCodePaths
	embedCodeConfig.DocumentationRoot = userArgs.BaseDocsPath
	slog.Info(fmt.Sprintf(
		"Preparing command line settings. Documentation folder: `%s`. %s.",
		logging.FileReference(userArgs.BaseDocsPath), sourceFoldersLabel(userArgs.BaseCodePaths),
	))

	embedCodeConfigs = append(embedCodeConfigs, embedCodeConfig)

	return embedCodeConfigs
}

// configFromEmbedding creates a new Configuration from one complete embedding config.
func configFromEmbedding(embedding EmbeddingConfig) configuration.Configuration {
	embedCodeConfig := configuration.NewConfiguration()
	embedCodeConfig.Name = embedding.Name
	embedCodeConfig.CodeRoots = embedding.CodePaths
	embedCodeConfig.DocumentationRoot = embedding.DocsPath

	if len(embedding.DocIncludes) > 0 {
		embedCodeConfig.DocIncludes = embedding.DocIncludes
	}
	if len(embedding.DocExcludes) > 0 {
		embedCodeConfig.DocExcludes = embedding.DocExcludes
	}
	if isNotEmpty(embedding.Separator) {
		embedCodeConfig.Separator = embedding.Separator
	}

	return embedCodeConfig
}

// configWithOptionalParams creates a new Configuration with optional properties from user args.
func configWithOptionalParams(userArgs Config) configuration.Configuration {
	embedCodeConfig := configuration.NewConfiguration()

	if len(userArgs.DocIncludes) > 0 {
		embedCodeConfig.DocIncludes = userArgs.DocIncludes
	}
	if len(userArgs.DocExcludes) > 0 {
		embedCodeConfig.DocExcludes = userArgs.DocExcludes
	}
	if isNotEmpty(userArgs.Separator) {
		embedCodeConfig.Separator = userArgs.Separator
	}

	return embedCodeConfig
}

// sourceFoldersLabel formats source folders for human-readable log messages.
func sourceFoldersLabel(paths _type.NamedPathList) string {
	if len(paths) == 0 {
		return "No source code folders configured"
	}

	var labels []string
	for _, path := range paths {
		label := fmt.Sprintf("`%s`", logging.FileReference(path.Path))
		if strings.TrimSpace(path.Name) != "" {
			label = fmt.Sprintf("`%s` as `%s`", logging.FileReference(path.Path), path.Name)
		}
		labels = append(labels, label)
	}

	if len(labels) == 1 {
		return "Source code folder: " + labels[0]
	}

	return "Source code folders: " + strings.Join(labels, ", ")
}

// parseListArgument splits a comma-separated command-line argument.
//
// Parameters:
// listArgument - provides a comma-separated string.
//
// Returns parsed non-empty values.
func parseListArgument(listArgument string) []string {
	splitArgs := strings.Split(listArgument, ",")
	parsedArgs := make([]string, 0)
	for _, v := range splitArgs {
		if trimmed := strings.TrimSpace(v); trimmed != "" {
			parsedArgs = append(parsedArgs, trimmed)
		}
	}

	return parsedArgs
}

// readConfigFields reads and parses a YAML configuration file.
//
// Parameters:
// configFilePath - provides the path to a YAML configuration file.
//
// Returns:
// Config - parsed configuration fields.
// error - when the file cannot be read or decoded.
func readConfigFields(configFilePath string) (Config, error) {
	content, err := os.ReadFile(configFilePath)
	if err != nil {
		return Config{}, err
	}

	configFields := Config{}
	err = yaml.Unmarshal(content, &configFields)
	if err != nil {
		return Config{}, err
	}

	return configFields, nil
}
