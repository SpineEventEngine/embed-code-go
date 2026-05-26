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

package parsing

import (
	"fmt"

	"embed-code/embed-code-go/configuration"
	"embed-code/embed-code-go/embedding/commentfilter"
	"embed-code/embed-code-go/fragmentation"
	"embed-code/embed-code-go/indent"
)

// Instruction specifies the code fragment to embed into a Markdown file, and the
// embedding parameters.
//
// Takes form of an XML processing instruction <embed-code file="..." fragment="..."/>.
//
// CodeFile — a path to a code file to embed. The path is relative to the corresponding code root.
//
// Fragment — name of the particular fragment in the code. If Fragment is empty, the whole file
// is embedded.
//
// StartPattern — an optional glob-like pattern. If specified, lines before the matching one
// are excluded.
//
// EndPattern — an optional glob-like pattern. If specified, lines after the matching one
// are excluded.
//
// LinePattern — an optional glob-like pattern. If specified, only the matching line is embedded.
//
// CommentMode — specifies which comments are retained in the embedded code.
//
// DocumentationFile — a documentation file containing the instruction.
//
// DocumentationLine — a line containing the start of the instruction.
//
// Configuration — a Configuration with all embed-code settings.
type Instruction struct {
	CodeFile          string
	Fragment          string
	StartPattern      *Pattern
	EndPattern        *Pattern
	LinePattern       *Pattern
	CommentMode       commentfilter.Mode
	DocumentationFile string
	DocumentationLine int
	Configuration     configuration.Configuration
}

// PatternNotFoundError reports that a start or end pattern did not match the code file.
type PatternNotFoundError struct {
	Line              int
	CodeFileReference string
	Kind              string
	Pattern           *Pattern
}

// Error returns a user-facing description of an unmatched start or end pattern.
func (e PatternNotFoundError) Error() string {
	pattern := ""
	if e.Pattern != nil {
		pattern = e.Pattern.sourceGlob
	}

	return fmt.Sprintf(
		"no line in code file `%s` matches the %s pattern `%s`",
		e.CodeFileReference,
		e.Kind,
		pattern,
	)
}

// NewInstruction creates an Instruction based on provided attributes and configuration.
//
// attributes — a map with string-typed both keys and values. Possible keys are:
//   - file — a mandatory relative path to the file with the code;
//   - fragment — an optional name of the particular fragment in the code. If no fragment
//     is specified, the whole file is embedded;
//   - start — an optional glob-like pattern. If specified, lines before the matching one
//     are excluded;
//   - end — an optional glob-like pattern. If specified, lines after the matching one are excluded.
//   - line — an optional glob-like pattern. If specified, only the matching line is embedded.
//   - comments — an optional comment filtering mode. If omitted, all comments are retained.
//
// config — a Configuration with all embed-code settings.
//
// Returns an error if the instruction is wrong.
func NewInstruction(
	attributes map[string]string, config configuration.Configuration) (Instruction, error) {
	codeFile := attributes["file"]
	fragment := attributes["fragment"]
	startValue := attributes["start"]
	endValue := attributes["end"]
	lineValue := attributes["line"]
	commentMode, err := commentfilter.ParseMode(attributes["comments"])
	if err != nil {
		return Instruction{}, err
	}

	if fragment != "" && (startValue != "" || endValue != "" || lineValue != "") {
		return Instruction{},
			fmt.Errorf("<embed-code> must NOT specify both a fragment name and start/end/line patterns")
	}
	if lineValue != "" && (startValue != "" || endValue != "") {
		return Instruction{},
			fmt.Errorf("<embed-code> must NOT specify both a line pattern and start/end patterns")
	}
	var end *Pattern
	var line *Pattern
	var start *Pattern

	if startValue != "" {
		startPattern := NewPattern(startValue)
		start = &startPattern
	}
	if endValue != "" {
		endPattern := NewPattern(endValue)
		end = &endPattern
	}
	if lineValue != "" {
		linePattern := NewPattern(lineValue)
		line = &linePattern
	}

	return Instruction{
		CodeFile:      codeFile,
		Fragment:      fragment,
		StartPattern:  start,
		EndPattern:    end,
		LinePattern:   line,
		CommentMode:   commentMode,
		Configuration: config,
	}, nil
}

// Content reads and returns the lines for specified fragment from the code.
//
// Returns an error if there was an error during reading the content.
func (e Instruction) Content() ([]string, error) {
	fileContent, err := fragmentation.ResolveContent(e.CodeFile, e.Fragment, e.Configuration)
	if err != nil {
		return nil, err
	}
	if e.StartPattern != nil || e.EndPattern != nil || e.LinePattern != nil {
		codeFileReference, err := fragmentation.ResolveCodeFileReference(e.CodeFile, e.Configuration)
		if err != nil {
			return nil, err
		}
		fileContent, err = e.matchingLines(fileContent, codeFileReference)
		if err != nil {
			return nil, err
		}
	}

	return commentfilter.Filter(
		fileContent,
		e.CodeFile,
		e.CommentMode,
		e.DocumentationFile,
		e.DocumentationLine,
	), nil
}

// Returns string representation of Instruction.
func (e Instruction) String() string {
	return fmt.Sprintf(
		"EmbeddingInstruction[file=`%s`, fragment=`%s`, start=`%s`, end=`%s`, line=`%s`, comments=`%s`]",
		e.CodeFile, e.Fragment, e.StartPattern, e.EndPattern, e.LinePattern, e.CommentMode,
	)
}

// Filters and returns a subset of input lines based on start, end, or line patterns.
//
// lines — a list of strings representing the input lines.
func (e Instruction) matchingLines(lines []string, codeFileReference string) ([]string, error) {
	if e.LinePattern != nil {
		linePosition, err := e.matchGlob(
			e.LinePattern, lines, 0, "line", codeFileReference,
		)
		if err != nil {
			return nil, err
		}
		requiredLines := []string{lines[linePosition]}
		indentation := indent.MaxCommonIndentation(requiredLines)

		return indent.CutIndent(requiredLines, indentation), nil
	}

	startPosition := 0
	if e.StartPattern != nil {
		var err error
		startPosition, err = e.matchGlob(
			e.StartPattern, lines, 0, "start", codeFileReference,
		)
		if err != nil {
			return nil, err
		}
	}
	endPosition := len(lines) - 1
	if e.EndPattern != nil {
		var err error
		endPosition, err = e.matchGlob(
			e.EndPattern, lines, startPosition, "end", codeFileReference,
		)
		if err != nil {
			return nil, err
		}
	}
	requiredLines := lines[startPosition : endPosition+1]
	indentation := indent.MaxCommonIndentation(requiredLines)

	return indent.CutIndent(requiredLines, indentation), nil
}

// Returns the index of a first line that matches given pattern.
//
// pattern — a pattern to search in lines for.
//
// lines — a list of lines to search in.
//
// startFrom — an index from which to start searching.
func (e Instruction) matchGlob(pattern *Pattern, lines []string, startFrom int,
	kind string, codeFileReference string) (int, error) {
	if kind != "line" && pattern.HasLineSeparator() {
		start, end, found := matchLineSequence(pattern, lines, startFrom)
		if found {
			if kind == "end" {
				return end, nil
			}
			return start, nil
		}
		return 0, PatternNotFoundError{
			Line:              e.DocumentationLine,
			CodeFileReference: codeFileReference,
			Kind:              kind,
			Pattern:           pattern,
		}
	}
	if line, found := matchSingleLine(pattern, lines, startFrom); found {
		return line, nil
	}
	return 0, PatternNotFoundError{
		Line:              e.DocumentationLine,
		CodeFileReference: codeFileReference,
		Kind:              kind,
		Pattern:           pattern,
	}
}

// matchSingleLine returns the first source line matching the pattern.
func matchSingleLine(pattern *Pattern, lines []string, startFrom int) (int, bool) {
	lineCount := len(lines)
	resultLine := startFrom
	for resultLine < lineCount {
		line := lines[resultLine]
		if pattern.Match(line) {
			return resultLine, true
		}
		resultLine++
	}

	return 0, false
}

// matchLineSequence returns the first source-line range matching an escaped-line pattern.
func matchLineSequence(pattern *Pattern, lines []string, startFrom int) (int, int, bool) {
	patternLines, _ := pattern.linePatterns()
	lineCount := len(patternLines)
	lastStart := len(lines) - lineCount
	for start := startFrom; start <= lastStart; start++ {
		end := start + lineCount
		if pattern.MatchLineSequence(lines[start:end]) {
			return start, end - 1, true
		}
	}

	return 0, 0, false
}
