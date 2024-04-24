// Copyright 2024, TeamDev. All rights reserved.
//
// Redistribution and use in source and/or binary forms, with or without
// modification, must retain the above copyright notice and the following
// disclaimer.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

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
//
// context — a context of the parsing process.
//
// config — a configuration of the embedding.
func (r RegularLine) Accept(context *ParsingContext, config configuration.Configuration) {
	line := context.CurrentLine()
	context.result = append(context.result, line)
	context.ToNextLine()
}
