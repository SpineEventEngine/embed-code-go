package embedding

type UnexpectedDiffError struct {
	changedFiles []string
}

func (m *UnexpectedDiffError) Error() string {
	return "boom"
}
