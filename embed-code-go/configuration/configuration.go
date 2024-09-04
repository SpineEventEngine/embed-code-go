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

// The configuration of the plugin.
package configuration

const (
	DefaultSeparator    = "..."
	DefaultFragmentsDir = "./build/fragments"
)

var DefaultInclude = []string{"**/*.*"}
var DefaultDocIncludes = []string{"**/*.md", "**/*.html"}

// The configuration for all the settings for the plugin to work.
//
// It is used to get data for scanning for and doc files, to receive fragments' dir and separator for partitions.
// The example of creating the Configuration with default values:
//
//	var config = configuration.NewConfiguration()
//
// If there's need to modify the default configuration,
// it can be done with just setting values to corresponding fields:
//
//	config.FragmentsDir = "foo/bar"
type Configuration struct {
	// A root directory of the source code to be embedded.
	CodeRoot string

	// A root directory of the documentation files.
	DocumentationRoot string

	// A list of patterns for filtering the code files to be considered.
	//
	// Directories are never matched by these patterns.
	//
	// For example, ["**/*.java", "**/*.gradle"].
	//
	// The default value is "**/*".
	CodeIncludes []string

	// A list of patterns for filtering files in which we should look for embedding instructions.
	//
	// The patterns are resolved relatively to the `documentation_root`.
	//
	// Directories are never matched by these patterns.
	//
	// For example, ["docs/**/*.md", "guides/*.html"].
	//
	// The default value is ["**/*.md", "**/*.html"].
	DocIncludes []string

	// A directory where fragmented code is stored. A temporary directory that should not be
	// tracked in VCS.
	//
	// The default value is: "./build/fragments".
	FragmentsDir string

	// A string that's inserted between multiple partitions of a single fragment.
	//
	// The default value is: "..." (three dots).
	Separator string
}

// Builds the default config.
func NewConfiguration() Configuration {
	return Configuration{
		CodeIncludes: DefaultInclude,
		DocIncludes:  DefaultDocIncludes,
		FragmentsDir: DefaultFragmentsDir,
		Separator:    DefaultSeparator,
	}
}
