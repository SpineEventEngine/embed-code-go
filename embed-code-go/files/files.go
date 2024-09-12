/*
 * Copyright 2024, TeamDev. All rights reserved.
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

// Package files holds common functions to operate with files and directories.
package files

import (
	"bufio"
	"embed-code/embed-code-go/configuration"
	"fmt"
	"os"
	"path/filepath"
)

const (
	ReadWriteExecPermission uint32 = 0777
	WritePermission         uint32 = 0600
)

// WriteLinesToFile writes lines to the file at given filePath.
func WriteLinesToFile(filepath string, lines []string) {
	file, err := os.Create(filepath)
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		if err = file.Close(); err != nil {
			panic(err)
		}
	}(file)

	for _, s := range lines {
		_, err := file.WriteString(s + "\n")
		if err != nil {
			panic(err)
		}
	}
}

// ReadFile reads and returns all lines from the file at given filePath.
func ReadFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	var lines []string
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	reader := bufio.NewReader(file)

	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			break
		}
		lines = append(lines, string(line))
	}

	return lines, nil
}

// BuildDocRelativePath builds a relative path for documentation file with a given config.
func BuildDocRelativePath(absolutePath string, config configuration.Configuration) string {
	relativePath, err := filepath.Rel(config.DocumentationRoot, absolutePath)
	if err != nil {
		panic(err)
	}

	return relativePath
}

// EnsureDirExists creates dir at given path if it doesn't exist. Does nothing if exists.
func EnsureDirExists(path string) error {
	exist, err := IsDirExist(path)
	if err != nil {
		return err
	}
	if !exist {
		err = os.MkdirAll(path, os.FileMode(ReadWriteExecPermission))
		if err != nil {
			return err
		}
	}

	return nil
}

// IsFileExist reports whether the given path to a file exists in the file system.
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

// IsDirExist reports whether the given directory exists in the file system.
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
// the path exists.
func validatePathExists(path string) (bool, os.FileInfo, error) {
	info, err := os.Stat(path)

	if err != nil {
		if os.IsNotExist(err) {
			return false, nil, nil
		}

		return false, nil, err
	}

	return true, info, nil
}
