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

package indent_test

import (
	"strings"
	"testing"

	"embed-code/embed-code-go/indent"
)

// FuzzCutIndent checks indentation trimming over arbitrary source blocks.
func FuzzCutIndent(f *testing.F) {
	f.Add("        System.out.println(\"Hi\");\n  \n        return;")
	f.Add("\n  foo\n    bar\n\n  baz")
	f.Add("\t\tfirst\n\tsecond\n")
	f.Add("plain\ntext")

	f.Fuzz(func(t *testing.T, source string) {
		lines := strings.Split(source, "\n")
		commonIndent := indent.MaxCommonIndentation(lines)

		changedLines := indent.CutIndent(lines, commonIndent)

		if len(changedLines) != len(lines) {
			t.Fatalf("CutIndent returned %d lines for %d input lines", len(changedLines), len(lines))
		}
		for lineIndex, line := range lines {
			if commonIndent <= len(line) && changedLines[lineIndex] != line[commonIndent:] {
				t.Fatalf(
					"line %d was cut incorrectly: got %q, want %q",
					lineIndex,
					changedLines[lineIndex],
					line[commonIndent:],
				)
			}
			if commonIndent > len(line) && changedLines[lineIndex] != "" {
				t.Fatalf(
					"short line %d was not fully removed: got %q",
					lineIndex,
					changedLines[lineIndex],
				)
			}
		}
	})
}
