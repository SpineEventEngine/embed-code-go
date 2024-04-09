package embedding

import (
	"embed-code/embed-code-go/configuration"
	"embed-code/embed-code-go/embedding_instruction"
	"fmt"
	"strings"
)

type Transition interface {
	accept(context *ParsingContext, config configuration.Configuration)
	recognize(context ParsingContext) bool
}

//
// Finish
//

type Finish struct{}

func (f Finish) recognize(context ParsingContext) bool {
	return context.ReachedEOF()
}

func (f Finish) accept(context *ParsingContext, config configuration.Configuration) {
}

//
// Code sample line
//

type CodeSampleLine struct{}

func (c CodeSampleLine) recognize(context ParsingContext) bool {
	return !context.ReachedEOF() && context.codeFenceStarted
}

func (c CodeSampleLine) accept(context *ParsingContext, config configuration.Configuration) {
	context.ToNextLine()
}

//
// Code fence end
//

type CodeFenceEnd struct{}

func (c CodeFenceEnd) recognize(context ParsingContext) bool {
	// Assuming context is of type `interface{}`.
	// Implement your logic here.
	if !context.ReachedEOF() {
		indentation := strings.Repeat(" ", context.codeFenceIndentation)
		return context.codeFenceStarted && strings.HasPrefix(context.CurrentLine(), indentation+"```")
	}
	return false
}

func (c CodeFenceEnd) accept(context *ParsingContext, config configuration.Configuration) {
	// Assuming the two arguments are of type `interface{}`.
	// Implement your logic here.
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

//
// Code fence end
//

type CodeFenceStart struct{}

func (c CodeFenceStart) recognize(context ParsingContext) bool {
	if !context.ReachedEOF() {
		return strings.HasPrefix(strings.TrimSpace(context.CurrentLine()), "```")
	}
	return false
}

func (c CodeFenceStart) accept(context *ParsingContext, config configuration.Configuration) {
	line := context.CurrentLine()
	context.result = append(context.result, line)
	context.codeFenceStarted = true
	leadingSpaces := len(line) - len(strings.TrimLeft(line, " "))
	context.codeFenceIndentation = leadingSpaces
	context.ToNextLine()
}

//
//
//

type BlankLine struct{}

func (b BlankLine) recognize(context ParsingContext) bool {
	if !context.ReachedEOF() && strings.TrimSpace(context.CurrentLine()) == "" {
		return !context.codeFenceStarted && context.embedding != nil
	}
	return false
}

func (b BlankLine) accept(context *ParsingContext, config configuration.Configuration) {
	line := context.CurrentLine()
	context.result = append(context.result, line)
	context.ToNextLine()
}

//
//
//

type RegularLine struct{}

func (r RegularLine) recognize(context ParsingContext) bool {
	return true
}

func (r RegularLine) accept(context *ParsingContext, config configuration.Configuration) {
	line := context.CurrentLine()
	context.result = append(context.result, line)
	context.ToNextLine()
}

//
//
//

type EmbedInstructionToken struct{}

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
