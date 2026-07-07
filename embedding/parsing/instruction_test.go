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
	_type "embed-code/embed-code-go/type"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"embed-code/embed-code-go/configuration"
	"embed-code/embed-code-go/embedding/parsing"
	"embed-code/embed-code-go/logging"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// TestInstructionParams contains instruction attributes used by parser tests.
type TestInstructionParams struct {
	// fragment is the optional fragment name.
	fragment string

	// startGlob is the optional start pattern.
	startGlob string

	// endGlob is the optional end pattern.
	endGlob string

	// lineGlob is the optional single-line pattern.
	lineGlob string

	// comments is the requested comment filtering mode.
	comments string

	// closeTag reports whether the instruction includes a closing tag.
	closeTag bool
}

// TestInstruction runs the instruction parsing test suite.
func TestInstruction(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Data Suite")
}

var _ = Describe("Instruction", func() {
	var config configuration.Configuration

	BeforeEach(func() {
		currentDir, err := os.Getwd()
		if err != nil {
			Fail("unexpected error during the test setup: " + err.Error())
		}
		err = os.Chdir(currentDir)
		if err != nil {
			Fail("unexpected error during the test setup: " + err.Error())
		}
		config = buildConfigWithSourceFiles()
	})

	It("should have an error while parsing malformed XML string", func() {
		xmlString := "<file=\"org/example/Hello.java\" fragment=\"Hello class\"/>"

		Expect(parsing.FromXML(xmlString, config)).Error().Should(HaveOccurred())
	})

	It("should successfully parse XML with no errors", func() {
		instructionParams := TestInstructionParams{fragment: "Hello class"}
		xmlString := buildInstruction("org/example/Hello.java", instructionParams)

		Expect(parsing.FromXML(xmlString, config)).Error().ShouldNot(HaveOccurred())
	})

	It("should successfully parse XML with closing tag and with no errors", func() {
		instructionParams := TestInstructionParams{
			fragment: "Hello class",
			closeTag: true,
		}
		xmlString := buildInstruction("org/example/Hello.java", instructionParams)

		Expect(parsing.FromXML(xmlString, config)).Error().ShouldNot(HaveOccurred())
	})

	It("should parse backslash-escaped quotes in XML attributes", func() {
		xmlString := `<embed-code file="org/example/Hello.java" line="println(\"Hello world\")"/>`

		attributes, err := parsing.ParseXMLLine(xmlString)

		Expect(err).ShouldNot(HaveOccurred())
		Expect(attributes["line"]).Should(Equal(`println("Hello world")`))
	})

	It("should have an error for unsupported comments mode", func() {
		instructionParams := TestInstructionParams{
			comments: "summary",
		}
		xmlString := buildInstruction("org/example/Comments.java", instructionParams)

		Expect(parsing.FromXML(xmlString, config)).Error().Should(HaveOccurred())
	})

	It("should have an error for an invalid glob pattern", func() {
		instructionParams := TestInstructionParams{
			startGlob: "[",
		}
		xmlString := buildInstruction("org/example/Hello.java", instructionParams)

		_, err := parsing.FromXML(xmlString, config)

		Expect(err).Should(MatchError(ContainSubstring("invalid start pattern `[`")))
	})

	It("should have an error for an unclosed alternative pattern", func() {
		instructionParams := TestInstructionParams{
			startGlob: "0{$",
		}
		xmlString := buildInstruction("org/example/Hello.java", instructionParams)

		_, err := parsing.FromXML(xmlString, config)

		Expect(err).Should(MatchError(ContainSubstring("invalid start pattern `0{$`")))
	})

	It("should successfully read source content", func() {
		instructionParams := TestInstructionParams{
			closeTag: true,
		}

		actualLines := getXMLExtractionContent(
			"org/example/Hello.java", instructionParams, config)

		expectedLength := 34
		checkedLine := 28
		expectedLine := "public class Hello {"

		Expect(actualLines).Should(HaveLen(expectedLength))
		Expect(actualLines[checkedLine]).Should(Equal(expectedLine))
	})

	It("should strip all recognized comments", func() {
		instructionParams := TestInstructionParams{
			comments: "none",
		}

		actualLines := getXMLExtractionContent(
			"org/example/Comments.java", instructionParams, config)

		Expect(actualLines).Should(Equal([]string{
			"package org.example;",
			"",
			"public interface Comments {",
			"    String marker = \"http://example.org/*not-comment*/\";",
			"",
			"    String create(String name); ",
			"}",
		}))
	})

	It("should keep documentation comments only", func() {
		instructionParams := TestInstructionParams{
			comments: "documentation",
		}

		actualLines := getXMLExtractionContent(
			"org/example/Comments.java", instructionParams, config)

		Expect(actualLines).Should(ContainElement("/**"))
		Expect(actualLines).Should(ContainElement(" * Documents the public API."))
		Expect(actualLines).ShouldNot(ContainElement("     * The block comment."))
		Expect(actualLines).ShouldNot(ContainElement("    // Full-line inline comment."))
	})

	It("should keep inline comments only", func() {
		instructionParams := TestInstructionParams{
			comments: "inline",
		}

		actualLines := getXMLExtractionContent(
			"org/example/Comments.java", instructionParams, config)

		Expect(actualLines).Should(ContainElement("    // Full-line inline comment."))
		Expect(actualLines).Should(ContainElement(
			"    String create(String name); // end-of-line inline comment.",
		))
		Expect(actualLines).ShouldNot(ContainElement("/**"))
		Expect(actualLines).ShouldNot(ContainElement("     * The block comment."))
	})

	It("should keep block comments only", func() {
		instructionParams := TestInstructionParams{
			comments: "block",
		}

		actualLines := getXMLExtractionContent(
			"org/example/Comments.java", instructionParams, config)

		Expect(actualLines).ShouldNot(ContainElement("/**"))
		Expect(actualLines).ShouldNot(ContainElement(" * Documents the public API."))
		Expect(actualLines).Should(ContainElement("     * The block comment."))
		Expect(actualLines).ShouldNot(ContainElement("    // Full-line inline comment."))
	})

	It("should keep regular comments only", func() {
		instructionParams := TestInstructionParams{
			comments: "regular",
		}

		actualLines := getXMLExtractionContent(
			"org/example/Comments.java", instructionParams, config)

		Expect(actualLines).ShouldNot(ContainElement("/**"))
		Expect(actualLines).ShouldNot(ContainElement(" * Documents the public API."))
		Expect(actualLines).Should(ContainElement("     * The block comment."))
		Expect(actualLines).Should(ContainElement("    // Full-line inline comment."))
		Expect(actualLines).Should(ContainElement(
			"    String create(String name); // end-of-line inline comment.",
		))
	})

	It("should have an error when parsing fragment with start glob", func() {
		instructionParams := TestInstructionParams{
			fragment:  "fragment",
			startGlob: "public void hello()",
		}
		xmlString := buildInstruction("org/example/Hello.java", instructionParams)

		Expect(parsing.FromXML(xmlString, config)).Error().Should(HaveOccurred())
	})

	It("should have an error when parsing fragment with end glob", func() {
		instructionParams := TestInstructionParams{
			fragment: "fragment",
			endGlob:  "}",
		}
		xmlString := buildInstruction("org/example/Hello.java", instructionParams)

		Expect(parsing.FromXML(xmlString, config)).Error().Should(HaveOccurred())
	})

	It("should have an error when parsing fragment with line glob", func() {
		instructionParams := TestInstructionParams{
			fragment: "fragment",
			lineGlob: "public void hello()",
		}
		xmlString := buildInstruction("org/example/Hello.java", instructionParams)

		Expect(parsing.FromXML(xmlString, config)).Error().Should(HaveOccurred())
	})

	It("should have an error when parsing line glob with start glob", func() {
		instructionParams := TestInstructionParams{
			startGlob: "public class*",
			lineGlob:  "public void hello()",
		}
		xmlString := buildInstruction("org/example/Hello.java", instructionParams)

		Expect(parsing.FromXML(xmlString, config)).Error().Should(HaveOccurred())
	})

	It("should have an error when parsing line glob with end glob", func() {
		instructionParams := TestInstructionParams{
			endGlob:  "*System.out*",
			lineGlob: "public void hello()",
		}
		xmlString := buildInstruction("org/example/Hello.java", instructionParams)

		Expect(parsing.FromXML(xmlString, config)).Error().Should(HaveOccurred())
	})

	It("should successfully parse XML from start to end glob", func() {
		instructionParams := TestInstructionParams{
			startGlob: "public class*",
			endGlob:   "*System.out*",
		}

		actualLines := getXMLExtractionContent(
			"org/example/Hello.java", instructionParams, config)

		expectedLength := 4
		expectedFirstLine := "public class Hello {"
		expectedLastLine := "System.out.println(\"Hello world\");"

		Expect(actualLines).Should(HaveLen(expectedLength))
		Expect(actualLines[0]).Should(Equal(expectedFirstLine))
		Expect(strings.TrimLeft(actualLines[expectedLength-1], " ")).
			Should(Equal(expectedLastLine))
	})

	It("should successfully parse XML from start to end glob", func() {
		instructionParams := TestInstructionParams{
			startGlob: "*public static void main*",
			endGlob:   "*}*",
		}

		actualLines := getXMLExtractionContent(
			"org/example/Hello.java", instructionParams, config)

		expectedLength := 3
		expectedFirstLine := "public static void main(String[] args) {"
		expectedPattern := "^    "

		Expect(actualLines).Should(HaveLen(expectedLength))
		Expect(actualLines[0]).Should(Equal(expectedFirstLine))
		Expect(actualLines[1]).Should(MatchRegexp(expectedPattern))
	})

	It("should successfully parse XML from only start glob", func() {
		instructionParams := TestInstructionParams{
			startGlob: "*class*",
		}

		actualLines := getXMLExtractionContent(
			"org/example/Hello.java", instructionParams, config)

		expectedLength := 6
		expectedFirstLine := "public class Hello {"

		Expect(actualLines).Should(HaveLen(expectedLength))
		Expect(actualLines[0]).Should(Equal(expectedFirstLine))
	})

	It("should embed only the matching line when line glob is specified", func() {
		instructionParams := TestInstructionParams{
			lineGlob: "*class*",
		}

		actualLines := getXMLExtractionContent(
			"org/example/Hello.java", instructionParams, config)

		Expect(actualLines).Should(Equal([]string{
			"public class Hello {",
		}))
	})

	It("should embed a line with an escaped asterisk pattern", func() {
		instructionParams := TestInstructionParams{
			lineGlob: `Use \* to multiply`,
		}

		actualLines := getXMLExtractionContent(
			"literal-patterns.txt", instructionParams, config)

		Expect(actualLines).Should(Equal([]string{
			"Use * to multiply",
		}))
	})

	It("should embed a line starting with a literal caret pattern", func() {
		instructionParams := TestInstructionParams{
			lineGlob: "^^ starts with caret",
		}

		actualLines := getXMLExtractionContent(
			"literal-patterns.txt", instructionParams, config)

		Expect(actualLines).Should(Equal([]string{
			"^ starts with caret",
		}))
	})

	It("should embed a line ending with a literal dollar pattern", func() {
		instructionParams := TestInstructionParams{
			lineGlob: "The value ends with $$",
		}

		actualLines := getXMLExtractionContent(
			"literal-patterns.txt", instructionParams, config)

		Expect(actualLines).Should(Equal([]string{
			"The value ends with $",
		}))
	})

	It("should preserve pattern spaces that are not adjacent to a line separator", func() {
		instructionParams := TestInstructionParams{
			lineGlob: "^  padded text  $ \\n ^Use \\* to multiply$",
		}

		actualLines := getXMLExtractionContent(
			"literal-patterns.txt", instructionParams, config)

		Expect(actualLines).Should(Equal([]string{
			"  padded text  ",
			"Use * to multiply",
		}))
	})

	It("should preserve UTF-8 after trimming spaces following a line separator", func() {
		instructionParams := TestInstructionParams{
			lineGlob: "^Use \\* to multiply$ \\n żółć$",
		}

		actualLines := getXMLExtractionContent(
			"literal-patterns.txt", instructionParams, config)

		Expect(actualLines).Should(Equal([]string{
			"Use * to multiply",
			"żółć",
		}))
	})

	It("should successfully parse XML by only end glob", func() {
		instructionParams := TestInstructionParams{
			endGlob: "package*",
		}

		actualLines := getXMLExtractionContent(
			"org/example/Hello.java", instructionParams, config)

		expectedLength := 27
		expectedFirstLine := "/*"
		expectedLastLine := "package org.example;"

		Expect(actualLines).Should(HaveLen(expectedLength))
		Expect(actualLines[0]).Should(Equal(expectedFirstLine))
		Expect(actualLines[expectedLength-1]).Should(Equal(expectedLastLine))
	})

	It("should successfully parse XML by equal start and end glob", func() {
		instructionParams := TestInstructionParams{
			startGlob: "*main*",
			endGlob:   "*main*",
		}

		actualLines := getXMLExtractionContent(
			"org/example/Hello.java", instructionParams, config)

		expectedLength := 1
		expectedFirstLine := "public static void main(String[] args) {"

		Expect(actualLines).Should(HaveLen(expectedLength))
		Expect(actualLines[0]).Should(Equal(expectedFirstLine))
	})

	It("should successfully parse XML by globs without asterisks", func() {
		instructionParams := TestInstructionParams{
			startGlob: "main",
			endGlob:   "world",
		}

		actualLines := getXMLExtractionContent(
			"org/example/Hello.java", instructionParams, config)

		expectedLength := 2
		expectedFirstLinePattern := "^public static void main"
		expectedLastLinePattern := "^    System.out.println"

		Expect(actualLines).Should(HaveLen(expectedLength))
		Expect(actualLines[0]).Should(MatchRegexp(expectedFirstLinePattern))
		Expect(actualLines[1]).Should(MatchRegexp(expectedLastLinePattern))
	})

	It("should embed one line when the start and end globs match the same line", func() {
		instructionParams := TestInstructionParams{
			startGlob: "*spine.enableJava()*",
			endGlob:   "*.server()",
		}

		actualLines := getXMLExtractionContent(
			"examples/hello/build.gradle", instructionParams, config)

		Expect(actualLines).Should(Equal([]string{
			"spine.enableJava().server()",
		}))
	})

	It("should successfully parse XML by globs with line starts", func() {
		instructionParams := TestInstructionParams{
			startGlob: "^foo",
			endGlob:   "^bar",
		}

		actualLines := getXMLExtractionContent(
			"plain-text-to-embed.txt", instructionParams, config)

		expectedLength := 4
		expectedFirstLinePattern := "foo — this line starts with it"
		expectedLastLinePattern := "bar — this line starts with it"

		Expect(actualLines).Should(HaveLen(expectedLength))
		Expect(actualLines[0]).Should(MatchRegexp(expectedFirstLinePattern))
		Expect(actualLines[3]).Should(MatchRegexp(expectedLastLinePattern))
	})

	It("should successfully parse XML by globs with line ends", func() {
		instructionParams := TestInstructionParams{
			startGlob: "foo$",
			endGlob:   "bar$",
		}

		actualLines := getXMLExtractionContent(
			"plain-text-to-embed.txt", instructionParams, config)

		expectedLength := 6
		expectedFirstLine := "This line ends with foo"
		expectedLastLine := "This line ends with bar"

		Expect(actualLines).Should(HaveLen(expectedLength))
		Expect(actualLines[0]).Should(Equal(expectedFirstLine))
		Expect(actualLines[5]).Should(Equal(expectedLastLine))
	})

	It("should report an error when start glob does not match", func() {
		instructionParams := TestInstructionParams{
			startGlob: "foo bar",
			endGlob:   "*main*",
		}
		xmlString := buildInstruction("org/example/Hello.java", instructionParams)
		instruction := createInstructionFromXML(xmlString, config)

		_, err := instruction.Content()

		Expect(err).Should(MatchError(
			fmt.Sprintf(
				"no line in code file `%s` matches the start pattern `foo bar`",
				logging.FileReference(absTestCodeFile("org/example/Hello.java")),
			),
		))
	})

	It("should report an error when end glob does not match", func() {
		instructionParams := TestInstructionParams{
			startGlob: "*main*",
			endGlob:   "foo bar",
		}
		xmlString := buildInstruction("org/example/Hello.java", instructionParams)
		instruction := createInstructionFromXML(xmlString, config)

		_, err := instruction.Content()

		Expect(err).Should(MatchError(
			fmt.Sprintf(
				"no line in code file `%s` matches the end pattern `foo bar`",
				logging.FileReference(absTestCodeFile("org/example/Hello.java")),
			),
		))
	})

	It("should use shared resolver as default", func() {
		sourceRoot := GinkgoT().TempDir()
		markdownPath := filepath.Join(GinkgoT().TempDir(), "doc.md")
		codePath := filepath.Join(sourceRoot, "Example.java")
		config.CodeRoots = _type.NamedPathList{_type.NamedPath{Path: sourceRoot}}
		Expect(os.WriteFile(markdownPath, []byte(""), 0600)).To(Succeed())
		Expect(os.WriteFile(
			codePath,
			[]byte("class Example { String version = \"first\"; }"),
			0600,
		)).To(Succeed())
		context, err := parsing.NewContextWithResolver(markdownPath, nil)
		Expect(err).ShouldNot(HaveOccurred())
		context.StartEmbedding(parsing.Instruction{
			CodeFile:      "Example.java",
			Configuration: config,
		})

		firstContent, err := context.EmbeddingInstruction.Content()
		Expect(err).ShouldNot(HaveOccurred())
		Expect(os.WriteFile(
			codePath,
			[]byte("class Example { String version = \"second\"; }"),
			0600,
		)).To(Succeed())
		secondContent, err := context.EmbeddingInstruction.Content()
		Expect(err).ShouldNot(HaveOccurred())

		Expect(secondContent).Should(Equal(firstContent))
		Expect(secondContent).Should(Equal([]string{
			"class Example { String version = \"first\"; }",
		}))
	})
})

// getXMLExtractionContent returns source lines selected by an XML instruction fixture.
func getXMLExtractionContent(fileName string, params TestInstructionParams,
	config configuration.Configuration) []string {
	xmlString := buildInstruction(fileName, params)
	instruction := createInstructionFromXML(xmlString, config)

	return readInstructionContent(instruction)
}

// buildConfigWithSourceFiles returns a configuration using parser source fixtures.
func buildConfigWithSourceFiles() configuration.Configuration {
	var config = configuration.NewConfiguration()
	config.DocumentationRoot = "../../test/resources/docs"
	config.CodeRoots = _type.NamedPathList{_type.NamedPath{Path: "../../test/resources/code/java"}}

	return config
}

// buildInstruction builds an XML instruction string from test parameters.
func buildInstruction(fileName string, params TestInstructionParams) string {
	fragmentAttr := xmlAttribute("fragment", params.fragment)
	instructionLine := fmt.Sprintf("<embed-code file=\"%s\" %s", fileName, fragmentAttr)

	if len(params.startGlob) > 0 {
		startAttr := xmlAttribute("start", params.startGlob)
		instructionLine += " " + startAttr
	}
	if len(params.endGlob) > 0 {
		endAttr := xmlAttribute("end", params.endGlob)
		instructionLine += " " + endAttr
	}
	if len(params.lineGlob) > 0 {
		lineAttr := xmlAttribute("line", params.lineGlob)
		instructionLine += " " + lineAttr
	}
	if len(params.comments) > 0 {
		commentsAttr := xmlAttribute("comments", params.comments)
		instructionLine += " " + commentsAttr
	}
	if params.closeTag {
		instructionLine += "></embed-code>"
	} else {
		instructionLine += "/>"
	}

	return instructionLine
}

// createInstructionFromXML parses an XML instruction string for tests.
func createInstructionFromXML(xmlString string,
	config configuration.Configuration) parsing.Instruction {
	instruction, err := parsing.FromXML(xmlString, config)
	if err != nil {
		Fail("unexpected error occurred during XML parsing: " + err.Error())
	}

	return instruction
}

// readInstructionContent returns instruction content and fails the spec on errors.
func readInstructionContent(instruction parsing.Instruction) []string {
	lines, err := instruction.Content()
	if err != nil {
		Fail("unexpected error occurred while retrieving content: " + err.Error())
	}

	return lines
}

// absTestCodeFile returns an absolute path to a parser source fixture.
func absTestCodeFile(path string) string {
	absolutePath, err := filepath.Abs(filepath.Join("../../test/resources/code/java", path))
	if err != nil {
		Fail("unexpected error while resolving test code file: " + err.Error())
	}

	return absolutePath
}

// xmlAttribute formats one XML attribute for instruction fixtures.
func xmlAttribute(name string, value string) string {
	return fmt.Sprintf("%s=\"%v\"", name, value)
}
