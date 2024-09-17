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

// Transitions maps state names to the list of possible next states.
//
// States are chosen considered the logical validity of their existence.
//
// The order of the next states is important. States are ordered by the level of their specificity,
// so the first state in the list is the most specific one.
// When the state is chosen, the latter ones are skipped.

type TransitionMap map[Transition][]Transition

var Transitions = TransitionMap{
	start:            {finish, embedInstruction, regularLine},
	regularLine:      {finish, embedInstruction, regularLine},
	embedInstruction: {codeFenceStart, blankLine},
	blankLine:        {codeFenceStart, blankLine},
	codeFenceStart:   {codeFenceEnd, codeSampleLine},
	codeSampleLine:   {codeFenceEnd, codeSampleLine},
	codeFenceEnd:     {finish, embedInstruction, regularLine},
}

var (
	start            = Start{"START"}
	regularLine      = RegularLine{"REGULAR_LINE"}
	embedInstruction = EmbedInstructionToken{"EMBEDDING_INSTRUCTION"}
	blankLine        = BlankLine{"BLANK_LINE"}
	codeFenceStart   = CodeFenceStart{"CODE_FENCE_START"}
	codeFenceEnd     = CodeFenceEnd{"CODE_FENCE_END"}
	codeSampleLine   = CodeSampleLine{"CODE_SAMPLE_LINE"}
	finish           = Finish{"FINISH"}
)
