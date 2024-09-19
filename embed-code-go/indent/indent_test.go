// Copyright 2024, TeamDev. All rights reserved.
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

package indent_test

import (
	"embed-code/embed-code-go/indent"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"testing"
)

func TestIndent(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Data Suite")
}

var _ = Describe("Indent", func() {

	It("should not find indentations", func() {
		testLines := []string{"", "foo", "bar", "", "baz", ""}

		Expect(indent.MaxCommonIndentation(testLines)).Should(BeZero())
	})

	It("should not find indentations as the given lines asre nil", func() {
		var testLines []string

		Expect(indent.MaxCommonIndentation(testLines)).Should(BeZero())
	})

	It("should properly find indentations", func() {
		testLines := []string{"", "  foo", "    bar", "", "", "  baz"}
		expectedIndents := 2

		Expect(indent.MaxCommonIndentation(testLines)).Should(Equal(expectedIndents))
	})

	It("should properly cut indentations", func() {
		testLines := []string{"", "  foo", "    bar", "", "", "  baz"}
		changedLines := indent.CutIndent(testLines, 2)

		Expect(changedLines).ShouldNot(Equal(testLines))
	})

})
