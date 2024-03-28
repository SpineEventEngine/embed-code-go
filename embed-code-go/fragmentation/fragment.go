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
	"embed-code/embed-code-go/configuration"
	"embed-code/embed-code-go/indent"
	"math"
	"strings"
)

const (
	DefaultFragmentName = "_default"
)

// A single fragment in a file.
type Fragment struct {
	Name       string      // A name of a Fragment.
	Partitions []Partition // A list of partitions found for a Fragment.
}

//
// Initializers
//

// Creates and returns Fragment with DefaultFragmentName.
func CreateDefaultFragment() Fragment {
	return Fragment{
		Name:       DefaultFragmentName,
		Partitions: []Partition{},
	}
}

//
// Public methods
//

// Takes given allLines, unites them into a text and writes it into given file.
//
// Creates the file if not exists and overwrites if exists.
func (fragment Fragment) WriteTo(
	file FragmentFile,
	lines []string,
	configuration configuration.Configuration,
) {
	text := fragment.text(lines, configuration)
	file.Write(text)
}

//
// Static functions
//

// Calculates and returns a list which contains corresponding lines for every partition.
//
// allLines — a list with every line of the file.
//
// partitions — a list with partitions to select lines from.
func calculatePartitionLines(allLines []string, partitions []Partition) [][]string {
	partitionLines := [][]string{}
	for _, part := range partitions {
		partitionText := part.Select(allLines)
		partitionLines = append(partitionLines, partitionText)
	}
	return partitionLines
}

// Calculates and returns common indentation on which it is possible to trim the lines
// without any harm.
//
// partitionLines — a list which contains corresponding lines for every parition.
func calculateCommonIndentation(partitionLines [][]string) int {
	commonIndentation := math.MaxInt32
	for _, partitionText := range partitionLines {
		indentation := indent.MaxCommonIndentation(partitionText)
		if indentation < commonIndentation {
			commonIndentation = indentation
		}
	}
	return commonIndentation
}

//
// Private methods
//

// Obtains the text for the fragment.
//
// The each partition of the fragment is separated with the Configuration.Separator.
//
// allLines — a list with every line of the file.
func (fragment Fragment) text(allLines []string, configuration configuration.Configuration) string {

	if fragment.isDefault() {
		return strings.Join(allLines, "")
	}

	partitionLines := calculatePartitionLines(allLines, fragment.Partitions)
	commonIndentation := calculateCommonIndentation(partitionLines)

	text := ""
	for index, line := range partitionLines {
		if index != 0 {
			text += configuration.Separator + "\n"
		}
		cutIndentLines := indent.CutIndent(line, commonIndentation)
		for _, cutIndentLine := range cutIndentLines {
			text += cutIndentLine
		}
	}
	return text
}

func (fragment Fragment) isDefault() bool {
	return fragment.Name == DefaultFragmentName
}
