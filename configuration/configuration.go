/*
 * Copyright 2026, TeamDev. All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 * Redistribution and use in source and/or binary forms, with or without
 * modification, must retain the above copyright notice and the following
 * disclaimer.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
 * "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
 * LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
 * A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
 * OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
 * SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
 * LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
 * DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
 * THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 * OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

// Package configuration contains normalized embed-code settings.
package configuration

import (
	_type "embed-code/embed-code-go/type"
)

const (
	// DefaultSeparator joins multiple partitions of a single fragment.
	DefaultSeparator = "..."
)

// DefaultDocIncludes contains the default documentation glob patterns.
var DefaultDocIncludes = []string{"**/*.md", "**/*.html"}

// Configuration contains embed-code processing settings.
//
// It is used to get data for scanning docs and resolving source files.
// The example of creating the Configuration with default values:
//
//	var config = configuration.NewConfiguration()
type Configuration struct {
	// Name identifies this configuration when it is built from an embeddings entry.
	Name string

	// CodeRoots is a list of directories with the source code to be embedded.
	CodeRoots _type.NamedPathList

	// DocumentationRoot is a root directory of the documentation files.
	DocumentationRoot string

	// DocIncludes is a list of patterns for filtering files in which we should look for embedding
	// instructions.
	//
	// The patterns are resolved relatively to the `documentation_root`.
	//
	// Directories are never matched by these patterns.
	//
	// For example, ["docs/**/*.md", "guides/*.html"].
	//
	// The default value is ["**/*.md", "**/*.html"].
	DocIncludes []string

	// DocExcludes is a list of patterns for documentation files that should not be
	// processed for embedding instructions.
	//
	// The patterns are resolved relatively to the `documentation_root`.
	//
	// A pattern can match both directories and files.
	//
	// For example, ["old-docs/**/*.md", "old-docs-v1/**/*"]
	//
	// By default, it is not set.
	DocExcludes []string

	// Separator is a string that's inserted between multiple partitions of a single fragment.
	//
	// The default value is: "..." (three dots).
	Separator string
}

// NewConfiguration builds the default config.
//
// Returns configuration with default include patterns and separator.
func NewConfiguration() Configuration {
	return Configuration{
		DocIncludes: DefaultDocIncludes,
		Separator:   DefaultSeparator,
	}
}
