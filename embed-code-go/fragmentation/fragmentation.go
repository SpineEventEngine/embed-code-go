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

// Splits the given file into fragments.
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
	"embed-code/embed-code-go/configuration"
	"fmt"
	"os"
	"path/filepath"
)

// Splits the given file into fragments and writes them into corresponding output files.
//
// SourcesRoot a full path of the root directory of the source code to be embedded.
//
// CodeFile a full path of a file to fragment.
type Fragmentation struct {
	Configuration configuration.Configuration
	SourcesRoot   string
	CodeFile      string
}

//
// Initializers
//

// Builds Fragmentation from given codeFileRelative and config.
//
// codeFileRelative — a relative path to a code file to fragment.
func NewFragmentation(
	codeFileRelative string,
	config configuration.Configuration,
) Fragmentation {

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

	return fragmentation
}

//
// Public methods
//

// Splits the file into fragments.
//
// Returns a refined content of the file to be cut into fragments, and the Fragments.
func (fragmentation Fragmentation) Fragmentize() ([]string, map[string]Fragment, error) {
	fragmentBuilders := make(map[string]*FragmentBuilder)
	var contentToRender []string

	file, err := os.Open(fragmentation.CodeFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text() + "\n"
		contentToRender, fragmentBuilders, err = fragmentation.parseLine(line, contentToRender, fragmentBuilders)
		if err != nil {
			return nil, nil, err
		}
	}

	fragments := make(map[string]Fragment)
	for k, v := range fragmentBuilders {
		fragments[k] = v.Build()
	}
	fragments[DefaultFragment] = CreateDefaultFragment()

	return contentToRender, fragments, nil

}

// Serializes fragments to the output directory.
//
// Keeps the original directory structure relative to the sources root dir.
// That is, `SRC/src/main` becomes `OUT/src/main`.
func (fragmentation Fragmentation) WriteFragments() error {
	allLines, fragments, err := fragmentation.Fragmentize()
	if err != nil {
		return err
	}

	EnsureDirExists(fragmentation.targetDirectory())

	for _, fragment := range fragments {
		fragmentFile := NewFragmentFileFromAbsolute(fragmentation.CodeFile, fragment.Name, fragmentation.Configuration)
		fragment.WriteTo(fragmentFile, allLines, fragmentation.Configuration)
	}
	return nil
}

//
// Static functions
//

// Writes fragments into a corresponding files.
//
// Searches for code files with patterns defined in configuration
// and fragmentizes them with creating fragmented files as a result.
func WriteFragmentFiles(configuration configuration.Configuration) error {
	includes := configuration.CodeIncludes
	codeRoot := configuration.CodeRoot
	for _, rule := range includes {
		pattern := fmt.Sprintf("%s/%s", codeRoot, rule)
		codeFiles, _ := filepath.Glob(pattern)
		for _, codeFile := range codeFiles {
			if ShouldFragmentize(codeFile) {
				fragmentation := NewFragmentation(codeFile, configuration)
				err := fragmentation.WriteFragments()
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// Reports whether the file stored at filePath exists, is file and is textual-encoded.
func ShouldFragmentize(filePath string) bool {
	info, err := os.Stat(filePath)
	if err != nil {
		panic(err)
	}

	isFile := !info.IsDir()
	isValidEncoding := IsEncodedAsText(filePath)

	return isFile && isValidEncoding
}

//
// Private methods
//

// This function parses a single line of input and performs the following actions:
//   - identifies fragment start and end markers within given line;
//   - updates fragmentBuilders based on the markers;
//   - appends non-fragment lines to contentToRender.
//
// Returns updated contentToRender, fragmentBuilders and error if there's any.
func (fragmentation Fragmentation) parseLine(
	line string, contentToRender []string,
	fragmentBuilders map[string]*FragmentBuilder,
) ([]string, map[string]*FragmentBuilder, error) {
	cursor := len(contentToRender)

	fragmentStarts := GetFragmentStarts(line)
	fragmentEnds := GetFragmentEnds(line)

	if len(fragmentStarts) > 0 {
		for _, fragmentName := range fragmentStarts {
			fragment, exists := fragmentBuilders[fragmentName]
			if !exists {
				builder := FragmentBuilder{FileName: fragmentation.CodeFile, Name: fragmentName}
				fragmentBuilders[fragmentName] = &builder
				fragment = fragmentBuilders[fragmentName]
			}
			fragment.AddStartPosition(cursor)
		}
	} else if len(fragmentEnds) > 0 {
		for _, fragmentName := range fragmentEnds {
			if fragment, exists := fragmentBuilders[fragmentName]; exists {
				fragment.AddEndPosition(cursor - 1)
			} else {
				return nil, nil, fmt.Errorf("cannot end the fragment `%s` as it wasn't started", fragmentName)
			}
		}
	} else {
		contentToRender = append(contentToRender, line)
	}
	return contentToRender, fragmentBuilders, nil
}

// This function calculates the target directory path based on the
// Configuration.FragmentsDir and the parent dir of Fragmentation.CodeFile.
func (fragmentation Fragmentation) targetDirectory() string {
	fragmentsDir := fragmentation.Configuration.FragmentsDir
	codeRoot, err := filepath.Abs(fragmentation.Configuration.CodeRoot)
	if err != nil {
		panic(fmt.Sprintf("error calculating absolute path: %v", err))
	}
	relativeFile, err := filepath.Rel(codeRoot, fragmentation.CodeFile)
	if err != nil {
		panic(fmt.Sprintf("error calculating relative path: %v", err))
	}
	subTree := filepath.Dir(relativeFile)
	return filepath.Join(fragmentsDir, subTree)
}
