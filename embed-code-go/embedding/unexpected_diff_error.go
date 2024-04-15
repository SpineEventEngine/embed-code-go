package embedding

import "fmt"

type UnexpectedDiffError struct {
	changedFiles []string
}

func (m *UnexpectedDiffError) Error() string {
	return fmt.Sprintf("unexpected diff: %v", m.changedFiles)
}
