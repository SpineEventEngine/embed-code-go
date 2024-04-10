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
