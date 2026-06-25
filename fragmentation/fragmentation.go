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

// Package fragmentation contains functions for splitting the given file into fragments.
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
	"bytes"
	"fmt"
	"os"
	"path/filepath"
)

// NamedPathPrefix is the prefix before a named code source.
const NamedPathPrefix = "$"

// Fragmentation splits the given file into fragments.
type Fragmentation struct {
	// codeFile is the absolute path of the source file being fragmented.
	codeFile string

	// fragmentBuilders collects fragment partitions by name while the source file is scanned.
	fragmentBuilders map[string]*FragmentBuilder
}

// NewFragmentation builds Fragmentation for a relative or absolute source path.
//
// Parameters:
// codeFile - provides the source file path.
//
// Returns:
// Fragmentation - source file fragmentation context.
// error - when codeFile cannot be made absolute.
func NewFragmentation(codeFile string) (Fragmentation, error) {
	absoluteCodeFile, err := filepath.Abs(codeFile)
	if err != nil {
		return Fragmentation{}, err
	}

	return Fragmentation{
		codeFile:         absoluteCodeFile,
		fragmentBuilders: make(map[string]*FragmentBuilder),
	}, nil
}

// DoFragmentation splits the source file into renderable content and named fragments.
//
// Returns:
// []string - renderable source lines.
// map[string]Fragment - parsed fragments by name.
// error - when the source file cannot be read, decoded, or parsed.
func (f Fragmentation) DoFragmentation() ([]string, map[string]Fragment, error) {
	var contentToRender []string

	content, err := os.ReadFile(f.codeFile)
	if err != nil {
		return nil, nil, err
	}
	if err := validateTextEncoding(content); err != nil {
		return nil, nil, err
	}

	scanner := bufio.NewScanner(bytes.NewReader(content))
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()
		contentToRender, err = f.parseLine(line, contentToRender)
		if err != nil {
			return nil, nil, fmt.Errorf(
				"failed to do fragmentation on file `file://%s:%d`: %w",
				f.codeFile, lineNumber, err,
			)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}

	fragments := make(map[string]Fragment)
	for k, v := range f.fragmentBuilders {
		fragments[k] = v.Build()
	}
	fragments[DefaultFragmentName] = CreateDefaultFragment()

	return contentToRender, fragments, nil
}

// parseLine parses one source line and updates fragment builders or renderable content.
//
// It identifies fragment start and end markers, updates fragment builders,
// and appends non-fragment lines to renderable content.
//
// Parameters:
// line - provides one source line to parse.
// contentToRender - provides accumulated renderable source lines.
//
// Returns:
// []string - updated renderable source lines.
// error - when fragment marker parsing fails.
func (f Fragmentation) parseLine(line string, contentToRender []string) ([]string, error) {
	cursor := len(contentToRender)

	docFragments, startErr := FindDocFragments(line)
	if startErr != nil {
		return nil, startErr
	}
	endDocFragments, endErr := FindEndDocFragments(line)
	if endErr != nil {
		return nil, endErr
	}

	switch {
	case len(docFragments) > 0:
		if err := f.parseStartDocFragments(docFragments, cursor); err != nil {
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

// parseStartDocFragments starts a new partition for each named fragment marker.
//
// It creates fragment builders when necessary.
func (f Fragmentation) parseStartDocFragments(docFragments []string, cursor int) error {
	for _, fragmentName := range docFragments {
		fragment, exists := f.fragmentBuilders[fragmentName]
		if !exists {
			builder := FragmentBuilder{
				CodeFilePath: f.codeFile,
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

// parseEndDocFragments closes the latest partition for each named fragment marker.
//
// It requires a matching fragment builder to have been started earlier.
func (f Fragmentation) parseEndDocFragments(endDocFragments []string, cursor int) error {
	for _, fragmentName := range endDocFragments {
		if fragment, exists := f.fragmentBuilders[fragmentName]; exists {
			if err := fragment.AddEndPosition(cursor - 1); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("cannot end the fragment `%s` of the file `file://%s` as it wasn't started",
				fragmentName, f.codeFile)
		}
	}

	return nil
}
