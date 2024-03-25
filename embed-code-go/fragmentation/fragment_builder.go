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

import "fmt"

// A single fragment builder.
type FragmentBuilder struct {
	FileName   string
	Partitions []Partition
	Name       string
}

//
// Public methods
//

// Adds a new partition with the given start position.
//
// Don't forget to call `add_end_position` when the end of the fragment is reached.
//
// @param [Integer] start_position a starting position of the fragment
func (fragmentBuilder *FragmentBuilder) AddStartPosition(startPosition int) {

	if len(fragmentBuilder.Partitions) > 0 {
		lastAddedPartition := fragmentBuilder.Partitions[len(fragmentBuilder.Partitions)-1]
		if lastAddedPartition.EndPosition == nil {
			panic("error: the last added partition has no end position")
		}
	}

	partition := Partition{StartPosition: &startPosition}
	fragmentBuilder.Partitions = append(fragmentBuilder.Partitions, partition)
}

// Completes previously created fragment partition with its end position.
//
// Should be called after `add_start_position`.
//
// @param [Integer] end_position an end position position of the fragment
func (fragmentBuilder *FragmentBuilder) AddEndPosition(endPosition int) {
	if len(fragmentBuilder.Partitions) == 0 {
		panic("error: the list of partitions is empty")
	}
	lastAddedPartition := &fragmentBuilder.Partitions[len(fragmentBuilder.Partitions)-1]
	if lastAddedPartition.EndPosition != nil {
		panic(fmt.Sprintf("unexpected #enddocfragment statement at %s:%d",
			fragmentBuilder.FileName,
			*lastAddedPartition.EndPosition),
		)
	}
	lastAddedPartition.EndPosition = &endPosition
}

func (fragmentBuilder FragmentBuilder) Build() Fragment {
	return Fragment{
		Name:       fragmentBuilder.Name,
		Partitions: fragmentBuilder.Partitions,
	}
}
