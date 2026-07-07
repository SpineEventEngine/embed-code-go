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

package _type_test

import (
	"testing"

	_type "embed-code/embed-code-go/type"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v3"
)

// TestTypes runs the type package specs.
func TestTypes(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Types Suite")
}

var _ = Describe("NamedPathList", func() {

	It("should unmarshal a single path string", func() {
		var config struct {
			Paths _type.NamedPathList `yaml:"paths"`
		}

		err := yaml.Unmarshal([]byte("paths: ' ../examples '\n"), &config)

		Expect(err).ShouldNot(HaveOccurred())
		Expect(config.Paths).Should(Equal(_type.NamedPathList{
			{Path: "../examples"},
		}))
	})

	It("should unmarshal a sequence of path strings", func() {
		var config struct {
			Paths _type.NamedPathList `yaml:"paths"`
		}

		err := yaml.Unmarshal([]byte("paths:\n  - ' ../examples '\n  - ../runtime\n"), &config)

		Expect(err).ShouldNot(HaveOccurred())
		Expect(config.Paths).Should(Equal(_type.NamedPathList{
			{Path: "../examples"},
			{Path: "../runtime"},
		}))
	})

	It("should unmarshal a sequence of named paths", func() {
		var config struct {
			Paths _type.NamedPathList `yaml:"paths"`
		}

		err := yaml.Unmarshal([]byte(
			"paths:\n"+
				"  - name: examples\n"+
				"    path: ../examples\n"+
				"  - name: runtime\n"+
				"    path: ../runtime\n",
		), &config)

		Expect(err).ShouldNot(HaveOccurred())
		Expect(config.Paths).Should(Equal(_type.NamedPathList{
			{Name: "examples", Path: "../examples"},
			{Name: "runtime", Path: "../runtime"},
		}))
	})

	It("should reject mapping values", func() {
		var config struct {
			Paths _type.NamedPathList `yaml:"paths"`
		}

		err := yaml.Unmarshal([]byte("paths:\n  name: examples\n  path: ../examples\n"), &config)

		Expect(err).Should(MatchError(ContainSubstring("invalid format for named paths")))
	})
})

var _ = Describe("StringList", func() {

	It("should unmarshal a comma-separated string", func() {
		var config struct {
			Patterns _type.StringList `yaml:"patterns"`
		}

		err := yaml.Unmarshal([]byte("patterns: ' docs/*.md, , guides/*.html '\n"), &config)

		Expect(err).ShouldNot(HaveOccurred())
		Expect(config.Patterns).Should(Equal(_type.StringList{"docs/*.md", "guides/*.html"}))
	})

	It("should unmarshal a sequence of strings", func() {
		var config struct {
			Patterns _type.StringList `yaml:"patterns"`
		}

		err := yaml.Unmarshal([]byte("patterns:\n  - ' docs/*.md '\n  - guides/*.html\n"), &config)

		Expect(err).ShouldNot(HaveOccurred())
		Expect(config.Patterns).Should(Equal(_type.StringList{"docs/*.md", "guides/*.html"}))
	})

	It("should reject mapping values", func() {
		var config struct {
			Patterns _type.StringList `yaml:"patterns"`
		}

		err := yaml.Unmarshal([]byte("patterns:\n  markdown: docs/*.md\n"), &config)

		Expect(err).Should(MatchError(ContainSubstring("invalid format for string list")))
	})
})
