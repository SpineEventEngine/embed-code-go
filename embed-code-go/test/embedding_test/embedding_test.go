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
	"fmt"
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
	copyDirRecursive("./test/resources/docs", config.DocumentationRoot)
}

func (testPreparator EmbeddingInstructionTestsPreparator) cleanup() {
	config := buildConfigWithPreparedFragments()
	cleanupDir(config.DocumentationRoot)
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

	isUpToDate := processor.CheckUpToDate()
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

	isUpToDate := processor.CheckUpToDate()
	assert.True(t, isUpToDate)
}

func TestNothingToUpdate(t *testing.T) {
	preparator := newEmbeddingInstructionTestsPreparator()
	preparator.setup()
	defer preparator.cleanup()

	config := buildConfigWithPreparedFragments()
	docPath := fmt.Sprintf("%s/no-embedding-doc.md", config.DocumentationRoot)
	processor := embedding.NewEmbeddingProcessor(docPath, config)
	assert.True(t, processor.CheckUpToDate())
}
