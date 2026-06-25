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

package fragmentation

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var quotedNamePattern = regexp.MustCompile("\"(.*)\"")

const (
	// FragmentStart marks the beginning of a named source fragment.
	FragmentStart = "#docfragment"

	// FragmentEnd marks the end of a named source fragment.
	FragmentEnd = "#enddocfragment"
)

// FindDocFragments finds fragment names declared with the start marker.
//
// For example, FindDocFragments("// #docfragment \"main\",\"sub-main\"\n")
// returns ["main", "sub-main"]
//
// Parameters:
// line - provides one source line.
//
// Returns:
// []string - fragment names declared on the line.
// error - when a declaration is malformed.
func FindDocFragments(line string) ([]string, error) {
	return lookup(line, FragmentStart)
}

// FindEndDocFragments finds fragment names declared with the end marker.
//
// For example, FindEndDocFragments("// #enddocfragment \"main\",\"sub-main\"\n")
// returns ["main", "sub-main"]
//
// Parameters:
// line - provides one source line.
//
// Returns:
// []string - fragment names closed on the line.
// error - when a declaration is malformed.
func FindEndDocFragments(line string) ([]string, error) {
	return lookup(line, FragmentEnd)
}

// lookup finds fragment names in line after the given fragment marker prefix.
//
// For example, lookup("// #enddocfragment \"main\",\"sub-main\"\n", "#enddocfragment")
// returns ["main", "sub-main"]
//
// Parameters:
// line - provides one source line to search in.
// prefix - provides the fragment marker prefix, for example "#docfragment".
//
// Returns:
// []string - fragment names found on the line.
// error - when prefix is found without valid names.
func lookup(line string, prefix string) ([]string, error) {
	var unquotedNames []string
	if strings.Contains(line, prefix) {
		// 1 for trailing space after the prefix.
		fragmentsStart := strings.Index(line, prefix) + len(prefix) + 1
		if len(line) < fragmentsStart {
			return unquotedNames, fmt.Errorf(
				"found `%s` pefix without any name", prefix,
			)
		}
		for _, fragmentName := range strings.Split(line[fragmentsStart:], ",") {
			quotedName := strings.Trim(fragmentName, "\n\t ")
			unquotedName, err := unquoteName(quotedName)
			if err != nil {
				return unquotedNames, err
			}
			unquotedNames = append(unquotedNames, unquotedName)
		}
	}

	return unquotedNames, nil
}

// unquoteName removes quotes from a fragment marker name.
//
// Parameters:
// quotedName - provides a quoted fragment name.
//
// Returns:
// string - unquoted fragment name.
// error - when quotedName cannot be unquoted.
func unquoteName(quotedName string) (string, error) {
	nameQuoted := quotedNamePattern.FindString(quotedName)
	nameCleaned, err := strconv.Unquote(nameQuoted)
	if err != nil {
		return "", fmt.Errorf("failed to unquote name `%s`: %w", quotedName, err)
	}

	return nameCleaned, nil
}
