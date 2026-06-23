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

package logging

import (
	"fmt"
	"strings"
)

// FormatError formats a single error inline and joined errors as a bullet list.
func FormatError(message string, err error) string {
	errs := flattenedErrors(err)
	if len(errs) <= 1 {
		return fmt.Sprintf("%s: %v", message, err)
	}

	var builder strings.Builder
	builder.WriteString(message)
	builder.WriteString(":")
	for _, nestedErr := range errs {
		builder.WriteString("\n  - ")
		builder.WriteString(nestedErr.Error())
	}

	return builder.String()
}

// flattenedErrors returns the leaf errors from a joined error.
func flattenedErrors(err error) []error {
	joined, ok := err.(interface {
		Unwrap() []error
	})
	if !ok {
		return []error{err}
	}

	var result []error
	for _, nestedErr := range joined.Unwrap() {
		result = append(result, flattenedErrors(nestedErr)...)
	}

	return result
}
