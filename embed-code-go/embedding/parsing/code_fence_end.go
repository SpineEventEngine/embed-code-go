package parsing

import (
	"embed-code/embed-code-go/configuration"
	"strings"
)

//
// Public methods
//

// Represents the end of a code fence.
type CodeFenceEnd struct{}

// Reports whether the current line is the end of a code fence.
//
// The line is a code fence end if:
//   - the end is not reached;
//   - the code fence has started;
//   - the current line starts with the appropriate indentation and "```"
//
// context — a context of the parsing process.
func (c CodeFenceEnd) Recognize(context ParsingContext) bool {
	if !context.ReachedEOF() {
		indentation := strings.Repeat(" ", context.codeFenceIndentation)
		return context.codeFenceStarted && strings.HasPrefix(context.CurrentLine(), indentation+"```")
	}
	return false
}

// Processes the end of a code fence by adding the current line to the result,
// resetting certain context variables, and moving to the next line.
//
// context — a context of the parsing process.
//
// config — a configuration of the embedding.
func (c CodeFenceEnd) Accept(context *ParsingContext, config configuration.Configuration) {
	line := context.CurrentLine()
	renderSample(context)
	context.result = append(context.result, line)
	context.SetEmbedding(nil)
	context.codeFenceStarted = false
	context.codeFenceIndentation = 0
	context.ToNextLine()
}

//
// Private methods
//

func renderSample(context *ParsingContext) {
	for _, line := range context.embedding.Content() {
		indentation := strings.Repeat(" ", context.codeFenceIndentation)
		context.result = append(context.result, indentation+line)
	}
}
