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

package fragmentation

import (
	"bufio"
	"os"
	"path/filepath"

	"embed-code/embed-code-go/configuration"
)

// Creates dir at given dirPath if it doesn't exist.
// Does nothing if exists.
func EnsureDirExists(dirPath string) {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		var readWriteExecPermission uint32 = 777
		err := os.MkdirAll(dirPath, os.FileMode(readWriteExecPermission))
		if err != nil {
			panic(err)
		}
	}
}

// Reports whether file exists at given filePath.
//
// Returns an error if any faced.
func IsFileExists(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false, nil
	}

	return err == nil, err
}

// Reads and returns all lines from the file at given filePath.
func ReadLines(filePath string) []string {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	lines := []string{}
	defer file.Close()

	r := bufio.NewReader(file)

	for {
		line, _, err := r.ReadLine()
		if err != nil {
			break
		}
		lines = append(lines, string(line))
	}

	return lines
}

// Writes lines to the file at given filePath.
func WriteLinesToFile(filepath string, lines []string) {
	file, err := os.Create(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	for _, s := range lines {
		_, err := file.WriteString(s + "\n")
		if err != nil {
			panic(err)
		}
	}
}

// Builds a relative path for documentation file with a given config.
func BuildDocRelativePath(absolutePath string, config configuration.Configuration) string {
	absolutePath, err := filepath.Rel(config.DocumentationRoot, absolutePath)
	if err != nil {
		panic(err)
	}

	return absolutePath
}
