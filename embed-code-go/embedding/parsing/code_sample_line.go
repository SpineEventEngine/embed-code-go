package parsing

import (
	"embed-code/embed-code-go/configuration"
)

// Represents a line of a code sample.
type CodeSampleLine struct{}

//
// Public methods
//

// Reports whether the current line is a code sample line.
//
// If codeFenceStarted is true and it's not the end of the file,
// the line is a code sample line.
//
// context — a context of the parsing process.
func (c CodeSampleLine) Recognize(context ParsingContext) bool {
	return !context.ReachedEOF() && context.codeFenceStarted
}

// Moves to the next line.
//
// context — a context of the parsing process.
//
// config — a configuration of the embedding.
func (c CodeSampleLine) Accept(context *ParsingContext, config configuration.Configuration) {
	context.ToNextLine()
}
