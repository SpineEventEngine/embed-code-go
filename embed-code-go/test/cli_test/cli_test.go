package cli_test

import (
	"embed-code/embed-code-go/cli"
	"embed-code/embed-code-go/configuration"
	"fmt"
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

func TestRequiredArgsFilled(t *testing.T) {
	args := cli.Args{
		DocsRoot: "docs",
		CodeRoot: "code",
		Mode:     "embed",
	}
	validation_message := cli.Validate(args)
	assert.Equal(t, "", validation_message)
}

func TestModeMissed(t *testing.T) {
	args := cli.Args{
		DocsRoot: "docs",
		CodeRoot: "code",
	}
	validation_message := cli.Validate(args)
	assert.Equal(t, "Mode must be set.", validation_message)
}

func TestDocsRootMissed(t *testing.T) {
	args := cli.Args{
		CodeRoot: "code",
		Mode:     "embed",
	}
	validation_message := cli.Validate(args)
	assert.Equal(t, "If one of code_root and docs_root is set, the another one must be set as well.",
		validation_message)
}

func TestConfigAndRootDirsSet(t *testing.T) {
	args := cli.Args{
		CodeRoot:       "code",
		DocsRoot:       "docs",
		Mode:           "embed",
		ConfigFilePath: "config.yaml",
	}
	validation_message := cli.Validate(args)
	assert.Equal(t, "Config path cannot be set when code_root, docs_root or optional params are set.",
		validation_message)
}

func TestCorrectConfigFile(t *testing.T) {
	preparator := newTestsPreparator()
	preparator.setup()
	defer preparator.cleanup()

	args := cli.Args{
		Mode:           "embed",
		ConfigFilePath: "./test/resources/config_files/correct_config.yml",
	}
	validation_message := cli.Validate(args)
	assert.Equal(t, "", validation_message)

	config_file_validation_message := cli.ValidateConfigFile(args.ConfigFilePath)
	assert.Equal(t, "", config_file_validation_message)
}

func TestConfigFileNotExist(t *testing.T) {
	args := cli.Args{
		Mode:           "embed",
		ConfigFilePath: "/some/path/to/config.yaml",
	}
	validation_message := cli.Validate(args)
	assert.Equal(t, "", validation_message)

	config_file_validation_message := cli.ValidateConfigFile(args.ConfigFilePath)
	assert.Equal(t, fmt.Sprintf("The file %s is not exists.", args.ConfigFilePath), config_file_validation_message)
}

func TestConfigFileWithoutDocsRoot(t *testing.T) {
	preparator := newTestsPreparator()
	preparator.setup()
	defer preparator.cleanup()

	args := cli.Args{
		Mode:           "embed",
		ConfigFilePath: "./test/resources/config_files/config_without_docs_root.yml",
	}
	validation_message := cli.Validate(args)
	assert.Equal(t, "", validation_message)

	config_file_validation_message := cli.ValidateConfigFile(args.ConfigFilePath)
	assert.Equal(t, "Config must include both code_root and docs_root fields.", config_file_validation_message)
}
