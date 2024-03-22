package fragmentation_test

import (
	"embed-code/embed-code-go/configuration"
	"embed-code/embed-code-go/fragmentation"
	"fmt"
	"os"
	"regexp"
	"testing"
)

func buildTestConfig() configuration.Configuration {
	var config = configuration.NewConfiguration()
	config.FragmentsDir = "./test/.fragments"
	config.DocumentationRoot = "./test/resources/docs"
	config.CodeRoot = "./test/resources/code"
	return config
}

func TestFragmentizeFile(t *testing.T) {
	// TODO: remove os.Chdir, it's just for vscode debugging
	os.Chdir(os.Getenv("WORKSPACE_DIR"))
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
	// TODO: remove os.Chdir, it's just for vscode debugging
	os.Chdir(os.Getenv("WORKSPACE_DIR"))

	var configuration = buildTestConfig()
	path := fmt.Sprintf("%s/org/example/Unopen.java", configuration.CodeRoot)
	fragmentation := fragmentation.NewFragmentation(path, configuration)
	err := fragmentation.WriteFragments()
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestFragmentWithoutEnd(t *testing.T) {
	// TODO: remove os.Chdir, it's just for vscode debugging
	os.Chdir(os.Getenv("WORKSPACE_DIR"))

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
	// TODO: remove os.Chdir, it's just for vscode debugging
	os.Chdir(os.Getenv("WORKSPACE_DIR"))

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
	// TODO: remove os.Chdir, it's just for vscode debugging
	os.Chdir(os.Getenv("WORKSPACE_DIR"))

	configuration := buildTestConfig()
	configuration.CodeIncludes = []string{"**.jar"}

	fragmentation.WriteFragmentFiles(configuration)
	if _, err := os.Stat(configuration.FragmentsDir); !os.IsNotExist(err) {
		t.Errorf("Expected file does not exist, got %v", err)
	}
}
