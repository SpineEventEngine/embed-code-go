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

package parsing_test

import (
	"strings"
	"testing"

	"embed-code/embed-code-go/embedding/parsing"
)

// FuzzPatternFindIn checks source-line pattern compilation and matching.
func FuzzPatternFindIn(f *testing.F) {
	f.Add("main", "public static void main(String[] args) {\n    return;\n}", 0)
	f.Add("^  padded text  $ \\n ^Use \\* to multiply$", "  padded text  \nUse * to multiply", 0)
	f.Add("^Use \\* to multiply$ \\n żółć$", "Use * to multiply\nżółć", 0)
	f.Add("        return;", "        System.out.println(\"Hi\");\n  \n        return;", 0)
	f.Add("[", "invalid pattern is rejected", 0)
	f.Add("0{$", "0", 0)
	f.Add("0{}", "0", 0)

	f.Fuzz(func(t *testing.T, sourceGlob string, source string, startFrom int) {
		if len(sourceGlob) > 4096 || len(source) > 8192 {
			t.Skip("generated input exceeds the bounded fuzz target size")
		}

		pattern, err := parsing.NewPattern(sourceGlob)
		if err != nil {
			return
		}
		lines := strings.Split(source, "\n")

		start, end, found := pattern.FindIn(lines, startFrom)

		if !found {
			return
		}
		assertValidPatternMatch(t, lines, startFrom, start, end)
	})
}

// assertValidPatternMatch checks the public range contract returned by FindIn.
func assertValidPatternMatch(t *testing.T, lines []string, startFrom int, start int, end int) {
	t.Helper()
	if start < 0 || start >= len(lines) {
		t.Fatalf("start index %d is outside %d source lines", start, len(lines))
	}
	if end < start || end >= len(lines) {
		t.Fatalf(
			"end index %d is outside matched range starting at %d in %d source lines",
			end,
			start,
			len(lines),
		)
	}
	if startFrom >= 0 && start < startFrom {
		t.Fatalf("match starts at %d before requested start %d", start, startFrom)
	}
}
