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

package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"embed-code/embed-code-go/configuration"
	"embed-code/embed-code-go/logging"
	_type "embed-code/embed-code-go/type"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// TestMainOrchestrator runs the main package specs.
func TestMainOrchestrator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Main Orchestrator Suite")
}

var _ = Describe("Main orchestrator", func() {
	It("should aggregate check errors while printing stale files", func() {
		staleConfig, staleDocPath := writeMainModeFixture("stale documentation")
		errorConfig, _ := writeMainModeFixture("documentation with a missing source")
		errorConfig.CodeRoots = _type.NamedPathList{_type.NamedPath{Path: GinkgoT().TempDir()}}

		output := captureStdout(func() {
			err := checkByConfigs([]configuration.Configuration{staleConfig, errorConfig})

			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(And(
				ContainSubstring("the documentation files are not up-to-date with code files"),
				ContainSubstring("code file `file://"),
				ContainSubstring("Example.java` not found"),
			))
		})

		Expect(output).Should(ContainSubstring("File to update:\n"))
		Expect(output).Should(ContainSubstring("- " + logging.FileReference(staleDocPath) + ".\n"))
	})

	It("should print updated files after embedding", func() {
		config, docPath := writeMainModeFixture("outdated documentation")

		output := captureStdout(func() {
			Expect(embedByConfigs([]configuration.Configuration{config})).Should(Succeed())
		})

		Expect(output).Should(ContainSubstring("File updated:\n"))
		Expect(output).Should(ContainSubstring("- " + logging.FileReference(docPath) + ".\n"))
		content, err := os.ReadFile(docPath)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(string(content)).Should(ContainSubstring("class Example {}"))
	})

	It("should aggregate embed errors from multiple configs", func() {
		firstConfig, firstDocPath := writeMainModeFixture("first documentation")
		secondConfig, secondDocPath := writeMainModeFixture("second documentation")
		firstConfig.CodeRoots = _type.NamedPathList{_type.NamedPath{Path: GinkgoT().TempDir()}}
		secondConfig.CodeRoots = _type.NamedPathList{_type.NamedPath{Path: GinkgoT().TempDir()}}

		output := captureStdout(func() {
			err := embedByConfigs([]configuration.Configuration{firstConfig, secondConfig})

			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(And(
				ContainSubstring(logging.FileReferenceWithLine(firstDocPath, 3)),
				ContainSubstring(logging.FileReferenceWithLine(secondDocPath, 3)),
				ContainSubstring("Example.java` not found"),
			))
		})

		Expect(output).Should(BeEmpty())
	})

	It("should print plural headings for multiple files", func() {
		firstPath := filepath.ToSlash(filepath.Join(GinkgoT().TempDir(), "first.md"))
		secondPath := filepath.ToSlash(filepath.Join(GinkgoT().TempDir(), "second.md"))

		output := captureStdout(func() {
			printFiles("One file:", "Many files:", []string{firstPath, secondPath})
		})

		Expect(output).Should(Equal(
			"Many files:\n" +
				"- " + logging.FileReference(firstPath) + ".\n" +
				"- " + logging.FileReference(secondPath) + ".\n",
		))
	})

	It("should capture output larger than the pipe buffer", func() {
		largeOutput := bytes.Repeat([]byte("x"), 128*1024)

		output := captureStdout(func() {
			_, err := os.Stdout.Write(largeOutput)
			Expect(err).ShouldNot(HaveOccurred())
		})

		Expect(output).Should(Equal(string(largeOutput)))
	})
})

// captureStdout runs action and returns text written to standard output.
func captureStdout(action func()) string {
	originalStdout := os.Stdout
	outputFile, err := os.CreateTemp("", "embed-code-stdout-*.txt")
	Expect(err).ShouldNot(HaveOccurred())
	os.Stdout = outputFile
	defer func() {
		os.Stdout = originalStdout
		_ = outputFile.Close()
		_ = os.Remove(outputFile.Name())
	}()

	action()

	_, err = outputFile.Seek(0, io.SeekStart)
	Expect(err).ShouldNot(HaveOccurred())
	output, err := io.ReadAll(outputFile)
	Expect(err).ShouldNot(HaveOccurred())

	return string(output)
}

// writeMainModeFixture creates one source file and one stale documentation file.
func writeMainModeFixture(docTitle string) (configuration.Configuration, string) {
	root := GinkgoT().TempDir()
	codeRoot := filepath.Join(root, "code")
	docsRoot := filepath.Join(root, "docs")
	Expect(os.MkdirAll(codeRoot, 0700)).To(Succeed())
	Expect(os.MkdirAll(docsRoot, 0700)).To(Succeed())
	Expect(os.WriteFile(
		filepath.Join(codeRoot, "Example.java"),
		[]byte("class Example {}\n"),
		0600,
	)).To(Succeed())
	docPath := filepath.ToSlash(filepath.Join(docsRoot, "doc.md"))
	Expect(os.WriteFile(
		docPath,
		[]byte("# "+docTitle+"\n\n<embed-code file=\"Example.java\"/>\n```java\nold source\n```\n"),
		0600,
	)).To(Succeed())

	config := configuration.NewConfiguration()
	config.DocumentationRoot = docsRoot
	config.CodeRoots = _type.NamedPathList{_type.NamedPath{Path: codeRoot}}
	config.DocIncludes = []string{"doc.md"}

	return config, docPath
}
