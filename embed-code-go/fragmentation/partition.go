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

// Partition a code fragment partition.
//
// A fragment may consist of a few partitions, collected from different points in the code file.
// In the resulting doc file, the partitions are joined by the Configuration.Separator.
// StartPosition and EndPosition are both set to -1 by default as the default int value for them
// is 0, which is wrong, because 0 is in the scope of possible values for them.
//
// StartPosition — an index from which the scope of partition exists.
//
// EndPosition — an index on which the scope of partition ends.
type Partition struct {
	StartPosition int
	EndPosition   int
}

// NewPartition returns a new Partition with both positions set to -1, as they should to be
// positive once set by a user.
func NewPartition() Partition {
	return Partition{
		-1,
		-1,
	}
}

// Select returns the partition-related lines from given lines.
// If EndPosition is not set, returns all the lines started from StartPosition.
func (p Partition) Select(lines []string) []string {
	startPosition := p.StartPosition
	endPosition := p.EndPosition

	// Verifying lines actually have those indexes.
	_, ok := safeAccess(lines, startPosition)
	if !ok {
		panic("an unexpected error occurred. the given lines don't have start position")
	}

	if endPosition < 0 {
		return lines[startPosition:]
	}

	_, ok = safeAccess(lines, endPosition)
	if !ok {
		panic("an unexpected error occurred. the given lines don't have end position")
	}

	return lines[startPosition : endPosition+1]
}

func safeAccess(slice []string, index int) (string, bool) {
	var ok bool
	defer func() {
		if r := recover(); r != nil {
			ok = false
		}
	}()
	value := slice[index]
	ok = true

	return value, ok
}
