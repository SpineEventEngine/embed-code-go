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

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type TestInstructionParams struct {
	fragment  string
	startGlob string
	endGlob   string
	comments  string
	closeTag  bool
}

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

	It("should have an error for unsupported comments mode", func() {
		instructionParams := TestInstructionParams{
			comments: "summary",
		}
		xmlString := buildInstruction("org/example/Comments.java", instructionParams)

		Expect(parsing.FromXML(xmlString, config)).Error().Should(HaveOccurred())
	})

	It("should successfully read source content", func() {
		instructionParams := TestInstructionParams{
			closeTag: true,
		}

		actualLines := getXMLExtractionContent(
			"org/example/Hello.java", instructionParams, config)

		expectedLength := 28
		checkedLine := 22
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
		Expect(actualLines).Should(ContainElement("    String create(String name); // end-of-line inline comment."))
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
		Expect(actualLines).Should(ContainElement("    String create(String name); // end-of-line inline comment."))
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

	It("should successfully parse XML by only end glob", func() {
		instructionParams := TestInstructionParams{
			endGlob: "package*",
		}

		actualLines := getXMLExtractionContent(
			"org/example/Hello.java", instructionParams, config)

		expectedLength := 21
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

	It("should embed one line when start and end globs match the same line", func() {
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
				"no line in code file `file://%s` matches the start pattern `foo bar`",
				absTestCodeFile("org/example/Hello.java"),
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
				"no line in code file `file://%s` matches the end pattern `foo bar`",
				absTestCodeFile("org/example/Hello.java"),
			),
		))
	})
})

func getXMLExtractionContent(fileName string, params TestInstructionParams,
	config configuration.Configuration) []string {
	xmlString := buildInstruction(fileName, params)
	instruction := createInstructionFromXML(xmlString, config)

	return readInstructionContent(instruction)
}

func buildConfigWithSourceFiles() configuration.Configuration {
	var config = configuration.NewConfiguration()
	config.DocumentationRoot = "../../test/resources/docs"
	config.CodeRoots = _type.NamedPathList{_type.NamedPath{Path: "../../test/resources/code/java"}}

	return config
}

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

func createInstructionFromXML(xmlString string,
	config configuration.Configuration) parsing.Instruction {
	instruction, err := parsing.FromXML(xmlString, config)
	if err != nil {
		Fail("unexpected error occurred during XML parsing: " + err.Error())
	}

	return instruction
}

func readInstructionContent(instruction parsing.Instruction) []string {
	lines, err := instruction.Content()
	if err != nil {
		Fail("unexpected error occurred while retrieving content: " + err.Error())
	}

	return lines
}

func absTestCodeFile(path string) string {
	absolutePath, err := filepath.Abs(filepath.Join("../../test/resources/code/java", path))
	if err != nil {
		Fail("unexpected error while resolving test code file: " + err.Error())
	}

	return absolutePath
}

func xmlAttribute(name string, value string) string {
	return fmt.Sprintf("%s=\"%v\"", name, value)
}
