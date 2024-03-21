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
	return lines[partition.StartPosition : partition.EndPosition+1]
}
