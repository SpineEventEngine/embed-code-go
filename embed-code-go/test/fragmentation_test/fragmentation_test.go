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

package fragmentation_test

import (
	"embed-code/embed-code-go/configuration"
	"embed-code/embed-code-go/fragmentation"
	"embed-code/embed-code-go/test/utils"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/stretchr/testify/suite"
)

func buildTestConfig() configuration.Configuration {
	var config = configuration.NewConfiguration()
	config.DocumentationRoot = "./test/resources/docs"
	config.CodeRoot = "./test/resources/code"
	return config
}

func buildTestFragmentation(
	testFileName string,
	config configuration.Configuration,
) fragmentation.Fragmentation {
	testFilePath := fmt.Sprintf("%s/org/example/%s", config.CodeRoot, testFileName)
	fragmentation := fragmentation.NewFragmentation(testFilePath, config)
	return fragmentation
}

type FragmentationTestSuite struct {
	suite.Suite
}

func (suite *FragmentationTestSuite) SetupSuite() {
	rootDir, err := filepath.Abs("../../")
	if err != nil {
		panic(err)
	}
	os.Chdir(rootDir)
}

func (suite *FragmentationTestSuite) TearDownTest() {
	var config = buildTestConfig()
	utils.CleanupDir(config.FragmentsDir)
}

func (suite *FragmentationTestSuite) TestFragmentizeFile() {
	var config = buildTestConfig()
	fileName := "Hello.java"
	frag := buildTestFragmentation(fileName, config)
	frag.WriteFragments()

	fragmentChildren, _ := os.ReadDir(config.FragmentsDir)
	suite.Len(fragmentChildren, 1)
	suite.Equal("org", fragmentChildren[0].Name())

	fragmentFiles, _ := os.ReadDir(fmt.Sprintf("%s/org/example", config.FragmentsDir))
	suite.Len(fragmentFiles, 4)

	defaultFragmentExists := false
	for _, file := range fragmentFiles {
		if file.Name() == fileName {
			defaultFragmentExists = true
		} else {
			suite.Regexp(`Hello-\w+\.java`, file.Name(), "File name does not match pattern")
		}
	}

	suite.True(defaultFragmentExists, "Default fragment '%s' does not exist", fileName)
}

func (suite *FragmentationTestSuite) TestFailNotOpenFragment() {
	var config = buildTestConfig()
	fileName := "Unopen.java"
	frag := buildTestFragmentation(fileName, config)
	err := frag.WriteFragments()
	suite.Error(err, "The file without opening tag should not be processed.")
}

func (suite *FragmentationTestSuite) TestFragmentWithoutEnd() {
	config := buildTestConfig()
	fileName := "Unclosed.java"
	frag := buildTestFragmentation(fileName, config)
	err := frag.WriteFragments()
	suite.Require().NoError(err, "Writing fragments went wrong")

	fragmentDir := fmt.Sprintf("%s/org/example", config.FragmentsDir)
	fragmentFiles, _ := os.ReadDir(fragmentDir)
	suite.Len(fragmentFiles, 2)

	var fragmentFileName string
	for _, file := range fragmentFiles {
		if file.Name() != fileName {
			fragmentFileName = file.Name()
			break
		}
	}

	fragmentContent, _ := os.ReadFile(fmt.Sprintf("%s/%s", fragmentDir, fragmentFileName))
	fragmentContentStr := string(fragmentContent)

	re, _ := regexp.Compile(`[.\n\s]+}\n}\n`)

	matched := re.FindStringSubmatch(fragmentContentStr)

	suite.Greater(len(matched), 0, "Fragment content does not match pattern", fragmentContentStr)
}

func (suite *FragmentationTestSuite) TestFragmentizeEmptyFile() {
	config := buildTestConfig()
	fileName := "Empty.java"
	frag := buildTestFragmentation(fileName, config)
	frag.WriteFragments()

	fragmentDir := fmt.Sprintf("%s/org/example", config.FragmentsDir)
	fragmentFiles, _ := os.ReadDir(fragmentDir)
	suite.Len(fragmentFiles, 1)

	fragmentContent, _ := os.ReadFile(fmt.Sprintf("%s/%s", fragmentDir, fragmentFiles[0].Name()))
	suite.Equal("", string(fragmentContent), "Expected empty string, got '%s'", string(fragmentContent))
}

func (suite *FragmentationTestSuite) TestIgnoreBinary() {
	configuration := buildTestConfig()
	configuration.CodeIncludes = []string{"**.jar"}

	fragmentation.WriteFragmentFiles(configuration)
	suite.NoDirExists(configuration.FragmentsDir, "Expected file does not exist")
}

func (suite *FragmentationTestSuite) TestManyPartitions() {
	config := buildTestConfig()

	fileName := "Complex.java"
	frag := buildTestFragmentation(fileName, config)
	frag.WriteFragments()

	fragmentDir := fmt.Sprintf("%s/org/example", config.FragmentsDir)
	fragmentFiles, _ := os.ReadDir(fragmentDir)
	suite.Len(fragmentFiles, 2)

	var fragmentFileName string
	for _, file := range fragmentFiles {
		if file.Name() != fileName {
			fragmentFileName = file.Name()
			break
		}
	}

	fragmentLines := fragmentation.ReadLines(fmt.Sprintf("%s/%s", fragmentDir, fragmentFileName))

	suite.Equal("public class Main {", fragmentLines[0])
	suite.Equal(config.Separator, fragmentLines[1])
	suite.Regexp(`\s{4}public.*`, fragmentLines[2])
	suite.Equal(config.Separator, fragmentLines[3])
	suite.Regexp(`\s{8}System.*`, fragmentLines[4])
	suite.Equal("", fragmentLines[5])
	suite.Equal("    }", fragmentLines[6])
	suite.Equal(config.Separator, fragmentLines[7])
	suite.Equal("}", fragmentLines[8])
}

func (suite *FragmentationTestSuite) TestFindFragmentOpenings() {
	testString := "// #docfragment \"main\",\"sub-main\"\n"
	foundedOpenings := fragmentation.FindFragmentOpenings(testString)

	suite.Len(foundedOpenings, 2)
	suite.Equal("main", foundedOpenings[0])
	suite.Equal("sub-main", foundedOpenings[1])
}

func (suite *FragmentationTestSuite) TestFindFragmentEndings() {
	testString := "// #enddocfragment \"main\",\"sub-main\"\n"
	foundedEndings := fragmentation.FindFragmentEndings(testString)

	suite.Len(foundedEndings, 2)
	suite.Equal("main", foundedEndings[0])
	suite.Equal("sub-main", foundedEndings[1])
}

func TestFragmentationTestSuite(t *testing.T) {
	suite.Run(t, new(FragmentationTestSuite))
}
