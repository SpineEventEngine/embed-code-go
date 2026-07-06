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
})
