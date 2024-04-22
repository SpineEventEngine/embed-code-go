// Copyright 2024, TeamDev. All rights reserved.
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

package embedding_instruction_test

import (
	"embed-code/embed-code-go/configuration"
	"embed-code/embed-code-go/embedding_instruction"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"
)

func buildConfigWithPreparedFragments() configuration.Configuration {
	var config = configuration.NewConfiguration()
	config.DocumentationRoot = "./test/resources/docs"
	config.CodeRoot = "./test/resources/code"
	config.FragmentsDir = "./test/resources/prepared-fragments"
	return config
}

func buildInstruction(fileName string, params buildInstructionParams) string {
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
	if params.closeTag {
		instructionLine += "></embed-code>"
	} else {
		instructionLine += "/>"
	}
	return instructionLine
}

func xmlAttribute(name string, value string) string {
	return fmt.Sprintf("%s=\"%v\"", name, value)
}

type buildInstructionParams struct {
	fragment  string
	startGlob string
	endGlob   string
	closeTag  bool
}

type EmbeddingInstructionTestSuite struct {
	suite.Suite
	config configuration.Configuration
}

func (suite *EmbeddingInstructionTestSuite) SetupSuite() {
	rootDir, err := filepath.Abs("../../")
	if err != nil {
		panic(err)
	}
	os.Chdir(rootDir)
	suite.config = buildConfigWithPreparedFragments()
}

func (suite *EmbeddingInstructionTestSuite) TestParsingMisformedXML() {
	xmlString := "<file=\"org/example/Hello.java\" fragment=\"Hello class\"/>"
	config := buildConfigWithPreparedFragments()

	_, err := embedding_instruction.FromXML(xmlString, config)
	suite.Require().Error(err, "Parsing misformed XML should cause an error.")
}

func (suite *EmbeddingInstructionTestSuite) TestParseFromXML() {
	instructionParams := buildInstructionParams{
		fragment: "Hello class",
	}
	xmlString := buildInstruction("org/example/Hello.java", instructionParams)
	config := buildConfigWithPreparedFragments()
	_, err := embedding_instruction.FromXML(xmlString, config)
	suite.NoError(err, "There was unexpected error during the XML parsing")
}

func (suite *EmbeddingInstructionTestSuite) TestParseWithClosingTag() {
	instructionParams := buildInstructionParams{
		fragment: "Hello class",
		closeTag: true,
	}
	xmlString := buildInstruction("org/example/Hello.java", instructionParams)
	config := buildConfigWithPreparedFragments()

	_, err := embedding_instruction.FromXML(xmlString, config)
	suite.NoError(err, "There was unexpected error during the XML parsing")
}

func (suite *EmbeddingInstructionTestSuite) TestReadFragmentDir() {
	instructionParams := buildInstructionParams{
		closeTag: true,
	}
	xmlString := buildInstruction("org/example/Hello.java", instructionParams)
	config := buildConfigWithPreparedFragments()

	instruction, err := embedding_instruction.FromXML(xmlString, config)
	suite.Require().NoError(err, "There was unexpected error during the XML parsing")

	lines := instruction.Content()

	suite.Len(lines, 28, "the length of the lines is incorrect")
	suite.Equal("public class Hello {", lines[22], "the line at the 22 index is incorrect")
}

func (suite *EmbeddingInstructionTestSuite) TestFragmentAndStart() {
	instructionParams := buildInstructionParams{
		fragment:  "fragment",
		startGlob: "public void hello()",
	}

	xmlString := buildInstruction("org/example/Hello.java", instructionParams)
	config := buildConfigWithPreparedFragments()

	_, err := embedding_instruction.FromXML(xmlString, config)
	suite.Error(err, "Instruction tag with both fragment and startGlob provided should cause an error.")
}

func (suite *EmbeddingInstructionTestSuite) TestFragmentAndEnd() {
	instructionParams := buildInstructionParams{
		fragment: "fragment",
		endGlob:  "}",
	}

	xmlString := buildInstruction("org/example/Hello.java", instructionParams)
	config := buildConfigWithPreparedFragments()

	_, err := embedding_instruction.FromXML(xmlString, config)
	suite.Error(err, "Instruction tag with both fragment and endGlob provided should cause an error.")
}

func (suite *EmbeddingInstructionTestSuite) TestExtractByGlob() {
	instructionParams := buildInstructionParams{
		startGlob: "public class*",
		endGlob:   "*System.out*",
	}

	xmlString := buildInstruction("org/example/Hello.java", instructionParams)
	config := buildConfigWithPreparedFragments()

	instruction, err := embedding_instruction.FromXML(xmlString, config)
	suite.Require().NoError(err, "There was unexpected error during the XML parsing")

	lines := instruction.Content()

	suite.Len(lines, 4, "the length of the lines should be 4")
	suite.Equal("public class Hello {", lines[0], "the line at the 0 index is not equal to \"public class Hello {\"")
	suite.Equal("        System.out.println(\"Hello world\");", lines[3], "the line at the 3 index is incorrect")
}

func (suite *EmbeddingInstructionTestSuite) TestMinIndentation() {
	instructionParams := buildInstructionParams{
		startGlob: "*public static void main*",
		endGlob:   "*}*",
	}
	xmlString := buildInstruction("org/example/Hello.java", instructionParams)
	config := buildConfigWithPreparedFragments()

	instruction, err := embedding_instruction.FromXML(xmlString, config)
	suite.Require().NoError(err, "There was unexpected error during the XML parsing")

	lines := instruction.Content()

	suite.Len(lines, 3)
	suite.Equal("public static void main(String[] args) {", lines[0])
	suite.Regexp("^    ", lines[1])
}

func (suite *EmbeddingInstructionTestSuite) TestStartWithoutEnd() {
	instructionParams := buildInstructionParams{
		startGlob: "*class*",
	}
	xmlString := buildInstruction("org/example/Hello.java", instructionParams)
	config := buildConfigWithPreparedFragments()

	instruction, err := embedding_instruction.FromXML(xmlString, config)
	suite.Require().NoError(err, "There was unexpected error during the XML parsing")

	lines := instruction.Content()

	suite.Len(lines, 6)
	suite.Equal("}", lines[5])
}

func (suite *EmbeddingInstructionTestSuite) TestEndWithoutStart() {
	instructionParams := buildInstructionParams{
		endGlob: "package*",
	}
	xmlString := buildInstruction("org/example/Hello.java", instructionParams)
	config := buildConfigWithPreparedFragments()

	instruction, err := embedding_instruction.FromXML(xmlString, config)
	suite.Require().NoError(err, "There was unexpected error during the XML parsing")

	lines := instruction.Content()

	suite.Len(lines, 21)
	suite.Equal("/*", lines[0])
	suite.Equal("package org.example;", lines[20])
}

func (suite *EmbeddingInstructionTestSuite) TestOneLine() {
	instructionParams := buildInstructionParams{
		startGlob: "*main*",
		endGlob:   "*main*",
	}
	xmlString := buildInstruction("org/example/Hello.java", instructionParams)
	config := buildConfigWithPreparedFragments()

	instruction, err := embedding_instruction.FromXML(xmlString, config)
	suite.Require().NoError(err, "There was unexpected error during the XML parsing")

	lines := instruction.Content()

	suite.Len(lines, 1)
	suite.Equal("public static void main(String[] args) {", lines[0])
}

func (suite *EmbeddingInstructionTestSuite) TestNoMatchStart() {
	instructionParams := buildInstructionParams{
		startGlob: "foo bar",
		endGlob:   "*main*",
	}

	xmlString := buildInstruction("org/example/Hello.java", instructionParams)
	config := buildConfigWithPreparedFragments()

	instruction, err := embedding_instruction.FromXML(xmlString, config)
	suite.Require().NoError(err, "There was unexpected error during the XML parsing")

	suite.Require().Panics(func() {
		instruction.Content()
	})
}

func (suite *EmbeddingInstructionTestSuite) TestNoMatchEnd() {
	instructionParams := buildInstructionParams{
		startGlob: "*main*",
		endGlob:   "foo bar",
	}

	xmlString := buildInstruction("org/example/Hello.java", instructionParams)
	config := buildConfigWithPreparedFragments()

	instruction, err := embedding_instruction.FromXML(xmlString, config)
	suite.Require().NoError(err, "There was unexpected error during the XML parsing")

	suite.Require().Panics(func() {
		instruction.Content()
	})
}

func (suite *EmbeddingInstructionTestSuite) TestImplyAsterisk() {
	instructionParams := buildInstructionParams{
		startGlob: "main",
		endGlob:   "world",
	}
	xmlString := buildInstruction("org/example/Hello.java", instructionParams)
	config := buildConfigWithPreparedFragments()

	instruction, err := embedding_instruction.FromXML(xmlString, config)
	suite.Require().NoError(err, "There was unexpected error during the XML parsing")

	lines := instruction.Content()

	suite.Len(lines, 2)
	suite.Regexp("^public static void main", lines[0])
	suite.Regexp("^    System.out.println", lines[1])
}

func (suite *EmbeddingInstructionTestSuite) TestExplicitLineStart() {
	instructionParams := buildInstructionParams{
		startGlob: "^foo",
		endGlob:   "^bar",
	}
	xmlString := buildInstruction("plain-text-to-embed.txt", instructionParams)
	config := buildConfigWithPreparedFragments()

	instruction, err := embedding_instruction.FromXML(xmlString, config)
	suite.Require().NoError(err, "There was unexpected error during the XML parsing")

	lines := instruction.Content()

	suite.Len(lines, 4)
	suite.Equal("foo — this line starts with it", lines[0])
	suite.Equal("bar — this line starts with it", lines[3])
}

func (suite *EmbeddingInstructionTestSuite) TestExplicitLineEnd() {
	instructionParams := buildInstructionParams{
		startGlob: "foo$",
		endGlob:   "bar$",
	}
	xmlString := buildInstruction("plain-text-to-embed.txt", instructionParams)
	config := buildConfigWithPreparedFragments()

	instruction, err := embedding_instruction.FromXML(xmlString, config)
	suite.Require().NoError(err, "There was unexpected error during the XML parsing")

	lines := instruction.Content()

	suite.Len(lines, 6)
	suite.Equal("This line ends with foo", lines[0])
	suite.Equal("This line ends with bar", lines[5])
}

func TestEmbeddingInstructionTestSuite(t *testing.T) {
	suite.Run(t, new(EmbeddingInstructionTestSuite))
}
