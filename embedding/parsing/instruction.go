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
	"log/slog"
	"strings"

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

	if err = validateExclusiveAttributes(fragment, startValue, endValue, lineValue); err != nil {
		return Instruction{}, err
	}

	startPattern, err := patternFromValue("start", startValue)
	if err != nil {
		return Instruction{}, err
	}
	endPattern, err := patternFromValue("end", endValue)
	if err != nil {
		return Instruction{}, err
	}
	linePattern, err := patternFromValue("line", lineValue)
	if err != nil {
		return Instruction{}, err
	}

	return Instruction{
		CodeFile:      codeFile,
		Fragment:      fragment,
		StartPattern:  startPattern,
		EndPattern:    endPattern,
		LinePattern:   linePattern,
		CommentMode:   commentMode,
		Configuration: config,
	}, nil
}

// validateExclusiveAttributes reports mutually exclusive instruction attributes.
func validateExclusiveAttributes(fragment string, start string, end string, line string) error {
	if fragment != "" && (start != "" || end != "" || line != "") {
		return fmt.Errorf(
			"<embed-code> must NOT specify both a fragment name and start/end/line patterns",
		)
	}
	if line != "" && (start != "" || end != "") {
		return fmt.Errorf(
			"<embed-code> must NOT specify both a line pattern and start/end patterns",
		)
	}

	return nil
}

// patternFromValue creates a Pattern pointer for a non-empty attribute value.
func patternFromValue(attribute string, value string) (*Pattern, error) {
	if value == "" {
		return nil, nil
	}
	pattern, err := NewPattern(value)
	if err != nil {
		return nil, fmt.Errorf("invalid %s pattern `%s`: %w", attribute, value, err)
	}

	return &pattern, nil
}

// Content reads and returns the lines for specified fragment from the code.
//
// Returns an error if there was an error during reading the content.
func (e Instruction) Content() ([]string, error) {
	fileContent, err := fragmentation.ResolveContent(e.CodeFile, e.Fragment, e.Configuration)
	if err != nil {
		return nil, err
	}
	codeFileReference, referenceErr := fragmentation.ResolveCodeFileReference(
		e.CodeFile,
		e.Configuration,
	)
	if e.StartPattern != nil || e.EndPattern != nil || e.LinePattern != nil {
		if referenceErr != nil {
			return nil, referenceErr
		}
		fileContent, err = e.matchingLines(fileContent, codeFileReference)
		if err != nil {
			return nil, err
		}
	}
	if referenceErr == nil {
		slog.Info(e.contentLogMessage(codeFileReference))
	}

	return commentfilter.Filter(
		fileContent,
		e.CodeFile,
		e.CommentMode,
		e.DocumentationFile,
		e.DocumentationLine,
	), nil
}

// contentLogMessage describes the source content selected by this instruction.
func (e Instruction) contentLogMessage(codeFileReference string) string {
	switch {
	case e.Fragment != "":
		return fmt.Sprintf("Extracted fragment `%s` from `%s`.",
			e.Fragment, codeFileReference)
	case e.LinePattern != nil:
		return fmt.Sprintf("Extracted line-pattern embedding from `%s` using %s.",
			codeFileReference, patternLabel("line", e.LinePattern))
	case e.StartPattern != nil || e.EndPattern != nil:
		return fmt.Sprintf("Extracted start/end-pattern embedding from `%s` using %s.",
			codeFileReference, rangePatternLabel(e.StartPattern, e.EndPattern))
	default:
		return fmt.Sprintf("Extracted source file `%s`.", codeFileReference)
	}
}

// rangePatternLabel formats the start and end patterns set on an instruction.
func rangePatternLabel(start *Pattern, end *Pattern) string {
	var labels []string
	if start != nil {
		labels = append(labels, patternLabel("start", start))
	}
	if end != nil {
		labels = append(labels, patternLabel("end", end))
	}

	return strings.Join(labels, " and ")
}

// patternLabel formats a pattern for human-readable logs.
func patternLabel(kind string, pattern *Pattern) string {
	return fmt.Sprintf("%s pattern `%s`", kind, pattern.sourceGlob)
}

// String returns a string representation of Instruction.
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
		if e.LinePattern.HasLineSeparator() {
			startPosition, endPosition, err := e.matchLineSequence(
				e.LinePattern, lines, 0, "line", codeFileReference,
			)
			if err != nil {
				return nil, err
			}
			requiredLines := lines[startPosition : endPosition+1]
			indentation := indent.MaxCommonIndentation(requiredLines)

			return indent.CutIndent(requiredLines, indentation), nil
		}
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
	if pattern.HasLineSeparator() {
		start, end, err := e.matchLineSequence(
			pattern, lines, startFrom, kind, codeFileReference,
		)
		if err != nil {
			return 0, err
		}
		if kind == "end" {
			return end, nil
		}

		return start, nil
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

// matchLineSequence returns the first line range matching the pattern or a not-found error.
func (e Instruction) matchLineSequence(pattern *Pattern, lines []string, startFrom int,
	kind string, codeFileReference string) (int, int, error) {
	start, end, found, err := matchLineSequence(pattern, lines, startFrom)
	if err != nil {
		return 0, 0, err
	}
	if found {
		return start, end, nil
	}

	return 0, 0, PatternNotFoundError{
		Line:              e.DocumentationLine,
		CodeFileReference: codeFileReference,
		Kind:              kind,
		Pattern:           pattern,
	}
}

// matchLineSequence returns the first source-line range matching an escaped-line pattern.
func matchLineSequence(pattern *Pattern, lines []string, startFrom int) (int, int, bool, error) {
	patterns, err := pattern.lineSequencePatterns()
	if err != nil {
		return 0, 0, false, err
	}
	lineCount := len(patterns)
	lastStart := len(lines) - lineCount
	for start := startFrom; start <= lastStart; start++ {
		end := start + lineCount
		if matchLineSequencePatterns(patterns, lines[start:end]) {
			return start, end - 1, true, nil
		}
	}

	return 0, 0, false, nil
}
