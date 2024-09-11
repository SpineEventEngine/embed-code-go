/*
 * Copyright 2024, TeamDev. All rights reserved.
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
	. "embed-code/embed-code-go/cli"
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"os"
	"path/filepath"
	"testing"
)

func TestCli(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Data Suite")
}

var _ = Describe("CLI package", func() {

	It("should set default values when they are not provided", func() {
		defaultConfig := ReadArgs()

		Expect(defaultConfig.CodeIncludes).Should(Equal("**/*.*"))
		Expect(defaultConfig.DocIncludes).Should(Equal("**/*.md,**/*.html"))
		Expect(defaultConfig.FragmentsPath).Should(Equal("./build/fragments"))
		Expect(defaultConfig.Separator).Should(Equal("..."))
	})

	Context("with valid config", func() {
		var config Config

		BeforeEach(func() {
			config = baseCliConfig()
		})

		DescribeTable("should pass validation when all required args are set",
			func(mode string) {
				config.Mode = mode
				Expect(ValidateConfig(config)).Error().ShouldNot(HaveOccurred())
			},

			Entry("with check mode", ModeCheck),
			Entry("with analyze mode", ModeAnalyze),
			Entry("with embed mode", ModeEmbed),
		)

		It("should pass validation when correct config file is set", func() {
			config := Config{
				Mode:       ModeCheck,
				ConfigPath: configFilePath(),
			}

			Expect(ValidateConfig(config)).Error().ShouldNot(HaveOccurred())
			Expect(ValidateConfigFile(config.ConfigPath)).Error().ShouldNot(HaveOccurred())
		})
	})

	Context("with invalid config", func() {

		DescribeTable("should fail validation when mode is invalid",
			func(mode string) {
				config := baseCliConfig()
				config.Mode = mode
				Expect(ValidateConfig(config)).Error().Should(HaveOccurred())
			},

			Entry("with random mode", "justarandomstring"),
			Entry("with numeric mode", "123123123123"),
			Entry("with symbols mode", "!@#$%^&*()"),
			Entry("with empty mode", "         "),
		)

		It("should fail validation when config file is not exist", func() {
			invalidConfig := Config{
				Mode:       ModeEmbed,
				ConfigPath: "/some/path/to/config.yaml",
			}

			Expect(ValidateConfigFile(invalidConfig.ConfigPath)).Error().Should(HaveOccurred())
			Expect(ValidateConfigFile(invalidConfig.ConfigPath).Error()).Should(Equal(
				fmt.Sprintf("the path %s is not exist", invalidConfig.ConfigPath)))
		})

		It("should fail validation when mode is not set", func() {
			invalidConfig := baseCliConfig()
			invalidConfig.Mode = ""

			Expect(ValidateConfig(invalidConfig)).Error().Should(HaveOccurred())
			Expect(ValidateConfig(invalidConfig).Error()).Should(Equal("mode must be set"))
		})

		It("should fail validation when docs path is missed", func() {
			invalidConfig := baseCliConfig()
			invalidConfig.DocsPath = ""

			Expect(ValidateConfig(invalidConfig)).Error().Should(HaveOccurred())
			Expect(ValidateConfig(invalidConfig).Error()).Should(Equal(
				"if one of code-path and docs-path is set, the another one must be set as well"))
		})

		It("should fail validation when config, code and docs paths are set at the same time", func() {
			invalidConfig := baseCliConfig()
			invalidConfig.ConfigPath = configFilePath()

			Expect(ValidateConfig(invalidConfig)).Error().Should(HaveOccurred())
			Expect(ValidateConfig(invalidConfig).Error()).Should(Equal(
				"config path cannot be set when code-path, docs-path or optional params are set"))
		})

	})

})

func baseCliConfig() Config {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	parentDir := filepath.Dir(currentDir)

	return Config{
		Mode:     ModeCheck,
		DocsPath: parentDir + "/test/resources/docs",
		CodePath: parentDir + "/test/resources/code",
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
