package parsing

import (
	"embed-code/embed-code-go/configuration"
	"strings"
)

// Represents a blank line of a markdown.
type BlankLine struct{}

// Reports whether the current line is a blank line.
//
// Checks if the current line is empty and not part of a code fence,
// and if there is an embedding. If these conditions are met, it returns true.
// Otherwise, it returns false.
func (b BlankLine) Recognize(context ParsingContext) bool {
	if !context.ReachedEOF() && strings.TrimSpace(context.CurrentLine()) == "" {
		return !context.codeFenceStarted && context.embedding != nil
	}
	return false
}

// Processes a blank line of a markdown.
//
// Appends the current line of the context to the result, and moves to the next line.
func (b BlankLine) Accept(context *ParsingContext, config configuration.Configuration) {
	line := context.CurrentLine()
	context.result = append(context.result, line)
	context.ToNextLine()
}
