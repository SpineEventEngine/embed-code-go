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
	_ "embed"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"embed-code/embed-code-go/cli"
	"embed-code/embed-code-go/configuration"
	"embed-code/embed-code-go/logging"
)

//go:embed VERSION
var versionFile string

// Version is the embed-code application version embedded from the VERSION file.
var Version = strings.TrimSpace(versionFile)

// The entry point for embed-code.
//
// There are two modes, which are chosen by 'mode' arg. If it is set to 'check',
// then the checking for up-to-date is performed. If it is set to 'embed',
// the embedding is performed.
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
//   - mode — string which represents the mode of embed-code execution. If it is set to 'check',
//     then the checking for up-to-date is performed. If it is set to 'embed', the embedding
//     is performed.
//   - doc-includes — a comma-separated string of glob patterns for docs files to include.
//     For example:
//     "docs/**/*.md,guides/*.html". Default value is "**/*.md,**/*.html";
//   - doc-excludes - a comma-separated string of glob patterns for docs files to exclude from
//     the embedding.
//     For example:
//     "old-docs/**/*.md,old-guides/*.html". It is not set by default;
//   - separator — a string which is used as a separator between code fragments. Default value
//     is "...".
func main() {
	fmt.Printf("Running embed-code v%s.\n", Version)
	userArgs := cli.ReadArgs()
	configureLogging(userArgs)
	defer logging.HandlePanic(userArgs.Stacktrace)
	source := "command line arguments"
	if cli.IsUsingConfigFile(userArgs) {
		source = fmt.Sprintf("configuration file `%s`", logging.FileReference(userArgs.ConfigPath))
	}
	slog.Info(fmt.Sprintf("Started embed-code in `%s` mode using %s.", userArgs.Mode, source))

	if cli.IsUsingConfigFile(userArgs) {
		err := cli.ValidateConfigFile(userArgs)
		if err != nil {
			exitWithError("The provided config file is not valid", err)
		}
		userArgs, err = cli.FillArgsFromConfigFile(userArgs)
		if err != nil {
			exitWithError("Received an issue while reading config file", err)
		}
	}
	err := cli.ValidateConfig(userArgs)
	if err != nil {
		exitWithError("User arguments are not valid", err)
	}
	configs := cli.BuildEmbedCodeConfiguration(userArgs)

	switch userArgs.Mode {
	case cli.ModeCheck:
		if err := checkByConfigs(configs); err != nil {
			exitWithError("Check failed", err)
		}
	case cli.ModeEmbed:
		if err := embedByConfigs(configs); err != nil {
			exitWithError("Embedding failed", err)
		}
		fmt.Println("Embedding process finished.")
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

// logError writes a user-facing error through the configured logger.
func logError(message string, err error) {
	slog.Error(formatError(message, err))
}

// formatError formats single errors inline and joined errors as a bullet list.
func formatError(message string, err error) string {
	errs := flattenedErrors(err)
	if len(errs) <= 1 {
		return fmt.Sprintf("%s: %v", message, err)
	}

	var builder strings.Builder
	builder.WriteString(message)
	builder.WriteString(":")
	for _, nestedErr := range errs {
		builder.WriteString("\n- ")
		builder.WriteString(nestedErr.Error())
	}

	return builder.String()
}

// flattenedErrors returns the leaf errors from a joined error joined.
func flattenedErrors(err error) []error {
	joined, ok := err.(interface {
		Unwrap() []error
	})
	if !ok {
		return []error{err}
	}

	var result []error
	for _, nestedErr := range joined.Unwrap() {
		result = append(result, flattenedErrors(nestedErr)...)
	}

	return result
}

// exitWithError writes a user-facing error and exits with a failing status.
func exitWithError(message string, err error) {
	logError(message, err)
	os.Exit(1)
}

// checkByConfigs runs check for all configs and reports documentation files that are outdated.
func checkByConfigs(configs []configuration.Configuration) error {
	var totalOutdatedFiles []string
	var checkErrors []error
	for _, config := range configs {
		outdatedFiles, err := cli.CheckCodeSamples(config)
		if err != nil {
			checkErrors = append(checkErrors, err)

			continue
		}
		totalOutdatedFiles = append(totalOutdatedFiles, outdatedFiles...)
	}
	if len(totalOutdatedFiles) > 0 {
		printFiles("File to update:", "Files to update:", totalOutdatedFiles)
		checkErrors = append(checkErrors,
			fmt.Errorf("the documentation files are not up-to-date with code files"))
	}
	if len(checkErrors) == 0 {
		fmt.Println("The documentation files are up-to-date with code files.")

		return nil
	}

	return errors.Join(checkErrors...)
}

// embedByConfigs runs embedding for all configs and logs the results.
func embedByConfigs(configs []configuration.Configuration) error {
	var totalEmbeddedFiles []string
	totalEmbeddings := 0
	var embeddingErrors []error
	for _, config := range configs {
		result, err := cli.EmbedCodeSamples(config)
		if err != nil {
			embeddingErrors = append(embeddingErrors, err)

			continue
		}
		totalEmbeddedFiles = append(totalEmbeddedFiles, result.UpdatedTargetFiles...)
		totalEmbeddings += result.TotalEmbeddings
	}
	if len(embeddingErrors) > 0 {
		return errors.Join(embeddingErrors...)
	}
	if len(totalEmbeddedFiles) == 0 && totalEmbeddings != 0 {
		fmt.Println("All documentation files are already up to date. Nothing to update.")
	}
	printFiles("File updated:", "Files updated:", totalEmbeddedFiles)

	return nil
}

// printFiles prints file paths with the singular or plural heading.
func printFiles(singularHeading string, pluralHeading string, files []string) {
	if len(files) == 1 {
		fmt.Println(singularHeading)
	}
	if len(files) > 1 {
		fmt.Println(pluralHeading)
	}
	for _, file := range files {
		fmt.Printf("- %s.\n", logging.FileReference(file))
	}
}
