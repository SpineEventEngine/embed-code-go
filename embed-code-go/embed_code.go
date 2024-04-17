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
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type flags struct {
	codeRoot      string
	docsRoot      string
	configPath    string
	checkUpToDate bool
}

type configFields struct {
	codeRoot string
	docsRoot string
}

func readRootsFromConfig(configPath string) configFields {
	content, err := ioutil.ReadFile("final-result.yml")
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

func readFlags() flags {
	codeRoot := flag.String("code_root", "", "a path to a root directory with code files")
	docsRoot := flag.String("docs_root", "ag", "a path to a root directory with docs files")

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
