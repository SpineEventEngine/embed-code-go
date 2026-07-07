// Copyright 2026, TeamDev. All rights reserved.
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

package commentfilter_test

import (
	. "embed-code/embed-code-go/embedding/commentfilter"

	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("Java-style languages", func() {
	It("should keep documentation comments", func() {
		lines := []string{
			"/** API docs. */",
			"// implementation note",
			"fun call() = \"// literal\"",
		}

		expected := []string{
			"/** API docs. */",
			"fun call() = \"// literal\"",
		}

		assertFiltered("api.kt", RetainDocumentation, lines, expected)
	})

	It("should keep block comments", func() {
		lines := []string{
			"/** API docs. */",
			"/* implementation note */",
			"String create();",
		}

		expected := []string{
			"/* implementation note */",
			"String create();",
		}

		assertFiltered("Api.java", RetainBlock, lines, expected)
	})

	It("should keep regular comments", func() {
		lines := []string{
			"/** API docs. */",
			"/* implementation note */",
			"String create(); // inline note",
		}

		expected := []string{
			"/* implementation note */",
			"String create(); // inline note",
		}

		assertFiltered("Api.java", RetainRegular, lines, expected)
	})

	It("should strip comments without treating text block content as comments", func() {
		lines := []string{
			"// header comment",
			"String help = \"\"\"",
			"    Keep this // text.",
			"    Keep this /* text */ too.",
			"    \"\"\";",
			"String value = \"not a // comment\"; // inline comment",
		}

		expected := []string{
			"String help = \"\"\"",
			"    Keep this // text.",
			"    Keep this /* text */ too.",
			"    \"\"\";",
			"String value = \"not a // comment\"; ",
		}

		assertFiltered("Api.java", RetainNone, lines, expected)
	})

	It("should not close text blocks on escaped triple quotes", func() {
		lines := []string{
			"String help = \"\"\"",
			`    Quote: \"""`,
			`    Escaped quote: \"`,
			"    Keep this // text.",
			"    \"\"\";",
			"String value = \"kept\"; // real comment",
		}

		expected := []string{
			"String help = \"\"\"",
			`    Quote: \"""`,
			`    Escaped quote: \"`,
			"    Keep this // text.",
			"    \"\"\";",
			"String value = \"kept\"; ",
		}

		assertFiltered("Api.java", RetainNone, lines, expected)
	})
})
