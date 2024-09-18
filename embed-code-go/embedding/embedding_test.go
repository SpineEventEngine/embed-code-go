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
	"fmt"
	"os"
	"testing"

	"embed-code/embed-code-go/configuration"
	"embed-code/embed-code-go/embedding"
	"embed-code/embed-code-go/embedding/parsing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestEmbedding(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Data Suite")
}

var _ = Describe("Embedding", func() {
	var config configuration.Configuration

	BeforeEach(func() {
		currentDir, err := os.Getwd()
		if err != nil {
			Fail("unexpected error during the test setup: " + err.Error())
		}
		err = os.Chdir(currentDir)
		if err != nil {
			Fail("unexpected error during the test setup: " + err.Error())
		}
		config = buildConfigWithPreparedFragments()
	})

	It("should be up to date", func() {
		docPath := fmt.Sprintf("%s/whole-file-fragment.md", config.DocumentationRoot)
		processor := embedding.NewProcessor(docPath, config)

		Expect(processor.Embed()).ShouldNot(HaveOccurred())
		Expect(processor.IsUpToDate()).Should(BeTrue())
	})

	It("should be up to date as there is nothing to update", func() {
		docPath := fmt.Sprintf("%s/no-embedding-doc.md", config.DocumentationRoot)
		processor := embedding.NewProcessor(docPath, config)

		Expect(processor.Embed()).ShouldNot(HaveOccurred())
		Expect(processor.IsUpToDate()).Should(BeTrue())
	})

	It("should have error as it has invalid transition map", func() {
		docPath := fmt.Sprintf("%s/split-lines.md", config.DocumentationRoot)

		falseTransitions := parsing.TransitionMap{
			parsing.Start: {parsing.RegularLine, parsing.Finish,
				parsing.EmbedInstruction},
			parsing.RegularLine: {parsing.Finish, parsing.EmbedInstruction,
				parsing.RegularLine},
			parsing.EmbedInstruction: {parsing.CodeFenceStart, parsing.BlankLine},
			parsing.BlankLine:        {parsing.CodeFenceStart, parsing.BlankLine},
			parsing.CodeFenceStart:   {parsing.CodeFenceEnd, parsing.CodeSampleLine},
			parsing.CodeSampleLine:   {parsing.CodeFenceEnd, parsing.CodeSampleLine},
			parsing.CodeFenceEnd: {parsing.Finish, parsing.EmbedInstruction,
				parsing.RegularLine},
		}

		falseProcessor := embedding.NewProcessorWithTransitions(docPath, config, falseTransitions)
		Expect(falseProcessor.Embed()).Error().Should(HaveOccurred())
	})

	It("should successfully embed with multi lined tag", func() {
		docPath := fmt.Sprintf("%s/multi-lined-tag.md", config.DocumentationRoot)
		processor := embedding.NewProcessor(docPath, config)
		Expect(processor.Embed()).Error().ShouldNot(HaveOccurred())

		Expect(processor.IsUpToDate()).Should(BeTrue())
	})
})

func buildConfigWithPreparedFragments() configuration.Configuration {
	var config = configuration.NewConfiguration()
	config.DocumentationRoot = "../test/resources/docs"
	config.CodeRoot = "../test/resources/code"
	config.FragmentsDir = "../test/resources/prepared-fragments"

	return config
}
