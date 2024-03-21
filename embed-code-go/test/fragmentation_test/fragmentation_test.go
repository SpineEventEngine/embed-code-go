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
