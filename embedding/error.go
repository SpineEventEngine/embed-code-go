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

package embedding

import (
	"fmt"

	"embed-code/embed-code-go/logging"
)

// ProcessingError wraps a parser or source-resolution error with documentation location.
type ProcessingError struct {
	// DocFilePath is the path to the documentation file being processed.
	DocFilePath string

	// Line is the one-based documentation line where processing failed.
	Line int

	// Err is the underlying parser or source-resolution error.
	Err error
}

// Error returns a user-facing description of the failed documentation processing operation.
//
// Returns formatted processing error text.
func (e ProcessingError) Error() string {
	return fmt.Sprintf(
		"failed to embed code fragment into doc file `%s`: %s",
		logging.FileReferenceWithLine(e.DocFilePath, e.Line),
		e.Err,
	)
}

// Unwrap returns the parser or source-resolution error that caused processing to fail.
//
// Returns the underlying parser or source-resolution error.
func (e ProcessingError) Unwrap() error {
	return e.Err
}
