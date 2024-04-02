package embedding_instruction

import (
	"embed-code/embed-code-go/configuration"
	"embed-code/embed-code-go/fragmentation"
	"embed-code/embed-code-go/indent"
	"fmt"
)

// EmbeddingInstruction represents the code fragment to embed into a Markdown file.
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

// NewEmbeddingInstruction creates a new EmbeddingInstruction.
func NewEmbeddingInstruction(values map[string]string, configuration configuration.Configuration) EmbeddingInstruction {
	codeFile := values["file"]
	fragment := values["fragment"]
	startValue := values["start"]
	endValue := values["end"]

	if fragment != "" && (startValue != "" || endValue != "") {
		panic("<embed-code> must NOT specify both a fragment name and start/end patterns.")
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
		Configuration: configuration,
	}
}

// FromXML reads the instruction from the '<embed-code>' XML tag.
func FromXML(line string, configuration configuration.Configuration) EmbeddingInstruction {
	fields := ParseXmlLine(line)
	return NewEmbeddingInstruction(fields, configuration)
}

//
// Public methods
//

// Content reads the specified fragment from the code.
func (e EmbeddingInstruction) Content() []string {
	fragmentName := e.Fragment
	if fragmentName == "" {
		fragmentName = fragmentation.DefaultFragmentName
	}
	file := fragmentation.FragmentFile{
		CodeFile:      e.CodeFile,
		FragmentName:  fragmentName,
		Configuration: e.Configuration,
	}
	if e.StartPattern != nil || e.EndPattern != nil {
		return e.matchingLines(file.Content())
	}
	return file.Content()
}

func (e *EmbeddingInstruction) String() string {
	return fmt.Sprintf("EmbeddingInstruction[file=`%s`, fragment=`%s`, start=`%s`, end=`%s`]",
		e.CodeFile, e.Fragment, e.StartPattern, e.EndPattern)
}

//
// Private methods
//

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
	panic(fmt.Sprintf("There is no line matching `%s`.", pattern))
}
