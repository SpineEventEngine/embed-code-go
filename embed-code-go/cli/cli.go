// Copyright 2024, TeamDev. All rights reserved.
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
// CodePath — a path to a root directory with code files.
//
// DocsPath — a path to a root directory with docs files.
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
type Config struct {
	CodePath      string `yaml:"code-path"`
	DocsPath      string `yaml:"docs-path"`
	CodeIncludes  string `yaml:"code-includes"`
	DocIncludes   string `yaml:"doc-includes"`
	FragmentsPath string `yaml:"fragments-path"`
	Separator     string `yaml:"separator"`
	ConfigPath    string
	Mode          string
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
	err := fragmentation.WriteFragmentFiles(config)
	if err != nil {
		panic(err)
	}
	embedding.CheckUpToDate(config)
}

// EmbedCodeSamples embeds code fragments in documentation files.
//
// config — a configuration for embedding.
func EmbedCodeSamples(config configuration.Configuration) {
	err := fragmentation.WriteFragmentFiles(config)
	if err != nil {
		panic(err)
	}
	embedding.EmbedAll(config)
}

// AnalyzeCodeSamples analyzes code fragments in documentation files.
//
// config — a configuration for embedding.
func AnalyzeCodeSamples(config configuration.Configuration) {
	err := fragmentation.WriteFragmentFiles(config)
	if err != nil {
		panic(err)
	}
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
	fragmentsPath := flag.String("fragments-path", "",
		"a path to a directory where fragmented code is stored")
	separator := flag.String("separator", "",
		"a string that's inserted between multiple partitions of a single fragment")
	configPath := flag.String("config-path", "", "a path to a yaml configuration file")
	mode := flag.String("mode", "",
		"a mode of embed-code execution, which can be 'check' or 'embed'")

	flag.Parse()

	return Config{
		CodePath:      *codePath,
		DocsPath:      *docsPath,
		CodeIncludes:  *codeIncludes,
		DocIncludes:   *docIncludes,
		FragmentsPath: *fragmentsPath,
		Separator:     *separator,
		ConfigPath:    *configPath,
		Mode:          *mode,
	}
}

// FillArgsFromConfigFile fills config with the values read from config file.
//
// args — Config struct with user-provided args.
//
// Returns filled Config.
func FillArgsFromConfigFile(args Config) (Config, error) {
	configFields := readConfigFields(args.ConfigPath)
	args.CodePath = configFields.CodePath
	args.DocsPath = configFields.DocsPath

	if isNotEmpty(configFields.CodeIncludes) {
		args.CodeIncludes = configFields.CodeIncludes
	}
	if isNotEmpty(configFields.DocIncludes) {
		args.DocIncludes = configFields.DocIncludes
	}
	if isNotEmpty(configFields.FragmentsPath) {
		args.FragmentsPath = configFields.FragmentsPath
	}
	if isNotEmpty(configFields.Separator) {
		args.Separator = configFields.Separator
	}

	return args, nil
}

// BuildEmbedCodeConfiguration generates and returns a configuration based on provided userArgs.
//
// userArgs — a Config with user-provided args.
func BuildEmbedCodeConfiguration(userArgs Config) configuration.Configuration {
	embedCodeConfig := configuration.NewConfiguration()
	embedCodeConfig.CodeRoot = userArgs.CodePath
	embedCodeConfig.DocumentationRoot = userArgs.DocsPath

	if isNotEmpty(userArgs.CodeIncludes) {
		embedCodeConfig.CodeIncludes = parseListArgument(userArgs.CodeIncludes)
	}
	if isNotEmpty(userArgs.DocIncludes) {
		embedCodeConfig.DocIncludes = parseListArgument(userArgs.DocIncludes)
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
