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

//nolint:testpackage // Covers package-private filesystem hooks used for deterministic errors.
package files

import (
	"errors"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Files internals", func() {
	It("should treat a path removed after glob resolution as missing", func() {
		filePath := writeTemporaryFile()

		withStatError(os.ErrNotExist, func() {
			exists, err := IsFileExist(filePath)

			Expect(err).ShouldNot(HaveOccurred())
			Expect(exists).Should(BeFalse())
		})
	})

	It("should report stat errors after glob resolution", func() {
		filePath := writeTemporaryFile()
		statErr := errors.New("stat failed")

		withStatError(statErr, func() {
			exists, err := IsFileExist(filePath)

			Expect(exists).Should(BeFalse())
			Expect(err).Should(MatchError(statErr))
		})
	})
})

// withStatError replaces path inspection with a deterministic error.
func withStatError(statErr error, action func()) {
	originalStatPath := statPath
	statPath = func(_ string) (os.FileInfo, error) {
		return nil, statErr
	}
	defer func() {
		statPath = originalStatPath
	}()

	action()
}

// writeTemporaryFile creates an existing path for glob resolution.
func writeTemporaryFile() string {
	filePath := filepath.Join(GinkgoT().TempDir(), "file.txt")
	Expect(os.WriteFile(filePath, []byte("content"), os.FileMode(WritePermission))).To(Succeed())

	return filePath
}
