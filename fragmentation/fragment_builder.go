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

package fragmentation

import (
	"errors"
	"fmt"
)

// FragmentBuilder is a single fragment builder.
type FragmentBuilder struct {
	// CodeFilePath is the path to the file being fragmented.
	CodeFilePath string

	// Partitions contains the partitions found for the fragment.
	Partitions []Partition

	// Name is the fragment name.
	Name string
}

// AddStartPosition adds a new fragment partition starting at startPosition.
//
// Parameters:
// startPosition - provides the zero-based source line where the partition starts.
//
// Returns an error when the previous partition is still open.
func (b *FragmentBuilder) AddStartPosition(startPosition int) error {
	if !b.isPartitionsEmpty() {
		lastPartition := b.lastAddedPartition()
		if lastPartition.EndPosition < 0 {
			return fmt.Errorf("error: for the fragment \"%s\" of the file \"%s\", "+
				"the last added partition has no end position", b.Name, b.CodeFilePath)
		}
	}

	partition := NewPartition()
	partition.StartPosition = startPosition
	b.Partitions = append(b.Partitions, partition)

	return nil
}

// AddEndPosition completes the latest fragment partition at endPosition.
//
// It is needed to be called when the end of the fragment is reached,
// or else it will be considered that the end of partition is in the end of the file.
//
// Parameters:
// endPosition - provides the zero-based source line where the partition ends.
//
// Returns an error when no partition is open or the latest partition already has an end.
func (b *FragmentBuilder) AddEndPosition(endPosition int) error {
	if b.isPartitionsEmpty() {
		return errors.New("the list of partitions is empty")
	}
	lastPartition := b.lastAddedPartition()
	if lastPartition.EndPosition < 0 {
		lastPartition.EndPosition = endPosition
	} else {
		return fmt.Errorf("unexpected #enddocfragment statement at %s:%d", b.CodeFilePath,
			lastPartition.EndPosition)
	}

	return nil
}

// Build creates a Fragment from the collected partition positions.
//
// Returns fragment with collected partitions.
func (b *FragmentBuilder) Build() Fragment {
	return Fragment{
		Name:       b.Name,
		Partitions: b.Partitions,
	}
}

// isPartitionsEmpty reports whether no partition positions have been collected.
func (b *FragmentBuilder) isPartitionsEmpty() bool {
	return len(b.Partitions) == 0
}

// lastAddedPartition returns the most recently collected partition.
func (b *FragmentBuilder) lastAddedPartition() *Partition {
	lastIndex := len(b.Partitions) - 1

	return &b.Partitions[lastIndex]
}
