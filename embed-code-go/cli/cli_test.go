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
	var config Config

	It("should set default values when they are not provided", func() {
		defaultConfig := ReadArgs()

		Expect(defaultConfig.CodeIncludes).Should(Equal("**/*.*"))
		Expect(defaultConfig.DocIncludes).Should(Equal("**/*.md,**/*.html"))
		Expect(defaultConfig.FragmentsPath).Should(Equal("./build/fragments"))
		Expect(defaultConfig.Separator).Should(Equal("..."))
	})

	Context("should pass validation", func() {
		BeforeEach(func() {
			config = validConfig()
		})

		DescribeTable("when all required args are set",
			func(mode string) {
				config.Mode = mode
				Expect(ValidateConfig(config)).Error().ShouldNot(HaveOccurred())
			},

			Entry("with check mode", ModeCheck),
			Entry("with analyze mode", ModeAnalyze),
			Entry("with embed mode", ModeEmbed),
		)

		It("when correct config file is set", func() {
			currentDir, _ := os.Getwd()
			parentDir := filepath.Dir(currentDir)

			config = Config{
				Mode:       "embed",
				ConfigPath: parentDir + "/test/resources/config_files/correct_config.yml",
			}

			Expect(ValidateConfig(config)).Error().ShouldNot(HaveOccurred())
			Expect(ValidateConfigFile(config.ConfigPath)).Error().ShouldNot(HaveOccurred())
		})
	})

	Context("should not pass validation", func() {

		DescribeTable("when mode is invalid",
			func(mode string) {
				config = validConfig()
				config.Mode = mode
				Expect(ValidateConfig(config)).Error().Should(HaveOccurred())
			},

			Entry("with random mode", "justarandomstring"),
			Entry("with numeric mode", "123123123123"),
			Entry("with symbols mode", "!@#$%^&*()"),
			Entry("with empty mode", "         "),
		)

	})

})

func validConfig() Config {
	currentDir, _ := os.Getwd()
	parentDir := filepath.Dir(currentDir)

	return Config{
		DocsPath: parentDir + "/test/resources/docs",
		CodePath: parentDir + "/test/resources/code",
	}
}
