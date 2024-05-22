// Copyright 2024, TeamDev. All rights reserved.
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

package embedding

import (
	"embed-code/embed-code-go/embedding/parsing"
	"fmt"
	"path/filepath"
)

// Describes an error which occurs if something goes wrong during embedding.
type EmbeddingError struct {
	Context       parsing.ParsingContext
	OriginalError error
}

func (err EmbeddingError) Error() string {
	relativeMarkdownPath, filepathErr := filepath.Rel(
		err.Context.Embedding.Configuration.DocumentationRoot,
		err.Context.MarkdownFile)

	if filepathErr != nil {
		panic(err)
	}

	return fmt.Sprintf("error: %s | %s — %s | %s",
		relativeMarkdownPath,
		err.Context.Embedding.CodeFile,
		err.Context.Embedding.Fragment,
		err.OriginalError.Error())

}
