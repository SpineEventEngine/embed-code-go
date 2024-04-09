package embedding

import (
	"embed-code/embed-code-go/configuration"
	"embed-code/embed-code-go/embedding_instruction"
	"fmt"
	"strings"
)

// Represent a transition of a single line in the parsing process.
type Transition interface {

	// Updates the parsing context based on the transition.
	accept(context *ParsingContext, config configuration.Configuration)

	// Reports whether the current line satisfies the transition.
	recognize(context ParsingContext) bool
}

// Represents the end of the file.
type Finish struct{}

// Reports whether the current line satisfies the transition.
func (f Finish) recognize(context ParsingContext) bool {
	return context.ReachedEOF()
}

// Accepts the end of the file.
func (f Finish) accept(context *ParsingContext, config configuration.Configuration) {
}

// Represents a line of a code sample.
type CodeSampleLine struct{}

// Reports whether the current line is a code sample line.
//
// If codeFenceStarted is true and it's not the end of the file,
// the line is a code sample line.
func (c CodeSampleLine) recognize(context ParsingContext) bool {
	return !context.ReachedEOF() && context.codeFenceStarted
}

// Moves to the next line.
func (c CodeSampleLine) accept(context *ParsingContext, config configuration.Configuration) {
	context.ToNextLine()
}

// Represents the end of a code fence.
type CodeFenceEnd struct{}

// Reports whether the current line is the end of a code fence.
//
// The line is a code fence end if:
//   - the end is not reached;
//   - the code fence has started;
//   - the current line starts with the appropriate indentation and "```"
func (c CodeFenceEnd) recognize(context ParsingContext) bool {
	if !context.ReachedEOF() {
		indentation := strings.Repeat(" ", context.codeFenceIndentation)
		return context.codeFenceStarted && strings.HasPrefix(context.CurrentLine(), indentation+"```")
	}
	return false
}

// Processes the end of a code fence by adding the current line to the result,
// resetting certain context variables, and moving to the next line.
func (c CodeFenceEnd) accept(context *ParsingContext, config configuration.Configuration) {
	line := context.CurrentLine()
	renderSample(context)
	context.result = append(context.result, line)
	context.SetEmbedding(nil)
	context.codeFenceStarted = false
	context.codeFenceIndentation = 0
	context.ToNextLine()
}

func renderSample(context *ParsingContext) {
	for _, line := range context.embedding.Content() {
		indentation := strings.Repeat(" ", context.codeFenceIndentation)
		context.result = append(context.result, indentation+line)
	}
}

// Represents the start of a code fence.
type CodeFenceStart struct{}

// Reports whether the current line is the start of a code fence.
//
// The line is a code fence start if the end is not reached and the current line starts with "```".
func (c CodeFenceStart) recognize(context ParsingContext) bool {
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
func (c CodeFenceStart) accept(context *ParsingContext, config configuration.Configuration) {
	line := context.CurrentLine()
	context.result = append(context.result, line)
	context.codeFenceStarted = true
	leadingSpaces := len(line) - len(strings.TrimLeft(line, " "))
	context.codeFenceIndentation = leadingSpaces
	context.ToNextLine()
}

// Represents a blank line of a markdown.
type BlankLine struct{}

// Reports whether the current line is a blank line.
//
// Checks if the current line is empty and not part of a code fence,
// and if there is an embedding. If these conditions are met, it returns true.
// Otherwise, it returns false.
func (b BlankLine) recognize(context ParsingContext) bool {
	if !context.ReachedEOF() && strings.TrimSpace(context.CurrentLine()) == "" {
		return !context.codeFenceStarted && context.embedding != nil
	}
	return false
}

// Processes a blank line of a markdown.
//
// Appends the current line of the context to the result, and moves to the next line.
func (b BlankLine) accept(context *ParsingContext, config configuration.Configuration) {
	line := context.CurrentLine()
	context.result = append(context.result, line)
	context.ToNextLine()
}

// Represents a regular line of a markdown.
type RegularLine struct{}

// Reports whether the current line is a regular line.
//
// Every line can be considered as a regular line.
func (r RegularLine) recognize(context ParsingContext) bool {
	return true
}

// Adds the current line from the parsing context to the result
// and moves to the next line in the context.
func (r RegularLine) accept(context *ParsingContext, config configuration.Configuration) {
	line := context.CurrentLine()
	context.result = append(context.result, line)
	context.ToNextLine()
}

// Represents an embedding instruction token of a markdown.
type EmbedInstructionToken struct{}

// Reports whether the current line in the parsing context starts with "<embed-code",
// and if there is no ongoing embedding and the end of the file is not reached, it returns true.
// Otherwise, it returns false.
func (e EmbedInstructionToken) recognize(context ParsingContext) bool {
	line := context.CurrentLine()
	isStatement := strings.HasPrefix(strings.TrimSpace(line), Statement)
	if context.embedding == nil && !context.ReachedEOF() && isStatement {
		return true
	}
	return false
}

func (e EmbedInstructionToken) accept(context *ParsingContext, config configuration.Configuration) {
	instructionBody := []string{}
	for !context.ReachedEOF() {
		instructionBody = append(instructionBody, context.CurrentLine())
		instruction := embedding_instruction.FromXML(strings.Join(instructionBody, ""), config)
		context.SetEmbedding(&instruction)
		context.result = append(context.result, context.CurrentLine())
		context.ToNextLine()
		if context.embedding != nil {
			break
		}
	}
	if context.embedding == nil {
		panic(fmt.Sprintf("Failed to parse an embedding instruction. Context: %v", context))
	}
}
