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
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// TODO: make checking for valid encoding
// TODO: handle the errors
func shouldFragmentize(file string) bool {
	info, err := os.Stat(file)
	if err != nil {
		fmt.Println(err)
	}

	isFile := !info.IsDir()
	isValidEncoding := isValidEncoding(file)

	return isFile && isValidEncoding
}

func isValidEncoding(file string) bool {
	return true
}

// TODO: handle the errors
func ensureDirExists(dirPath string) {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err := os.Mkdir(dirPath, os.ModeDir)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func isFileExists(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false, nil
	}
	return err == nil, err
}

func unquoteNameAndClean(name string) string {
	r, _ := regexp.Compile("\"(.*)\"")
	nameQuoted := r.FindString(name)
	nameCleaned, _ := strconv.Unquote(nameQuoted)
	return nameCleaned
}

func lookup(line string, prefix string) []string {
	if strings.Contains(line, prefix) {
		fragmentsStart := strings.Index(line, prefix) + len(prefix) + 1 // 1 for trailing space after the prefix
		unquotedFragmentNames := []string{}
		for _, fragmentName := range strings.Split(line[fragmentsStart:len(line)-1], ",") {
			quotedFragmentName := strings.Trim(fragmentName, "\n\t ")
			unquotedFragmentName := unquoteNameAndClean(quotedFragmentName)
			unquotedFragmentNames = append(unquotedFragmentNames, unquotedFragmentName)
		}
		return unquotedFragmentNames
	} else {
		return []string{}
	}
}

func getFragmentStarts(line string) []string {
	return lookup(line, FragmentStart)
}

func getFragmentEnds(line string) []string {
	return lookup(line, FragmentEnd)
}

func readLines(filePath string) []string {
	file, _ := os.Open(filePath)
	lines := []string{}
	defer file.Close()

	r := bufio.NewReader(file)

	for {
		line, _, err := r.ReadLine()
		if len(line) > 0 {
			lines = append(lines, string(line))
		}
		if err != nil {
			break
		}
	}
	return lines
}
