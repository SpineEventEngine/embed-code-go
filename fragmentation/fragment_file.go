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

package fragmentation

import (
	"crypto/sha256"
	_type "embed-code/embed-code-go/type"
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
// CodePath — a relative path to a code file. The path is relative to the corresponding code root,
// and starts with the code root name if it's provided.
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
// codeRoot - a _type.NamedPath to the code root.
//
// fragmentName — a name of the fragment in the code file.
//
// configuration — configuration for embedding.
//
// Returns composed fragment.
func NewFragmentFileFromAbsolute(
	codeFile string,
	codeRoot _type.NamedPath,
	fragmentName string,
	config config.Configuration,
) FragmentFile {
	absolutePath, err := filepath.Abs(codeRoot.Path)
	if err != nil {
		panic(err)
	}

	relativePath, err := filepath.Rel(absolutePath, codeFile)
	if err != nil {
		panic(err)
	}

	if strings.TrimSpace(codeRoot.Name) != "" {
		relativePath = filepath.Join(NamedPathPrefix+codeRoot.Name, relativePath)
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

	if !isPathFileExits {
		codeFileReference, err := f.codeFileReference()
		if err != nil {
			return nil, err
		}

		if f.FragmentName != "" {
			if f.FragmentName == "_default" {
				return nil, fmt.Errorf("code file `%s` not found", codeFileReference)
			}
			return nil, fmt.Errorf(
				"fragment `%s` from code file `%s` not found", f.FragmentName, codeFileReference,
			)
		}
		return nil, fmt.Errorf(
			"code file %s fragment not found",
			codeFileReference,
		)
	}

	return files.ReadFile(path)
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

// Builds a user-facing reference to the original code file for error messages.
//
// If the source file exists, returns its absolute `file://` URL. If the file path uses a named
// code root and the resolved file does not exist, returns both the prefixed path and the expanded
// absolute path.
func (f FragmentFile) codeFileReference() (string, error) {
	originalCodePath, isPrefixed, err := f.originalCodePath()
	if err != nil {
		return "", err
	}
	if originalCodePath == "" {
		return f.CodePath, nil
	}

	exists, err := files.IsFileExist(originalCodePath)
	if err != nil {
		return "", err
	}
	if exists {
		return "file://" + originalCodePath, nil
	}
	if isPrefixed {
		return fmt.Sprintf("%s (%s)", f.CodePath, originalCodePath), nil
	}

	return originalCodePath, nil
}

// Resolves the original source file path from the fragment's code path.
//
// Returns the absolute path to the source file and reports whether the input path used a named
// code-root prefix such as `$runtime/...`.
//
// Returns an error if the path uses a named code root that is not present in the configuration.
func (f FragmentFile) originalCodePath() (string, bool, error) {
	normalizedPath := filepath.ToSlash(filepath.Clean(f.CodePath))

	if strings.HasPrefix(normalizedPath, NamedPathPrefix) {
		withoutPrefix := strings.TrimPrefix(normalizedPath, NamedPathPrefix)
		codeRootName, relativePath, _ := strings.Cut(withoutPrefix, "/")

		for _, codeRoot := range f.Configuration.CodeRoots {
			if strings.TrimSpace(codeRoot.Name) != codeRootName {
				continue
			}

			return filepath.Join(resolvedRootPath(codeRoot.Path), filepath.FromSlash(relativePath)),
				true, nil
		}

		return "", true, fmt.Errorf("code root with name `%s` not found for path `%s`",
			codeRootName, f.CodePath)
	}

	if len(f.Configuration.CodeRoots) == 1 {
		return filepath.Join(
			resolvedRootPath(f.Configuration.CodeRoots[0].Path),
			filepath.FromSlash(normalizedPath),
		), false, nil
	}

	return "", false, nil
}

// Resolves the given path to an absolute path when possible.
//
// If absolute-path resolution fails, returns the original path.
func resolvedRootPath(path string) string {
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return path
	}

	return absolutePath
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
