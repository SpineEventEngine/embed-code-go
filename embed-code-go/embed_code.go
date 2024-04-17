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
	"fmt"
)

// The entry point for embed-code.
//
// There are two modes, which are chosen by 'up_to_date' arg. If it is set to 'true',
// then the check for up-to-date is performed. Otherwise, the embedding is performed.
//
// There are two options to set the roots:
//   - code_root and docs_root args, in this case roots are read directly from provided paths;
//   - config_path arg, in this case roots are read from the given config file.
//
// If both options are missed, the embedding fails.
// If both options are set, the embedding fails as well.
//
// All possible args:
//   - code_root — a path to a root directory with code files;
//   - docs_root — a path to a root directory with docs files;
//   - config_path — a path to a yaml configuration file. It must contain 'code_root' and 'docs_root' fields;
//   - up_to_date — true to check for code embeddings to be up-to-date. Otherwise, the embedding is performed.
func main() {

	userArgs := cli.ReadArgs()

	validationMessage := cli.Validate(userArgs)
	if validationMessage != "" {
		fmt.Println("Validation error:")
		fmt.Println(validationMessage)
		return
	}

	config := cli.BuildEmbedCodeConfiguration(userArgs)

	if userArgs.CheckUpToDate {
		cli.CheckCodeSamples(config)
	} else {
		cli.EmbedCodeSamples(config)
		cli.CheckCodeSamples(config)
	}
}
