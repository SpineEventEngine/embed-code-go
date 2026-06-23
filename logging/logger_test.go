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

package logging //nolint:testpackage // Tests OS-specific normalization in an unexported helper.

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// TestLogging runs the logging package specs.
func TestLogging(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Logging Suite")
}

var _ = Describe("File URL formatting", func() {

	It("should format a Unix path", func() {
		Expect(fileURLFromAbsolutePath("/Users/me/project/file.go")).To(
			Equal("file:///Users/me/project/file.go"),
		)
	})

	It("should format a Windows drive path", func() {
		Expect(fileURLFromAbsolutePath(`C:\Users\me\project\file.go`)).To(
			Equal("file:///C:/Users/me/project/file.go"),
		)
	})

	It("should escape spaces in a Windows drive path", func() {
		Expect(fileURLFromAbsolutePath(`C:\Users\me\my project\file.go`)).To(
			Equal("file:///C:/Users/me/my%20project/file.go"),
		)
	})

	It("should format a Windows UNC path", func() {
		Expect(fileURLFromAbsolutePath(`\\server\share\project\file.go`)).To(
			Equal("file://server/share/project/file.go"),
		)
	})
})
