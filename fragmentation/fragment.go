/*
 * Copyright 2026, TeamDev. All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 * Redistribution and use in source and/or binary forms, with or without
 * modification, must retain the above copyright notice and the following
 * disclaimer.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
 * "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
 * LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
 * A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
 * OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
 * SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
 * LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
 * DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
 * THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 * OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

package fragmentation

import (
	"strings"

	"embed-code/embed-code-go/indent"
)

// DefaultFragmentName identifies the whole-file fragment.
const DefaultFragmentName = "_default"

// Fragment is a single fragment in a file.
type Fragment struct {
	// Name is the fragment name.
	Name string

	// Partitions contains the source partitions that form the fragment.
	Partitions []Partition
}

// partitionText contains selected source lines and their indentation group.
type partitionText struct {
	// lines contains the source lines selected by one partition.
	lines []string

	// indentGroup identifies partitions whose common indentation is normalized together.
	indentGroup string
}

// CreateDefaultFragment creates a whole-file fragment.
//
// Returns whole-file fragment.
func CreateDefaultFragment() Fragment {
	return Fragment{
		Name:       DefaultFragmentName,
		Partitions: []Partition{},
	}
}

// isDefault reports whether this fragment represents the whole source file.
func (f Fragment) isDefault() bool {
	return f.Name == DefaultFragmentName
}

// text returns source text selected by the fragment.
//
// Parameters:
// lines - provides every source line in the file.
// separator - provides text inserted between multiple partitions of one fragment.
//
// Returns:
// string - rendered fragment text.
// error - when a partition cannot select its lines.
func (f Fragment) text(lines []string, separator string) (string, error) {
	if f.isDefault() {
		return strings.Join(lines, "\n"), nil
	}
	partitionsTexts, err := f.obtainPartitionTexts(lines)
	if err != nil {
		return "", err
	}
	indentations := commonIndentations(partitionsTexts)

	text := ""
	for index, partition := range partitionsTexts {
		indentation := indentations[partition.indentGroup]
		cutIndentLines := indent.CutIndent(partition.lines, indentation)

		if index > 0 {
			separatorIndentation := separatorIndent(cutIndentLines)
			text += separatorIndentation + separator + "\n"
		}

		text += strings.Join(cutIndentLines, "\n") + "\n"
	}

	return text, nil
}

// obtainPartitionTexts returns source lines selected for every fragment partition.
//
// Parameters:
// lines - provides every source line in the file.
//
// Returns:
// []partitionText - selected lines and indentation group for every partition.
// error - when a partition cannot select its lines.
func (f Fragment) obtainPartitionTexts(lines []string) ([]partitionText, error) {
	var partitions []partitionText
	for _, part := range f.Partitions {
		selectedLines, err := part.Select(lines)
		if err != nil {
			return nil, err
		}
		partitions = append(partitions, partitionText{
			lines:       selectedLines,
			indentGroup: part.IndentGroup,
		})
	}

	return partitions, nil
}

// commonIndentations calculates common indentation for every partition group.
//
// Parameters:
// partitions - provides selected source lines in source order.
//
// Returns indentation width by group name.
func commonIndentations(partitions []partitionText) map[string]int {
	groupLines := make(map[string][]string)
	for _, partition := range partitions {
		groupLines[partition.indentGroup] = append(
			groupLines[partition.indentGroup],
			partition.lines...,
		)
	}

	indentations := make(map[string]int, len(groupLines))
	for indentGroup, lines := range groupLines {
		indentations[indentGroup] = indent.MaxCommonIndentation(lines)
	}

	return indentations
}

// separatorIndent returns the indentation to use before a partition separator.
func separatorIndent(lines []string) string {
	if len(lines) > 0 {
		firstLine := lines[0]
		leadingSpaces := len(firstLine) - len(strings.TrimLeft(firstLine, " "))

		return strings.Repeat(" ", leadingSpaces)
	}

	return ""
}
