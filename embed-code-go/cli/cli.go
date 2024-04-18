package cli

import (
	"embed-code/embed-code-go/configuration"
	"embed-code/embed-code-go/embedding"
	"embed-code/embed-code-go/fragmentation"
	"flag"
	"fmt"
	"strings"

	"os"

	"gopkg.in/yaml.v3"
)

// User-specified embed-code Args.
//
// codeRoot — a path to a root directory with code files.
//
// docsRoot — a path to a root directory with docs files.
//
// codeIncludes — a string with coma-separated patterns for filtering the code files to be considered.
// Directories are never matched by these patterns.
// For example, "**/*.java,**/*.gradle".
// The default value is "**/*.*".
//
// docIncludes — a string with coma-separated patterns for filtering files
// in which we should look for embedding instructions.
// The patterns are resolved relatively to the `documentation_root`.
// Directories are never matched by these patterns.
// For example, "docs/**/*.md,guides/*.html".
// The default value is "**/*.md,**/*.html".
//
// fragmentsDir — a directory where fragmented code is stored. A temporary directory that should not be
// tracked in VCS. The default value is: "./build/fragments".
//
// separator — a string that's inserted between multiple partitions of a single fragment.
// The default value is "...".
//
// configPath — a path to a yaml configuration file which contains the roots.
//
// checkUpToDate — true to check for code embeddings to be up-to-date. Otherwise, the embedding is performed.
type Args struct {
	CodeRoot      string
	DocsRoot      string
	CodeIncludes  string
	DocIncludes   string
	FragmentsDir  string
	Separator     string
	ConfigPath    string
	CheckUpToDate bool
}

// Needed for yaml.Unmarshal to parse into.
type ConfigFields struct {
	CodeRoot     string `yaml:"code_root"`
	DocsRoot     string `yaml:"docs_root"`
	CodeIncludes string `yaml:"code_includes"`
	DocIncludes  string `yaml:"doc_includes"`
	FragmentsDir string `yaml:"fragments_dir"`
	Separator    string `yaml:"separator"`
}

// Checks documentation to be up-to-date with code files. Raises UnexpectedDiffError if not.
//
// config — a configuration for checking code samples.
func CheckCodeSamples(config configuration.Configuration) {
	fragmentation.WriteFragmentFiles(config)
	embedding.CheckUpToDate(config)
}

// Embeds code fragments in documentation files.
//
// config — a configuration for embedding.
func EmbedCodeSamples(config configuration.Configuration) {
	fragmentation.WriteFragmentFiles(config)
	embedding.EmbedAll(config)
}

// Reads user-specified args from the command line.
//
// Returns an Args struct filled with the corresponding args.
func ReadArgs() Args {
	codeRoot := flag.String("code_root", "", "a path to a root directory with code files")
	docsRoot := flag.String("docs_root", "", "a path to a root directory with docs files")
	codeIncludes := flag.String("code_includes", "", "a coma-separated list of glob patterns for code files to include")
	docIncludes := flag.String("doc_includes", "", "a coma-separated list of glob patterns for docs files to include")
	fragmentsDir := flag.String("fragments_dir", "", "a path to a directory where fragmented code is stored")
	separator := flag.String("separator", "", "a string that's inserted between multiple partitions of a single fragment")
	configPath := flag.String("config_path", "", "a path to a yaml configuration file")
	checkUpToDate := flag.Bool("up_to_date", false,
		"true to check for code embeddings to be up-to-date, false to perform embedding")

	flag.Parse()

	return Args{
		CodeRoot:      *codeRoot,
		DocsRoot:      *docsRoot,
		CodeIncludes:  *codeIncludes,
		DocIncludes:   *docIncludes,
		FragmentsDir:  *fragmentsDir,
		Separator:     *separator,
		ConfigPath:    *configPath,
		CheckUpToDate: *checkUpToDate,
	}

}

// Checks the validity of provided userArgs and returns an error message if any of the validation rules are broken.
// If everything is ok, returns an empty string.
//
// userArgs — a struct with user-provided args.
func Validate(userArgs Args) string {
	isRootsSet := userArgs.CodeRoot != "" && userArgs.DocsRoot != ""
	isOneOfRootsSet := userArgs.CodeRoot != "" || userArgs.DocsRoot != ""
	isConfigSet := userArgs.ConfigPath != ""
	isOptionalParamsSet := userArgs.CodeIncludes != "" || userArgs.DocIncludes != "" ||
		userArgs.FragmentsDir != "" || userArgs.Separator != ""

	validationMessage := ""

	if isConfigSet && (isOneOfRootsSet || isOptionalParamsSet) {
		return "Config path cannot be set when code_root, docs_root or optional params are set."
	}
	if isOneOfRootsSet && !isRootsSet {
		return "If one of code_root and docs_root is set, the another one must be set as well."
	}
	if !(isRootsSet || isConfigSet) {
		return "Embed code should be used with either config_path or both code_root and docs_root being set."
	}

	return validationMessage
}

// Performs several checks to ensure that the necessary configuration values are present.
// Also checks for the existence of the config file.
//
// configPath — a path to a yaml configuration file.
//
// Returns validation message. If everything is ok, returns an empty string.
func ValidateConfig(configPath string) string {
	validationMessage := ""

	stat, err := os.Stat(configPath)
	if os.IsNotExist(err) {
		return fmt.Sprintf("The file %s is not exists.", configPath)
	}
	if stat.IsDir() {
		return fmt.Sprintf("%s is a dir, not a file.", configPath)
	}
	configFields := readConfigFields(configPath)
	if configFields.CodeRoot == "" || configFields.DocsRoot == "" {
		return "Config must include both code_root and docs_root fields."
	}
	return validationMessage
}

// Fills args with the values read from config file.
//
// configPath — a path to a yaml configuration file.
//
// args — an Args struct with user-provided args.
//
// Returns filled Args.
func FillArgsFromConfig(args Args) Args {
	configFields := readConfigFields(args.ConfigPath)
	args.CodeRoot = configFields.CodeRoot
	args.DocsRoot = configFields.DocsRoot

	if configFields.CodeIncludes != "" {
		args.CodeIncludes = configFields.CodeIncludes
	}
	if configFields.DocIncludes != "" {
		args.DocIncludes = configFields.DocIncludes
	}
	if configFields.FragmentsDir != "" {
		args.FragmentsDir = configFields.FragmentsDir
	}
	if configFields.Separator != "" {
		args.Separator = configFields.Separator
	}
	return args
}

// Generates and returns a configuration based on provided userArgs.
//
// userArgs — a struct with user-provided args.
func BuildEmbedCodeConfiguration(userArgs Args) configuration.Configuration {

	embedCodeConfig := configuration.NewConfiguration()
	embedCodeConfig.CodeRoot = userArgs.CodeRoot
	embedCodeConfig.DocumentationRoot = userArgs.DocsRoot

	if userArgs.CodeIncludes != "" {
		embedCodeConfig.CodeIncludes = parseListArgument(userArgs.CodeIncludes)
	}
	if userArgs.DocIncludes != "" {
		embedCodeConfig.DocIncludes = parseListArgument(userArgs.DocIncludes)
	}
	if userArgs.FragmentsDir != "" {
		embedCodeConfig.FragmentsDir = userArgs.FragmentsDir
	}
	if userArgs.Separator != "" {
		embedCodeConfig.Separator = userArgs.Separator
	}
	return embedCodeConfig
}

// Returns a list of strings from given coma-separated string listArgument.
func parseListArgument(listArgument string) []string {
	extractedArgs := strings.Split(listArgument, ",")
	for i, v := range extractedArgs {
		extractedArgs[i] = strings.TrimSpace(v)
	}
	return extractedArgs
}

// Reads the file from provided configPath and returns a ConfigFields struct.
//
// configPath — a path to a yaml configuration file.
//
// Returns a filled ConfigFields struct.
func readConfigFields(configPath string) ConfigFields {
	content, err := os.ReadFile(configPath)
	if err != nil {
		panic(err)
	}

	configFields := ConfigFields{}
	err = yaml.Unmarshal(content, &configFields)
	if err != nil {
		panic(err)
	}

	return configFields
}
