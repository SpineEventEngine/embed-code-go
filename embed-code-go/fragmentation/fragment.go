package fragmentation

import (
	"embed-code/embed-code-go/configuration"
	"embed-code/embed-code-go/indent"
	"math"
)

const (
	DefaultFragment = "_default"
)

type Fragment struct {
	Name       string
	Partitions []Partition
}

// Initializers
func CreateDefaultFragment() Fragment {
	return Fragment{
		Name:       DefaultFragment,
		Partitions: []Partition{},
	}
}

// Public methods
func (fragment Fragment) Text(allLines []string, configuration configuration.Configuration) string {
	commonIndentetaion := math.MaxInt32
	partitionLines := [][]string{}

	for _, part := range fragment.Partitions {
		partitionText := part.Select(allLines)
		partitionLines = append(partitionLines, partitionText)
		indentetaion := indent.MaxCommonIndentation(partitionText)
		if indentetaion < commonIndentetaion {
			commonIndentetaion = indentetaion
		}
	}

	text := ""
	for index, line := range partitionLines {
		if index != 0 {
			text += configuration.Separator + "\n"
		}
		cutIndentLines := indent.CutIndent(line, commonIndentetaion)
		for _, cutIndentLine := range cutIndentLines {
			text += cutIndentLine
		}
	}
	return text
}

func (fragment Fragment) WriteTo(file FragmentFile,
	allLines []string,
	configuration configuration.Configuration,
) {
	text := fragment.Text(allLines, configuration)
	file.Write(text)
}
