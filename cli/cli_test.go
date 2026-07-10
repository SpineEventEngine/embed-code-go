/*
 * Copyright 2026, TeamDev. All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * https://www.apache.org/licenses/LICENSE-2.0
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
	"flag"
	"io"
	"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// TestCli runs the CLI test suite.
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
			Entry("with embed mode", cli.ModeEmbed),
		)

		It("should pass validation when correct config file is set", func() {
			config := cli.Config{
				Mode:       cli.ModeCheck,
				ConfigPath: configFilePath(),
			}

			Expect(cli.IsUsingConfigFile(config)).To(BeTrue())
			Expect(cli.ValidateConfig(config)).Error().ShouldNot(HaveOccurred())
			Expect(cli.ValidateConfigFile(config)).Error().ShouldNot(HaveOccurred())
		})

		It("should pass validation when no roots are set", func() {
			config := cli.Config{
				Mode: cli.ModeCheck,
			}

			Expect(cli.IsUsingConfigFile(config)).To(BeFalse())
			Expect(cli.ValidateConfig(config)).Error().ShouldNot(HaveOccurred())
		})

		It("should pass validation when embeddings are set", func() {
			config := cli.Config{
				Mode:       cli.ModeCheck,
				Embeddings: []cli.EmbeddingConfig{baseEmbeddingConfig()},
			}

			Expect(cli.ValidateConfig(config)).Error().ShouldNot(HaveOccurred())
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

		It("should reject embeddings with root optional params", func() {
			invalidConfig := cli.Config{
				Mode:        cli.ModeCheck,
				DocIncludes: []string{"**/*.md"},
				Embeddings:  []cli.EmbeddingConfig{baseEmbeddingConfig()},
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

		It("should fail validation when embedding name contains illegal folder characters", func() {
			invalidConfig := cli.Config{
				Mode:       cli.ModeCheck,
				Embeddings: []cli.EmbeddingConfig{baseEmbeddingConfig()},
			}
			invalidConfig.Embeddings[0].Name = "bad/name"

			Expect(cli.ValidateConfig(invalidConfig)).Error().Should(HaveOccurred())
			Expect(cli.ValidateConfig(invalidConfig).Error()).Should(Equal(
				"embedding `bad/name`: `name` `bad/name` is not valid, " +
					"those characters are not allowed `/\\ *?:\"<>|`"))
		})

		It("should fail validation when embedding roots are incomplete", func() {
			invalidConfig := cli.Config{
				Mode:       cli.ModeCheck,
				Embeddings: []cli.EmbeddingConfig{baseEmbeddingConfig()},
			}
			invalidConfig.Embeddings[0].DocsPath = ""

			Expect(cli.ValidateConfig(invalidConfig)).Error().Should(HaveOccurred())
			Expect(cli.ValidateConfig(invalidConfig).Error()).Should(Equal(
				"embedding `docs`: `code-path` and `docs-path` must both be set"))
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

		It("should fail validation when source code path names are duplicated", func() {
			invalidConfig := baseCliConfig()
			invalidConfig.BaseCodePaths = _type.NamedPathList{
				_type.NamedPath{Name: "samples", Path: codeResourcePath("java")},
				_type.NamedPath{Name: "samples", Path: codeResourcePath("kotlin")},
			}

			Expect(cli.ValidateConfig(invalidConfig)).Error().Should(HaveOccurred())
			Expect(cli.ValidateConfig(invalidConfig).Error()).Should(Equal(
				"duplicate source code path names detected:\n- samples"))
		})

		It("should fail validation when code path name contains illegal folder characters", func() {
			invalidConfig := baseCliConfig()
			invalidConfig.BaseCodePaths = _type.NamedPathList{
				_type.NamedPath{Name: "bad/name", Path: codeResourcePath("java")},
			}

			Expect(cli.ValidateConfig(invalidConfig)).Error().Should(HaveOccurred())
			Expect(cli.ValidateConfig(invalidConfig).Error()).Should(Equal(
				"the given code path name `bad/name` is not a valid name for the folder, " +
					"those characters are not allowed `/\\ *?:\"<>|`"))
		})

		It("should fail validation when multiple unnamed sources code paths are configured", func() {
			invalidConfig := baseCliConfig()
			invalidConfig.BaseCodePaths = _type.NamedPathList{
				_type.NamedPath{Path: codeResourcePath("java")},
				_type.NamedPath{Path: codeResourcePath("kotlin")},
			}

			Expect(cli.ValidateConfig(invalidConfig)).Error().Should(HaveOccurred())
			Expect(cli.ValidateConfig(invalidConfig).Error()).Should(Equal(
				"only one unnamed source code path is allowed"))
		})

		It("should fail validation when named and unnamed source code paths are mixed", func() {
			invalidConfig := baseCliConfig()
			invalidConfig.BaseCodePaths = _type.NamedPathList{
				_type.NamedPath{Name: "java", Path: codeResourcePath("java")},
				_type.NamedPath{Path: codeResourcePath("kotlin")},
			}

			Expect(cli.ValidateConfig(invalidConfig)).Error().Should(HaveOccurred())
			Expect(cli.ValidateConfig(invalidConfig).Error()).Should(Equal(
				"named and unnamed source code paths cannot be mixed"))
		})

		It("should pass validation and warn when duplicate docs paths are configured", func() {
			first := baseEmbeddingConfig()
			first.Name = "first"
			second := baseEmbeddingConfig()
			second.Name = "second"
			config := cli.Config{
				Mode:       cli.ModeCheck,
				Embeddings: []cli.EmbeddingConfig{first, second},
			}

			Expect(cli.ValidateConfig(config)).Error().ShouldNot(HaveOccurred())
		})

		It("should pass validation and warn when duplicate named code paths are configured", func() {
			config := baseCliConfig()
			config.BaseCodePaths = _type.NamedPathList{
				_type.NamedPath{Name: "first", Path: codeResourcePath("java")},
				_type.NamedPath{Name: "second", Path: codeResourcePath("java")},
			}

			Expect(cli.ValidateConfig(config)).Error().ShouldNot(HaveOccurred())
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
			Expect(embedConfigs[1].Name).To(Equal("kotlin"))
			Expect(embedConfigs[1].CodeRoots[0].Path).To(Equal("test/resources/code/kotlin"))
			Expect(embedConfigs[1].DocumentationRoot).To(Equal("test/resources/docs/nested-dir-1"))
			Expect(embedConfigs[2].Name).To(Equal("nested-java"))
			Expect(embedConfigs[2].DocumentationRoot).To(
				Equal("test/resources/docs/nested-dir-1/nested-dir-3"))
			Expect(embedConfigs[2].Separator).To(Equal("---"))
		})

		It("should copy command line doc excludes to the runtime config", func() {
			config := baseCliConfig()
			config.DocExcludes = []string{"old-docs/**/*.md", "drafts/**/*"}

			embedConfigs := cli.BuildEmbedCodeConfiguration(config)

			Expect(embedConfigs).To(HaveLen(1))
			Expect(embedConfigs[0].DocExcludes).To(Equal([]string(config.DocExcludes)))
		})

	})

})

var _ = Describe("CLI arguments", func() {

	It("should read command-line arguments", func() {
		config := readArgs(
			"-mode=embed",
			"-code-path=/code",
			"-docs-path=/docs",
			"-doc-includes=**/*.md, guides/*.html, ,",
			"-doc-excludes=archive/**/*, drafts/**/*.md",
			"-separator=---",
			"-config-path=config.yml",
			"-info=true",
			"-stacktrace=true",
		)

		Expect(config.Mode).To(Equal(cli.ModeEmbed))
		Expect(config.BaseCodePaths).To(Equal(_type.NamedPathList{
			_type.NamedPath{Path: "/code"},
		}))
		Expect(config.BaseDocsPath).To(Equal("/docs"))
		Expect(config.DocIncludes).To(Equal(_type.StringList{"**/*.md", "guides/*.html"}))
		Expect(config.DocExcludes).To(Equal(_type.StringList{"archive/**/*", "drafts/**/*.md"}))
		Expect(config.Separator).To(Equal("---"))
		Expect(config.ConfigPath).To(Equal("config.yml"))
		Expect(config.Info).To(BeTrue())
		Expect(config.Stacktrace).To(BeTrue())
	})

})

var _ = Describe("CLI configuration building", func() {

	It("should fill args from a config file with optional root settings", func() {
		config := cli.Config{
			Mode:       cli.ModeCheck,
			ConfigPath: "../test/resources/config_files/optional_root_config.yml",
		}

		fileConfig, err := cli.FillArgsFromConfigFile(config)

		Expect(err).ToNot(HaveOccurred())
		Expect(fileConfig.BaseCodePaths).To(Equal(_type.NamedPathList{
			_type.NamedPath{Name: "java", Path: "test/resources/code/java"},
		}))
		Expect(fileConfig.BaseDocsPath).To(Equal("test/resources/docs"))
		Expect(fileConfig.DocIncludes).To(Equal(_type.StringList{"**/*.md"}))
		Expect(fileConfig.DocExcludes).To(Equal(_type.StringList{"archive/**/*", "drafts/**/*.md"}))
		Expect(fileConfig.Separator).To(Equal("---"))
		Expect(fileConfig.Info).To(BeTrue())
		Expect(fileConfig.Stacktrace).To(BeTrue())
	})

	It("should return an error when config file YAML is invalid", func() {
		configPath := writeTempConfigFile("doc-includes: [")
		config := cli.Config{
			Mode:       cli.ModeCheck,
			ConfigPath: configPath,
		}

		_, err := cli.FillArgsFromConfigFile(config)

		Expect(err).To(HaveOccurred())
	})

	It("should build default command-line configuration without roots", func() {
		configs := cli.BuildEmbedCodeConfiguration(cli.Config{})

		Expect(configs).To(HaveLen(1))
		Expect(configs[0]).To(Equal(configuration.NewConfiguration()))
	})

	It("should build embedding configuration with optional settings", func() {
		embedding := baseEmbeddingConfig()
		embedding.CodePaths = _type.NamedPathList{
			_type.NamedPath{Name: "java", Path: codeResourcePath("java")},
			_type.NamedPath{Name: "kotlin", Path: codeResourcePath("kotlin")},
		}
		embedding.DocIncludes = []string{"guides/**/*.md"}
		embedding.DocExcludes = []string{"archive/**/*"}
		embedding.Separator = "---"
		config := cli.Config{
			Mode:       cli.ModeCheck,
			Embeddings: []cli.EmbeddingConfig{embedding},
		}

		configs := cli.BuildEmbedCodeConfiguration(config)

		Expect(configs).To(HaveLen(1))
		Expect(configs[0].Name).To(Equal("docs"))
		Expect(configs[0].CodeRoots).To(Equal(embedding.CodePaths))
		Expect(configs[0].DocumentationRoot).To(Equal(embedding.DocsPath))
		Expect(configs[0].DocIncludes).To(Equal([]string{"guides/**/*.md"}))
		Expect(configs[0].DocExcludes).To(Equal([]string{"archive/**/*"}))
		Expect(configs[0].Separator).To(Equal("---"))
	})

})

var _ = Describe("CLI processing wrappers", func() {

	It("should check selected docs through the public wrapper", func() {
		config := noEmbeddingInstructionsConfig()

		staleDocs, err := cli.CheckCodeSamples(config)

		Expect(err).ToNot(HaveOccurred())
		Expect(staleDocs).To(BeEmpty())
	})

	It("should embed selected docs through the public wrapper", func() {
		config := noEmbeddingInstructionsConfig()

		result, err := cli.EmbedCodeSamples(config)

		Expect(err).ToNot(HaveOccurred())
		Expect(result.TotalEmbeddings).To(Equal(0))
		Expect(result.UpdatedTargetFiles).To(BeEmpty())
	})

})

// baseCliConfig returns the default valid CLI config used by validation specs.
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

// baseEmbeddingConfig returns the default valid multi-target embedding config.
func baseEmbeddingConfig() cli.EmbeddingConfig {
	baseConfig := baseCliConfig()

	return cli.EmbeddingConfig{
		Name:      "docs",
		CodePaths: baseConfig.BaseCodePaths,
		DocsPath:  baseConfig.BaseDocsPath,
	}
}

// configFilePath returns the path to a valid YAML config fixture.
func configFilePath() string {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	parentDir := filepath.Dir(currentDir)

	return parentDir + "/test/resources/config_files/correct_config.yml"
}

// codeResourcePath builds an absolute path to a test source-code fixture directory.
func codeResourcePath(name string) string {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	parentDir := filepath.Dir(currentDir)

	return filepath.Join(parentDir, "test/resources/code", name)
}

// docsResourcePath builds an absolute path to a test documentation fixture path.
func docsResourcePath(name string) string {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	parentDir := filepath.Dir(currentDir)

	return filepath.Join(parentDir, "test/resources/docs", name)
}

// readArgs runs CLI argument parsing with isolated global flag state.
func readArgs(args ...string) cli.Config {
	oldArgs := os.Args
	oldCommandLine := flag.CommandLine
	defer func() {
		os.Args = oldArgs
		flag.CommandLine = oldCommandLine
	}()

	os.Args = append([]string{"embed-code"}, args...)
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	flag.CommandLine.SetOutput(io.Discard)

	return cli.ReadArgs()
}

// writeTempConfigFile writes a YAML config fixture and returns its path.
func writeTempConfigFile(content string) string {
	configFile, err := os.CreateTemp("", "embed-code-cli-*.yml")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = configFile.Close(); err != nil {
			panic(err)
		}
	}()
	if _, err = configFile.WriteString(content); err != nil {
		panic(err)
	}

	return configFile.Name()
}

// noEmbeddingInstructionsConfig builds a config selecting a document without embed-code tags.
func noEmbeddingInstructionsConfig() configuration.Configuration {
	config := configuration.NewConfiguration()
	config.CodeRoots = _type.NamedPathList{_type.NamedPath{Path: codeResourcePath("java")}}
	config.DocumentationRoot = docsResourcePath("")
	config.DocIncludes = []string{"no-embedding-doc.md"}

	return config
}
