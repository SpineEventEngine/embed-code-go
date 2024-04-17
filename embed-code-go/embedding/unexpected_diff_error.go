package embedding

import "fmt"

// Describes an error which occurs if not up-to-dated files are found during the checking.
type UnexpectedDiffError struct {
	changedFiles []string
}

func (m *UnexpectedDiffError) Error() string {
	return fmt.Sprintf("unexpected diff: %v", m.changedFiles)
}
