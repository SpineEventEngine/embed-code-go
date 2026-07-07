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

package parsing_test

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"embed-code/embed-code-go/configuration"
	"embed-code/embed-code-go/embedding/parsing"
	_type "embed-code/embed-code-go/type"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Parser states", func() {
	It("should consume a multiline instruction before opening an embedding fence", func() {
		config := configuration.NewConfiguration()
		context := newStateContext(
			"<embed-code",
			"    file=\"Example.java\"/>",
			"    ````java",
			"old source",
			"    ````",
		)

		Expect(parsing.EmbedInstruction.Recognize(context)).Should(BeTrue())
		Expect(parsing.EmbedInstruction.Accept(&context, config)).Should(Succeed())

		Expect(context.CurrentIndex()).Should(Equal(3))
		Expect(context.EmbeddingInstruction).ShouldNot(BeNil())
		Expect(context.EmbeddingInstruction.CodeFile).Should(Equal("Example.java"))
		Expect(context.EmbeddingInstruction.DocumentationLine).Should(Equal(1))
		Expect(context.GetResult()).Should(Equal([]string{
			"<embed-code",
			"    file=\"Example.java\"/>",
		}))

		Expect(parsing.CodeFenceStart.Recognize(context)).Should(BeTrue())
		Expect(parsing.CodeFenceStart.Accept(&context, config)).Should(Succeed())

		Expect(context.CurrentIndex()).Should(Equal(4))
		Expect(context.CodeFenceStarted).Should(BeTrue())
		Expect(context.CodeFenceMarker).Should(Equal("````"))
		Expect(context.CodeFenceIndentation).Should(Equal(4))
		Expect(context.CurrentEmbedding().SourceStartIndex).Should(Equal(3))
	})

	It("should recognize only matching embedding fence endings", func() {
		context := newStateContext(
			"    ```",
			"```",
			"    ````",
			"    ``` language",
		)
		context.CodeFenceStarted = true
		context.CodeFenceMarker = "```"
		context.CodeFenceIndentation = 4

		Expect(parsing.CodeFenceEnd.Recognize(context)).Should(BeTrue())

		context.ToNextLine()
		Expect(parsing.CodeFenceEnd.Recognize(context)).Should(BeFalse())

		context.ToNextLine()
		Expect(parsing.CodeFenceEnd.Recognize(context)).Should(BeTrue())

		context.ToNextLine()
		Expect(parsing.CodeFenceEnd.Recognize(context)).Should(BeFalse())
	})

	It("should report a malformed instruction when EOF is reached before the tag closes", func() {
		config := configuration.NewConfiguration()
		context := newStateContext(
			"preface",
			"<embed-code",
			"    file=\"Example.java\"",
		)
		Expect(parsing.RegularLine.Accept(&context, config)).Should(Succeed())

		err := parsing.EmbedInstruction.Accept(&context, config)

		Expect(err).Should(HaveOccurred())
		var parseErr parsing.InstructionParseError
		Expect(errors.As(err, &parseErr)).Should(BeTrue())
		Expect(parseErr.Line).Should(Equal(2))
		Expect(parseErr.Reason).Should(Equal(
			"the opening `<embed-code>` tag is not closed; add `>` or `/>` before the code fence",
		))
		Expect(context.ReachedEOF()).Should(BeTrue())
		Expect(context.GetResult()).Should(Equal([]string{
			"preface",
			"<embed-code",
			"    file=\"Example.java\"",
		}))
	})

	It("should report a missing opening tag end before a code fence", func() {
		config := configuration.NewConfiguration()
		context := newStateContext(
			"<embed-code",
			"file=\"$test/JunitTest.java\"",
			"start=\"Test\"",
			"end=\"assertEquals(2, value)\\n\\n\\n\"",
			"",
			"",
			"```",
			"old source",
			"```",
			"",
			"<embed-code",
			"file=\"$test/Hello.kt\"",
			"line=\"public void subscribe\">",
			"</embed-code>",
			"```",
			"```",
		)

		err := parsing.EmbedInstruction.Accept(&context, config)

		Expect(err).Should(HaveOccurred())
		var parseErr parsing.InstructionParseError
		Expect(errors.As(err, &parseErr)).Should(BeTrue())
		Expect(parseErr.Line).Should(Equal(1))
		Expect(parseErr.Reason).Should(Equal(
			"the opening `<embed-code>` tag is not closed; add `>` or `/>` before the code fence",
		))
	})

	It("should render source and close the embedding fence when the end state is accepted", func() {
		sourceRoot := GinkgoT().TempDir()
		Expect(os.WriteFile(
			filepath.Join(sourceRoot, "Example.java"),
			[]byte("class Example {}\n"),
			0600,
		)).To(Succeed())
		config := configuration.NewConfiguration()
		config.CodeRoots = _type.NamedPathList{_type.NamedPath{Path: sourceRoot}}
		context := newStateContext(
			"<embed-code file=\"Example.java\"/>",
			"```java",
			"old source",
			"```",
		)
		Expect(parsing.EmbedInstruction.Accept(&context, config)).Should(Succeed())
		Expect(parsing.CodeFenceStart.Accept(&context, config)).Should(Succeed())
		Expect(parsing.CodeSampleLine.Accept(&context, config)).Should(Succeed())

		Expect(parsing.CodeFenceEnd.Recognize(context)).Should(BeTrue())
		Expect(parsing.CodeFenceEnd.Accept(&context, config)).Should(Succeed())

		Expect(context.CodeFenceStarted).Should(BeFalse())
		Expect(context.EmbeddingInstruction).Should(BeNil())
		Expect(context.CurrentEmbedding().SourceEndIndex).Should(Equal(3))
		Expect(context.GetResult()).Should(Equal([]string{
			"<embed-code file=\"Example.java\"/>",
			"```java",
			"class Example {}",
			"```",
		}))
	})

	It("should report changed content when generated result is shorter than processed source", func() {
		config := configuration.NewConfiguration()
		context := newStateContext(
			"<embed-code file=\"Example.java\"/>",
			"```java",
			"old source",
		)
		Expect(parsing.EmbedInstruction.Accept(&context, config)).Should(Succeed())
		Expect(parsing.CodeFenceStart.Accept(&context, config)).Should(Succeed())
		Expect(parsing.CodeSampleLine.Accept(&context, config)).Should(Succeed())

		Expect(context.GetResult()).Should(Equal([]string{
			"<embed-code file=\"Example.java\"/>",
			"```java",
		}))
		Expect(context.IsContentChanged()).Should(BeTrue())
	})

	It("should report changed content when generated result is longer than processed source", func() {
		config := configuration.NewConfiguration()
		context := newStateContext("original source")
		Expect(parsing.RegularLine.Accept(&context, config)).Should(Succeed())
		context.Result = append(context.Result, "extra generated line")

		Expect(context.IsContentChanged()).Should(BeTrue())
	})
})

// newStateContext builds a parser context from in-memory source lines.
func newStateContext(lines ...string) parsing.Context {
	docPath := filepath.Join(GinkgoT().TempDir(), "doc.md")
	err := os.WriteFile(docPath, []byte(strings.Join(lines, "\n")), 0600)
	if err != nil {
		Fail("unexpected error while writing parser state fixture: " + err.Error())
	}
	context, err := parsing.NewContext(docPath)
	if err != nil {
		Fail("unexpected error while creating parser state context: " + err.Error())
	}

	return context
}
