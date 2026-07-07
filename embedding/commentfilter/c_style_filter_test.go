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

var _ = Describe("C and C++", func() {
	It("should strip all comments without treating literals as comments", func() {
		lines := []string{
			"// header comment",
			"#include <stdio.h>",
			"",
			"/* block comment */",
			"const char slash = '/';",
			"const char* url = \"http://example.org\";",
			"int create() { return 1; } // inline comment",
		}

		expected := []string{
			"#include <stdio.h>",
			"",
			"const char slash = '/';",
			"const char* url = \"http://example.org\";",
			"int create() { return 1; } ",
		}

		assertFiltered("sample.cpp", RetainNone, lines, expected)
	})

	It("should keep inline comments", func() {
		lines := []string{
			"// header comment",
			"int create();",
			"/* block comment */",
			"int count(); // inline comment",
		}

		expected := []string{
			"// header comment",
			"int create();",
			"int count(); // inline comment",
		}

		assertFiltered("sample.cpp", RetainInline, lines, expected)
	})

	It("should keep block comments", func() {
		lines := []string{
			"// header comment",
			"int create();",
			"/* block comment */",
			"int count(); // inline comment",
		}

		expected := []string{
			"int create();",
			"/* block comment */",
			"int count(); ",
		}

		assertFiltered("sample.hpp", RetainBlock, lines, expected)
	})

	It("should not close block comments on an overlapping end marker", func() {
		lines := []string{
			"int before = 1; /*/ still comment */ int after = 2;",
		}

		expected := []string{
			"int before = 1;  int after = 2;",
		}

		assertFiltered("sample.cpp", RetainNone, lines, expected)
	})
})
