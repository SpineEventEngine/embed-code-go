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
	"bytes"
	"log/slog"
	"testing"

	. "embed-code/embed-code-go/embedding/commentfilter"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// TestCommentFilter runs the comment filter test suite.
func TestCommentFilter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Comment Filter Suite")
}

var _ = Describe("Comment filter", func() {
	Describe("unsupported extensions", func() {
		It("should return unsupported files unchanged", func() {
			lines := []string{
				"# docs",
				"sub call { } # inline",
			}

			assertFiltered("service.pl", RetainAll, lines, lines)
		})

		It("should warn about unsupported comment modes", func() {
			output := captureWarnings(func() {
				Filter([]string{"# comment"}, "service.pl", RetainNone, "docs/guide.md", 12)
			})

			Expect(output).Should(ContainSubstring(
				"comment filtering is not supported for this file extension",
			))
			Expect(output).Should(ContainSubstring("file://"))
			Expect(output).Should(ContainSubstring("guide.md:12"))
		})
	})

	Describe("warnings", func() {
		It("should warn about modes without language-specific meaning", func() {
			output := captureWarnings(func() {
				Filter([]string{"<!-- comment -->"}, "layout.xml", RetainDocumentation, "docs/guide.md", 12)
			})

			Expect(output).Should(ContainSubstring("documentation"))
			Expect(output).Should(ContainSubstring("layout.xml"))
			Expect(output).Should(ContainSubstring("file://"))
			Expect(output).Should(ContainSubstring("guide.md:12"))
			Expect(output).Should(ContainSubstring("does not have a distinct meaning"))
		})

		It("should leave content unchanged when mode is not supported for file type", func() {
			lines := []string{
				"<root>",
				"  <!-- hidden -->",
				"</root>",
			}

			assertFiltered("layout.xml", RetainDocumentation, lines, lines)
		})
	})
})

// assertFiltered verifies filtering output for one file path and mode.
func assertFiltered(
	filePath string,
	mode Mode,
	lines []string,
	expected []string,
) {
	got := Filter(lines, filePath, mode, "docs/guide.md", 12)

	Expect(got).Should(Equal(expected))
}

// captureWarnings runs action and returns slog warning output.
func captureWarnings(action func()) string {
	var output bytes.Buffer
	previous := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&output, &slog.HandlerOptions{
		Level: slog.LevelWarn,
	})))
	defer slog.SetDefault(previous)

	action()

	return output.String()
}
