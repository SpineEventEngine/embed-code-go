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

var _ = Describe("Visual Basic", func() {
	It("should strip all comments", func() {
		lines := []string{
			"' file comment",
			"REM module comment",
			"Dim text = \"REM not a comment\"",
			"Dim quotedRem = \"\"\"REM not a comment\"",
			"Dim quotedApostrophe = \"\"\"' not a comment\"",
			"Dim escapedQuote = \"Say \"\"REM\"\" and keep going\" ' inline",
			"Dim value = 1 ' inline",
			"Dim ready = True : Rem after statement separator",
			"Dim reminder = 1",
		}

		expected := []string{
			"Dim text = \"REM not a comment\"",
			"Dim quotedRem = \"\"\"REM not a comment\"",
			"Dim quotedApostrophe = \"\"\"' not a comment\"",
			"Dim escapedQuote = \"Say \"\"REM\"\" and keep going\" ",
			"Dim value = 1 ",
			"Dim ready = True : ",
			"Dim reminder = 1",
		}

		assertFiltered("Module.vb", RetainNone, lines, expected)
	})

	It("should keep regular comments", func() {
		lines := []string{
			"''' <summary>Creates a value.</summary>",
			"' file comment",
			"REM module comment",
			"Dim value = 1 ' inline",
		}

		expected := []string{
			"' file comment",
			"REM module comment",
			"Dim value = 1 ' inline",
		}

		assertFiltered("Module.vb", RetainRegular, lines, expected)
	})

	It("should keep documentation comments", func() {
		lines := []string{
			"''' <summary>Creates a value.</summary>",
			"' implementation note",
			"REM module comment",
			"Public Function Create() As String",
		}

		expected := []string{
			"''' <summary>Creates a value.</summary>",
			"Public Function Create() As String",
		}

		assertFiltered("Module.vb", RetainDocumentation, lines, expected)
	})
})
