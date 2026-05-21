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
	_type "embed-code/embed-code-go/type"
	"flag"
	"os"
	"strings"

	"embed-code/embed-code-go/configuration"
	"embed-code/embed-code-go/embedding"

	"gopkg.in/yaml.v3"
)

// Config — user-specified embed-code configurations.
//
// BaseCodePaths — a NamedPathList to directories with code files.
//
// BaseDocsPath — a path to a root directory with docs files.
//
// DocIncludes — a StringList with patterns for filtering files
// in which we should look for embedding instructions.
// The patterns are resolved relatively to the `documentation_root`.
// Directories are never matched by these patterns.
// For example, "docs/**/*.md,guides/*.html".
// The default value is "**/*.md,**/*.html".
//
// DocExcludes - a StringList with patterns for filtering documentation files
// which should be excluded from the embedding process.
//
// Separator — a string that's inserted between multiple partitions of a single fragment.
// The default value is "...".
//
// Embeddings — independent configurations for embedding multiple documentation targets.
//
// Info - specifies whether info-level logs should be shown.
//
// Stacktrace - specifies whether error stack traces should be shown.
//
// ConfigPath — a path to a yaml configuration file which contains roots or embeddings.
//
// Mode — defines the mode of embed-code execution.
type Config struct {
	BaseCodePaths _type.NamedPathList `yaml:"code-path"`
	BaseDocsPath  string              `yaml:"docs-path"`
	DocIncludes   _type.StringList    `yaml:"doc-includes"`
	DocExcludes   _type.StringList    `yaml:"doc-excludes"`
	Separator     string              `yaml:"separator"`
	Embeddings    []EmbeddingConfig   `yaml:"embeddings"`
	Info          bool                `yaml:"info"`
	Stacktrace    bool                `yaml:"stacktrace"`
	ConfigPath    string
	Mode          string
}

// EmbeddingConfig contains a complete configuration for one embedding target.
type EmbeddingConfig struct {
	Name        string              `yaml:"name"`
	CodePaths   _type.NamedPathList `yaml:"code-path"`
	DocsPath    string              `yaml:"docs-path"`
	DocIncludes _type.StringList    `yaml:"doc-includes"`
	DocExcludes _type.StringList    `yaml:"doc-excludes"`
	Separator   string              `yaml:"separator"`
}

// EmbedCodeSamplesResult is result of the EmbedCodeSamples method.
//
// EmbedAllResult the result of embedding code fragments in the documentation.
type EmbedCodeSamplesResult struct {
	embedding.EmbedAllResult
}

const (
	ModeCheck = "check"
	ModeEmbed = "embed"
)

// CheckCodeSamples returns documentation files that are not up-to-date with code files.
//
// config — a configuration for checking code samples.
func CheckCodeSamples(config configuration.Configuration) []string {
	return embedding.CheckUpToDate(config)
}

// EmbedCodeSamples embeds code fragments in documentation files.
//
// config — a configuration for embedding.
func EmbedCodeSamples(config configuration.Configuration) EmbedCodeSamplesResult {
	embeddingResult := embedding.EmbedAll(config)
	return EmbedCodeSamplesResult{
		embeddingResult,
	}
}

// ReadArgs reads user-specified args from the command line.
//
// Returns Config struct filled with the corresponding args.
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

// FillArgsFromConfigFile fills config with the values read from config file.
//
// args — Config struct with user-provided args.
//
// Returns filled Config.
func FillArgsFromConfigFile(args Config) (Config, error) {
	configFields := readConfigFields(args.ConfigPath)
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

// BuildEmbedCodeConfiguration generates and returns a configuration based on provided userArgs.
//
// userArgs — a Config with user-provided args.
func BuildEmbedCodeConfiguration(userArgs Config) []configuration.Configuration {
	embedCodeConfigs := make([]configuration.Configuration, 0)

	if len(userArgs.Embeddings) > 0 {
		for _, embedding := range userArgs.Embeddings {
			embedCodeConfigs = append(embedCodeConfigs, configFromEmbedding(embedding))
		}
		return embedCodeConfigs
	}

	embedCodeConfig := configWithOptionalParams(userArgs)
	embedCodeConfig.CodeRoots = userArgs.BaseCodePaths
	embedCodeConfig.DocumentationRoot = userArgs.BaseDocsPath

	embedCodeConfigs = append(embedCodeConfigs, embedCodeConfig)

	return embedCodeConfigs
}

// Creates a new Configuration from one complete embedding config.
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

// Creates a new Configuration with the filled optional properties from the user args.
func configWithOptionalParams(userArgs Config) configuration.Configuration {
	embedCodeConfig := configuration.NewConfiguration()

	if len(userArgs.DocIncludes) > 0 {
		embedCodeConfig.DocIncludes = userArgs.DocIncludes
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
