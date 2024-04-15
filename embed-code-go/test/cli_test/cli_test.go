package cli_test

import (
	"embed-code/embed-code-go/cli"
	"embed-code/embed-code-go/configuration"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Removes directory and all it's subdirectories if exists, does nothing if not exists.
func cleanupDir(dir string) {
	if _, err := os.Stat(dir); err == nil {
		err = os.RemoveAll(dir)
		if err != nil {
			panic(err)
		}
	}
}

// Copies directory from source path to target path with all subdirs and children.
func copyDirRecursive(source string, target string) {
	info, err := os.Stat(source)
	if err != nil {
		panic(err)
	}

	err = os.MkdirAll(target, info.Mode())
	if err != nil {
		panic(err)
	}

	entries, err := os.ReadDir(source)
	if err != nil {
		panic(err)
	}

	for _, entry := range entries {
		sourcePath := filepath.Join(source, entry.Name())
		targetPath := filepath.Join(target, entry.Name())

		if entry.IsDir() {
			copyDirRecursive(sourcePath, targetPath)
		} else {
			err = copyFile(sourcePath, targetPath)
			if err != nil {
				panic(err)
			}
		}
	}
}

func copyFile(source string, target string) (err error) {
	sourceFile, err := os.Open(source)
	if err != nil {
		return
	}
	defer sourceFile.Close()

	targetFile, err := os.Create(target)
	if err != nil {
		return
	}
	defer func() {
		cerr := targetFile.Close()
		if err == nil {
			err = cerr
		}
	}()

	if _, err = io.Copy(targetFile, sourceFile); err != nil {
		return
	}

	err = os.Chmod(target, os.FileMode(0666))
	return
}

type TestsPreparator struct {
	rootDir  string
	testsDir string
}

func newTestsPreparator() TestsPreparator {
	rootDir, err := filepath.Abs("../../")
	if err != nil {
		panic(err)
	}
	testsDir, err := filepath.Abs(".")
	if err != nil {
		panic(err)
	}
	return TestsPreparator{
		rootDir:  rootDir,
		testsDir: testsDir,
	}
}

func (testPreparator TestsPreparator) setup() {
	config := buildConfig()
	os.Chdir(testPreparator.rootDir)
	copyDirRecursive("./test/resources/docs", config.DocumentationRoot)
}

func (testPreparator TestsPreparator) cleanup() {
	config := buildConfig()
	cleanupDir(config.DocumentationRoot)
	os.Chdir(testPreparator.testsDir)
}

func buildConfig() configuration.Configuration {
	var config = configuration.NewConfiguration()
	config.DocumentationRoot = "./test/.docs"
	config.CodeRoot = "./test/resources/code"
	return config
}

func TestEmbedding(t *testing.T) {
	preparator := newTestsPreparator()
	preparator.setup()
	defer preparator.cleanup()

	config := buildConfig()

	assert.Panics(t, assert.PanicTestFunc(func() {
		cli.CheckCodeSamples(config)
	}))

	cli.EmbedCodeSamples(config)

	assert.NotPanics(t, assert.PanicTestFunc(func() {
		cli.CheckCodeSamples(config)
	}))
}
