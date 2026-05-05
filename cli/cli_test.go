/*
 * Copyright 2026, TeamDev. All rights reserved.
 *
 * Redistribution and use in source and/or binary forms, with or without
 * modification, must retain the above copyright notice and the following
 * disclaimer.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
 * "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
 * LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
 * A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
 * OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
 * SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
 * LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
 * DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
 * THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 * OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

package cli_test

import (
	"embed-code/embed-code-go/cli"
	"embed-code/embed-code-go/configuration"
	_type "embed-code/embed-code-go/type"
	"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCli(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Data Suite")
}

var _ = Describe("CLI validation", func() {

	Context("with valid config", func() {
		var config cli.Config

		BeforeEach(func() {
			config = baseCliConfig()
		})

		DescribeTable("should pass validation when all required args are set",
			func(mode string) {
				config.Mode = mode
				Expect(cli.ValidateConfig(config)).Error().ShouldNot(HaveOccurred())
			},

			Entry("with check mode", cli.ModeCheck),
			Entry("with analyze mode", cli.ModeAnalyze),
			Entry("with embed mode", cli.ModeEmbed),
		)

		It("should pass validation when correct config file is set", func() {
			config := cli.Config{
				Mode:       cli.ModeCheck,
				ConfigPath: configFilePath(),
			}

			Expect(cli.ValidateConfig(config)).Error().ShouldNot(HaveOccurred())
			Expect(cli.ValidateConfigFile(config)).Error().ShouldNot(HaveOccurred())
		})

		It("should pass validation when embeddings are set", func() {
			config := cli.Config{
				Mode:       cli.ModeCheck,
				Embeddings: []cli.EmbeddingConfig{baseEmbeddingConfig()},
			}

			Expect(cli.ValidateConfig(config)).Error().ShouldNot(HaveOccurred())
		})

		It("should store embedding fragments under a named subfolder", func() {
			embedding := baseEmbeddingConfig()
			embedding.FragmentsPath = "/tmp/fragments"
			config := cli.Config{
				Mode:       cli.ModeCheck,
				Embeddings: []cli.EmbeddingConfig{embedding},
			}

			embedConfigs := cli.BuildEmbedCodeConfiguration(config)

			Expect(embedConfigs).To(HaveLen(1))
			Expect(embedConfigs[0].Name).To(Equal("docs"))
			Expect(embedConfigs[0].FragmentsDir).To(Equal(filepath.Join("/tmp/fragments", "docs")))
		})
	})

	Context("with invalid config", func() {

		DescribeTable("should fail validation when mode is invalid",
			func(mode string) {
				config := baseCliConfig()
				config.Mode = mode
				Expect(cli.ValidateConfig(config)).Error().Should(HaveOccurred())
			},

			Entry("with random mode", "justarandomstring"),
			Entry("with numeric mode", "123123123123"),
			Entry("with symbols mode", "!@#$%^&*()"),
			Entry("with empty mode", "         "),
		)

		It("should fail validation when config file is not exist", func() {
			invalidConfig := cli.Config{
				Mode:       cli.ModeEmbed,
				ConfigPath: "/some/path/to/config.yaml",
			}

			Expect(cli.ValidateConfigFile(invalidConfig)).Error().Should(HaveOccurred())
			Expect(cli.ValidateConfigFile(invalidConfig).Error()).Should(
				Equal("expected to use config file, but it does not exist"))
		})

		It("should fail validation when mode is not set", func() {
			invalidConfig := baseCliConfig()
			invalidConfig.Mode = ""

			Expect(cli.ValidateConfig(invalidConfig)).Error().Should(HaveOccurred())
			Expect(cli.ValidateConfig(invalidConfig).Error()).Should(Equal("mode must be set"))
		})

		It("should fail validation when docs path is missed", func() {
			invalidConfig := baseCliConfig()
			invalidConfig.BaseDocsPath = ""

			Expect(cli.ValidateConfig(invalidConfig)).Error().Should(HaveOccurred())
			Expect(cli.ValidateConfig(invalidConfig).Error()).Should(Equal(
				"`code-path` and `docs-path` must both be set"))
		})

		It("should fail validation when config, code and docs paths are set at the same time", func() {
			invalidConfig := baseCliConfig()
			invalidConfig.ConfigPath = configFilePath()

			Expect(cli.ValidateConfigFile(invalidConfig)).Error().Should(HaveOccurred())
			Expect(cli.ValidateConfigFile(invalidConfig).Error()).Should(Equal(
				"config path cannot be set when code-path, docs-path or optional params are set"))
		})

		It("should fail validation when embeddings and root paths are set at the same time", func() {
			invalidConfig := baseCliConfig()
			invalidConfig.Embeddings = []cli.EmbeddingConfig{baseEmbeddingConfig()}

			Expect(cli.ValidateConfig(invalidConfig)).Error().Should(HaveOccurred())
			Expect(cli.ValidateConfig(invalidConfig).Error()).Should(Equal(
				"`code-path` and `docs-path` cannot be set when `embeddings` are set"))
		})

		It("should fail validation when embeddings and root optional params are set at the same time", func() {
			invalidConfig := cli.Config{
				Mode:         cli.ModeCheck,
				CodeIncludes: []string{"**/*.java"},
				Embeddings:   []cli.EmbeddingConfig{baseEmbeddingConfig()},
			}

			Expect(cli.ValidateConfig(invalidConfig)).Error().Should(HaveOccurred())
			Expect(cli.ValidateConfig(invalidConfig).Error()).Should(Equal(
				"root optional embedding options cannot be set when `embeddings` are set"))
		})

		It("should fail validation when embedding name is missed", func() {
			invalidConfig := cli.Config{
				Mode:       cli.ModeCheck,
				Embeddings: []cli.EmbeddingConfig{baseEmbeddingConfig()},
			}
			invalidConfig.Embeddings[0].Name = ""

			Expect(cli.ValidateConfig(invalidConfig)).Error().Should(HaveOccurred())
			Expect(cli.ValidateConfig(invalidConfig).Error()).Should(Equal(
				"embedding #1: `name` must be set"))
		})

		It("should fail validation when embedding names are duplicated", func() {
			embedding := baseEmbeddingConfig()
			duplicateEmbedding := baseEmbeddingConfig()
			invalidConfig := cli.Config{
				Mode:       cli.ModeCheck,
				Embeddings: []cli.EmbeddingConfig{embedding, duplicateEmbedding},
			}

			Expect(cli.ValidateConfig(invalidConfig)).Error().Should(HaveOccurred())
			Expect(cli.ValidateConfig(invalidConfig).Error()).Should(Equal(
				"duplicate embedding names detected:\n- docs"))
		})

		It("should correctly convert embeddings to a few configs", func() {
			config := cli.Config{
				Mode:       cli.ModeCheck,
				ConfigPath: "../test/resources/config_files/embeddings_config.yml",
			}

			fileConfig, err := cli.FillArgsFromConfigFile(config)
			embedConfigs := cli.BuildEmbedCodeConfiguration(fileConfig)

			Expect(err).ToNot(HaveOccurred())
			Expect(embedConfigs).To(HaveLen(3))
			Expect(embedConfigs[0].Name).To(Equal("java"))
			Expect(embedConfigs[0].CodeRoots[0].Path).To(Equal("test/resources/code/java"))
			Expect(embedConfigs[0].DocumentationRoot).To(Equal("test/resources/docs"))
			Expect(embedConfigs[0].FragmentsDir).To(Equal(
				filepath.Join(configuration.DefaultFragmentsDir, "java")))
			Expect(embedConfigs[1].Name).To(Equal("kotlin"))
			Expect(embedConfigs[1].CodeRoots[0].Path).To(Equal("test/resources/code/kotlin"))
			Expect(embedConfigs[1].DocumentationRoot).To(Equal("test/resources/docs/nested-dir-1"))
			Expect(embedConfigs[1].FragmentsDir).To(Equal(
				filepath.Join(configuration.DefaultFragmentsDir, "kotlin")))
			Expect(embedConfigs[2].Name).To(Equal("nested-java"))
			Expect(embedConfigs[2].DocumentationRoot).To(
				Equal("test/resources/docs/nested-dir-1/nested-dir-3"))
			Expect(embedConfigs[2].FragmentsDir).To(Equal(
				filepath.Join(configuration.DefaultFragmentsDir, "nested-java")))
			Expect(embedConfigs[2].Separator).To(Equal("---"))
		})

	})

})

func baseCliConfig() cli.Config {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	parentDir := filepath.Dir(currentDir)

	return cli.Config{
		Mode:          cli.ModeCheck,
		BaseDocsPath:  parentDir + "/test/resources/docs",
		BaseCodePaths: _type.NamedPathList{_type.NamedPath{Path: parentDir + "/test/resources/code"}},
	}
}

func baseEmbeddingConfig() cli.EmbeddingConfig {
	baseConfig := baseCliConfig()
	return cli.EmbeddingConfig{
		Name:      "docs",
		CodePaths: baseConfig.BaseCodePaths,
		DocsPath:  baseConfig.BaseDocsPath,
	}
}

func configFilePath() string {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	parentDir := filepath.Dir(currentDir)

	return parentDir + "/test/resources/config_files/correct_config.yml"
}
