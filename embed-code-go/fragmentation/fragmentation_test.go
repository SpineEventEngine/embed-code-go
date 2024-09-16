/*
 * Copyright 2024, TeamDev. All rights reserved.
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

package fragmentation_test

import (
	"embed-code/embed-code-go/configuration"
	"embed-code/embed-code-go/files"
	"embed-code/embed-code-go/fragmentation"
	"embed-code/embed-code-go/test/filesystem"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const (
	correctFragmentsFileName = "Hello.java"
	unclosedFragmentFileName = "Unclosed.java"
	unopenedFragmentFileName = "Unopen.java"
	complexFragmentsFileName = "Complex.java"
	emptyFileName            = "Empty.java"
)

func TestFragmentation(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Data Suite")
}

var _ = Describe("Fragmentation", func() {
	var config configuration.Configuration

	BeforeEach(func() {
		config = configuration.NewConfiguration()
		config.DocumentationRoot = "../test/resources/docs"
		config.CodeRoot = "../test/resources/code"
	})

	AfterEach(func() {
		filesystem.CleanupDir(config.FragmentsDir)
	})

	It("should do file fragmentation successfully", func() {
		frag := buildTestFragmentation(correctFragmentsFileName, config)
		Expect(frag.WriteFragments()).Error().ShouldNot(HaveOccurred())

		fragmentChild, _ := os.ReadDir(config.FragmentsDir)
		Expect(fragmentChild).Should(HaveLen(1))
		Expect(fragmentChild[0].Name()).Should(Equal("org"))

		fragmentFiles := readFragmentsDir(config)
		Expect(fragmentFiles).Should(HaveLen(4))

		var isDefaultFragmentExist bool
		for _, file := range fragmentFiles {
			if file.Name() == correctFragmentsFileName {
				isDefaultFragmentExist = true
			} else {
				Expect(file.Name()).Should(MatchRegexp(`Hello-\w+\.java`))
			}
		}

		Expect(isDefaultFragmentExist).Should(BeTrue())
	})

	It("should do fragmentation of a fragment without end", func() {
		frag := buildTestFragmentation(unclosedFragmentFileName, config)
		Expect(frag.WriteFragments()).Error().ShouldNot(HaveOccurred())

		fragmentFiles := readFragmentsDir(config)
		Expect(fragmentFiles).Should(HaveLen(2))

		fragmentFileName := findFragmentFile(fragmentFiles, unclosedFragmentFileName)
		fragmentsDir := fragmentsDirPath(config.FragmentsDir)
		content, err := os.ReadFile(fmt.Sprintf("%s/%s", fragmentsDir, fragmentFileName))
		if err != nil {
			Fail(err.Error())
		}

		re := regexp.MustCompile(`[.\n\s]+}\n}\n`)
		matchedStrings := re.FindStringSubmatch(string(content))

		Expect(matchedStrings).Should(Not(BeEmpty()))
	})

	It("should not do fragmentation of an empty file", func() {
		frag := buildTestFragmentation(emptyFileName, config)
		Expect(frag.WriteFragments()).Error().ShouldNot(HaveOccurred())

		fragmentFiles := readFragmentsDir(config)
		Expect(fragmentFiles).Should(HaveLen(1))
		fragmentsFilePath := fragmentsDirPath(config.FragmentsDir) + "/" + fragmentFiles[0].Name()

		content, err := os.ReadFile(fragmentsFilePath)
		if err != nil {
			Fail(err.Error())
		}

		Expect(content).Should(BeEmpty())
	})

	It("should not do fragmentation of a binary file", func() {
		config.CodeIncludes = []string{"**.jar"}

		Expect(fragmentation.WriteFragmentFiles(config)).Error().ShouldNot(HaveOccurred())
		Expect(files.IsDirExist(config.FragmentsDir)).Should(BeFalse())
	})

	It("should not do fragmentation of an unopened fragment", func() {
		frag := buildTestFragmentation(unopenedFragmentFileName, config)

		Expect(frag.WriteFragments()).Error().Should(HaveOccurred())
	})

	Context("fragments parsing", func() {
		mainFragment := "main"
		subMainFragment := "sub-main"

		It("should correctly find fragment openings", func() {
			docFragment := fmt.Sprintf(
				"// #docfragment \"%s\",\"%s\"", mainFragment, subMainFragment)

			openings := fragmentation.FindDocFragments(docFragment)
			Expect(openings).Should(HaveLen(2))
			Expect(openings[0]).Should(Equal(mainFragment))
			Expect(openings[1]).Should(Equal(subMainFragment))
		})

		It("should correctly find fragment endings", func() {
			endDocFragment := fmt.Sprintf(
				"// #enddocfragment \"%s\",\"%s\"", mainFragment, subMainFragment)

			endings := fragmentation.FindEndDocFragments(endDocFragment)
			Expect(endings).Should(HaveLen(2))
			Expect(endings[0]).Should(Equal(mainFragment))
			Expect(endings[1]).Should(Equal(subMainFragment))
		})

		It("should not find fragment endings as there are openings", func() {
			docFragment := fmt.Sprintf(
				"// #docfragment \"%s\",\"%s\"", mainFragment, subMainFragment)

			openings := fragmentation.FindEndDocFragments(docFragment)
			Expect(openings).Should(BeEmpty())
		})

		It("should not find fragment openings as there are endings", func() {
			endDocFragment := fmt.Sprintf(
				"// #enddocfragment \"%s\",\"%s\"", mainFragment, subMainFragment)

			endings := fragmentation.FindDocFragments(endDocFragment)
			Expect(endings).Should(BeEmpty())
		})
	})

	It("should correctly parse file into many partitions", func() {
		frag := buildTestFragmentation(complexFragmentsFileName, config)
		err := frag.WriteFragments()
		if err != nil {
			Fail(err.Error())
		}

		fragmentFiles := readFragmentsDir(config)
		Expect(fragmentFiles).Should(HaveLen(2))

		fragmentFileName := findFragmentFile(fragmentFiles, complexFragmentsFileName)
		fragmentDir := fragmentsDirPath(config.FragmentsDir)

		content, err := files.ReadFile(fmt.Sprintf("%s/%s", fragmentDir, fragmentFileName))
		if err != nil {
			Fail(err.Error())
		}

		expected := []string{
			"public class Main {",
			config.Separator,
			"public static void main(String[] args) {",
			config.Separator,
			"System.out.println(helperMethod());",
			"",
			"}",
			config.Separator,
			"}",
		}

		for index, line := range content {
			Expect(strings.TrimLeft(line, " ")).Should(Equal(expected[index]))
		}
	})
})

func buildTestFragmentation(testFileName string,
	config configuration.Configuration) fragmentation.Fragmentation {
	testFilePath := fmt.Sprintf("%s/org/example/%s", config.CodeRoot, testFileName)

	return fragmentation.NewFragmentation(testFilePath, config)
}

func readFragmentsDir(config configuration.Configuration) []os.DirEntry {
	fragmentFiles, err := os.ReadDir(fragmentsDirPath(config.FragmentsDir))
	if err != nil {
		Fail(err.Error())
	}

	return fragmentFiles
}

func fragmentsDirPath(path string) string {
	return fmt.Sprintf("%s/org/example", path)
}

func findFragmentFile(files []os.DirEntry, fileName string) string {
	for _, file := range files {
		if file.Name() != fileName {
			return file.Name()
		}
	}

	return ""
}
