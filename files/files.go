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

// Package files holds common functions to operate with files and directories.
package files

import (
	"fmt"
	"os"

	"github.com/bmatcuk/doublestar/v4"
)

const (
	// DocumentationFilePermission is the mode used when writing documentation data files.
	DocumentationFilePermission uint32 = 0644
	// WritePermission is the mode used for private writable files in tests and fixtures.
	WritePermission uint32 = 0600
)

// IsFileExist reports whether the given path exists as a file.
//
// Parameters:
// filePath - provides a file path or glob pattern.
//
// Returns:
// bool - whether the first match exists as a file.
// error - when the path cannot be inspected or points to a directory.
func IsFileExist(filePath string) (bool, error) {
	exists, info, err := validatePathExists(filePath)
	if err != nil {
		return false, err
	}
	if exists {
		if info.IsDir() {
			return false, fmt.Errorf("%s is a directory, the file was expected", filePath)
		}

		return true, nil
	}

	return false, nil
}

// IsDirExist reports whether the given path exists as a directory.
//
// Parameters:
// path - provides a directory path or glob pattern.
//
// Returns:
// bool - whether the first match exists as a directory.
// error - when the path cannot be inspected or points to a file.
func IsDirExist(path string) (bool, error) {
	exists, info, err := validatePathExists(path)
	if err != nil {
		return false, err
	}
	if exists {
		if info.IsDir() {
			return true, nil
		}

		return false, fmt.Errorf("%s is a file, the directory was expected", path)
	}

	return false, nil
}

// Reports whether the given path is valid and exist in the file system. Also returns a FileInfo if
// the path (relative or absolute) exists.
func validatePathExists(path string) (bool, os.FileInfo, error) {
	// Getting matches for the given path if it is a glob format. Otherwise, does nothing.
	matches, err := doublestar.FilepathGlob(path)
	if err != nil {
		return false, nil, err
	}

	if len(matches) == 0 {
		return false, nil, nil
	}

	firstMatch := matches[0]
	info, err := os.Stat(firstMatch)

	if err != nil {
		if os.IsNotExist(err) {
			return false, nil, nil
		}

		return false, nil, err
	}

	return true, info, nil
}
