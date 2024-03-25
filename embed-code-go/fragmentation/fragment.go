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
	DefaultFragment = "_default"
)

// A single fragment in a file.
type Fragment struct {
	Name       string      // A name of a Fragment
	Partitions []Partition // A list of partitions found for a Fragment
}

//
// Initializers
//

func CreateDefaultFragment() Fragment {
	return Fragment{
		Name:       DefaultFragment,
		Partitions: []Partition{},
	}
}

//
// Public methods
//

// WriteTo takes the given allLines,
// unites them into a text
// and writes it into the given file.
func (fragment Fragment) WriteTo(
	file FragmentFile,
	allLines []string,
	configuration configuration.Configuration,
) {
	text := fragment.text(allLines, configuration)
	file.Write(text)
}

//
// Private methods
//

func (fragment Fragment) text(allLines []string, configuration configuration.Configuration) string {

	if fragment.isDefault() {
		return strings.Join(allLines, "")
	}

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

func (fragment Fragment) isDefault() bool {
	return fragment.Name == DefaultFragment
}
