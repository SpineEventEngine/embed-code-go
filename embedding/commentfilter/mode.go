// Copyright 2026, TeamDev. All rights reserved.
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

package commentfilter

import "fmt"

// CommentFilterMode controls which source comments are retained in embedded snippets.
type CommentFilterMode string

const (
	// RetainAll keeps all comments in the embedded source.
	RetainAll CommentFilterMode = "all"
	// RetainNone removes all comments recognized for the source language.
	RetainNone CommentFilterMode = "none"
	// RetainDocumentation keeps only API documentation comments.
	RetainDocumentation CommentFilterMode = "documentation"
	// RetainRegular keeps inline and block comments that are not documentation comments.
	RetainRegular CommentFilterMode = "regular"
	// RetainInline keeps only inline comments such as `//` and `#`.
	RetainInline CommentFilterMode = "inline"
	// RetainBlock keeps only block comments such as `/* */`.
	RetainBlock CommentFilterMode = "block"
)

// ParseMode converts an embed-code `comments` attribute value into a CommentFilterMode.
func ParseMode(value string) (CommentFilterMode, error) {
	switch CommentFilterMode(value) {
	case "":
		return RetainAll, nil
	case RetainAll, RetainNone, RetainDocumentation, RetainRegular, RetainInline, RetainBlock:
		return CommentFilterMode(value), nil
	default:
		return "", fmt.Errorf("unsupported comments value `%s`; expected one of "+
			"`all`, `none`, `documentation`, `regular`, `inline`, or `block`", value)
	}
}
