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

// NamedPathPrefix the prefix before the named code source.
const NamedPathPrefix = "$"

// Fragmentation splits the given file into fragments.
type Fragmentation struct {
	// codeFile is the absolute path of the source file being fragmented.
	codeFile string

	// fragmentBuilders collects fragment partitions by name while the source file is scanned.
	fragmentBuilders map[string]*FragmentBuilder
}

// NewFragmentation builds Fragmentation for the given code file.
//
// codeFile — a relative or absolute path to a code file to fragment.
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

// DoFragmentation splits the file into fragments.
//
// Returns a refined content of the file to be cut into fragments, and the Fragments.
// Also returns an error if the fragmentation couldn't be done.
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

// Parses a single line of input and performs the following actions:
//   - identifies fragment start and end markers within given line;
//   - updates fragmentBuilders based on the markers;
//   - appends non-fragment lines to contentToRender.
//
// line — a string to parse.
//
// contentToRender — a list of strings which meant to be rendered. It fills up here.
//
// Returns updated contentToRender, and error if there's any.
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

// Iterates through the fragments` starts, creates fragments builders (if necessary), and adds a
// new partition to the fragment.
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

// Iterates through the fragments` ends, creates fragments builders (if necessary), and adds a
// new partition to the fragment.
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
