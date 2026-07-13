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

//nolint:testpackage // Covers package-private cache and resolver branches.
package fragmentation

import (
	"embed-code/embed-code-go/configuration"
	_type "embed-code/embed-code-go/type"
	"errors"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Fragmentation internals", func() {
	Describe("cache", func() {
		It("should replace an already stored key without adding a duplicate entry", func() {
			cache := newCache[int, string](2, func(_ int) (string, error) {
				return "loaded", nil
			})

			cache.storeLoaded(1, "first")
			cache.storeLoaded(1, "second")

			Expect(cache.values).Should(HaveKeyWithValue(1, "second"))
			Expect(cache.order).Should(Equal([]int{1}))
		})

	})

	It("should propagate partition selection errors while rendering fragments", func() {
		fragment := Fragment{
			Name: "broken",
			Partitions: []Partition{
				{StartPosition: 1, EndPosition: 1},
			},
		}

		text, err := fragment.text([]string{"only line"}, "...")

		Expect(text).Should(BeEmpty())
		Expect(err).Should(MatchError(
			"fragment partition start position 1 is outside source lines",
		))
	})

	It("should propagate fragment rendering errors through fragment line selection", func() {
		fragment := Fragment{
			Name: "broken",
			Partitions: []Partition{
				{StartPosition: 1, EndPosition: 1},
			},
		}

		lines, err := fragmentLines(fragment, []string{"only line"}, "...")

		Expect(lines).Should(BeNil())
		Expect(err).Should(MatchError(
			"fragment partition start position 1 is outside source lines",
		))
	})

	It("should propagate reference errors from unresolved source errors", func() {
		config := configuration.Configuration{
			CodeRoots: _type.NamedPathList{
				_type.NamedPath{Name: "library", Path: "/tmp"},
			},
		}

		err := unresolvedSourceError("$missing/Example.java", DefaultFragmentName, config)

		Expect(err).Should(MatchError(
			"code root with name `missing` not found for path `$missing/Example.java`",
		))
	})

	It("should describe missing default fragments as source load failures", func() {
		message := missingFragmentLogMessage(DefaultFragmentName, "/tmp/Source.java")

		Expect(message).Should(And(
			ContainSubstring("Could not load source file"),
			ContainSubstring("file://"),
			ContainSubstring("Source.java"),
		))
	})

	Describe("resolver error propagation", func() {
		It("should propagate cached source reload errors while resolving content", func() {
			sourceRoot := GinkgoT().TempDir()
			sourcePath := filepath.Join(sourceRoot, "Example.java")
			Expect(os.WriteFile(sourcePath, []byte("class Example {}"), 0600)).To(Succeed())
			loads := 0
			resolver := Resolver{
				cache: newCache[absolutePath, fragmentedFile](0,
					func(_ absolutePath) (fragmentedFile, error) {
						loads++
						if loads > 1 {
							return fragmentedFile{}, errors.New("source reload failed")
						}

						return fragmentedFile{
							lines: []string{"class Example {}"},
							fragments: map[string]Fragment{
								DefaultFragmentName: CreateDefaultFragment(),
							},
						}, nil
					},
				),
			}
			config := configuration.Configuration{
				CodeRoots: _type.NamedPathList{
					_type.NamedPath{Path: sourceRoot},
				},
			}

			content, err := resolver.ResolveContent("Example.java", DefaultFragmentName, config)

			Expect(content).Should(BeNil())
			Expect(err).Should(MatchError("source reload failed"))
			Expect(loads).Should(Equal(2))
		})

		It("should propagate absolute path errors while resolving a source in a root", func() {
			resolver, err := NewResolver(DefaultResolverCacheLimit)
			Expect(err).ShouldNot(HaveOccurred())

			withAbsolutePathError(func() {
				source, found, err := resolver.resolveSourceInRoot(
					_type.NamedPath{Path: "relative-root"},
					"Example.java",
				)

				Expect(source).Should(BeEmpty())
				Expect(found).Should(BeFalse())
				Expect(err).Should(MatchError("absolute path failed"))
			})
		})

		It("should propagate absolute path errors while loading source fragments", func() {
			withAbsolutePathError(func() {
				content, err := loadSourceFragments("Example.java")

				Expect(content).Should(Equal(fragmentedFile{}))
				Expect(err).Should(MatchError("absolute path failed"))
			})
		})

		It("should propagate absolute path errors while building code file references", func() {
			config := configuration.Configuration{
				CodeRoots: _type.NamedPathList{
					_type.NamedPath{Path: "relative-root"},
				},
			}

			withAbsolutePathError(func() {
				reference, err := codeFileReference("Example.java", config)

				Expect(reference).Should(BeEmpty())
				Expect(err).Should(MatchError("absolute path failed"))
			})
		})
	})
})

// withAbsolutePathError replaces absolute path resolution with a deterministic failure.
func withAbsolutePathError(action func()) {
	originalMakeAbsolutePath := makeAbsolutePath
	makeAbsolutePath = func(_ string) (string, error) {
		return "", errors.New("absolute path failed")
	}
	defer func() {
		makeAbsolutePath = originalMakeAbsolutePath
	}()

	action()
}
