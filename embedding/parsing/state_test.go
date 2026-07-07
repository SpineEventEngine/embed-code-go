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
		Expect(parseErr.Reason).Should(Equal("the `<embed-code>` tag is not closed"))
		Expect(context.ReachedEOF()).Should(BeTrue())
		Expect(context.GetResult()).Should(Equal([]string{
			"preface",
			"<embed-code",
			"    file=\"Example.java\"",
		}))
	})

	// The following specs reproduce
	// https://github.com/SpineEventEngine/embed-code-go/issues/19.
	//
	// A self-closed `<embed-code .../>` instruction that fails validation is a
	// complete, single-line instruction, yet `Accept` keeps re-parsing the
	// joined body and consumes every following line up to EOF before reporting
	// the failure. That is why a single malformed instruction near the top of a
	// document surfaces its error many lines later (line 76 -> line 404 in the
	// original report) and destroys the parse of everything after it.
	//
	// The parser should report the InstructionParseError but stop right after
	// the instruction's `/>` terminator, leaving the code fence and the rest of
	// the document intact. These specs assert that bounded behavior and
	// therefore FAIL until the over-consumption is fixed.
	It("should not consume the rest of the document when a start pattern is invalid", func() {
		assertBoundedMalformedInstruction(`<embed-code file="Example.java" start="["/>`)
	})

	It("should not consume the rest of the document when the comments mode is invalid", func() {
		assertBoundedMalformedInstruction(`<embed-code file="Example.java" comments="bogus"/>`)
	})

	It("should not consume the rest of the document when attributes are mutually exclusive", func() {
		assertBoundedMalformedInstruction(`<embed-code file="Example.java" fragment="f" line="l"/>`)
	})

	// Reproduces https://github.com/SpineEventEngine/embed-code-go/issues/19.
	//
	// A self-closed instruction whose start/end/line pattern contains raw XML
	// metacharacters (`<`, `&`) is legitimate user input (e.g. matching a Java
	// generic or comparison), yet `quoteEscapedXMLLine` only escapes `\"`, so
	// `xml.Unmarshal` rejects it. Because the tag is self-closed, this is NOT
	// the "tag is not closed" case; instead `Accept` keeps re-parsing and
	// consumes every following line up to EOF, then fails with
	// InstructionParseError -- swallowing the code fence and the rest of the
	// document, exactly the behavior reported in the issue.
	//
	// This test asserts the desired behavior and therefore FAILS until the XML
	// metacharacters in attribute values are escaped before unmarshalling.
	It("should parse an instruction whose pattern contains XML metacharacters", func() {
		config := configuration.NewConfiguration()
		context := newStateContext(
			"<embed-code file=\"Example.java\" line=\"if (a < b)\"/>",
			"```java",
			"old source",
			"```",
			"text after the fence",
		)

		Expect(parsing.EmbedInstruction.Recognize(context)).Should(BeTrue())

		err := parsing.EmbedInstruction.Accept(&context, config)

		Expect(err).ShouldNot(HaveOccurred())
		Expect(context.EmbeddingInstruction).ShouldNot(BeNil())
		Expect(context.EmbeddingInstruction.LinePattern).ShouldNot(BeNil())
		// The parser must stop right after the instruction line rather than
		// swallowing the code fence and the rest of the document.
		Expect(context.ReachedEOF()).Should(BeFalse())
		Expect(context.GetResult()).Should(Equal([]string{
			"<embed-code file=\"Example.java\" line=\"if (a < b)\"/>",
		}))
	})

	// Regresses https://github.com/SpineEventEngine/embed-code-go/issues/19.
	//
	// The `/>` inside the `line` value must not be mistaken for the tag
	// terminator: the instruction spans several lines and only closes on the
	// last one. Recognizing the inner `/>` would stop accumulation early and
	// fail with an "unexpected EOF" before the real closing `/>`.
	It("should accept a multiline instruction whose value contains a slash-close sequence", func() {
		config := configuration.NewConfiguration()
		context := newStateContext(
			"<embed-code",
			`    line="<br/>"`,
			`    file="Example.java"/>`,
			"```java",
			"old source",
			"```",
		)

		Expect(parsing.EmbedInstruction.Recognize(context)).Should(BeTrue())
		Expect(parsing.EmbedInstruction.Accept(&context, config)).Should(Succeed())

		Expect(context.EmbeddingInstruction).ShouldNot(BeNil())
		Expect(context.EmbeddingInstruction.CodeFile).Should(Equal("Example.java"))
		Expect(context.EmbeddingInstruction.LinePattern).ShouldNot(BeNil())
		Expect(context.GetResult()).Should(Equal([]string{
			"<embed-code",
			`    line="<br/>"`,
			`    file="Example.java"/>`,
		}))
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

// assertBoundedMalformedInstruction drives EmbedInstruction.Accept over a
// document whose first line is a self-closed but invalid instruction, followed
// by a code fence and trailing content.
//
// It asserts that the parser reports the failure as an InstructionParseError
// without consuming past the instruction line. See issue #19.
func assertBoundedMalformedInstruction(instruction string) {
	config := configuration.NewConfiguration()
	context := newStateContext(
		instruction,
		"```java",
		"old source",
		"```",
		"text after the fence",
	)

	Expect(parsing.EmbedInstruction.Recognize(context)).Should(BeTrue())

	err := parsing.EmbedInstruction.Accept(&context, config)

	var parseErr parsing.InstructionParseError
	Expect(errors.As(err, &parseErr)).Should(BeTrue())
	Expect(context.ReachedEOF()).Should(BeFalse())
	Expect(context.GetResult()).Should(Equal([]string{instruction}))
}

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
