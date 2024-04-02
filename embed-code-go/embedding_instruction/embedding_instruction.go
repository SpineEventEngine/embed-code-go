package embedding_instruction

import (
	"embed-code/embed-code-go/configuration"
	"embed-code/embed-code-go/fragmentation"
	"embed-code/embed-code-go/indent"
	"fmt"
)

// Specifies the code fragment to embed into a Markdown file, and the embedding parameters.
//
// Takes form of an XML processing instruction <embed-code file="..." fragment="..."/>.
//
// CodeFile — a path to a code file to embed. The path is relative to
// Configuration.CodeRoot dir.
//
// Fragment — name of the particular fragment in the code. If Fragment is empty,
// the whole file is embedded.
//
// StartPattern — an optional glob-like pattern. If specified, lines before the matching one are excluded.
//
// EndPattern — an optional glob-like pattern. If specified, lines after the matching one are excluded.
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

// Creates a new EmbeddingInstruction based on the provided values and configuration.
//
// values — a map with string-typed both keys and values. Possible keys are:
//   - file — a mandatory relative path to the file with the code;
//   - fragment — an optional name of the particular fragment in the code. If no fragment is specified,
//     the whole file is embedded;
//   - start — an optional glob-like pattern. If specified, lines before the matching one are excluded;
//   - end — an optional glob-like pattern. If specified, lines after the matching one are excluded.
//
// config — a Configuration with all embed-code settings.
func NewEmbeddingInstruction(values map[string]string, config configuration.Configuration) EmbeddingInstruction {
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
		Configuration: config,
	}
}

// Reads the instruction from the '<embed-code>' XML tag and creates new EmbeddingInstruction.
//
// line — a line which contains '<embed-code>' XML tag.
// For example: '<embed-code file="org/example/Hello.java" fragment="Hello class"/>'.
// The line can also contain closing tag: '<embed-code file=\"org/example/Hello.java\" fragment=\"Hello class\"></embed-code>'.
// The following parameters are currently supported:
//   - file — a mandatory relative path to the file with the code;
//   - fragment — an optional name of the particular fragment in the code. If no fragment is specified,
//     the whole file is embedded;
//   - start — an optional glob-like pattern. If specified, lines before the matching one are excluded;
//   - end — an optional glob-like pattern. If specified, lines after the matching one are excluded.
//
// config — a Configuration with all embed-code settings.
func FromXML(line string, config configuration.Configuration) EmbeddingInstruction {
	fields := ParseXmlLine(line)
	return NewEmbeddingInstruction(fields, config)
}

//
// Public methods
//

// Reads and returns the specified fragment from the code.
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

// Returns the index of a line that matches given pattern.
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
	panic(fmt.Sprintf("There is no line matching `%s`.", pattern))
}
