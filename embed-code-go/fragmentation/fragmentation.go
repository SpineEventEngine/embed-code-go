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

// Package fragmentation splits the given file into fragments.
//
// The fragments are named parts of the file that are surrounded by "fragment brackets":
//
//	class HelloWorld {
//	    // #docfragment main_method
//	    public static void main(String[] argv) {
//	        // #docfragment printing
//	        System.out.println("Hello World");
//	        // #enddocfragment printing
//	    }
//	    // #enddocfragment main_method
//	}
//
// Fragments with the same name may appear multiple times in the same document.
//
// Even if no fragments are defined explicitly, the whole file is always a fragment on its own.
package fragmentation

import (
	"bufio"
	"embed-code/embed-code-go/files"
	"fmt"
	"os"
	"path/filepath"

	config "embed-code/embed-code-go/configuration"

	"github.com/bmatcuk/doublestar/v4"
)

// Fragmentation splits the given file into fragments and writes them into corresponding
// output files.
//
// Configuration — a configuration for embedding.
//
// SourcesRoot — a full path of the root directory of the source code to be embedded.
//
// CodeFile — a full path of a file to fragment.
type Fragmentation struct {
	Configuration    config.Configuration
	SourcesRoot      string
	CodeFile         string
	fragmentBuilders map[string]*FragmentBuilder
}

// NewFragmentation builds Fragmentation from given codeFileRelative and config.
//
// codeFileRelative — a relative path to a code file to fragment.
//
// config — a configuration for embedding.
func NewFragmentation(codeFileRelative string, config config.Configuration) Fragmentation {
	fragmentation := Fragmentation{}

	sourcesRootRelative := config.CodeRoot

	absoluteSourcesRoot, err := filepath.Abs(sourcesRootRelative)
	fragmentation.SourcesRoot = absoluteSourcesRoot
	if err != nil {
		panic(err)
	}

	absoluteCodeFile, err := filepath.Abs(codeFileRelative)
	fragmentation.CodeFile = absoluteCodeFile
	if err != nil {
		panic(err)
	}

	fragmentation.Configuration = config
	fragmentation.fragmentBuilders = make(map[string]*FragmentBuilder)

	return fragmentation
}

// Fragmentize splits the file into fragments.
//
// Returns a refined content of the file to be cut into fragments, and the Fragments.
// Also returns an error if the fragmentation couldn't be done.
func (f Fragmentation) Fragmentize() ([]string, map[string]Fragment, error) {
	var contentToRender []string

	file, err := os.Open(f.CodeFile)
	if err != nil {
		panic(err)
	}

	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		contentToRender, err = f.parseLine(line, contentToRender)
		if err != nil {
			return nil, nil, err
		}
	}

	fragments := make(map[string]Fragment)
	for k, v := range f.fragmentBuilders {
		fragments[k] = v.Build()
	}
	fragments[DefaultFragmentName] = CreateDefaultFragment()

	return contentToRender, fragments, nil
}

// WriteFragments serializes fragments to the output directory.
//
// Keeps the original directory structure relative to the sources root dir.
// That is, `SRC/src/main` becomes `OUT/src/main`.
//
// Returns an error if the fragmentation couldn't be done.
func (f Fragmentation) WriteFragments() error {
	allLines, fragments, err := f.Fragmentize()
	if err != nil {
		return err
	}

	err = files.EnsureDirExists(f.targetDirectory())
	if err != nil {
		return err
	}

	for _, fragment := range fragments {
		fragmentFile := NewFragmentFileFromAbsolute(f.CodeFile, fragment.Name, f.Configuration)
		fragment.WriteTo(fragmentFile, allLines, f.Configuration)
	}

	return nil
}

// WriteFragmentFiles writes each fragment into a corresponding file.
//
// Searches for code files with patterns defined in configuration and fragmentizes them with
// creating fragmented files as a result.
//
// All fragments are placed inside Configuration.FragmentsDir with keeping the original directory
// structure relative to the sources root dir.
// That is, `SRC/src/main` becomes `OUT/src/main`.
//
// configuration — a configuration for embedding.
//
// Returns an error if any of the fragments couldn't be written.
func WriteFragmentFiles(configuration config.Configuration) error {
	includes := configuration.CodeIncludes
	codeRoot := configuration.CodeRoot
	for _, rule := range includes {
		pattern := fmt.Sprintf("%s/%s", codeRoot, rule)
		codeFiles, _ := doublestar.FilepathGlob(pattern)
		for _, codeFile := range codeFiles {
			err := writeFragments(configuration, codeFile)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// CleanFragmentFiles deletes Configuration.FragmentsDir if it exists.
func CleanFragmentFiles(config config.Configuration) {
	exists, err := files.IsDirExist(config.FragmentsDir)
	if err != nil {
		panic(err)
	}
	if !exists {
		panic(fmt.Errorf("%s directory is not exist", config.FragmentsDir))
	}
	if err = os.RemoveAll(config.FragmentsDir); err != nil {
		panic(err)
	}
}

func writeFragments(configuration config.Configuration, codeFile string) error {
	if shouldFragmentize(codeFile) {
		fragmentation := NewFragmentation(codeFile, configuration)
		err := fragmentation.WriteFragments()
		if err != nil {
			return err
		}
	}

	return nil
}

// shouldFragmentize reports whether if the file stored at filePath:
//   - exists
//   - is a file (not a dir)
//   - is textual-encoded.
func shouldFragmentize(filePath string) bool {
	exists, err := files.IsFileExist(filePath)
	if err != nil {
		return false
	}
	if exists {
		return IsEncodedAsText(filePath)
	}

	return false
}

// Parses a single line of input and performs the following actions:
//   - identifies fragment start and end markers within given line;
//   - updates fragmentBuilders based on the markers;
//   - appends non-fragment lines to contentToRender.
//
// line — a string to parse.
//
// contentToRender — a list of strings which meant to be rendered. It fills up here.
//
// fragmentBuilders — a list of FragmentBuilder. This list fills up here and gets start/end
// positions of it's items updated.
//
// Returns updated contentToRender, fragmentBuilders and error if there's any.
func (f Fragmentation) parseLine(line string, contentToRender []string) ([]string, error) {
	cursor := len(contentToRender)

	docFragments := FindDocFragments(line)
	endDocFragments := FindEndDocFragments(line)

	switch {
	case len(docFragments) > 0:
		if err := f.parseDocFragments(docFragments, cursor); err != nil {
			return nil, err
		}
	case len(endDocFragments) > 0:
		if err := f.parseEndDocFragments(endDocFragments, cursor); err != nil {
			return nil, err
		}
	default:
		contentToRender = append(contentToRender, line)
	}

	return contentToRender, nil
}

func (f Fragmentation) parseDocFragments(docFragments []string, cursor int) error {
	for _, fragmentName := range docFragments {
		fragment, exists := f.fragmentBuilders[fragmentName]
		if !exists {
			builder := FragmentBuilder{
				CodeFilePath: f.CodeFile,
				Name:         fragmentName,
			}
			f.fragmentBuilders[fragmentName] = &builder
			fragment = f.fragmentBuilders[fragmentName]
		}
		if err := fragment.AddStartPosition(cursor); err != nil {
			return err
		}
	}

	return nil
}

func (f Fragmentation) parseEndDocFragments(endDocFragments []string, cursor int) error {
	for _, fragmentName := range endDocFragments {
		if fragment, exists := f.fragmentBuilders[fragmentName]; exists {
			if err := fragment.AddEndPosition(cursor - 1); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("cannot end the fragment `%s` of the file `%s` as it wasn't started",
				fragmentName, f.CodeFile)
		}
	}

	return nil
}

// Calculates the target directory path based on the
// Configuration.FragmentsDir and the parent dir of Fragmentation.CodeFile.
func (f Fragmentation) targetDirectory() string {
	fragmentsDir := f.Configuration.FragmentsDir
	codeRoot, err := filepath.Abs(f.Configuration.CodeRoot)
	if err != nil {
		panic(fmt.Sprintf("error calculating absolute path: %v", err))
	}
	relativeFile, err := filepath.Rel(codeRoot, f.CodeFile)
	if err != nil {
		panic(fmt.Sprintf("error calculating relative path: %v", err))
	}
	subTree := filepath.Dir(relativeFile)

	return filepath.Join(fragmentsDir, subTree)
}
