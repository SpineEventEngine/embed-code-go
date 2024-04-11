package parsing

import (
	"embed-code/embed-code-go/configuration"
	"strings"
)

// Represents the start of a code fence.
type CodeFenceStart struct{}

//
// Public methods
//

// Reports whether the current line is the start of a code fence.
//
// The line is a code fence start if the end is not reached and the current line starts with "```".
//
// context — a context of the parsing process.
func (c CodeFenceStart) Recognize(context ParsingContext) bool {
	if !context.ReachedEOF() {
		return strings.HasPrefix(strings.TrimSpace(context.CurrentLine()), "```")
	}
	return false
}

// Processes the start of a code fence.
//
// Appends the current line from the parsing context to the result,
// sets a flag to indicate that a code fence has started,
// calculates the indentation level of the code fence, and moves to the next line in the context.
//
// context — a context of the parsing process.
//
// config — a configuration of the embedding.
func (c CodeFenceStart) Accept(context *ParsingContext, config configuration.Configuration) {
	line := context.CurrentLine()
	context.result = append(context.result, line)
	context.codeFenceStarted = true
	leadingSpaces := len(line) - len(strings.TrimLeft(line, " "))
	context.codeFenceIndentation = leadingSpaces
	context.ToNextLine()
}
