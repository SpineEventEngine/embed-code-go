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
	"bufio"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

func ShouldFragmentize(file string) bool {
	info, err := os.Stat(file)
	if err != nil {
		panic(err)
	}

	isFile := !info.IsDir()
	isValidEncoding := IsValidEncoding(file)

	return isFile && isValidEncoding
}

func IsFileUTF8Encoded(filename string) (bool, error) {
	// Read the entire file into memory
	content, err := os.ReadFile(filename)
	if err != nil {
		return false, err
	}

	// Check if the content contains valid UTF-8 characters
	isUTF8 := utf8.Valid(content)

	return isUTF8, nil
}

// If all the characters fall within the ASCII range (0 to 127), it’s likely an ASCII-encoded file.
func IsFileASCIIEncoded(filename string) (bool, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return false, err
	}

	for _, char := range content {
		if char > 127 {
			return false, nil
		}
	}

	return true, nil
}

func IsValidEncoding(file string) bool {
	isUTF8Encoded, err := IsFileUTF8Encoded(file)
	if err != nil {
		panic(err)
	}

	isASCIIEncoded, err := IsFileASCIIEncoded(file)
	if err != nil {
		panic(err)
	}

	return isUTF8Encoded || isASCIIEncoded
}

func EnsureDirExists(dirPath string) {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err := os.MkdirAll(dirPath, os.ModeDir)
		if err != nil {
			panic(err)
		}
	}
}

func IsFileExists(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false, nil
	}
	return err == nil, err
}

func UnquoteNameAndClean(name string) string {
	r, _ := regexp.Compile("\"(.*)\"")
	nameQuoted := r.FindString(name)
	nameCleaned, _ := strconv.Unquote(nameQuoted)
	return nameCleaned
}

func Lookup(line string, prefix string) []string {
	if strings.Contains(line, prefix) {
		fragmentsStart := strings.Index(line, prefix) + len(prefix) + 1 // 1 for trailing space after the prefix
		unquotedFragmentNames := []string{}
		for _, fragmentName := range strings.Split(line[fragmentsStart:], ",") {
			quotedFragmentName := strings.Trim(fragmentName, "\n\t ")
			unquotedFragmentName := UnquoteNameAndClean(quotedFragmentName)
			unquotedFragmentNames = append(unquotedFragmentNames, unquotedFragmentName)
		}
		return unquotedFragmentNames
	} else {
		return []string{}
	}
}

func GetFragmentStarts(line string) []string {
	return Lookup(line, FragmentStart)
}

func GetFragmentEnds(line string) []string {
	return Lookup(line, FragmentEnd)
}

func ReadLines(filePath string) []string {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	lines := []string{}
	defer file.Close()

	r := bufio.NewReader(file)

	for {
		line, _, err := r.ReadLine()
		lines = append(lines, string(line))
		if err != nil {
			break
		}
	}
	return lines
}
