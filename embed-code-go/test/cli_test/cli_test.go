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

// Removes directory and all its subdirectories if exists, does nothing if not exists.
//
// dir_path — a path (full or relative) of the directory to be removed.
func cleanupDir(dir_path string) {
	if _, err := os.Stat(dir_path); err == nil {
		err = os.RemoveAll(dir_path)
		if err != nil {
			panic(err)
		}
	}
}

// Copies directory from source path to target path with all subdirs and children.
//
// source_dir_path — a path (full or relative) of the directory to be copied.
//
// target_dir_path — a path (full or relative) of the directory to be copied to.
func copyDirRecursive(source_dir_path string, target_dir_path string) {
	info, err := os.Stat(source_dir_path)
	if err != nil {
		panic(err)
	}

	err = os.MkdirAll(target_dir_path, info.Mode())
	if err != nil {
		panic(err)
	}

	entries, err := os.ReadDir(source_dir_path)
	if err != nil {
		panic(err)
	}

	for _, entry := range entries {
		sourcePath := filepath.Join(source_dir_path, entry.Name())
		targetPath := filepath.Join(target_dir_path, entry.Name())

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

// Copies file from source_file_path to target_file_path.
func copyFile(source_file_path string, target_file_path string) (err error) {
	sourceFile, err := os.Open(source_file_path)
	if err != nil {
		return
	}
	defer sourceFile.Close()

	targetFile, err := os.Create(target_file_path)
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

	err = os.Chmod(target_file_path, os.FileMode(0666))
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
	cleanupDir(config.FragmentsDir)
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
