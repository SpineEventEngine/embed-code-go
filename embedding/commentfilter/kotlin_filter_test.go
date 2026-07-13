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

package commentfilter_test

import (
	. "embed-code/embed-code-go/embedding/commentfilter"

	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("Kotlin", func() {
	It("should keep all comments", func() {
		lines := []string{
			"/** API docs. */",
			"/* implementation note */",
			"val value = 1 // inline note",
		}

		assertFiltered("Sample.kt", RetainAll, lines, lines)
	})

	It("should strip comments without treating raw string text as comments", func() {
		lines := []string{
			"/* outer /* nested */ still comment */",
			"val text = \"\"\"",
			"    This is not a /* comment */.",
			"    This is not a // comment either.",
			"    This removes a real comment: ${render(/* real raw argument */ value)}",
			"\"\"\"",
			"val message = \"value = ${render(/* real argument */ value)}\"",
		}

		expected := []string{
			"val text = \"\"\"",
			"    This is not a /* comment */.",
			"    This is not a // comment either.",
			"    This removes a real comment: ${render( value)}",
			"\"\"\"",
			"val message = \"value = ${render( value)}\"",
		}

		assertFiltered("Sample.kt", RetainNone, lines, expected)
	})

	It("should continue raw string interpolation after line comments", func() {
		lines := []string{
			"val text = \"\"\"",
			"    ${render(",
			"        value, // real line comment",
			"        /* real block comment */ nextValue",
			"    )}",
			"    Keep // raw text.",
			"\"\"\"",
		}

		expected := []string{
			"val text = \"\"\"",
			"    ${render(",
			"        value, ",
			"         nextValue",
			"    )}",
			"    Keep // raw text.",
			"\"\"\"",
		}

		assertFiltered("Sample.kt", RetainNone, lines, expected)
	})

	It("should keep KDoc comments", func() {
		lines := []string{
			"/** API docs. */",
			"/* implementation note */",
			"val value = 1 // inline note",
		}

		expected := []string{
			"/** API docs. */",
			"val value = 1 ",
		}

		assertFiltered("Sample.kt", RetainDocumentation, lines, expected)
	})

	It("should keep regular comments", func() {
		lines := []string{
			"/** API docs. */",
			"/* implementation note */",
			"val value = 1 // inline note",
		}

		expected := []string{
			"/* implementation note */",
			"val value = 1 // inline note",
		}

		assertFiltered("Sample.kt", RetainRegular, lines, expected)
	})

	It("should keep inline comments", func() {
		lines := []string{
			"/** API docs. */",
			"/* implementation note */",
			"val value = 1 // inline note",
		}

		expected := []string{
			"val value = 1 // inline note",
		}

		assertFiltered("Sample.kt", RetainInline, lines, expected)
	})

	It("should keep nested block comments", func() {
		lines := []string{
			"val before = 1",
			"/* outer",
			"   /* nested */",
			"   still outer */",
			"val after = 2 // inline",
		}

		expected := []string{
			"val before = 1",
			"/* outer",
			"   /* nested */",
			"   still outer */",
			"val after = 2 ",
		}

		assertFiltered("Sample.kts", RetainBlock, lines, expected)
	})

	It("should close empty documentation block comments", func() {
		lines := []string{
			"/**/",
			"val a = 1",
			"val b = 2 /**/ val c = 3",
		}

		expected := []string{
			"val a = 1",
			"val b = 2  val c = 3",
		}

		assertFiltered("Sample.kt", RetainNone, lines, expected)
	})

	It("should preserve escaped strings, characters, and nested interpolation braces", func() {
		lines := []string{
			`val quote = '\'' // comment`,
			`val text = "escaped \\" // literal" // comment`,
			`val nested = "${if (ready) { "// literal" } else { value }}" // comment`,
			`val raw = "${"""raw"""}" // comment`,
			`val quoted = "${"text"}" // comment`,
		}
		expected := []string{
			`val quote = '\'' `,
			`val text = "escaped \\" `,
			`val nested = "${if (ready) { "// literal" } else { value }}" `,
			`val raw = "${"""raw"""}" `,
			`val quoted = "${"text"}" `,
		}

		assertFiltered("Sample.kt", RetainNone, lines, expected)
	})
})
