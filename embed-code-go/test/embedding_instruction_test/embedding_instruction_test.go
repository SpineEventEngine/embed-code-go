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
	"strings"
	"testing"
)

type EmbeddingInstructionTestsPreparator struct {
	rootDir  string
	testsDir string
}

func newEmbeddingInstructionTestsPreparator() EmbeddingInstructionTestsPreparator {
	rootDir, err := filepath.Abs("../../")
	if err != nil {
		panic(err)
	}
	testsDir, err := filepath.Abs(".")
	if err != nil {
		panic(err)
	}
	return EmbeddingInstructionTestsPreparator{
		rootDir:  rootDir,
		testsDir: testsDir,
	}
}

func (testPreparator EmbeddingInstructionTestsPreparator) setup() {
	os.Chdir(testPreparator.rootDir)
}

func (testPreparator EmbeddingInstructionTestsPreparator) cleanup() {
	os.Chdir(testPreparator.testsDir)
}

func buildConfigWithPreparedFragments() configuration.Configuration {
	var config = configuration.NewConfiguration()
	config.DocumentationRoot = "./test/resources/docs"
	config.CodeRoot = "./test/resources/code"
	config.FragmentsDir = "./test/resources/prepared-fragments"
	return config
}

type buildInstructionParams struct {
	fragment  string
	startGlob string
	endGlob   string
	closeTag  bool
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

func TestFalseXML(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("FromXML had to raise a panic, but it didn't.")
		}
	}()

	xmlString := "<file=\"org/example/Hello.java\" fragment=\"Hello class\"/>"
	config := buildConfigWithPreparedFragments()

	embedding_instruction.FromXML(xmlString, config)
}

func TestParseFromXML(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Error: an exception occured on FromXML.")
		}
	}()

	instructionParams := buildInstructionParams{}
	instructionParams.fragment = "Hello class"
	xmlString := buildInstruction("org/example/Hello.java", instructionParams)
	config := buildConfigWithPreparedFragments()
	embedding_instruction.FromXML(xmlString, config)
}

func TestParseWithClosingTag(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Error: an exception occured on FromXML.")
		}
	}()

	instructionParams := buildInstructionParams{}
	instructionParams.fragment = "Hello class"
	instructionParams.closeTag = true
	xmlString := buildInstruction("org/example/Hello.java", instructionParams)
	config := buildConfigWithPreparedFragments()

	embedding_instruction.FromXML(xmlString, config)
}

func TestReadFragmentDir(t *testing.T) {
	preparator := newEmbeddingInstructionTestsPreparator()
	preparator.setup()

	instructionParams := buildInstructionParams{}
	instructionParams.closeTag = true
	xmlString := buildInstruction("org/example/Hello.java", instructionParams)
	config := buildConfigWithPreparedFragments()
	instruction := embedding_instruction.FromXML(xmlString, config)
	lines := instruction.Content()

	if len(lines) != 29 {
		t.Errorf("Error: the length of the lines is %d, but have to be 29", len(lines))
	}
	if lines[22] != "public class Hello {" {
		t.Errorf(
			"Error: the line at the 22 index is %s, but have to be \"public class Hello {\"",
			lines[22])
	}

	preparator.cleanup()
}

func TestFragmentAndStart(t *testing.T) {
	preparator := newEmbeddingInstructionTestsPreparator()
	preparator.setup()

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("FromXML had to raise a panic, but it didn't.")
		}
		preparator.cleanup()
	}()

	instructionParams := buildInstructionParams{fragment: "fragment", startGlob: "public void hello()"}

	xmlString := buildInstruction("org/example/Hello.java", instructionParams)
	config := buildConfigWithPreparedFragments()
	embedding_instruction.FromXML(xmlString, config)

	preparator.cleanup()
}

func TestFragmentAndEnd(t *testing.T) {
	preparator := newEmbeddingInstructionTestsPreparator()
	preparator.setup()

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("FromXML had to raise a panic, but it didn't.")
		}
		preparator.cleanup()
	}()

	instructionParams := buildInstructionParams{fragment: "fragment", endGlob: "}"}

	xmlString := buildInstruction("org/example/Hello.java", instructionParams)
	config := buildConfigWithPreparedFragments()
	embedding_instruction.FromXML(xmlString, config)

	preparator.cleanup()
}

func TestExtractByGlob(t *testing.T) {
	preparator := newEmbeddingInstructionTestsPreparator()
	preparator.setup()

	instructionParams := buildInstructionParams{
		startGlob: "public class*",
		endGlob:   "*System.out*",
	}

	xmlString := buildInstruction("org/example/Hello.java", instructionParams)
	config := buildConfigWithPreparedFragments()
	instruction := embedding_instruction.FromXML(xmlString, config)
	lines := instruction.Content()

	if len(lines) != 4 {
		t.Errorf("Error: the length of the lines is %d, but have to be 4", len(lines))
	}
	if lines[0] != "public class Hello {" {
		t.Errorf(
			"Error: the line at the 0 index is %s, but have to be \"public class Hello {\"",
			lines[0])
	}
	if lines[3] != "        System.out.println(\"Hello world\");" {
		t.Errorf(
			"Error: the line at the 3 index is %s, but have to be \"        System.out.println(\"Hello world\");\"",
			lines[3])
	}

	preparator.cleanup()
}

func TestMinIndentation(t *testing.T) {
	preparator := newEmbeddingInstructionTestsPreparator()
	preparator.setup()

	instructionParams := buildInstructionParams{
		startGlob: "*public static void main*",
		endGlob:   "*}*",
	}
	xmlString := buildInstruction("org/example/Hello.java", instructionParams)
	config := buildConfigWithPreparedFragments()
	instruction := embedding_instruction.FromXML(xmlString, config)
	lines := instruction.Content()

	if len(lines) != 3 {
		t.Errorf("Error: the length of the lines is %d, but have to be 3", len(lines))
	}
	if lines[0] != "public static void main(String[] args) {" {
		t.Errorf(
			"Error: the line at the 1 index is %s, but have to be \"public static void main(String[] args) {\"",
			lines[0])
	}
	if !strings.HasPrefix(lines[1], "    ") {
		t.Errorf(
			"Error: the line at the 1 index is %s, but it have to have \"    \" prefix.",
			lines[1])
	}

	preparator.cleanup()
}

func TestStartWithoutEnd(t *testing.T) {
	preparator := newEmbeddingInstructionTestsPreparator()
	preparator.setup()

	instructionParams := buildInstructionParams{
		startGlob: "*class*",
	}
	xmlString := buildInstruction("org/example/Hello.java", instructionParams)
	config := buildConfigWithPreparedFragments()
	instruction := embedding_instruction.FromXML(xmlString, config)
	lines := instruction.Content()

	if len(lines) != 7 {
		t.Errorf("Error: the length of the lines is %d, but have to be 7", len(lines))
	}
	if lines[5] != "}" {
		t.Errorf(
			"Error: the line at the 5 index is %s, but have to be \"}\"",
			lines[5])
	}

	preparator.cleanup()
}

func TestEndWithoutStart(t *testing.T) {
	preparator := newEmbeddingInstructionTestsPreparator()
	preparator.setup()

	instructionParams := buildInstructionParams{
		endGlob: "package*",
	}
	xmlString := buildInstruction("org/example/Hello.java", instructionParams)
	config := buildConfigWithPreparedFragments()
	instruction := embedding_instruction.FromXML(xmlString, config)
	lines := instruction.Content()

	if len(lines) != 21 {
		t.Errorf("Error: the length of the lines is %d, but have to be 21", len(lines))
	}
	if lines[0] != "/*" {
		t.Errorf(
			"Error: the line at the 0 index is %s, but have to be \"/*\"",
			lines[0])
	}
	if lines[20] != "package org.example;" {
		t.Errorf(
			"Error: the line at the 20 index is %s, but have to be \"package org.example;\"",
			lines[20])
	}

	preparator.cleanup()
}

func TestOneLine(t *testing.T) {
	preparator := newEmbeddingInstructionTestsPreparator()
	preparator.setup()

	instructionParams := buildInstructionParams{
		startGlob: "*main*",
		endGlob:   "*main*",
	}
	xmlString := buildInstruction("org/example/Hello.java", instructionParams)
	config := buildConfigWithPreparedFragments()
	instruction := embedding_instruction.FromXML(xmlString, config)
	lines := instruction.Content()

	if len(lines) != 1 {
		t.Errorf("Error: the length of the lines is %d, but have to be 1", len(lines))
	}
	if lines[0] != "public static void main(String[] args) {" {
		t.Errorf(
			"Error: the line at the 0 index is %s, but have to be \"public static void main(String[] args) {\"",
			lines[0])
	}

	preparator.cleanup()
}

func TestNoMatchStart(t *testing.T) {
	preparator := newEmbeddingInstructionTestsPreparator()
	preparator.setup()

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("instruction.Content() had to raise a panic, but it didn't.")
		}
		preparator.cleanup()
	}()

	instructionParams := buildInstructionParams{
		startGlob: "foo bar",
		endGlob:   "*main*",
	}

	xmlString := buildInstruction("org/example/Hello.java", instructionParams)
	config := buildConfigWithPreparedFragments()
	instruction := embedding_instruction.FromXML(xmlString, config)
	instruction.Content()

	preparator.cleanup()
}

func TestNoMatchEnd(t *testing.T) {
	preparator := newEmbeddingInstructionTestsPreparator()
	preparator.setup()

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("instruction.Content() had to raise a panic, but it didn't.")
		}
		preparator.cleanup()
	}()

	instructionParams := buildInstructionParams{
		startGlob: "*main*",
		endGlob:   "foo bar",
	}

	xmlString := buildInstruction("org/example/Hello.java", instructionParams)
	config := buildConfigWithPreparedFragments()
	instruction := embedding_instruction.FromXML(xmlString, config)
	instruction.Content()

	preparator.cleanup()
}

func TestImplyAsterisk(t *testing.T) {
	preparator := newEmbeddingInstructionTestsPreparator()
	preparator.setup()

	instructionParams := buildInstructionParams{
		startGlob: "main",
		endGlob:   "world",
	}
	xmlString := buildInstruction("org/example/Hello.java", instructionParams)
	config := buildConfigWithPreparedFragments()
	instruction := embedding_instruction.FromXML(xmlString, config)
	lines := instruction.Content()

	if len(lines) != 2 {
		t.Errorf("Error: the length of the lines is %d, but have to be 1", len(lines))
	}
	if !strings.HasPrefix(lines[0], "public static void main") {
		t.Errorf(
			"Error: the line %s has to have \"public static void main\" at the beginning.",
			lines[0])
	}
	if !strings.HasPrefix(lines[1], "    System.out.println") {
		t.Errorf(
			"Error: the line %s has to have \"    System.out.println\" at the beginning.",
			lines[1])
	}

	preparator.cleanup()
}

func TestExplicitLineStart(t *testing.T) {
	preparator := newEmbeddingInstructionTestsPreparator()
	preparator.setup()

	instructionParams := buildInstructionParams{
		startGlob: "^foo",
		endGlob:   "^bar",
	}
	xmlString := buildInstruction("plain-text-to-embed.txt", instructionParams)
	config := buildConfigWithPreparedFragments()
	instruction := embedding_instruction.FromXML(xmlString, config)
	lines := instruction.Content()

	if len(lines) != 4 {
		t.Errorf("Error: the length of the lines is %d, but have to be 4", len(lines))
	}
	if lines[0] != "foo — this line starts with it" {
		t.Errorf(
			"Error: the line at the 0 index is %s, but have to be \"foo — this line starts with it\"",
			lines[0])
	}
	if lines[3] != "bar — this line starts with it" {
		t.Errorf(
			"Error: the line at the 0 index is %s, but have to be \"bar — this line starts with it\"",
			lines[3])
	}

	preparator.cleanup()
}

func TestExplicitLineEnd(t *testing.T) {
	preparator := newEmbeddingInstructionTestsPreparator()
	preparator.setup()

	instructionParams := buildInstructionParams{
		startGlob: "foo$",
		endGlob:   "bar$",
	}
	xmlString := buildInstruction("plain-text-to-embed.txt", instructionParams)
	config := buildConfigWithPreparedFragments()
	instruction := embedding_instruction.FromXML(xmlString, config)
	lines := instruction.Content()

	if len(lines) != 6 {
		t.Errorf("Error: the length of the lines is %d, but have to be 4", len(lines))
	}
	if lines[0] != "This line ends with foo" {
		t.Errorf(
			"Error: the line at the 0 index is %s, but have to be \"This line ends with foo\"",
			lines[0])
	}
	if lines[5] != "This line ends with bar" {
		t.Errorf(
			"Error: the line at the 0 index is %s, but have to be \"This line ends with bar\"",
			lines[5])
	}

	preparator.cleanup()
}
