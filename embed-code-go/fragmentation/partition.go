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

// A code fragment partition.
//
// A fragment may consist of a few partitions, collected from different points in the code file. In
// the resulting doc file, the partitions are joined by the Configuration.Separator.
//
// StartPosition is an index from which the scope of partition exists.
//
// EndPosition is an index on which the scope of partition ends.
type Partition struct {
	StartPosition *int
	EndPosition   *int
}

//
// Public methods
//

// Returns the partition-related lines from the given lines.
// If EndPosition is nil, returns all the lines started from StartPosition.
func (partition Partition) Select(lines []string) []string {
	if partition.EndPosition == nil {
		// This part is for emulating the behaviour of the original embed code.
		// In Ruby, with unsetted EndPosition, it used to return all the lines from the StartPosition.
		return lines[*partition.StartPosition:]
	} else {
		return lines[*partition.StartPosition : *partition.EndPosition+1]
	}
}
