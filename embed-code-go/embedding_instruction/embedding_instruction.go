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

package embedding_instruction

import (
	"fmt"

	"embed-code/embed-code-go/configuration"
	"embed-code/embed-code-go/fragmentation"
	"embed-code/embed-code-go/indent"
)

// Specifies the code fragment to embed into a Markdown file, and the embedding parameters.
//
// Takes form of an XML processing instruction <embed-code file="..." fragment="..."/>.
//
// CodeFile — a path to a code file to embed. The path is relative to Configuration.CodeRoot dir.
//
// Fragment — name of the particular fragment in the code. If Fragment is empty,
// the whole file is embedded.
//
// StartPattern — an optional glob-like pattern. If specified, lines before the matching one
// are excluded.
//
// EndPattern — an optional glob-like pattern. If specified, lines after the matching one
// are excluded.
//
// Configuration — a Configuration with all embed-code settings.
type EmbeddingInstruction struct {
	CodeFile      string
	Fragment      string
	StartPattern  *Pattern
	EndPattern    *Pattern
	Configuration configuration.Configuration
}

//
// Initializers
//

// Creates a new EmbeddingInstruction based on provided attributes and configuration.
//
// attributes — a map with string-typed both keys and values. Possible keys are:
//   - file — a mandatory relative path to the file with the code;
//   - fragment — an optional name of the particular fragment in the code. If no fragment
//     is specified, the whole file is embedded;
//   - start — an optional glob-like pattern. If specified, lines before the matching one
//     are excluded;
//   - end — an optional glob-like pattern. If specified, lines after the matching one are excluded.
//
// config — a Configuration with all embed-code settings.
//
// Returns an error if the instruction is wrong.
func NewEmbeddingInstruction(
	attributes map[string]string, config configuration.Configuration) (EmbeddingInstruction, error) {
	codeFile := attributes["file"]
	fragment := attributes["fragment"]
	startValue := attributes["start"]
	endValue := attributes["end"]

	if fragment != "" && (startValue != "" || endValue != "") {
		return EmbeddingInstruction{},
			fmt.Errorf("<embed-code> must NOT specify both a fragment name and start/end patterns")
	}
	var end *Pattern
	var start *Pattern

	if startValue != "" {
		startPattern := NewPattern(startValue)
		start = &startPattern
	}
	if endValue != "" {
		endPattern := NewPattern(endValue)
		end = &endPattern
	}

	return EmbeddingInstruction{
		CodeFile:      codeFile,
		Fragment:      fragment,
		StartPattern:  start,
		EndPattern:    end,
		Configuration: config,
	}, nil
}

// Reads the instruction from the '<embed-code>' XML tag and creates new EmbeddingInstruction.
//
// line — a line which contains '<embed-code>' XML tag.
// For example: '<embed-code file="org/example/Hello.java" fragment="Hello class"/>'.
// The line can also contain closing tag:
// '<embed-code file=\"org/example/Hello.java\" fragment=\"Hello class\"></embed-code>'.
// The following parameters are currently supported:
//   - file — a mandatory relative path to the file with the code;
//   - fragment — an optional name of the particular fragment in the code. If no fragment
//     is specified, the whole file is embedded;
//   - start — an optional glob-like pattern. If specified, lines before the matching one
//     are excluded;
//   - end — an optional glob-like pattern. If specified, lines after the matching one are excluded.
//
// config — a Configuration with all embed-code settings.
//
// Returns an error if the paring of XML instruction failed.
func FromXML(line string, config configuration.Configuration) (EmbeddingInstruction, error) {
	fields, err := ParseXMLLine(line)
	if err != nil {
		return EmbeddingInstruction{}, err
	}

	return NewEmbeddingInstruction(fields, config)
}

//
// Public methods
//

// Reads and returns the lines for specified fragment from the code.
//
// Returns an error if there was an error during reading the content.
func (e EmbeddingInstruction) Content() ([]string, error) {
	fragmentName := e.Fragment
	if fragmentName == "" {
		fragmentName = fragmentation.DefaultFragmentName
	}
	file := fragmentation.FragmentFile{
		CodePath:      e.CodeFile,
		FragmentName:  fragmentName,
		Configuration: e.Configuration,
	}
	if e.StartPattern != nil || e.EndPattern != nil {
		fileContent, err := file.Content()
		if err != nil {
			return nil, err
		}

		return e.matchingLines(fileContent), nil
	}

	return file.Content()
}

// Returns string representation of EmbeddingInstruction.
func (e EmbeddingInstruction) String() string {
	return fmt.Sprintf("EmbeddingInstruction[file=`%s`, fragment=`%s`, start=`%s`, end=`%s`]",
		e.CodeFile, e.Fragment, e.StartPattern, e.EndPattern)
}

//
// Private methods
//

// Filters and returns a subset of input lines based on start and end patterns.
//
// lines — a list of strings representing the input lines.
func (e EmbeddingInstruction) matchingLines(lines []string) []string {
	startPosition := 0
	if e.StartPattern != nil {
		startPosition = e.matchGlob(e.StartPattern, lines, 0)
	}
	endPosition := len(lines) - 1
	if e.EndPattern != nil {
		endPosition = e.matchGlob(e.EndPattern, lines, startPosition)
	}
	requiredLines := lines[startPosition : endPosition+1]
	indentation := indent.MaxCommonIndentation(requiredLines)

	return indent.CutIndent(requiredLines, indentation)
}

// Returns the index of a first line that matches given pattern.
//
// pattern — a pattern to search in lines for.
//
// lines — a list of lines to search in.
//
// startFrom — an index from which to start searching.
func (e EmbeddingInstruction) matchGlob(pattern *Pattern, lines []string, startFrom int) int {
	lineCount := len(lines)
	resultLine := startFrom
	for resultLine < lineCount {
		line := lines[resultLine]
		if pattern.Match(line) {
			return resultLine
		}
		resultLine++
	}
	panic(fmt.Sprintf("there is no line matching `%s`", pattern))
}
