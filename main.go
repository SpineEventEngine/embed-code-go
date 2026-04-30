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

package main

import (
	"embed-code/embed-code-go/cli"
	"embed-code/embed-code-go/configuration"
	"embed-code/embed-code-go/logging"
	"fmt"
	"log/slog"
	"path/filepath"
)

// Version of the embed-code application.
const Version = "1.1.0"

// The entry point for embed-code.
//
// There are three modes, which are chosen by 'mode' arg. If it is set to 'check',
// then the checking for up-to-date is performed. If it is set to 'embed', the embedding is
// performed. If it is set to 'analyze', the analyzing is performed.
//
// EmbeddingInstruction is the process that consists of the following steps:
//   - the code fragments are extracted from the code files;
//   - the docs files are scanned for <embed-code> tags;
//   - for each tag, the code fragments are embedded into the docs. The embedding
//     is parametrized with the tag attributes.
//
// Checking for up-to-date is the process that consists of the following steps:
//   - the code fragments are extracted from the code files;
//   - the docs files are scanned for <embed-code> tags;
//   - for each tag, the code fragments are compared to the code which is already embedded
//     into the docs;
//   - if there is a difference, the error is reported.
//
// The 'mode' arg is required.
//
// Embed code also needs root directories to be set.
// There are two options to set them:
//   - code-path and docs-path args, in this case roots are read directly from provided paths;
//   - config-path arg, in this case roots are read from the given config file.
//
// If both options are missed, the embedding fails.
// If both options are set, the embedding fails as well.
// If config file does not exist, or contains neither root 'code-path' and 'docs-path' fields nor
// 'embeddings' entries, the embedding fails.
//
// All possible args:
//   - code-path — a path to a root directory with code files;
//   - docs-path — a path to a root directory with docs files;
//   - config-path — a path to a yaml configuration file;
//   - mode — string which represents the mode of embed-code execution. if it is set to 'check',
//     then the checking for up-to-date is performed. If it is set to 'embed', the embedding
//     is performed.
//     If it is set to 'analyze', the analyzing is performed;
//   - code-includes — a comma-separated string of glob patterns for code files to include.
//     For example:
//     "**/*.java,**/*.gradle". Default value is "**/*.*";
//   - doc-includes — a comma-separated string of glob patterns for docs files to include.
//     For example:
//     "docs/**/*.md,guides/*.html". Default value is "**/*.md,**/*.html";
//   - doc-excludes - a comma-separated string of glob patterns for docs files to exclude from
//     the embedding.
//     For example:
//     "old-docs/**/*.md,old-guides/*.html". It is not set by default;
//   - fragments-path — a path to a directory with code fragments. Default value is
//     "./build/fragments";
//   - separator — a string which is used as a separator between code fragments. Default value
//     is "...".
func main() {
	fmt.Println(fmt.Sprintf("Running embed-code v%s.", Version))
	userArgs := cli.ReadArgs()
	configureLogging(userArgs)
	defer logging.HandlePanic(userArgs.Stacktrace)

	if cli.IsUsingConfigFile(userArgs) {
		err := cli.ValidateConfigFile(userArgs)
		if err != nil {
			slog.Error("The provided config file is not valid.", "error", err)

			return
		}
		userArgs, err = cli.FillArgsFromConfigFile(userArgs)
		if err != nil {
			slog.Error("Received an issue while reading config file: ", "error", err)

			return
		}
	}
	err := cli.ValidateConfig(userArgs)
	if err != nil {
		slog.Error("User arguments are not valid.", "error", err)

		return
	}
	configs := cli.BuildEmbedCodeConfiguration(userArgs)

	switch userArgs.Mode {
	case cli.ModeCheck:
		for _, config := range configs {
			cli.CheckCodeSamples(config)
		}
		fmt.Println("The documentation files are up-to-date with code files.")
	case cli.ModeEmbed:
		embedByConfigs(configs)
		fmt.Println("Embedding process finished.")
	case cli.ModeAnalyze:
		for _, config := range configs {
			cli.AnalyzeCodeSamples(config)
		}
		fmt.Println("Analysis is completed, analytics files can be found in /build/analytics folder.")
	}
}

// configureLogging configures the slog logging.
func configureLogging(config cli.Config) {
	level := slog.LevelWarn
	if config.Info {
		level = slog.LevelInfo
	}
	logger := slog.New(&logging.Handler{Level: level})
	slog.SetDefault(logger)
}

// embedByConfig runs the embedByConfig for all configs and logs the results.
func embedByConfigs(configs []configuration.Configuration) {
	var totalEmbeddedFiles []string
	totalEmbeddings := 0
	totalFragments := 0
	for _, config := range configs {
		result := cli.EmbedCodeSamples(config)
		totalEmbeddedFiles = append(totalEmbeddedFiles, result.UpdatedTargetFiles...)
		totalEmbeddings += result.TotalEmbeddings
		totalFragments += result.TotalFragments
	}
	if len(totalEmbeddedFiles) == 0 &&
		totalEmbeddings != 0 &&
		totalFragments != 0 {
		fmt.Println("All documentation files are already up to date. Nothing to update.")
	}
	if len(totalEmbeddedFiles) == 1 {
		fmt.Println("File updated:")
	}
	if len(totalEmbeddedFiles) > 1 {
		fmt.Println("Files updated:")
	}
	for _, updatedDocFile := range totalEmbeddedFiles {
		absPath, err := filepath.Abs(updatedDocFile)
		if err != nil {
			panic(err)
		}
		fmt.Printf("- file://%s.\n", absPath)
	}
}
