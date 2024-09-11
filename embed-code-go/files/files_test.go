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

package files_test

import (
	"embed-code/embed-code-go/files"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestFiles(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Data Suite")
}

var _ = Describe("Files actions", func() {

	Context("testing directories", func() {

		It("should return true as the referenced directory exists", func() {
			currentDir, _ := os.Getwd()

			Expect(files.IsDirExist(currentDir)).Error().ShouldNot(HaveOccurred())
			Expect(files.IsDirExist(currentDir)).Should(BeTrue())
		})

		It("should return false as the referenced path does not exist", func() {
			filePath := "/a/path/to/nowhere"

			Expect(files.IsDirExist(filePath)).Error().ShouldNot(HaveOccurred())
			Expect(files.IsDirExist(filePath)).Should(BeFalse())
		})

		It("should return error as the referenced path is a file", func() {
			currentDir, _ := os.Getwd()
			path := filepath.Dir(currentDir) + "/test/resources/config_files/correct_config.yml"

			Expect(files.IsDirExist(path)).Error().Should(HaveOccurred())
			_, err := files.IsDirExist(path)
			Expect(err.Error()).Should(
				Equal(fmt.Sprintf("%s is a file, the directory was expected", path)))
		})
	})

	Context("testing files", func() {

		It("should return true as the referenced file exists", func() {
			currentDir, _ := os.Getwd()
			filePath := filepath.Dir(currentDir) + "/test/resources/config_files/correct_config.yml"

			Expect(files.IsFileExist(filePath)).Error().ShouldNot(HaveOccurred())
			Expect(files.IsFileExist(filePath)).Should(BeTrue())
		})

		It("should return false as the referenced file does not exist", func() {
			filePath := "/path/to/nowhere/file.txt"

			Expect(files.IsFileExist(filePath)).Error().ShouldNot(HaveOccurred())
			Expect(files.IsFileExist(filePath)).Should(BeFalse())
		})

		It("should return false as the referenced path point to directory", func() {
			currentDir, _ := os.Getwd()

			Expect(files.IsFileExist(currentDir)).Error().Should(HaveOccurred())
			_, err := files.IsFileExist(currentDir)
			Expect(err.Error()).Should(
				Equal(fmt.Sprintf("%s is a directory, the file was expected", currentDir)))
		})
	})
})
