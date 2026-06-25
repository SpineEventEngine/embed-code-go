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

// Instruction specifies the code fragment to embed into a Markdown file.
//
// It is parsed from an XML-like `<embed-code>` instruction such as
// `<embed-code file="..." fragment="..."/>`.
type Instruction struct {
	// CodeFile is the path to the source file relative to its code root.
	CodeFile string

	// Fragment identifies a named fragment; an empty value selects the whole file.
	Fragment string

	// StartPattern excludes lines before its first match when set.
	StartPattern *Pattern

	// EndPattern excludes lines after its first match when set.
	EndPattern *Pattern

	// LinePattern selects only its matching line when set.
	LinePattern *Pattern

	// CommentMode selects which comments are retained in embedded code.
	CommentMode commentfilter.Mode

	// DocumentationFile is the path to the documentation containing the instruction.
	DocumentationFile string

	// DocumentationLine is the line containing the start of the instruction.
	DocumentationLine int

	// Configuration contains the embedding settings.
	Configuration configuration.Configuration
}

// PatternNotFoundError reports that an instruction pattern did not match the code file.
type PatternNotFoundError struct {
	// Line is the documentation line containing the instruction.
	Line int

	// CodeFileReference is the user-facing reference to the searched source file.
	CodeFileReference string

	// Kind identifies the unmatched pattern as start or end.
	Kind string

	// Pattern is the source-line pattern that did not match.
	Pattern *Pattern
}

// Error returns a user-facing description of an unmatched start or end pattern.
//
// Returns formatted pattern error text.
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

// NewInstruction builds an instruction from parsed `<embed-code>` attributes and configuration.
//
// Parameters:
// attributes - provides embed-code tag attributes. Supported keys are:
//   - file - mandatory relative path to the source file;
//   - fragment - optional source fragment name. When omitted, the whole file is embedded;
//   - start - optional glob-like pattern. Matching lines before it are excluded;
//   - end - optional glob-like pattern. Matching lines after it are excluded;
//   - line - optional glob-like pattern. Only the matching line is embedded;
//   - comments - optional comment filtering mode. When omitted, all comments are retained.
//
// config - provides embedding configuration.
//
// Returns:
// Instruction - parsed embedding instruction.
// error - when instruction attributes are invalid.
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

	patterns, err := parseInstructionPatterns(startValue, endValue, lineValue)
	if err != nil {
		return Instruction{}, err
	}

	return Instruction{
		CodeFile:      codeFile,
		Fragment:      fragment,
		StartPattern:  patterns.start,
		EndPattern:    patterns.end,
		LinePattern:   patterns.line,
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

// instructionPatterns holds the optional source-line patterns from instruction attributes.
type instructionPatterns struct {
	// start is the optional start pattern.
	start *Pattern

	// end is the optional end pattern.
	end *Pattern

	// line is the optional single-line pattern.
	line *Pattern
}

// parseInstructionPatterns parses all optional source-line pattern attributes.
func parseInstructionPatterns(start string, end string, line string) (instructionPatterns, error) {
	var patterns instructionPatterns
	if start != "" {
		pattern, err := parseInstructionPattern("start", start)
		if err != nil {
			return instructionPatterns{}, err
		}
		patterns.start = &pattern
	}
	if end != "" {
		pattern, err := parseInstructionPattern("end", end)
		if err != nil {
			return instructionPatterns{}, err
		}
		patterns.end = &pattern
	}
	if line != "" {
		pattern, err := parseInstructionPattern("line", line)
		if err != nil {
			return instructionPatterns{}, err
		}
		patterns.line = &pattern
	}

	return patterns, nil
}

// parseInstructionPattern parses one non-empty source-line pattern attribute.
func parseInstructionPattern(attribute string, value string) (Pattern, error) {
	pattern, err := NewPattern(value)
	if err != nil {
		return Pattern{}, fmt.Errorf("invalid %s pattern `%s`: %w", attribute, value, err)
	}

	return pattern, nil
}

// Content returns source lines selected and filtered by this instruction.
//
// It reads source content for the configured file and fragment before applying
// optional source-line patterns and comment filtering.
//
// Returns:
// []string - selected and filtered source lines.
// error - when source resolution or pattern matching fails.
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
//
// Returns diagnostic instruction text.
func (e Instruction) String() string {
	return fmt.Sprintf(
		"EmbeddingInstruction[file=`%s`, fragment=`%s`, start=`%s`, end=`%s`, line=`%s`, comments=`%s`]",
		e.CodeFile, e.Fragment, e.StartPattern, e.EndPattern, e.LinePattern, e.CommentMode,
	)
}

// matchingLines filters and returns input lines based on start, end, or line patterns.
func (e Instruction) matchingLines(lines []string, codeFileReference string) ([]string, error) {
	var selectedLines []string
	var err error
	if e.LinePattern != nil {
		selectedLines, err = e.matchLinePattern(lines, codeFileReference)
	} else {
		selectedLines, err = e.matchRangePattern(lines, codeFileReference)
	}
	if err != nil {
		return nil, err
	}

	return removeCommonIndent(selectedLines), nil
}

// matchLinePattern returns the source lines matched by the instruction line pattern.
func (e Instruction) matchLinePattern(
	lines []string,
	codeFileReference string,
) ([]string, error) {
	startPosition, endPosition, err := e.matchPattern(
		e.LinePattern, lines, 0, "line", codeFileReference,
	)
	if err != nil {
		return nil, err
	}

	return lines[startPosition : endPosition+1], nil
}

// matchRangePattern returns the source lines matched by the instruction start/end patterns.
func (e Instruction) matchRangePattern(
	lines []string,
	codeFileReference string,
) ([]string, error) {
	startPosition := 0
	if e.StartPattern != nil {
		var err error
		startPosition, _, err = e.matchPattern(
			e.StartPattern, lines, 0, "start", codeFileReference,
		)
		if err != nil {
			return nil, err
		}
	}
	endPosition := len(lines) - 1
	if e.EndPattern != nil {
		var err error
		_, endPosition, err = e.matchPattern(
			e.EndPattern, lines, startPosition, "end", codeFileReference,
		)
		if err != nil {
			return nil, err
		}
	}

	return lines[startPosition : endPosition+1], nil
}

// removeCommonIndent removes shared indentation from the selected source lines.
func removeCommonIndent(lines []string) []string {
	indentation := indent.MaxCommonIndentation(lines)

	return indent.CutIndent(lines, indentation)
}

// matchPattern returns the first source-line range matching pattern.
//
// Parameters:
// pattern - provides the source-line pattern to search for.
// lines - provides source lines to search in.
// startFrom - provides the first index to search.
// kind - identifies the pattern kind for errors.
// codeFileReference - identifies the searched file for errors.
//
// Returns:
// int - inclusive start index.
// int - inclusive end index.
// error - when pattern does not match.
func (e Instruction) matchPattern(
	pattern *Pattern, lines []string, startFrom int, kind string, codeFileReference string,
) (int, int, error) {
	if start, end, found := pattern.FindIn(lines, startFrom); found {
		return start, end, nil
	}

	return 0, 0, PatternNotFoundError{
		Line:              e.DocumentationLine,
		CodeFileReference: codeFileReference,
		Kind:              kind,
		Pattern:           pattern,
	}
}
