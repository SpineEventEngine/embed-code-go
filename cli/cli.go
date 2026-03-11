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

package cli

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"embed-code/embed-code-go/analyzing"
	"embed-code/embed-code-go/configuration"
	"embed-code/embed-code-go/embedding"
	"embed-code/embed-code-go/fragmentation"

	"gopkg.in/yaml.v3"
)

// Config — user-specified embed-code configurations.
//
// BaseCodePath — a path to a root directory with code files.
//
// BaseDocsPath — a path to a root directory with docs files.
//
// CodeIncludes — a string with comma-separated patterns for filtering the code files
// to be considered.
// Directories are never matched by these patterns.
// For example, "**/*.java,**/*.gradle".
// The default value is "**/*.*".
//
// DocIncludes — a string with comma-separated patterns for filtering files
// in which we should look for embedding instructions.
// The patterns are resolved relatively to the `documentation_root`.
// Directories are never matched by these patterns.
// For example, "docs/**/*.md,guides/*.html".
// The default value is "**/*.md,**/*.html".
//
// FragmentsPath — a directory where fragmented code is stored. A temporary directory that should
// not be tracked in VCS. The default value is: "./build/fragments".
//
// Separator — a string that's inserted between multiple partitions of a single fragment.
// The default value is "...".
//
// ConfigPath — a path to a yaml configuration file which contains the roots.
//
// Mode — defines the mode of embed-code execution.
//
// EmbedMappings — an additional optional list of configs, which will be executed together with the
// main one. A config written here has higher priority and may overwrite the base one.
type Config struct {
	CodeIncludes  StringList     `yaml:"code-includes"`
	DocIncludes   StringList     `yaml:"doc-includes"`
	DocExcludes   StringList     `yaml:"doc-excludes"`
	FragmentsPath string         `yaml:"fragments-path"`
	Separator     string         `yaml:"separator"`
	BaseCodePath  string         `yaml:"code-path"`
	BaseDocsPath  string         `yaml:"docs-path"`
	EmbedMappings []EmbedMapping `yaml:"embed-mappings"`
	Info          bool           `yaml:"info"`
	Stacktrace    bool           `yaml:"stacktrace"`
	ConfigPath    string
	Mode          string
}

// EmbedMapping is a pair of a source code path and a destination docs path to perform an embedding.
type EmbedMapping struct {
	CodePath string `yaml:"code-path"`
	DocsPath string `yaml:"docs-path"`
}

// EmbedCodeSamplesResult is result of the EmbedCodeSamples method.
//
// WriteFragmentFilesResult the result of code fragmentation.
//
// EmbedAllResult the result of embedding code fragments in the documentation.
type EmbedCodeSamplesResult struct {
	fragmentation.WriteFragmentFilesResult
	embedding.EmbedAllResult
}

const (
	ModeCheck   = "check"
	ModeEmbed   = "embed"
	ModeAnalyze = "analyze"
)

// CheckCodeSamples checks documentation to be up-to-date with code files. Raises
// UnexpectedDiffError if not.
//
// config — a configuration for checking code samples.
func CheckCodeSamples(config configuration.Configuration) {
	fragmentation.WriteFragmentFiles(config)
	embedding.CheckUpToDate(config)
}

// EmbedCodeSamples embeds code fragments in documentation files.
//
// config — a configuration for embedding.
func EmbedCodeSamples(config configuration.Configuration) EmbedCodeSamplesResult {
	fragmentationResult := fragmentation.WriteFragmentFiles(config)
	embeddingResult := embedding.EmbedAll(config)
	embedding.CheckUpToDate(config)
	return EmbedCodeSamplesResult{
		fragmentationResult,
		embeddingResult,
	}
}

// AnalyzeCodeSamples analyzes code fragments in documentation files.
//
// config — a configuration for embedding.
func AnalyzeCodeSamples(config configuration.Configuration) {
	fragmentation.WriteFragmentFiles(config)
	analyzing.AnalyzeAll(config)
	fragmentation.CleanFragmentFiles(config)
}

// ReadArgs reads user-specified args from the command line.
//
// Returns Config struct filled with the corresponding args.
func ReadArgs() Config {
	codePath := flag.String("code-path", "", "a path to a root directory with code files")
	docsPath := flag.String("docs-path", "", "a path to a root directory with docs files")
	codeIncludes := flag.String("code-includes", "",
		"a comma-separated string of glob patterns for code files to include")
	docIncludes := flag.String("doc-includes", "",
		"a comma-separated string of glob patterns for docs files to include")
	docExcludes := flag.String("doc-excludes", "",
		"a comma-separated string of glob patterns for docs files to exclude")
	fragmentsPath := flag.String("fragments-path", "",
		"a path to a directory where fragmented code is stored")
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
		BaseCodePath:  *codePath,
		BaseDocsPath:  *docsPath,
		CodeIncludes:  parseListArgument(*codeIncludes),
		DocIncludes:   parseListArgument(*docIncludes),
		DocExcludes:   parseListArgument(*docExcludes),
		FragmentsPath: *fragmentsPath,
		Separator:     *separator,
		ConfigPath:    *configPath,
		Mode:          *mode,
		Info:          *info,
		Stacktrace:    *stacktrace,
	}
}

// FillArgsFromConfigFile fills config with the values read from config file.
//
// args — Config struct with user-provided args.
//
// Returns filled Config.
func FillArgsFromConfigFile(args Config) (Config, error) {
	configFields := readConfigFields(args.ConfigPath)
	args.BaseDocsPath = configFields.BaseDocsPath
	args.BaseCodePath = configFields.BaseCodePath

	if len(configFields.CodeIncludes) > 0 {
		args.CodeIncludes = configFields.CodeIncludes
	}
	if len(configFields.EmbedMappings) > 0 {
		args.EmbedMappings = configFields.EmbedMappings
	}
	if len(configFields.DocIncludes) > 0 {
		args.DocIncludes = configFields.DocIncludes
	}
	if len(configFields.DocExcludes) > 0 {
		args.DocExcludes = configFields.DocExcludes
	}
	if isNotEmpty(configFields.FragmentsPath) {
		args.FragmentsPath = configFields.FragmentsPath
	}
	if isNotEmpty(configFields.Separator) {
		args.Separator = configFields.Separator
	}
	args.Info = configFields.Info
	args.Stacktrace = configFields.Stacktrace

	return args, nil
}

// BuildEmbedCodeConfiguration generates and returns a configuration based on provided userArgs.
//
// userArgs — a Config with user-provided args.
func BuildEmbedCodeConfiguration(userArgs Config) []configuration.Configuration {
	embedCodeConfigs := make([]configuration.Configuration, 0)
	excludedConfigs := make([]string, 0)

	if len(userArgs.EmbedMappings) > 0 {
		for _, mapping := range userArgs.EmbedMappings {
			embedCodeConfig := configWithOptionalParams(userArgs)
			embedCodeConfig.CodeRoot = mapping.CodePath
			embedCodeConfig.DocumentationRoot = mapping.DocsPath

			// As the top config may overwrite those files, we need to exclude it from the embedding
			excludedConfigs = append(excludedConfigs, fmt.Sprintf("%s**/*.*", mapping.DocsPath))
			embedCodeConfigs = append(embedCodeConfigs, embedCodeConfig)
		}
	}

	embedCodeConfig := configWithOptionalParams(userArgs)
	embedCodeConfig.CodeRoot = userArgs.BaseCodePath
	embedCodeConfig.DocumentationRoot = userArgs.BaseDocsPath

	if len(userArgs.DocExcludes) > 0 {
		embedCodeConfig.DocExcludes = append(embedCodeConfig.DocExcludes, excludedConfigs...)
	} else {
		embedCodeConfig.DocExcludes = excludedConfigs
	}
	embedCodeConfigs = append(embedCodeConfigs, embedCodeConfig)

	return embedCodeConfigs
}

// Creates a new Configuration with the filled optional properties from the user args.
func configWithOptionalParams(userArgs Config) configuration.Configuration {
	embedCodeConfig := configuration.NewConfiguration()

	if len(userArgs.CodeIncludes) > 0 {
		embedCodeConfig.CodeIncludes = userArgs.CodeIncludes
	}
	if len(userArgs.DocIncludes) > 0 {
		embedCodeConfig.DocIncludes = userArgs.DocIncludes
	}
	if isNotEmpty(userArgs.FragmentsPath) {
		embedCodeConfig.FragmentsDir = userArgs.FragmentsPath
	}
	if isNotEmpty(userArgs.Separator) {
		embedCodeConfig.Separator = userArgs.Separator
	}

	return embedCodeConfig
}

// Returns a list of strings from given comma-separated string listArgument.
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

// Reads the file from provided configFilePath and returns a ConfigFields struct.
//
// configFilePath — a path to a yaml configuration file.
//
// Returns a filled ConfigFields struct.
func readConfigFields(configFilePath string) Config {
	content, err := os.ReadFile(configFilePath)
	if err != nil {
		panic(err)
	}

	configFields := Config{}
	err = yaml.Unmarshal(content, &configFields)
	if err != nil {
		panic(err)
	}

	return configFields
}
