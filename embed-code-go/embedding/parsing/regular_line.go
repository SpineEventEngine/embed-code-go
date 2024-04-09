package parsing

import (
	"embed-code/embed-code-go/configuration"
)

// Represents a regular line of a markdown.
type RegularLine struct{}

//
// Public methods
//

// Reports whether the current line is a regular line.
//
// Every line can be considered as a regular line.
//
// context — a context of the parsing process.
func (r RegularLine) Recognize(context ParsingContext) bool {
	return true
}

// Adds the current line from the parsing context to the result
// and moves to the next line in the context.
func (r RegularLine) Accept(context *ParsingContext, config configuration.Configuration) {
	line := context.CurrentLine()
	context.result = append(context.result, line)
	context.ToNextLine()
}
