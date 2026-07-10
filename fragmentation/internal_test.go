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

package fragmentation

import (
	"embed-code/embed-code-go/configuration"
	_type "embed-code/embed-code-go/type"

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
			Expect(cache.entries).Should(HaveLen(1))
			Expect(cache.order.Len()).Should(Equal(1))
		})

		It("should ignore eviction when the usage order is empty", func() {
			cache := newCache[int, string](1, func(_ int) (string, error) {
				return "loaded", nil
			})

			cache.evictOldest()

			Expect(cache.values).Should(BeEmpty())
			Expect(cache.entries).Should(BeEmpty())
			Expect(cache.order.Len()).Should(Equal(0))
		})

		It("should discard an invalid usage-order entry", func() {
			cache := newCache[int, string](1, func(_ int) (string, error) {
				return "loaded", nil
			})
			cache.order.PushBack("not an int key")

			cache.evictOldest()

			Expect(cache.values).Should(BeEmpty())
			Expect(cache.entries).Should(BeEmpty())
			Expect(cache.order.Len()).Should(Equal(0))
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
})
