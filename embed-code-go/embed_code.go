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

package main

import (
	"embed-code/embed-code-go/cli"
	"embed-code/embed-code-go/configuration"
	"flag"
	"fmt"

	"os"

	"gopkg.in/yaml.v2"
)

// Struct with user-specified flags.
//
// codeRoot — a path to a root directory with code files.
//
// docsRoot — a path to a root directory with docs files.
//
// configPath — a path to a yaml configuration file which contains the roots.
//
// checkUpToDate — true to check for code embeddings to be up to date. Otherwise, the embedding is performed.
type flags struct {
	codeRoot      string
	docsRoot      string
	configPath    string
	checkUpToDate bool
}

// Struct with roots that contained in a yaml configuration file.
//
// codeRoot — a path to a root directory with code files.
//
// docsRoot — a path to a root directory with docs files.
type configFields struct {
	codeRoot string
	docsRoot string
}

// Reads the roots from the provided configPath and returns a configFields struct.
//
// configPath — a path to a yaml configuration file which contains the roots.
//
// Returns a configFields struct filled with the roots.
func readRootsFromConfig(configPath string) configFields {
	content, err := os.ReadFile(configPath)
	if err != nil {
		panic(err)
	}

	configFields := configFields{}
	err = yaml.Unmarshal(content, &configFields)
	if err != nil {
		panic(err)
	}

	return configFields
}

// Reads user-specified flags from the command line.
//
// Returns a flags struct filled with the corresponding flags.
func readFlags() flags {
	codeRoot := flag.String("code_root", "", "a path to a root directory with code files")
	docsRoot := flag.String("docs_root", "", "a path to a root directory with docs files")

	configPath := flag.String("config_path", "", "a path to a configuration file")

	checkUpToDate := flag.Bool("up_to_date", false, "true to check for code embeddings to be up to date")

	flag.Parse()

	return flags{
		codeRoot:      *codeRoot,
		docsRoot:      *docsRoot,
		configPath:    *configPath,
		checkUpToDate: *checkUpToDate,
	}

}

// Checks the validity of user-provided flags and returns an error message if any of the validation rules are broken.
// If everything is ok, returns an empty string.
//
// flagsSet — a struct with user-provided flags.
func validate(flagsSet flags) string {
	isRootsSet := flagsSet.codeRoot != "" && flagsSet.docsRoot != ""
	isOneOfRootsSet := flagsSet.codeRoot != "" || flagsSet.docsRoot != ""
	isConfigSet := flagsSet.configPath != ""

	validationMessage := ""

	if isConfigSet && isOneOfRootsSet {
		return "Config path cannot be set when code_root and docs_root are set."
	}
	if isOneOfRootsSet && !isRootsSet {
		return "If one of code_root and docs_root is set, the another one must be set as well."
	}
	if !(isRootsSet || isConfigSet) {
		return "Embed code should be used with either config_path or both code_root and docs_root being set."
	}

	return validationMessage
}

// Generates and returns a configuration based on the provided flags.
//
// flagsSet — a struct with user-provided flags.
func buildEmbedCodeConfiguration(flagsSet flags) configuration.Configuration {
	codeRoot := flagsSet.codeRoot
	docsRoot := flagsSet.docsRoot
	if flagsSet.configPath != "" {
		configFields := readRootsFromConfig(flagsSet.configPath)
		codeRoot = configFields.codeRoot
		docsRoot = configFields.docsRoot
	}

	return configuration.NewConfigurationWithRoots(codeRoot, docsRoot)
}

// The entry point for embed-code.
//
// There are two modes, which are chosen by 'up_to_date' flag. If it is set to 'true',
// then the check for up-to-date is performed. Otherwise, the embedding is performed.
//
// There are two options to set the roots:
//   - code_root and docs_root flags, in this case roots are read directly from provided paths;
//   - config_path flag, in this case roots are read from the given config file.
//
// If both options are missed, the embedding fails.
// If both options are set, the embedding fails as well.
//
// All possible flags:
//   - code_root — a path to a root directory with code files;
//   - docs_root — a path to a root directory with docs files;
//   - config_path — a path to a yaml configuration file;
//   - up_to_date — true to check for code embeddings to be up-to-date. Otherwise, the embedding is performed.
func main() {

	flagsSet := readFlags()

	validationMessage := validate(flagsSet)
	if validationMessage != "" {
		fmt.Println("Validation error:")
		fmt.Println(validationMessage)
		return
	}

	config := buildEmbedCodeConfiguration(flagsSet)

	if flagsSet.checkUpToDate {
		cli.CheckCodeSamples(config)
	} else {
		cli.EmbedCodeSamples(config)
		cli.CheckCodeSamples(config)
	}
}
