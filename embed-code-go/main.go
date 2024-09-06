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
	"fmt"

	"embed-code/embed-code-go/cli"
)

// The entry point for embed-code.
//
// There are three modes, which are chosen by 'mode' arg. If it is set to 'check',
// then the checking for up-to-date is performed. If it is set to 'embed', the embedding is
// performed. If it is set to 'analyze', the analyzing is performed.
//
// Embedding is the process that consists of the following steps:
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
// If config file is not exists or does not contain 'code-path' and 'docs-path' fields, the
// embedding fails.
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
//   - fragments-path — a path to a directory with code fragments. Default value is
//     "./build/fragments";
//   - separator — a string which is used as a separator between code fragments. Default value
//     is "...".
func main() {
	userArgs := cli.ReadArgs()

	err := cli.ValidateConfig(userArgs)
	if err != nil {
		fmt.Println("Validation error:")
		fmt.Println(err.Error())

		return
	}

	if userArgs.ConfigPath != "" {
		err = cli.ValidateConfigFile(userArgs.ConfigPath)
		if err != nil {
			fmt.Println("Configuration file validation error:")
			fmt.Println(err.Error())

			return
		}
		userArgs = cli.FillArgsFromConfigFile(userArgs)
	}

	config := cli.BuildEmbedCodeConfiguration(userArgs)

	switch userArgs.Mode {
	case cli.ModeCheck:
		cli.CheckCodeSamples(config)
	case cli.ModeEmbed:
		cli.EmbedCodeSamples(config)
		cli.CheckCodeSamples(config)
	case cli.ModeAnalyze:
		cli.AnalyzeCodeSamples(config)
	}
}
