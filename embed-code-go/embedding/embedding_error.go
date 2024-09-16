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
	"fmt"

	"embed-code/embed-code-go/embedding/parsing"
)

// EmbeddingError Describes an error which occurs if something goes wrong during embedding.
type EmbeddingError struct {
	Context parsing.Context
}

func (err EmbeddingError) Error() string {
	errorString := fmt.Sprintf("embedding error for file `%s`.", err.Context.MarkdownFile)

	if len(err.Context.EmbeddingsNotFound) > 0 {
		embeddingsNotFoundStr := "\nMissing embeddings: \n"
		for _, emb := range err.Context.EmbeddingsNotFound {
			embeddingsNotFoundStr += fmt.Sprintf(
				"%s — %s\n",
				emb.CodeFile,
				emb.Fragment)
		}
		errorString += embeddingsNotFoundStr
	}

	if len(err.Context.UnacceptedEmbeddings) > 0 {
		unacceptedEmbeddingStr := "\nUnaccepted embeddings: \n"
		for _, emb := range err.Context.UnacceptedEmbeddings {
			unacceptedEmbeddingStr += fmt.Sprintf(
				"%s — %s\n",
				emb.CodeFile,
				emb.Fragment)
		}
		errorString += unacceptedEmbeddingStr
	}

	return errorString
}
