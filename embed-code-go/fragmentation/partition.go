package fragmentation

// A code fragment partition.
//
// A fragment may consist of a few partitions, collected from different points in the code file. In
// the resulting doc file, the partitions are joined by the +Configuration::separator+.
type Partition struct {
	StartPosition int
	EndPosition   int
}

func (partition Partition) Select(lines []string) []string {
	if partition.EndPosition != 0 {
		return lines[partition.StartPosition : partition.EndPosition+1]
	} else {
		// This part is for emulating the behaviour of the original embed code.
		// In Ruby, with unsetted EndPosition, it used to return all the lines from the StartPosition.
		return lines[partition.StartPosition:]
	}

}
