// Copyright 2020, TeamDev. All rights reserved.
//
// Redistribution and use in source and/or binary forms, with or without
// modification, must retain the above copyright notice and the following
// disclaimer.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package embedding_test

import (
	"embed-code/embed-code-go/configuration"
	"embed-code/embed-code-go/embedding"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNotUpToDate(t *testing.T) {
	config := configuration.PrepareConfiguration("./test/resources/docs")
	assert.False(t, embedding.IsUpToDate(config, "./test/resources/docs/whole-file-fragment.md"))
}

func TestUpToDate(t *testing.T) {
	config := configuration.PrepareConfiguration("./test/resources/docs")
	embedding.EmbedCode(config, "./test/resources/docs/whole-file-fragment.md")
	assert.True(t, embedding.IsUpToDate(config, "./test/resources/docs/whole-file-fragment.md"))
}

func TestNothingToUpdate(t *testing.T) {
	config := configuration.PrepareConfiguration("./test/resources/docs")
	assert.True(t, embedding.IsUpToDate(config, "./test/resources/docs/no-embedding-doc.md"))
}

func TestNonExistingFile(t *testing.T) {
	config := configuration.PrepareConfiguration("./test/resources/docs")
	assert.PanicsWithError(t, errors.FileNotFound("./test/resources/docs/non-existing-file.md"), func() {
		embedding.EmbedCode(config, "./test/resources/docs/non-existing-file.md")
	})
}

func TestDirectoryInsteadOfFile(t *testing.T) {
	config := configuration.PrepareConfiguration("./test/resources/docs")
	assert.PanicsWithError(t, errors.IsDirectory("./test/resources/docs"), func() {
		embedding.EmbedCode(config, "./test/resources/docs")
	})
}
