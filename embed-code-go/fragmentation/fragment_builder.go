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
//
// CodeFilePath — a path to a file to fragment.
//
// Partitions — a list of partitions of a file to fragment.
//
// Name — a name of a Fragment.
type FragmentBuilder struct {
	CodeFilePath string
	Partitions   []Partition
	Name         string
}

// AddStartPosition adds a new partition with given startPosition.
//
// AddEndPosition is need to be called when the end of the fragment is reached,
// or else it will be considered that the end of partition is in the end of the file.
//
// startPosition — starting position of the fragment.
func (b *FragmentBuilder) AddStartPosition(startPosition int) error {
	if !b.isPartitionsEmpty() {
		lastPartition := b.lastAddedPartition()
		if lastPartition.EndPosition == 0 {
			return fmt.Errorf("error: for the fragment \"%s\" of the file \"%s\", "+
				"the last added partition has no end position", b.Name, b.CodeFilePath)
		}
	}

	partition := Partition{StartPosition: startPosition}
	b.Partitions = append(b.Partitions, partition)

	return nil
}

// AddEndPosition completes previously created fragment partition with its endPosition.
// It should be called after AddStartPosition.
//
// endPosition — end position of the fragment.
func (b *FragmentBuilder) AddEndPosition(endPosition int) error {
	if b.isPartitionsEmpty() {
		return errors.New("error: the list of partitions is empty")
	}
	lastPartition := b.lastAddedPartition()
	if lastPartition.EndPosition != 0 {
		return fmt.Errorf("unexpected #enddocfragment statement at %s:%d", b.CodeFilePath,
			lastPartition.EndPosition)
	}
	lastPartition.EndPosition = endPosition

	return nil
}

// Build creates and returns new Fragment with the previously added and filled Partitions.
func (b *FragmentBuilder) Build() Fragment {
	return Fragment{
		Name:       b.Name,
		Partitions: b.Partitions,
	}
}

func (b *FragmentBuilder) isPartitionsEmpty() bool {
	return len(b.Partitions) == 0
}

func (b *FragmentBuilder) lastAddedPartition() Partition {
	lastIndex := len(b.Partitions) - 1

	return b.Partitions[lastIndex]
}
