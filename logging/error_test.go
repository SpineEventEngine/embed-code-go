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

package logging_test

import (
	"errors"
	"fmt"

	"embed-code/embed-code-go/logging"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Error formatting", func() {

	It("should format a single error inline", func() {
		err := errors.New("first failure")

		Expect(logging.FormatError("operation failed", err)).To(
			Equal("operation failed: first failure"),
		)
	})

	It("should format a nested error inline", func() {
		err := fmt.Errorf("outer context: %w", errors.New("first failure"))

		Expect(logging.FormatError("operation failed", err)).To(
			Equal("operation failed: outer context: first failure"),
		)
	})

	It("should format joined errors as a bullet list", func() {
		err := errors.Join(
			errors.New("first failure"),
			errors.New("second failure"),
		)

		Expect(logging.FormatError("operation failed", err)).To(
			Equal("operation failed:\n" +
				"  - first failure\n" +
				"  - second failure"),
		)
	})

	It("should flatten nested joined errors into a bullet list", func() {
		err := errors.Join(
			errors.New("first failure"),
			errors.Join(
				errors.New("second failure"),
				errors.New("third failure"),
			),
		)

		Expect(logging.FormatError("operation failed", err)).To(
			Equal("operation failed:\n" +
				"  - first failure\n" +
				"  - second failure\n" +
				"  - third failure"),
		)
	})
})
