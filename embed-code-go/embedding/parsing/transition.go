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

const Statement = "<embed-code"

// Transition represent a transition of a single line in the parsing process.
type Transition interface {

	// Accept updates the parsing context based on the transition.
	//
	// Switches context to the next line if EOF is not reached yet.
	//
	// context — a context of the parsing process.
	//
	// config — a configuration of the embedding.
	//
	// Returns an error if any occurs.
	Accept(context *Context, config configuration.Configuration) error

	// Recognize reports whether the current line satisfies the transition.
	//
	// context — a context of the parsing process.
	Recognize(context Context) bool

	// State returns a state name for this transition.
	State() string
}
