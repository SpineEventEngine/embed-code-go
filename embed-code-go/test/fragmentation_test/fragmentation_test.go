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
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"testing"
)

type FragmentationTestsPreparator struct {
	rootDir  string
	testsDir string
}

func newFragmentationTestsPreparator() FragmentationTestsPreparator {
	rootDir, err := filepath.Abs("../../")
	if err != nil {
		panic(err)
	}
	testsDir, err := filepath.Abs(".")
	if err != nil {
		panic(err)
	}
	return FragmentationTestsPreparator{
		rootDir:  rootDir,
		testsDir: testsDir,
	}
}

func (testPreparator FragmentationTestsPreparator) setup() {
	os.Chdir(testPreparator.rootDir)
}

func (testPreparator FragmentationTestsPreparator) cleanup() {
	os.Chdir(testPreparator.rootDir)
	var config = buildTestConfig()
	err := os.RemoveAll(config.FragmentsDir)
	if err != nil {
		panic(err)
	}
	os.Chdir(testPreparator.testsDir)
}

func buildTestConfig() configuration.Configuration {
	var config = configuration.NewConfiguration()
	config.FragmentsDir = "./test/.fragments"
	config.DocumentationRoot = "./test/resources/docs"
	config.CodeRoot = "./test/resources/code"
	return config
}

func TestFragmentizeFile(t *testing.T) {
	testPreparator := newFragmentationTestsPreparator()
	testPreparator.setup()
	defer testPreparator.cleanup()

	var config = buildTestConfig()
	fileName := "Hello.java"
	path := fmt.Sprintf("%s/org/example/%s", config.CodeRoot, fileName)
	fragmentation := fragmentation.NewFragmentation(path, config)
	fragmentation.WriteFragments()

	fragmentChildren, _ := os.ReadDir(config.FragmentsDir)
	if len(fragmentChildren) != 1 {
		t.Errorf("Expected 1, got %d", len(fragmentChildren))
	}
	if fragmentChildren[0].Name() != "org" {
		t.Errorf("Expected 'org', got '%s'", fragmentChildren[0].Name())
	}

	fragmentFiles, _ := os.ReadDir(fmt.Sprintf("%s/org/example", config.FragmentsDir))
	if len(fragmentFiles) != 4 {
		t.Errorf("Expected 4, got %d", len(fragmentFiles))
	}

	defaultFragmentExists := false
	for _, file := range fragmentFiles {
		if file.Name() == fileName {
			defaultFragmentExists = true
		} else if matched, _ := regexp.MatchString(`Hello-\w+\.java`, file.Name()); !matched {
			t.Errorf("File name does not match pattern: %s", file.Name())
		}
	}

	if !defaultFragmentExists {
		t.Errorf("Default fragment '%s' does not exist", fileName)
	}
}

func TestFailNotOpenFragment(t *testing.T) {
	testPreparator := newFragmentationTestsPreparator()
	testPreparator.setup()
	defer testPreparator.cleanup()

	var configuration = buildTestConfig()
	path := fmt.Sprintf("%s/org/example/Unopen.java", configuration.CodeRoot)
	fragmentation := fragmentation.NewFragmentation(path, configuration)
	err := fragmentation.WriteFragments()
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestFragmentWithoutEnd(t *testing.T) {
	testPreparator := newFragmentationTestsPreparator()
	testPreparator.setup()
	defer testPreparator.cleanup()

	configuration := buildTestConfig()
	fileName := "Unclosed.java"
	path := fmt.Sprintf("%s/org/example/%s", configuration.CodeRoot, fileName)
	fragmentation := fragmentation.NewFragmentation(path, configuration)
	err := fragmentation.WriteFragments()
	if err != nil {
		t.Errorf("Writing fragments went wrong: %d", err)
	}

	fragmentDir := fmt.Sprintf("%s/org/example", configuration.FragmentsDir)
	fragmentFiles, _ := os.ReadDir(fragmentDir)
	if len(fragmentFiles) != 2 {
		t.Errorf("Expected 2, got %d", len(fragmentFiles))
	}

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

	if len(matched) == 0 {
		t.Errorf("Fragment content does not match pattern: %s", fragmentContentStr)
	}
}

func TestFragmentizeEmptyFile(t *testing.T) {
	testPreparator := newFragmentationTestsPreparator()
	testPreparator.setup()
	defer testPreparator.cleanup()

	configuration := buildTestConfig()
	fileName := "Empty.java"
	path := fmt.Sprintf("%s/org/example/%s", configuration.CodeRoot, fileName)
	fragmentation := fragmentation.NewFragmentation(path, configuration)
	fragmentation.WriteFragments()

	fragmentDir := fmt.Sprintf("%s/org/example", configuration.FragmentsDir)
	fragmentFiles, _ := os.ReadDir(fragmentDir)
	if len(fragmentFiles) != 1 {
		t.Errorf("Expected 1, got %d", len(fragmentFiles))
	}

	fragmentContent, _ := os.ReadFile(fmt.Sprintf("%s/%s", fragmentDir, fragmentFiles[0].Name()))
	if string(fragmentContent) != "" {
		t.Errorf("Expected empty string, got '%s'", string(fragmentContent))
	}
}

func TestIgnoreBinary(t *testing.T) {
	testPreparator := newFragmentationTestsPreparator()
	testPreparator.setup()
	defer testPreparator.cleanup()

	configuration := buildTestConfig()
	configuration.CodeIncludes = []string{"**.jar"}

	fragmentation.WriteFragmentFiles(configuration)
	if _, err := os.Stat(configuration.FragmentsDir); !os.IsNotExist(err) {
		t.Errorf("Expected file does not exist, got %v", err)
	}
}

func TestManyPartitions(t *testing.T) {
	testPreparator := newFragmentationTestsPreparator()
	testPreparator.setup()
	defer testPreparator.cleanup()

	configuration := buildTestConfig()

	fileName := "Complex.java"
	path := fmt.Sprintf("%s/org/example/%s", configuration.CodeRoot, fileName)
	frag := fragmentation.NewFragmentation(path, configuration)
	frag.WriteFragments()

	fragmentDir := fmt.Sprintf("%s/org/example", configuration.FragmentsDir)
	fragmentFiles, _ := os.ReadDir(fragmentDir)
	if len(fragmentFiles) != 2 {
		t.Errorf("Expected 2, got %d", len(fragmentFiles))
	}

	var fragmentFileName string
	for _, file := range fragmentFiles {
		if file.Name() != fileName {
			fragmentFileName = file.Name()
			break
		}
	}

	fragmentLines := fragmentation.ReadLines(fmt.Sprintf("%s/%s", fragmentDir, fragmentFileName))

	if fragmentLines[0] != "public class Main {" {
		t.Errorf("Expected 'public class Main {', got '%s'", fragmentLines[0])
	}
	if fragmentLines[1] != configuration.Separator {
		t.Errorf("Expected '%s', got '%s'", configuration.Separator, fragmentLines[1])
	}
	if matched, _ := regexp.MatchString(`\s{4}public.*`, fragmentLines[2]); !matched {
		t.Errorf("Line does not match pattern: %s", fragmentLines[2])
	}
	if fragmentLines[3] != configuration.Separator {
		t.Errorf("Expected '%s', got '%s'", configuration.Separator, fragmentLines[3])
	}
	if matched, _ := regexp.MatchString(`\s{8}System.*`, fragmentLines[4]); !matched {
		t.Errorf("Line does not match pattern: %s", fragmentLines[4])
	}
	if fragmentLines[5] != "" {
		t.Errorf("Expected empty string, got '%s'", fragmentLines[5])
	}
	if fragmentLines[6] != "    }" {
		t.Errorf("Expected '    }', got '%s'", fragmentLines[6])
	}
	if fragmentLines[7] != configuration.Separator {
		t.Errorf("Expected '%s', got '%s'", configuration.Separator, fragmentLines[7])
	}
	if fragmentLines[8] != "}" {
		t.Errorf("Expected '}', got '%s'", fragmentLines[8])
	}
}
