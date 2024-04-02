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

package fragmentation

import (
	"os"
	"unicode/utf8"
)

// Reports whether given bytes are UTF8-encoded.
func areUTF8Encoded(bytes []byte) bool {
	return utf8.Valid(bytes)
}

// Reports whether given bytes are ASCII-encoded.
//
// If all the characters fall within the ASCII range (0 to 127), it’s likely an ASCII-encoded file.
func areASCIIEncoded(bytes []byte) bool {
	for _, char := range bytes {
		if char > 127 {
			return false
		}
	}

	return true
}

// Reports whether the file stored at filePath is encoded as a text.
//
// If file encoded in ASCII or UTF-8, it is meant to be a text file.
func IsEncodedAsText(filePath string) bool {

	// Read the entire file into memory.
	content, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	isUTF8Encoded := areUTF8Encoded(content)
	isASCIIEncoded := areASCIIEncoded(content)
	return isUTF8Encoded || isASCIIEncoded
}
