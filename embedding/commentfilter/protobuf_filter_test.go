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

var _ = Describe("Protobuf", func() {
	It("should strip all comments without treating literals as comments", func() {
		lines := []string{
			"// file comment",
			"syntax = \"proto3\";",
			"",
			"/* message comment */",
			"message Sample {",
			"  string url = 1 [default = 'http://example.org'];",
			"  int32 count = 2; // inline comment",
			"}",
		}

		expected := []string{
			"syntax = \"proto3\";",
			"",
			"message Sample {",
			"  string url = 1 [default = 'http://example.org'];",
			"  int32 count = 2; ",
			"}",
		}

		assertFiltered("sample.proto", RetainNone, lines, expected)
	})

	It("should keep inline comments", func() {
		lines := []string{
			"// file comment",
			"syntax = \"proto3\";",
			"/* message comment */",
			"message Sample {} // inline comment",
		}

		expected := []string{
			"// file comment",
			"syntax = \"proto3\";",
			"message Sample {} // inline comment",
		}

		assertFiltered("sample.proto", RetainInline, lines, expected)
	})

	It("should keep block comments", func() {
		lines := []string{
			"// file comment",
			"syntax = \"proto3\";",
			"/* message comment */",
			"message Sample {} // inline comment",
		}

		expected := []string{
			"syntax = \"proto3\";",
			"/* message comment */",
			"message Sample {} ",
		}

		assertFiltered("sample.proto", RetainBlock, lines, expected)
	})
})
