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

var _ = Describe("JavaScript and TypeScript", func() {
	It("should strip comments without treating regex and template text as comments", func() {
		lines := []string{
			"// module comment",
			"const url = `http://example.org/*not-comment*/`;",
			"const pattern = /https?:\\/\\/example\\.com\\/docs/;",
			"const help = `Keep // and /* markers */ in template text`;",
			"const nested = `${format(value /* remove this real comment */)}`;",
			"const value = 42; // inline comment",
		}

		expected := []string{
			"const url = `http://example.org/*not-comment*/`;",
			"const pattern = /https?:\\/\\/example\\.com\\/docs/;",
			"const help = `Keep // and /* markers */ in template text`;",
			"const nested = `${format(value )}`;",
			"const value = 42; ",
		}

		assertFiltered("sample.ts", RetainNone, lines, expected)
	})

	It("should preserve nested template literals inside template interpolations", func() {
		lines := []string{
			"const msg = `${items.map(i => `// ${i}`).join()}`;",
			"const braces = `${items.map(i => `}`).join()}`; // real comment",
			"const multiline = `${items.map(i => `// text",
			"still } /* text */ ${i}`).join()}`; // real comment",
		}

		expected := []string{
			"const msg = `${items.map(i => `// ${i}`).join()}`;",
			"const braces = `${items.map(i => `}`).join()}`; ",
			"const multiline = `${items.map(i => `// text",
			"still } /* text */ ${i}`).join()}`; ",
		}

		assertFiltered("sample.ts", RetainNone, lines, expected)
	})

	It("should filter comments inside multi-line nested template interpolations", func() {
		lines := []string{
			"const nestedExpression = `${items.map(i => `value ${format(",
			"i, // real nested expression comment",
			")}`).join()}`;",
		}

		expected := []string{
			"const nestedExpression = `${items.map(i => `value ${format(",
			"i, ",
			")}`).join()}`;",
		}

		assertFiltered("sample.ts", RetainNone, lines, expected)
	})

	It("should preserve multi-line template literal text", func() {
		lines := []string{
			"const help = `Keep // marker",
			"and /* marker */ text`; // real comment",
		}

		expected := []string{
			"const help = `Keep // marker",
			"and /* marker */ text`; ",
		}

		assertFiltered("sample.ts", RetainNone, lines, expected)
	})

	It("should preserve regex literals after expression-starting keywords", func() {
		lines := []string{
			"function parse() { return /\"/; } // real comment",
			"case /\"/.source: // real comment",
			"const type = typeof /\"/; // real comment",
			"const match = await /\"/; // real comment",
			"const hasValue = name in /\"/; // real comment",
			"const isPattern = value instanceof /\"/; // real comment",
			"if (missing) {} else /\"/.test(value); // real comment",
			"function skipComment() { return /* stripped */ /\\/\\//; } // real comment",
			"const ratio = value++ / 2; // real comment",
		}

		expected := []string{
			"function parse() { return /\"/; } ",
			"case /\"/.source: ",
			"const type = typeof /\"/; ",
			"const match = await /\"/; ",
			"const hasValue = name in /\"/; ",
			"const isPattern = value instanceof /\"/; ",
			"if (missing) {} else /\"/.test(value); ",
			"function skipComment() { return  /\\/\\//; } ",
			"const ratio = value++ / 2; ",
		}

		assertFiltered("sample.ts", RetainNone, lines, expected)
	})

	It("should honor JavaScript comment retention modes around literals", func() {
		lines := []string{
			"/** API docs. */",
			"/* setup block */",
			"const regex = /\\/\\/literal\\/\\*not-comment\\*\\//; // inline note",
			"const template = `Keep // and /* markers */ " +
				"${format(value /* inner block */)}`; // trailing note",
		}
		cases := []struct {
			mode     Mode
			expected []string
		}{
			{
				mode: RetainDocumentation,
				expected: []string{
					"/** API docs. */",
					"const regex = /\\/\\/literal\\/\\*not-comment\\*\\//; ",
					"const template = `Keep // and /* markers */ ${format(value )}`; ",
				},
			},
			{
				mode: RetainRegular,
				expected: []string{
					"/* setup block */",
					"const regex = /\\/\\/literal\\/\\*not-comment\\*\\//; // inline note",
					"const template = `Keep // and /* markers */ " +
						"${format(value /* inner block */)}`; // trailing note",
				},
			},
			{
				mode: RetainInline,
				expected: []string{
					"const regex = /\\/\\/literal\\/\\*not-comment\\*\\//; // inline note",
					"const template = `Keep // and /* markers */ ${format(value )}`; // trailing note",
				},
			},
			{
				mode: RetainBlock,
				expected: []string{
					"/* setup block */",
					"const regex = /\\/\\/literal\\/\\*not-comment\\*\\//; ",
					"const template = `Keep // and /* markers */ ${format(value /* inner block */)}`; ",
				},
			},
		}

		for _, tc := range cases {
			By(string(tc.mode))
			assertFiltered("sample.ts", tc.mode, lines, tc.expected)
		}
	})

	It("should not close block comments on an overlapping end marker", func() {
		lines := []string{
			"const before = 1; /*/ still comment */ const after = 2;",
		}

		expected := []string{
			"const before = 1;  const after = 2;",
		}

		assertFiltered("sample.ts", RetainNone, lines, expected)
	})

	It("should preserve JavaScript regex, string, and interpolation variants", func() {
		lines := []string{
			`const text = "escaped \\" // literal"; // comment`,
			`const classPattern = /[\\/]/gi; // comment`,
			`/start/.test(value); // comment`,
			`const open = /unterminated`,
			`const division = identifier / 2; // comment`,
			"const template = `escaped \\` text ${value}`; // comment",
			"const nested = `${{ value: /[}]/g, text: \"}\" }}`; // comment",
			`const unusual = */ /regex/; // comment`,
			`*/ /regex/; // comment`,
		}
		expected := []string{
			`const text = "escaped \\" `,
			`const classPattern = /[\\/]/gi; `,
			`/start/.test(value); `,
			`const open = /unterminated`,
			`const division = identifier / 2; `,
			"const template = `escaped \\` text ${value}`; ",
			"const nested = `${{ value: /[}]/g, text: \"}\" }}`; ",
			`const unusual = */ /regex/; `,
			`*/ /regex/; `,
		}

		assertFiltered("sample.ts", RetainNone, lines, expected)
	})
})
