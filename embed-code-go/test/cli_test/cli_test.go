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

package cli_test

import (
	"embed-code/embed-code-go/cli"
	"embed-code/embed-code-go/configuration"
	"embed-code/embed-code-go/test/filesystem"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"
)

func buildConfig() configuration.Configuration {
	var config = configuration.NewConfiguration()
	config.DocumentationRoot = "./test/.docs"
	config.CodeRoot = "./test/resources/code"
	config.CodeIncludes = []string{"**/Hello.java", "**/Hello.kt"}
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
	filesystem.CopyDirRecursive("./test/resources/docs", suite.config.DocumentationRoot)
}

func (suite *CLITestSuite) TearDownTest() {
	filesystem.CleanupDir(suite.config.DocumentationRoot)
	filesystem.CleanupDir(suite.config.FragmentsDir)
}

func (suite *CLITestSuite) TestEmbedding() {
	suite.Panics(func() {
		cli.CheckCodeSamples(suite.config)
	})

	cli.EmbedCodeSamples(suite.config)

	suite.NotPanics(func() {
		cli.CheckCodeSamples(suite.config)
	})
}

func (suite *CLITestSuite) TestRequiredArgsFilled() {
	args := cli.Config{
		DocsPath: "docs",
		CodePath: "code",
		Mode:     "embed",
	}
	validation_message := cli.ValidateConfig(args)
	suite.Equal(nil, validation_message)
}

func (suite *CLITestSuite) TestModeMissed() {
	args := cli.Config{
		DocsPath: "docs",
		CodePath: "code",
	}
	validation_message := cli.ValidateConfig(args).Error()
	suite.Equal("mode must be set", validation_message)
}

func (suite *CLITestSuite) TestDocsRootMissed() {
	args := cli.Config{
		CodePath: "code",
		Mode:     "embed",
	}
	validation_message := cli.ValidateConfig(args).Error()
	suite.Equal("if one of code-path and docs-path is set, the another one must be set as well",
		validation_message)
}

func (suite *CLITestSuite) TestConfigAndRootDirsSet() {
	args := cli.Config{
		CodePath:   "code",
		DocsPath:   "docs",
		Mode:       "embed",
		ConfigPath: "config.yaml",
	}
	validation_message := cli.ValidateConfig(args).Error()
	suite.Equal("config path cannot be set when code-path, docs-path or optional params are set",
		validation_message)
}

func (suite *CLITestSuite) TestCorrectConfigFile() {
	args := cli.Config{
		Mode:       "embed",
		ConfigPath: "./test/resources/config_files/correct_config.yml",
	}
	validation_message := cli.ValidateConfig(args)
	suite.Equal(nil, validation_message)

	config_file_validation_message := cli.ValidateConfigFile(args.ConfigPath)
	suite.Equal("", config_file_validation_message)
}

func (suite *CLITestSuite) TestConfigFileNotExist() {
	args := cli.Config{
		Mode:       "embed",
		ConfigPath: "/some/path/to/config.yaml",
	}
	validation_message := cli.ValidateConfig(args)
	suite.Equal(nil, validation_message)

	config_file_validation_message := cli.ValidateConfigFile(args.ConfigPath)
	suite.Equal(fmt.Sprintf("The file %s is not exists.", args.ConfigPath), config_file_validation_message)
}

func (suite *CLITestSuite) TestConfigFileWithoutDocsRoot() {
	args := cli.Config{
		Mode:       "embed",
		ConfigPath: "./test/resources/config_files/config_without_docs_root.yml",
	}
	validation_message := cli.ValidateConfig(args)
	suite.Equal(nil, validation_message)

	config_file_validation_message := cli.ValidateConfigFile(args.ConfigPath)
	suite.Equal("Config must include both code-path and docs-path fields.", config_file_validation_message)
}

func TestCLITestSuite(t *testing.T) {
	suite.Run(t, new(CLITestSuite))
}
