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

package logging //nolint:testpackage // Tests OS-specific normalization in an unexported helper.

import "testing"

// TestFileURLFromAbsolutePath verifies OS-specific path shapes are valid file URLs.
func TestFileURLFromAbsolutePath(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{
			name: "unix path",
			path: "/Users/me/project/file.go",
			want: "file:///Users/me/project/file.go",
		},
		{
			name: "windows drive path",
			path: `C:\Users\me\project\file.go`,
			want: "file:///C:/Users/me/project/file.go",
		},
		{
			name: "windows drive path with spaces",
			path: `C:\Users\me\my project\file.go`,
			want: "file:///C:/Users/me/my%20project/file.go",
		},
		{
			name: "windows unc path",
			path: `\\server\share\project\file.go`,
			want: "file://server/share/project/file.go",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := fileURLFromAbsolutePath(test.path)
			if got != test.want {
				t.Fatalf("fileURLFromAbsolutePath(%q) = %q, want %q",
					test.path, got, test.want)
			}
		})
	}
}
