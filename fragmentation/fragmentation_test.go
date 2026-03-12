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

package fragmentation_test

import (
	"embed-code/embed-code-go/configuration"
	"embed-code/embed-code-go/files"
	"embed-code/embed-code-go/fragmentation"
	_type "embed-code/embed-code-go/type"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const (
	correctFragmentsFileName     = "Hello.java"
	unclosedFragmentFileName     = "Unclosed.java"
	unopenedFragmentFileName     = "Unopen.java"
	complexFragmentsFileName     = "Complex.java"
	twoFragmentsFileName         = "TwoFragments.java"
	overlappingFragmentsFileName = "OverlappingFragments.java"
	emptyFileName                = "Empty.java"
	indent                       = "    "
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
		config.CodeRoots = _type.NamedPathList{_type.NamedPath{Path: "../test/resources/code/java"}}
	})

	AfterEach(func() {
		cleanupDir(config.FragmentsDir)
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

	It("should do multi-source fragmentation successfully", func() {
		config := configuration.NewConfiguration()
		config.DocumentationRoot = "../test/resources/docs"
		javaCodePathName := "java-code"
		kotlinCodePathName := "kotlin-code"
		config.CodeRoots = _type.NamedPathList{
			_type.NamedPath{
				Name: javaCodePathName,
				Path: "../test/resources/code/java/org/example/multitest",
			},
			_type.NamedPath{
				Name: kotlinCodePathName,
				Path: "../test/resources/code/kotlin/org/example/multitest",
			},
		}
		result := fragmentation.WriteFragmentFiles(config)
		Expect(result.TotalSourceFiles).Should(Equal(2))
		javaFragments, _ := os.ReadDir(
			filepath.Join(config.FragmentsDir, fragmentation.NamedPathPrefix+javaCodePathName),
		)
		kotlinFragments, _ := os.ReadDir(
			filepath.Join(config.FragmentsDir, fragmentation.NamedPathPrefix+kotlinCodePathName),
		)
		Expect(javaFragments).Should(HaveLen(2))
		Expect(kotlinFragments).Should(HaveLen(2))
	})

	It("should do fragmentation of a fragment without end", func() {
		frag := buildTestFragmentation(unclosedFragmentFileName, config)
		Expect(frag.WriteFragments()).Error().ShouldNot(HaveOccurred())

		fragmentFiles := readFragmentsDir(config)
		Expect(fragmentFiles).Should(HaveLen(2))

		fragmentFileName := findFragmentFile(fragmentFiles, unclosedFragmentFileName)
		fragmentsDir := fragmentsDirPath(config.FragmentsDir)
		content, err := os.ReadFile(filepath.Join(fragmentsDir, fragmentFileName))
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

		Expect(fragmentation.WriteFragmentFiles(config).TotalFragments).Should(Equal(0))
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

			openings, _ := fragmentation.FindDocFragments(docFragment)
			Expect(openings).Should(HaveLen(2))
			Expect(openings[0]).Should(Equal(mainFragment))
			Expect(openings[1]).Should(Equal(subMainFragment))
		})

		It("should correctly find fragment endings", func() {
			endDocFragment := fmt.Sprintf(
				"// #enddocfragment \"%s\",\"%s\"", mainFragment, subMainFragment)

			endings, _ := fragmentation.FindEndDocFragments(endDocFragment)
			Expect(endings).Should(HaveLen(2))
			Expect(endings[0]).Should(Equal(mainFragment))
			Expect(endings[1]).Should(Equal(subMainFragment))
		})

		It("should not find fragment endings as there are openings", func() {
			docFragment := fmt.Sprintf(
				"// #docfragment \"%s\",\"%s\"", mainFragment, subMainFragment)

			openings, _ := fragmentation.FindEndDocFragments(docFragment)
			Expect(openings).Should(BeEmpty())
		})

		It("should not find fragment openings as there are endings", func() {
			endDocFragment := fmt.Sprintf(
				"// #enddocfragment \"%s\",\"%s\"", mainFragment, subMainFragment)

			endings, _ := fragmentation.FindDocFragments(endDocFragment)
			Expect(endings).Should(BeEmpty())
		})
	})

	It("should correctly parse file into many partitions", func() {
		frag := buildTestFragmentation(complexFragmentsFileName, config)
		_, err := frag.WriteFragments()
		Expect(err).ToNot(HaveOccurred())

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
			indent + config.Separator,
			indent + "public static void main(String[] args) {",
			indent + indent + config.Separator,
			indent + indent + "System.out.println(helperMethod());",
			"",
			indent + "}",
			config.Separator,
			"}",
		}
		Expect(content).Should(Equal(expected))
	})

	It("should correctly parse file with several different fragments", func() {
		frag := buildTestFragmentation(twoFragmentsFileName, config)
		_, err := frag.WriteFragments()
		Expect(err).ToNot(HaveOccurred())

		fragmentFiles := readFragmentsDir(config)
		Expect(fragmentFiles).Should(HaveLen(3))

		fragmentDir := fragmentsDirPath(config.FragmentsDir)
		actual := readFragmentsContent(fragmentDir, fragmentFiles, twoFragmentsFileName)

		expected := [][]string{
			{
				"public class TwoFragments {",
				indent + config.Separator,
				indent + "public static void main(String[] args) {",
				indent + indent + config.Separator,
				indent + indent + "System.out.println(helperMethod());",
				"",
				indent + "}",
				config.Separator,
				"}",
			},
			{
				"public static void hello(String[] args) {",
				indent + config.Separator,
				indent + "var coolText = \"Cool Text\";",
				indent + "System.out.println(coolText);",
				"}",
			},
		}

		Expect(actual).Should(ConsistOf(expected))
	})

	It("should correctly parse file with several overlapping fragments", func() {
		frag := buildTestFragmentation(overlappingFragmentsFileName, config)
		_, err := frag.WriteFragments()
		Expect(err).ToNot(HaveOccurred())

		fragmentFiles := readFragmentsDir(config)
		Expect(fragmentFiles).Should(HaveLen(3))

		fragmentDir := fragmentsDirPath(config.FragmentsDir)
		actual := readFragmentsContent(fragmentDir, fragmentFiles, overlappingFragmentsFileName)

		expected := [][]string{
			{
				"public class OverlappingFragments {",
				indent + config.Separator,
				indent + "public static void main(String[] args) {",
				indent + indent + config.Separator,
				indent + indent + "System.out.println(helperMethod());",
				"",
				indent + "}",
				config.Separator,
				"}",
			},
			{
				"public class OverlappingFragments {",
				indent + config.Separator,
				indent + "public static void hello(String[] args) {",
				indent + indent + config.Separator,
				indent + indent + "var coolText = \"Cool Text\";",
				indent + indent + "System.out.println(coolText);",
				indent + "}",
				config.Separator,
				"}",
			},
		}

		Expect(actual).Should(ConsistOf(expected))
	})
})

func buildTestFragmentation(testFileName string,
	config configuration.Configuration) fragmentation.Fragmentation {
	codeRoot := config.CodeRoots[0]
	testFilePath := fmt.Sprintf("%s/org/example/%s", codeRoot.Path, testFileName)

	return fragmentation.NewFragmentation(testFilePath, codeRoot, config)
}

func readFragmentsDir(config configuration.Configuration) []os.DirEntry {
	fragmentFiles, err := os.ReadDir(fragmentsDirPath(config.FragmentsDir))
	if err != nil {
		Fail(err.Error())
	}

	return fragmentFiles
}

// readFragmentsContent reads the contents of fragment files from a given directory.
//
// fragmentDir — path to the directory containing fragment files.
// fragmentFiles — list of directory entries representing the fragment files.
// skipFile — file name to skip (the default fragment file).
//
// Returns a slice of string slices, where each inner slice contains the lines
// of a fragment file.
//
// This function fails the test immediately if any file cannot be read.
func readFragmentsContent(
	fragmentDir string, fragmentFiles []os.DirEntry, skipFile string,
) [][]string {
	var result [][]string
	for _, file := range fragmentFiles {
		if file.Name() == skipFile {
			continue
		}

		content, err := files.ReadFile(fmt.Sprintf("%s/%s", fragmentDir, file.Name()))
		Expect(err).ShouldNot(HaveOccurred())

		result = append(result, content)
	}

	return result
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

func cleanupDir(dirPath string) {
	err := os.RemoveAll(dirPath)
	if err != nil {
		Fail(err.Error())
	}
}
