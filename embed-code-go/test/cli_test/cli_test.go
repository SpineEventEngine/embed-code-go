package cli_test

import (
	"embed-code/embed-code-go/cli"
	"embed-code/embed-code-go/configuration"
	"embed-code/embed-code-go/test/utils"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
	utils.CopyDirRecursive("./test/resources/docs", config.DocumentationRoot)
}

func (testPreparator TestsPreparator) cleanup() {
	config := buildConfig()
	utils.CleanupDir(config.DocumentationRoot)
	utils.CleanupDir(config.FragmentsDir)
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
	assert.Equal(t, "The file /some/path/to/config.yaml is not exists.", config_file_validation_message)
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
