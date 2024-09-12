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
	"strings"

	config "embed-code/embed-code-go/configuration"
	"embed-code/embed-code-go/indent"
)

const DefaultFragmentName = "_default"

// Fragment is a single fragment in a file.
//
// Name — a name of a Fragment.
//
// Partitions — a list of partitions found for a Fragment.
type Fragment struct {
	Name       string
	Partitions []Partition
}

// CreateDefaultFragment creates and returns Fragment with DefaultFragmentName.
func CreateDefaultFragment() Fragment {
	return Fragment{
		Name:       DefaultFragmentName,
		Partitions: []Partition{},
	}
}

// WriteTo takes given lines, unites them into a text and writes it into given file.
//
// file — a FragmentFile to write the lines to.
//
// lines — a list of strings to write.
//
// configuration — a Configuration with all embed-code settings.
//
// Creates the file if not exists and overwrites if exists.
func (f Fragment) WriteTo(file FragmentFile, lines []string, config config.Configuration) {
	text := f.text(lines, config)
	file.Write(text)
}

func (f Fragment) isDefault() bool {
	return f.Name == DefaultFragmentName
}

// Obtains the text for the fragment.
//
// The partition of the fragment is separated with the Configuration.Separator.
//
// lines — a list with every line of the file.
//
// configuration — a configuration for embedding.
func (f Fragment) text(lines []string, config config.Configuration) string {
	if f.isDefault() {
		return strings.Join(lines, "\n")
	}

	partitionsTexts := calculatePartitionsTexts(lines, f.Partitions)

	text := ""
	for index, partitionText := range partitionsTexts {
		indentation := indent.MaxCommonIndentation(partitionText)
		cutIndentLines := indent.CutIndent(partitionText, indentation)

		if index != 0 {
			separatorIndentation := separatorIndent(cutIndentLines)
			text += separatorIndentation + config.Separator + "\n"
		}

		text += strings.Join(cutIndentLines, "\n") + "\n"
	}

	return text
}

// Returns string indent for separator.
func separatorIndent(lines []string) string {
	if len(lines) > 0 {
		firstLine := lines[0]
		leadingSpaces := len(firstLine) - len(strings.TrimLeft(firstLine, " "))

		return strings.Repeat(" ", leadingSpaces)
	}

	return ""
}

// Calculates and returns a list which contains corresponding lines for every partition.
//
// lines — a list with every line of the file.
//
// partitions — a list with partitions to select lines from.
func calculatePartitionsTexts(lines []string, partitions []Partition) [][]string {
	var partitionLines [][]string
	for _, part := range partitions {
		partitionText := part.Select(lines)
		partitionLines = append(partitionLines, partitionText)
	}

	return partitionLines
}
