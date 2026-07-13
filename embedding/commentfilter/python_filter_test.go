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

var _ = Describe("Python", func() {
	It("should strip all comments", func() {
		lines := []string{
			"# module comment",
			"name = 'hash # literal'",
			"value = 1 # inline comment",
		}

		expected := []string{
			"name = 'hash # literal'",
			"value = 1 ",
		}

		assertFiltered("module.py", RetainNone, lines, expected)
	})

	It("should strip comments without treating triple-quoted string content as comments", func() {
		lines := []string{
			"# module comment",
			"def message():",
			"    \"\"\"",
			"    Keep this # docstring text.",
			"    \"\"\"",
			`    escaped = """a\"""" # inline comment`,
			"    value = '''",
			"    Keep this # multiline text.",
			"    ''' # inline comment",
			"    return value # real comment",
		}

		expected := []string{
			"def message():",
			"    \"\"\"",
			"    Keep this # docstring text.",
			"    \"\"\"",
			`    escaped = """a\"""" `,
			"    value = '''",
			"    Keep this # multiline text.",
			"    ''' ",
			"    return value ",
		}

		assertFiltered("module.py", RetainNone, lines, expected)
	})

	It("should strip hash comments inside multi-line f-string expressions", func() {
		lines := []string{
			"message = f\"\"\"",
			"Keep this # f-string text.",
			"Total: {",
			"    compute_total(items)  # expression comment",
			"}",
			"\"\"\"",
			`formatted = f"{value:#x}" # inline comment`,
			"slice = f\"\"\"",
			"{",
			"    values[",
			"        :  # slice comment",
			"    ]",
			"}",
			"\"\"\"",
			`braces = rf"{{ # literal }} {value # expression comment`,
			`}"`,
			"value = 1 # real comment",
		}

		expected := []string{
			"message = f\"\"\"",
			"Keep this # f-string text.",
			"Total: {",
			"    compute_total(items)  ",
			"}",
			"\"\"\"",
			`formatted = f"{value:#x}" `,
			"slice = f\"\"\"",
			"{",
			"    values[",
			"        :  ",
			"    ]",
			"}",
			"\"\"\"",
			`braces = rf"{{ # literal }} {value `,
			`}"`,
			"value = 1 ",
		}

		assertFiltered("module.py", RetainNone, lines, expected)
	})

	It("should preserve strings nested inside f-string expressions", func() {
		cases := []struct {
			lines    []string
			expected []string
		}{
			{
				lines:    []string{`value = f"{format('escaped \\' quote')}" # comment`},
				expected: []string{`value = f"{format('escaped \\' quote')}" # comment`},
			},
			{
				lines: []string{`value = f"{format("""raw`, `text`, `""")}" # comment`},
				expected: []string{
					`value = f"{format("""raw`, `text`, `""")}" `,
				},
			},
			{
				lines:    []string{`value = f"{number:04`},
				expected: []string{`value = f"{number:04`},
			},
			{
				lines:    []string{`value = f"{format('unterminated`},
				expected: []string{`value = f"{format('unterminated`},
			},
			{
				lines:    []string{`value = f"{ {'nested': value} }" # comment`},
				expected: []string{`value = f"{ {'nested': value} }" `},
			},
			{
				lines:    []string{`value = "unterminated`},
				expected: []string{`value = "unterminated`},
			},
		}

		for _, tc := range cases {
			assertFiltered("module.py", RetainNone, tc.lines, tc.expected)
		}
	})
})
