package fragmentation

import "fmt"

type FragmentBuilder struct {
	FileName   string
	Partitions []Partition
	Name       string
}

// Public methods

// Adds a new partition with the given start position.
//
// Don't forget to call `add_end_position` when the end of the fragment is reached.
//
// @param [Integer] start_position a starting position of the fragment
func (fragmentBuilder *FragmentBuilder) AddStartPosition(startPosition int) {

	if len(fragmentBuilder.Partitions) > 0 {
		lastAddedPartition := fragmentBuilder.Partitions[len(fragmentBuilder.Partitions)-1]
		if lastAddedPartition.EndPosition == 0 {
			fmt.Println("Error: the last added partition has no end position.")
			return
		}
	}

	partition := Partition{StartPosition: startPosition}
	fragmentBuilder.Partitions = append(fragmentBuilder.Partitions, partition)
}

// Completes previously created fragment partition with its end position.
//
// Should be called after `add_start_position`.
//
// @param [Integer] end_position an end position position of the fragment
func (fragmentBuilder *FragmentBuilder) AddEndPosition(endPosition int) {
	if len(fragmentBuilder.Partitions) == 0 {
		fmt.Println("Error: the list of partitions is empty.")
		return
	}
	lastAddedPartition := &fragmentBuilder.Partitions[len(fragmentBuilder.Partitions)-1]
	if lastAddedPartition.EndPosition != 0 {
		fmt.Printf("Unexpected #enddocfragment statement at %s:%s.\n",
			fragmentBuilder.FileName,
			fmt.Sprint((lastAddedPartition.EndPosition)),
		)
		return
	}
	lastAddedPartition.EndPosition = endPosition
}

func (fragmentBuilder FragmentBuilder) Build() Fragment {
	return Fragment{
		Name:       fragmentBuilder.Name,
		Partitions: fragmentBuilder.Partitions,
	}
}
