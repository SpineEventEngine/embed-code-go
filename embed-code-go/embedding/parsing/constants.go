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

// Maps state names to the list of possible next states.
//
// States are chosen considered the logical validity of their existence.
//
// The order of the next states is important. States are ordered by the level of their specificity,
// so the first state in the list is the most specific one.
// When the state is chosen, the latter ones are skipped.
var Transitions = map[string][]string{
	"START":                 {"FINISH", "EMBEDDING_INSTRUCTION", "REGULAR_LINE"},
	"REGULAR_LINE":          {"FINISH", "EMBEDDING_INSTRUCTION", "REGULAR_LINE"},
	"EMBEDDING_INSTRUCTION": {"CODE_FENCE_START", "BLANK_LINE"},
	"BLANK_LINE":            {"CODE_FENCE_START", "BLANK_LINE"},
	"CODE_FENCE_START":      {"CODE_FENCE_END", "CODE_SAMPLE_LINE"},
	"CODE_SAMPLE_LINE":      {"CODE_FENCE_END", "CODE_SAMPLE_LINE"},
	"CODE_FENCE_END":        {"FINISH", "EMBEDDING_INSTRUCTION", "REGULAR_LINE"},
}

// Maps a state name to a Transition.
var StateToTransition = map[string]Transition{
	"REGULAR_LINE":          RegularLine{},
	"EMBEDDING_INSTRUCTION": EmbedInstructionToken{},
	"BLANK_LINE":            BlankLine{},
	"CODE_FENCE_START":      CodeFenceStart{},
	"CODE_FENCE_END":        CodeFenceEnd{},
	"CODE_SAMPLE_LINE":      CodeSampleLine{},
	"FINISH":                Finish{},
}
