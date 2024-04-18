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
	"crypto/sha1"
	"embed-code/embed-code-go/configuration"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// A file storing a single fragment from the file.
//
// CodeFile — a relative path to a code file. The path is relative to Configuration.CodeRoot.
//
// FragmentName — a name of the fragment in the code file.
//
// Configuration — a configuration for embedding.
type FragmentFile struct {
	CodeFile      string
	FragmentName  string
	Configuration configuration.Configuration
}

// Iniitalizers

// Composes a FragmentFile for the given fragment in given codeFile.
//
// codeFile — an absolute path to a code file.
//
// fragmentName — a name of the fragment in the code file.
//
// configuration — configuration for embedding.
//
// Returns composed fragment.
func NewFragmentFileFromAbsolute(
	codeFile string,
	fragmentName string,
	configuration configuration.Configuration,
) FragmentFile {

	absoluteCodeRoot, err := filepath.Abs(configuration.CodeRoot)
	if err != nil {
		panic(err)
	}
	relativeCodeFile, err := filepath.Rel(absoluteCodeRoot, codeFile)
	if err != nil {
		panic(err)
	}

	return FragmentFile{
		CodeFile:      relativeCodeFile,
		FragmentName:  fragmentName,
		Configuration: configuration,
	}
}

//
// Public methods
//

// Writes text to the file.
//
// Creates the file if not exists and overwrites if exists.
func (fragmentFile FragmentFile) Write(text string) {
	byteStr := []byte(text)
	filePath := fragmentFile.absolutePath()
	os.WriteFile(filePath, byteStr, 0777)
}

// Reads content of the file.
//
// Returns contents of the file as a list of strings, or raises an error if it doesn't exists.
func (fragmentFile FragmentFile) Content() []string {
	path := fragmentFile.absolutePath()
	isPathFileExits, err := IsFileExists(path)
	if isPathFileExits {
		return ReadLines(path)
	} else if err != nil {
		panic(err)
	} else {
		panic(fmt.Sprintf("file %s doesn't exist", path))
	}
}

// Returns string representation of FragmentFile.
func (fragmentFile FragmentFile) String() string {
	return fragmentFile.absolutePath()
}

//
// Private methods
//

// Obtains the absolute path to this fragment file.
func (fragmentFile FragmentFile) absolutePath() string {

	fileExtension := filepath.Ext(fragmentFile.CodeFile)
	fragmentsAbsDir, err := filepath.Abs(fragmentFile.Configuration.FragmentsDir)
	if err != nil {
		panic(err)
	}

	if fragmentFile.FragmentName == DefaultFragmentName {
		return filepath.Join(fragmentsAbsDir, fragmentFile.CodeFile)
	} else {
		withoutExtension := strings.TrimSuffix(fragmentFile.CodeFile, fileExtension)
		filename := fmt.Sprintf("%s-%s", withoutExtension, fragmentFile.calculateFragmentHash())
		return filepath.Join(fragmentsAbsDir, filename+fileExtension)
	}
}

// Calculates and returns a hash string for FragmentFile.
//
// Since fragments which have the same name unite into one
// fragment with multiple partitions, the name of a fragment is unique.
func (fragmentFile FragmentFile) calculateFragmentHash() string {
	hash := sha1.New()
	hash.Write([]byte(fragmentFile.FragmentName))
	sha1_hash := hex.EncodeToString(hash.Sum(nil))[:8]
	return sha1_hash
}
