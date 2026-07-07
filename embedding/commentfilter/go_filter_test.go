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

var _ = Describe("Go", func() {
	It("should strip all comments without treating literals as comments", func() {
		lines := []string{
			"// package comment",
			"package sample",
			"",
			"/* block comment */",
			"const slash = '/'",
			"const url = \"http://example.org\"",
			"const raw = `Keep // and /* markers */ in raw strings`",
			"const path = `C:\\Users\\`",
			"const multi = `",
			"Keep // and /* markers */ across lines",
			"`",
			"value := 1 // remove this real comment",
			"func create() {} // inline comment",
		}

		expected := []string{
			"package sample",
			"",
			"const slash = '/'",
			"const url = \"http://example.org\"",
			"const raw = `Keep // and /* markers */ in raw strings`",
			"const path = `C:\\Users\\`",
			"const multi = `",
			"Keep // and /* markers */ across lines",
			"`",
			"value := 1 ",
			"func create() {} ",
		}

		assertFiltered("sample.go", RetainNone, lines, expected)
	})

	It("should keep inline comments", func() {
		lines := []string{
			"// package comment",
			"package sample",
			"/* block comment */",
			"func create() {} // inline comment",
		}

		expected := []string{
			"// package comment",
			"package sample",
			"func create() {} // inline comment",
		}

		assertFiltered("sample.go", RetainInline, lines, expected)
	})

	It("should keep block comments", func() {
		lines := []string{
			"// package comment",
			"package sample",
			"/* block comment */",
			"func create() {} // inline comment",
		}

		expected := []string{
			"package sample",
			"/* block comment */",
			"func create() {} ",
		}

		assertFiltered("sample.go", RetainBlock, lines, expected)
	})
})
