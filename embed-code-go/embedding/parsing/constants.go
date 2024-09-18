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

// EmbeddingTag is a StartState of a tag where it requires to embed the code.
const EmbeddingTag = "<embed-code"

// TransitionMap is a type for mapping one State to a list of possible next States.
type TransitionMap map[State][]State

// Transitions maps State to the list of possible next states.
//
// States are chosen considered the logical validity of their existence.
//
// The order of the next states is important. States are ordered by the level of their specificity,
// so the first state in the list is the most specific one.
// When the state is chosen, the latter ones are skipped.
var Transitions = TransitionMap{
	Start:            {Finish, EmbedInstruction, RegularLine},
	RegularLine:      {Finish, EmbedInstruction, RegularLine},
	EmbedInstruction: {CodeFenceStart, BlankLine},
	BlankLine:        {CodeFenceStart, BlankLine},
	CodeFenceStart:   {CodeFenceEnd, CodeSampleLine},
	CodeSampleLine:   {CodeFenceEnd, CodeSampleLine},
	CodeFenceEnd:     {Finish, EmbedInstruction, RegularLine},
}

var (
	Start            = StartState{}
	RegularLine      = RegularLineState{}
	EmbedInstruction = EmbedInstructionTokenState{}
	BlankLine        = BlankLineState{}
	CodeFenceStart   = CodeFenceStartState{}
	CodeFenceEnd     = CodeFenceEndState{}
	CodeSampleLine   = CodeSampleLineState{}
	Finish           = FinishState{}
)
