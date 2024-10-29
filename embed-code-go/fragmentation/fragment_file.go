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
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	config "embed-code/embed-code-go/configuration"
	"embed-code/embed-code-go/files"
)

// FragmentFile is a file storing a single fragment from the file.
//
// CodePath — a relative path to a code file. The path is relative to Configuration.CodeRoot.
//
// FragmentName — a name of the fragment in the code file.
//
// Configuration — a configuration for embedding.
type FragmentFile struct {
	CodePath      string
	FragmentName  string
	Configuration config.Configuration
}

// NewFragmentFileFromAbsolute composes a FragmentFile for the given fragment in given codeFile.
//
// codeFile — an absolute path to a code file.
//
// fragmentName — a name of the fragment in the code file.
//
// configuration — configuration for embedding.
//
// Returns composed fragment.
func NewFragmentFileFromAbsolute(codeFile string, fragmentName string,
	config config.Configuration) FragmentFile {
	absolutePath, err := filepath.Abs(config.CodeRoot)
	if err != nil {
		panic(err)
	}
	relativePath, err := filepath.Rel(absolutePath, codeFile)
	if err != nil {
		panic(err)
	}

	return FragmentFile{
		CodePath:      relativePath,
		FragmentName:  fragmentName,
		Configuration: config,
	}
}

// Writes the given text to the file.
//
// Creates the file if not exists and overwrites if exists.
func (f FragmentFile) Write(text string) {
	byteStr := []byte(text)
	filePath := f.absolutePath()
	err := os.WriteFile(filePath, byteStr, os.FileMode(files.WritePermission))
	if err != nil {
		panic(err)
	}
}

// Content reads content of the file.
//
// Returns contents of the file as a list of strings, or returns an error if it doesn't exist.
func (f FragmentFile) Content() ([]string, error) {
	path := f.absolutePath()
	isPathFileExits, err := files.IsFileExist(path)

	if err != nil {
		return nil, err
	}

	if isPathFileExits {
		return files.ReadFile(path)
	}

	return nil, fmt.Errorf("file %s doesn't exist", path)
}

// Returns string representation of FragmentFile.
func (f FragmentFile) String() string {
	return f.absolutePath()
}

// Obtains the absolute path to this fragment file.
func (f FragmentFile) absolutePath() string {
	fileExtension := filepath.Ext(f.CodePath)
	fragmentsAbsDir, err := filepath.Abs(f.Configuration.FragmentsDir)
	if err != nil {
		panic(err)
	}

	if f.FragmentName == DefaultFragmentName {
		return filepath.Join(fragmentsAbsDir, f.CodePath)
	}

	withoutExtension := strings.TrimSuffix(f.CodePath, fileExtension)
	filename := fmt.Sprintf("%s-%s", withoutExtension, f.fragmentHash())

	return filepath.Join(fragmentsAbsDir, filename+fileExtension)
}

// Calculates and returns a hash string for FragmentFile.
//
// Since fragments which have the same name unite into one fragment with multiple partitions,
// the name of a fragment is unique.
func (f FragmentFile) fragmentHash() string {
	hash := sha256.New()
	hash.Write([]byte(f.FragmentName))

	return hex.EncodeToString(hash.Sum(nil))[:8]
}
