package embedding

import "strings"

// Finish
type Finish struct{}

func (f Finish) recognize(context ParsingContext) {
	context.reachedEOF()
}

func (f Finish) accept(_, _ interface{}) {
}

// Code sample line
type CodeSampleLine struct{}

func (c *CodeSampleLine) recognize(context ParsingContext) bool {
	return !context.reachedEOF() && context.codeFenceStarted
}

func (c *CodeSampleLine) accept(context ParsingContext, _ interface{}) {
	context.toNextLine()
}

// Code fence end
type CodeFenceEnd struct{}

func (c *CodeFenceEnd) recognize(context ParsingContext) bool {
	// Assuming context is of type `interface{}`.
	// Implement your logic here.
	if !context.reachedEOF() {
		indentation := strings.Repeat(" ", context.codeFenceIndentation)
		return context.codeFenceStarted && strings.HasPrefix(context.currentLine(), indentation+"```")
	}
	return false
}

func (c *CodeFenceEnd) accept(context ParsingContext, _ interface{}) {
	// Assuming the two arguments are of type `interface{}`.
	// Implement your logic here.
	line := context.currentLine
	renderSample(context)
	context.result = append(context.result, line)
	context.embedding = nil
	context.codeFenceStarted = false
	context.codeFenceIndentation = 0
	context.toNextLine()
}

func renderSample(context ParsingContext) {
	for _, line := range context.embedding.content {
		indentation := strings.Repeat(" ", context.codeFenceIndentation)
		context.result = append(context.result, indentation+line)
	}
}
