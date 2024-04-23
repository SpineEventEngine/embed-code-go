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
// СodeRoot — a path to a root directory with code files.
//
// DocsRoot — a path to a root directory with docs files.
//
// CodeIncludes — a string with comma-separated patterns for filtering the code files to be considered.
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
// FragmentsDir — a directory where fragmented code is stored. A temporary directory that should not be
// tracked in VCS. The default value is: "./build/fragments".
//
// Separator — a string that's inserted between multiple partitions of a single fragment.
// The default value is "...".
//
// ConfigFilePath — a path to a yaml configuration file which contains the roots.
//
// Mode — defines the mode of embed-code execution.
type Args struct {
	CodeRoot       string
	DocsRoot       string
	CodeIncludes   string
	DocIncludes    string
	FragmentsDir   string
	Separator      string
	ConfigFilePath string
	Mode           string
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

//
// Public functions
//

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
	codeIncludes := flag.String("code_includes", "", "a comma-separated string of glob patterns for code files to include")
	docIncludes := flag.String("doc_includes", "", "a comma-separated string of glob patterns for docs files to include")
	fragmentsDir := flag.String("fragments_dir", "", "a path to a directory where fragmented code is stored")
	separator := flag.String("separator", "", "a string that's inserted between multiple partitions of a single fragment")
	configFilePath := flag.String("config_file_path", "", "a path to a yaml configuration file")
	mode := flag.String("mode", "",
		"a mode of embed-code execution, which can be 'check' or 'embed'")

	flag.Parse()

	return Args{
		CodeRoot:       *codeRoot,
		DocsRoot:       *docsRoot,
		CodeIncludes:   *codeIncludes,
		DocIncludes:    *docIncludes,
		FragmentsDir:   *fragmentsDir,
		Separator:      *separator,
		ConfigFilePath: *configFilePath,
		Mode:           *mode,
	}

}

// Checks the validity of provided userArgs and returns an error message if any of the validation rules are broken.
// If everything is ok, returns an empty string.
//
// userArgs — a struct with user-provided args.
func Validate(userArgs Args) string {
	isModeSet := userArgs.Mode != ""
	isRootsSet := userArgs.CodeRoot != "" && userArgs.DocsRoot != ""
	isOneOfRootsSet := userArgs.CodeRoot != "" || userArgs.DocsRoot != ""
	isConfigSet := userArgs.ConfigFilePath != ""
	isOptionalParamsSet := userArgs.CodeIncludes != "" || userArgs.DocIncludes != "" ||
		userArgs.FragmentsDir != "" || userArgs.Separator != ""

	validationMessage := ""

	if !isModeSet {
		return "Mode must be set."
	}
	if isConfigSet && (isOneOfRootsSet || isOptionalParamsSet) {
		return "Config path cannot be set when code_root, docs_root or optional params are set."
	}
	if isOneOfRootsSet && !isRootsSet {
		return "If one of code_root and docs_root is set, the another one must be set as well."
	}
	if !(isRootsSet || isConfigSet) {
		return "Embed code should be used with either config_file_path or both code_root and docs_root being set."
	}

	return validationMessage
}

// Performs several checks to ensure that the necessary configuration values are present.
// Also checks for the existence of the config file.
//
// configFilePath — a path to a yaml configuration file.
//
// Returns validation message. If everything is ok, returns an empty string.
func ValidateConfigFile(configFilePath string) string {
	validationMessage := ""

	stat, err := os.Stat(configFilePath)
	if os.IsNotExist(err) {
		return fmt.Sprintf("The file %s is not exists.", configFilePath)
	}
	if stat.IsDir() {
		return fmt.Sprintf("%s is a dir, not a file.", configFilePath)
	}
	configFields := readConfigFields(configFilePath)
	if configFields.CodeRoot == "" || configFields.DocsRoot == "" {
		return "Config must include both code_root and docs_root fields."
	}
	return validationMessage
}

// Fills args with the values read from config file.
//
// args — an Args struct with user-provided args.
//
// Returns filled Args.
func FillArgsFromConfigFile(args Args) Args {
	configFields := readConfigFields(args.ConfigFilePath)
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

//
// Private functions
//

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
func readConfigFields(configFilePath string) ConfigFields {
	content, err := os.ReadFile(configFilePath)
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
