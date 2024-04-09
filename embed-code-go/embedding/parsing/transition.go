package parsing

import (
	"embed-code/embed-code-go/configuration"
)

const Statement = "<embed-code"

// Represent a transition of a single line in the parsing process.
type Transition interface {

	// Updates the parsing context based on the transition.
	Accept(context *ParsingContext, config configuration.Configuration)

	// Reports whether the current line satisfies the transition.
	//
	// context — a context of the parsing process.
	Recognize(context ParsingContext) bool
}
