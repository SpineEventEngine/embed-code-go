package cli_test

import (
	"embed-code/embed-code-go/cli"
	"embed-code/embed-code-go/configuration"
	"embed-code/embed-code-go/test/utils"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func buildConfig() configuration.Configuration {
	var config = configuration.NewConfiguration()
	config.DocumentationRoot = "./test/.docs"
	config.CodeRoot = "./test/resources/code"
	return config
}

type CLITestSuite struct {
	suite.Suite
	config configuration.Configuration
}

func (suite *CLITestSuite) SetupSuite() {
	rootDir, err := filepath.Abs("../../")
	if err != nil {
		panic(err)
	}
	os.Chdir(rootDir)
	suite.config = buildConfig()
}

func (suite *CLITestSuite) SetupTest() {
	utils.CopyDirRecursive("./test/resources/docs", suite.config.DocumentationRoot)
}

func (suite *CLITestSuite) TearDownTest() {
	utils.CleanupDir(suite.config.DocumentationRoot)
	utils.CleanupDir(suite.config.FragmentsDir)
}

func (suite *CLITestSuite) TestEmbedding() {
	suite.Panics(assert.PanicTestFunc(func() {
		cli.CheckCodeSamples(suite.config)
	}))

	cli.EmbedCodeSamples(suite.config)

	suite.NotPanics(assert.PanicTestFunc(func() {
		cli.CheckCodeSamples(suite.config)
	}))
}

func (suite *CLITestSuite) TestRequiredArgsFilled() {
	args := cli.Args{
		DocsRoot: "docs",
		CodeRoot: "code",
		Mode:     "embed",
	}
	validation_message := cli.Validate(args)
	suite.Equal("", validation_message)
}

func (suite *CLITestSuite) TestModeMissed() {
	args := cli.Args{
		DocsRoot: "docs",
		CodeRoot: "code",
	}
	validation_message := cli.Validate(args)
	suite.Equal("Mode must be set.", validation_message)
}

func (suite *CLITestSuite) TestDocsRootMissed() {
	args := cli.Args{
		CodeRoot: "code",
		Mode:     "embed",
	}
	validation_message := cli.Validate(args)
	suite.Equal("If one of code_root and docs_root is set, the another one must be set as well.",
		validation_message)
}

func (suite *CLITestSuite) TestConfigAndRootDirsSet() {
	args := cli.Args{
		CodeRoot:       "code",
		DocsRoot:       "docs",
		Mode:           "embed",
		ConfigFilePath: "config.yaml",
	}
	validation_message := cli.Validate(args)
	suite.Equal("Config path cannot be set when code_root, docs_root or optional params are set.",
		validation_message)
}

func (suite *CLITestSuite) TestCorrectConfigFile() {
	args := cli.Args{
		Mode:           "embed",
		ConfigFilePath: "./test/resources/config_files/correct_config.yml",
	}
	validation_message := cli.Validate(args)
	suite.Equal("", validation_message)

	config_file_validation_message := cli.ValidateConfigFile(args.ConfigFilePath)
	suite.Equal("", config_file_validation_message)
}

func (suite *CLITestSuite) TestConfigFileNotExist() {
	args := cli.Args{
		Mode:           "embed",
		ConfigFilePath: "/some/path/to/config.yaml",
	}
	validation_message := cli.Validate(args)
	suite.Equal("", validation_message)

	config_file_validation_message := cli.ValidateConfigFile(args.ConfigFilePath)
	suite.Equal("The file /some/path/to/config.yaml is not exists.", config_file_validation_message)
}

func (suite *CLITestSuite) TestConfigFileWithoutDocsRoot() {
	args := cli.Args{
		Mode:           "embed",
		ConfigFilePath: "./test/resources/config_files/config_without_docs_root.yml",
	}
	validation_message := cli.Validate(args)
	suite.Equal("", validation_message)

	config_file_validation_message := cli.ValidateConfigFile(args.ConfigFilePath)
	suite.Equal("Config must include both code_root and docs_root fields.", config_file_validation_message)
}

func TestCLITestSuite(t *testing.T) {
	suite.Run(t, new(CLITestSuite))
}
