// Copyright 2020, TeamDev. All rights reserved.
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

package embedding_test

import (
	"embed-code/embed-code-go/configuration"
	"embed-code/embed-code-go/embedding"
	"embed-code/embed-code-go/test/filesystem"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"
)

func buildConfigWithPreparedFragments() configuration.Configuration {
	var config = configuration.NewConfiguration()
	config.DocumentationRoot = "./test/.docs"
	config.CodeRoot = "./test/resources/code"
	config.FragmentsDir = "./test/resources/prepared-fragments"
	return config
}

type EmbeddingTestSuite struct {
	suite.Suite
	config configuration.Configuration
}

func (suite *EmbeddingTestSuite) SetupSuite() {
	rootDir, err := filepath.Abs("../../")
	if err != nil {
		panic(err)
	}
	os.Chdir(rootDir)
	suite.config = buildConfigWithPreparedFragments()
}

func (suite *EmbeddingTestSuite) SetupTest() {
	filesystem.CopyDirRecursive("./test/resources/docs", suite.config.DocumentationRoot)
}

func (suite *EmbeddingTestSuite) TearDownTest() {
	filesystem.CleanupDir(suite.config.DocumentationRoot)
}

func (suite *EmbeddingTestSuite) TestNotUpToDate() {
	docPath := fmt.Sprintf("%s/whole-file-fragment.md", suite.config.DocumentationRoot)
	processor := embedding.NewEmbeddingProcessor(docPath, suite.config)

	isUpToDate := processor.IsUpToDate()
	suite.False(isUpToDate)
}

func (suite *EmbeddingTestSuite) TestUpToDate() {
	docPath := fmt.Sprintf("%s/whole-file-fragment.md", suite.config.DocumentationRoot)
	processor := embedding.NewEmbeddingProcessor(docPath, suite.config)
	processor.Embed()

	isUpToDate := processor.IsUpToDate()
	suite.True(isUpToDate)
}

func (suite *EmbeddingTestSuite) TestNothingToUpdate() {
	docPath := fmt.Sprintf("%s/no-embedding-doc.md", suite.config.DocumentationRoot)
	processor := embedding.NewEmbeddingProcessor(docPath, suite.config)
	suite.True(processor.IsUpToDate())
}

func (suite *EmbeddingTestSuite) TestFalseTransitions() {
	docPath := fmt.Sprintf("%s/split-lines.md", suite.config.DocumentationRoot)

	falseTransitions := map[string][]string{
		"START":                 {"REGULAR_LINE", "FINISH", "EMBEDDING_INSTRUCTION"},
		"REGULAR_LINE":          {"FINISH", "EMBEDDING_INSTRUCTION", "REGULAR_LINE"},
		"EMBEDDING_INSTRUCTION": {"CODE_FENCE_START", "BLANK_LINE"},
		"BLANK_LINE":            {"CODE_FENCE_START", "BLANK_LINE"},
		"CODE_FENCE_START":      {"CODE_FENCE_END", "CODE_SAMPLE_LINE"},
		"CODE_SAMPLE_LINE":      {"CODE_FENCE_END", "CODE_SAMPLE_LINE"},
		"CODE_FENCE_END":        {"FINISH", "EMBEDDING_INSTRUCTION", "REGULAR_LINE"},
	}

	falseProcessor := embedding.NewEmbeddingProcessorWithTransitions(docPath, suite.config, falseTransitions)

	suite.Require().Panics(func() {
		falseProcessor.Embed()
	})
}

func (suite *EmbeddingTestSuite) TestMultiLinedTag() {
	docPath := fmt.Sprintf("%s/multi-lined-tag.md", suite.config.DocumentationRoot)
	processor := embedding.NewEmbeddingProcessor(docPath, suite.config)
	processor.Embed()

	isUpToDate := processor.IsUpToDate()
	suite.True(isUpToDate)
}

func TestEmbeddingTestSuite(t *testing.T) {
	suite.Run(t, new(EmbeddingTestSuite))
}
