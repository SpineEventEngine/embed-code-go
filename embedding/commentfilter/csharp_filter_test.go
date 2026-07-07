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

var _ = Describe("C#", func() {
	It("should keep XML documentation comments", func() {
		lines := []string{
			"/// <summary>Creates a value.</summary>",
			"// implementation note",
			"public string Create() => \"// literal\";",
		}

		expected := []string{
			"/// <summary>Creates a value.</summary>",
			"public string Create() => \"// literal\";",
		}

		assertFiltered("Api.cs", RetainDocumentation, lines, expected)
	})

	It("should keep inline comments", func() {
		lines := []string{
			"/// <summary>Creates a value.</summary>",
			"// implementation note",
			"public string Create() => \"// literal\";",
		}

		expected := []string{
			"// implementation note",
			"public string Create() => \"// literal\";",
		}

		assertFiltered("Api.cs", RetainInline, lines, expected)
	})

	It("should strip comments after character literals", func() {
		lines := []string{
			`var quote = '"'; // inline comment`,
			`var slash = '/'; // inline comment`,
		}

		expected := []string{
			`var quote = '"'; `,
			`var slash = '/'; `,
		}

		assertFiltered("Api.cs", RetainNone, lines, expected)
	})

	It("should strip comments without treating verbatim string text as comments", func() {
		lines := []string{
			"// header comment",
			`var uri = @"https://example.com/*not-comment*/"; // inline comment`,
			`var block = @"Keep // marker`,
			`and /* marker */ text"; /* trailing block */`,
		}

		expected := []string{
			`var uri = @"https://example.com/*not-comment*/"; `,
			`var block = @"Keep // marker`,
			`and /* marker */ text"; `,
		}

		assertFiltered("Api.cs", RetainNone, lines, expected)
	})

	It("should strip comments without treating interpolated string text as comments", func() {
		lines := []string{
			`var message = $"Keep // and /* markers */ {Format(value /* real comment */)}";`,
			`var escaped = $"Keep {{ // text }} and {value}"; // inline comment`,
			`var nested = $"Value {Format("/* not comment */")} // still text"; // inline comment`,
			`var url = $"{scheme://example.com/*path*/}"; // inline comment`,
		}

		expected := []string{
			`var message = $"Keep // and /* markers */ {Format(value )}";`,
			`var escaped = $"Keep {{ // text }} and {value}"; `,
			`var nested = $"Value {Format("/* not comment */")} // still text"; `,
			`var url = $"{scheme://example.com/*path*/}"; `,
		}

		assertFiltered("Api.cs", RetainNone, lines, expected)
	})

	It("should strip comments without treating verbatim interpolated string text as comments", func() {
		lines := []string{
			`var path = $@"C:\Temp\// not comment {name /* real comment */}";`,
			`var template = @$"Keep /* marker */ and ""// marker""`,
			`with {Format(value // real comment`,
			`)}"; // trailing comment`,
		}

		expected := []string{
			`var path = $@"C:\Temp\// not comment {name }";`,
			`var template = @$"Keep /* marker */ and ""// marker""`,
			`with {Format(value `,
			`)}"; `,
		}

		assertFiltered("Api.cs", RetainNone, lines, expected)
	})

	It("should strip comments without treating raw string text as comments", func() {
		lines := []string{
			`var raw = """`,
			`Keep // text`,
			`Keep /* text */`,
			`"""; // trailing comment`,
			`var interpolated = $"""`,
			`Keep // text {Format(value /* real comment */)}`,
			`"""; // trailing comment`,
		}

		expected := []string{
			`var raw = """`,
			`Keep // text`,
			`Keep /* text */`,
			`"""; `,
			`var interpolated = $"""`,
			`Keep // text {Format(value )}`,
			`"""; `,
		}

		assertFiltered("Api.cs", RetainNone, lines, expected)
	})

	It("should not close block comments on an overlapping end marker", func() {
		lines := []string{
			"var before = 1; /*/ still comment */ var after = 2;",
		}

		expected := []string{
			"var before = 1;  var after = 2;",
		}

		assertFiltered("Api.cs", RetainNone, lines, expected)
	})
})
