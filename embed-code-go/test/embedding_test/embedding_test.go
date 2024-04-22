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
	"embed-code/embed-code-go/test/utils"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

type EmbeddingInstructionTestsPreparator struct {
	rootDir  string
	testsDir string
}

func newEmbeddingInstructionTestsPreparator() EmbeddingInstructionTestsPreparator {
	rootDir, err := filepath.Abs("../../")
	if err != nil {
		panic(err)
	}
	testsDir, err := filepath.Abs(".")
	if err != nil {
		panic(err)
	}
	return EmbeddingInstructionTestsPreparator{
		rootDir:  rootDir,
		testsDir: testsDir,
	}
}

func (testPreparator EmbeddingInstructionTestsPreparator) setup() {
	config := buildConfigWithPreparedFragments()
	os.Chdir(testPreparator.rootDir)
	utils.CopyDirRecursive("./test/resources/docs", config.DocumentationRoot)
}

func (testPreparator EmbeddingInstructionTestsPreparator) cleanup() {
	config := buildConfigWithPreparedFragments()
	utils.CleanupDir(config.DocumentationRoot)
	os.Chdir(testPreparator.testsDir)
}

func buildConfigWithPreparedFragments() configuration.Configuration {
	var config = configuration.NewConfiguration()
	config.DocumentationRoot = "./test/.docs"
	config.CodeRoot = "./test/resources/code"
	config.FragmentsDir = "./test/resources/prepared-fragments"
	return config
}

func TestNotUpToDate(t *testing.T) {
	preparator := newEmbeddingInstructionTestsPreparator()
	preparator.setup()
	defer preparator.cleanup()

	config := buildConfigWithPreparedFragments()
	docPath := fmt.Sprintf("%s/whole-file-fragment.md", config.DocumentationRoot)
	processor := embedding.NewEmbeddingProcessor(docPath, config)

	isUpToDate := processor.IsUpToDate()
	assert.False(t, isUpToDate)
}

func TestUpToDate(t *testing.T) {
	preparator := newEmbeddingInstructionTestsPreparator()
	preparator.setup()
	defer preparator.cleanup()

	config := buildConfigWithPreparedFragments()
	docPath := fmt.Sprintf("%s/whole-file-fragment.md", config.DocumentationRoot)
	processor := embedding.NewEmbeddingProcessor(docPath, config)
	processor.Embed()

	isUpToDate := processor.IsUpToDate()
	assert.True(t, isUpToDate)
}

func TestNothingToUpdate(t *testing.T) {
	preparator := newEmbeddingInstructionTestsPreparator()
	preparator.setup()
	defer preparator.cleanup()

	config := buildConfigWithPreparedFragments()
	docPath := fmt.Sprintf("%s/no-embedding-doc.md", config.DocumentationRoot)
	processor := embedding.NewEmbeddingProcessor(docPath, config)
	assert.True(t, processor.IsUpToDate())
}

func TestFalseTransitions(t *testing.T) {
	preparator := newEmbeddingInstructionTestsPreparator()
	preparator.setup()
	defer preparator.cleanup()

	config := buildConfigWithPreparedFragments()
	docPath := fmt.Sprintf("%s/split-lines.md", config.DocumentationRoot)

	falseTransitions := map[string][]string{
		"START":                 {"REGULAR_LINE", "FINISH", "EMBEDDING_INSTRUCTION"},
		"REGULAR_LINE":          {"FINISH", "EMBEDDING_INSTRUCTION", "REGULAR_LINE"},
		"EMBEDDING_INSTRUCTION": {"CODE_FENCE_START", "BLANK_LINE"},
		"BLANK_LINE":            {"CODE_FENCE_START", "BLANK_LINE"},
		"CODE_FENCE_START":      {"CODE_FENCE_END", "CODE_SAMPLE_LINE"},
		"CODE_SAMPLE_LINE":      {"CODE_FENCE_END", "CODE_SAMPLE_LINE"},
		"CODE_FENCE_END":        {"FINISH", "EMBEDDING_INSTRUCTION", "REGULAR_LINE"},
	}

	falseProcessor := embedding.NewEmbeddingProcessorWithTransitions(docPath, config, falseTransitions)

	assert.Panics(t, assert.PanicTestFunc(func() {
		falseProcessor.Embed()
	}))
}

func TestMultiLinedTag(t *testing.T) {
	preparator := newEmbeddingInstructionTestsPreparator()
	preparator.setup()
	defer preparator.cleanup()

	config := buildConfigWithPreparedFragments()
	docPath := fmt.Sprintf("%s/multi-lined-tag.md", config.DocumentationRoot)
	processor := embedding.NewEmbeddingProcessor(docPath, config)
	processor.Embed()

	isUpToDate := processor.IsUpToDate()
	assert.True(t, isUpToDate)
}
