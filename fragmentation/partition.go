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

import "fmt"

// Partition a code fragment partition.
//
// A fragment may consist of a few partitions, collected from different points in the code file.
// In the resulting doc file, the partitions are joined by the Configuration.Separator.
// StartPosition and EndPosition are both set to -1 by default as the default int value for them
// is 0, which is wrong, because 0 is in the scope of possible values for them.
type Partition struct {
	// StartPosition is the first source-line index included in the partition.
	StartPosition int

	// EndPosition is the last source-line index included in the partition.
	EndPosition int
}

// NewPartition returns a Partition with both positions unset as -1.
//
// Returns empty partition ready to receive start and end positions.
func NewPartition() Partition {
	return Partition{
		-1,
		-1,
	}
}

// Select returns the partition-related lines from given lines.
//
// Parameters:
// lines - provides source lines indexed by StartPosition and EndPosition.
//
// Returns:
// []string - selected source lines.
// error - when configured positions are outside lines.
func (p Partition) Select(lines []string) ([]string, error) {
	startPosition := p.StartPosition
	endPosition := p.EndPosition

	if !hasLineIndex(lines, startPosition) {
		return nil, fmt.Errorf(
			"fragment partition start position %d is outside source lines",
			startPosition,
		)
	}

	if endPosition < 0 {
		return lines[startPosition:], nil
	}

	if endPosition < startPosition-1 {
		return nil, fmt.Errorf(
			"fragment partition end position %d is before start position %d",
			endPosition,
			startPosition,
		)
	}

	if !hasLineIndex(lines, endPosition) {
		return nil, fmt.Errorf(
			"fragment partition end position %d is outside source lines",
			endPosition,
		)
	}

	return lines[startPosition : endPosition+1], nil
}

// hasLineIndex reports whether lines contain index.
func hasLineIndex(lines []string, index int) bool {
	return index >= 0 && index < len(lines)
}
