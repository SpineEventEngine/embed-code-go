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

package fragmentation

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	// FragmentStart marks the beginning of a named source fragment.
	FragmentStart = "#docfragment"

	// FragmentEnd marks the end of a named source fragment.
	FragmentEnd = "#enddocfragment"
)

// attributePrefixPattern matches an attribute-shaped token at the start of marker text.
var attributePrefixPattern = regexp.MustCompile(`^[A-Za-z][A-Za-z_-]*[ \t]*=`)

// fragmentDeclaration describes one source fragment opening marker.
type fragmentDeclaration struct {
	// names contains the fragments opened by the marker.
	names []string

	// indentGroup identifies partitions that share common indentation.
	indentGroup string
}

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
	declaration, err := findDocFragmentDeclaration(line)
	if err != nil {
		return nil, err
	}

	return declaration.names, nil
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
	declaration, err := parseFragmentDeclaration(line, FragmentEnd, false)
	if err != nil {
		return nil, err
	}

	return declaration.names, nil
}

// findDocFragmentDeclaration finds fragment names and indentation metadata on an opening marker.
//
// Parameters:
// line - provides one source line.
//
// Returns:
// fragmentDeclaration - parsed opening marker, or an empty declaration when no marker exists.
// error - when the declaration is malformed.
func findDocFragmentDeclaration(line string) (fragmentDeclaration, error) {
	return parseFragmentDeclaration(line, FragmentStart, true)
}

// parseFragmentDeclaration parses names and optional indentation metadata after a marker.
//
// Attribute-shaped text after the declaration is rejected, while other text is ignored
// for compatibility with source-language comment syntax and inline marker annotations.
//
// Parameters:
// line - provides one source line to search in.
// prefix - provides the fragment marker prefix.
// allowIndentGroup - allows an indent-group attribute on the marker.
//
// Returns:
// fragmentDeclaration - parsed marker, or an empty declaration when no marker exists.
// error - when the declaration is malformed.
func parseFragmentDeclaration(
	line string,
	prefix string,
	allowIndentGroup bool,
) (fragmentDeclaration, error) {
	var declaration fragmentDeclaration
	markerPosition := strings.Index(line, prefix)
	if markerPosition < 0 {
		return declaration, nil
	}

	remainder := strings.TrimLeft(line[markerPosition+len(prefix):], "\t ")
	if remainder == "" {
		return declaration, fmt.Errorf("found `%s` prefix without any name", prefix)
	}

	for {
		name, rest, err := consumeQuotedName(remainder)
		if err != nil {
			return fragmentDeclaration{}, err
		}
		declaration.names = append(declaration.names, name)
		remainder = strings.TrimLeft(rest, "\t ")
		if !strings.HasPrefix(remainder, ",") {
			break
		}
		remainder = strings.TrimLeft(strings.TrimPrefix(remainder, ","), "\t ")
	}

	remainder = strings.TrimLeft(remainder, "\t ")
	if strings.HasPrefix(remainder, "indent-group") {
		if !allowIndentGroup {
			return fragmentDeclaration{}, fmt.Errorf(
				"indent-group is only supported by %s", FragmentStart,
			)
		}
		indentGroup, rest, err := parseIndentGroup(remainder)
		if err != nil {
			return fragmentDeclaration{}, err
		}
		declaration.indentGroup = indentGroup
		remainder = rest
	}
	if err := validateDeclarationRemainder(remainder, prefix); err != nil {
		return fragmentDeclaration{}, err
	}

	return declaration, nil
}

// consumeQuotedName parses the next quoted fragment name.
//
// Parameters:
// source - provides marker text beginning with a fragment name.
//
// Returns:
// string - unquoted fragment name.
// string - unconsumed marker text.
// error - when the name is not a valid quoted string.
func consumeQuotedName(source string) (string, string, error) {
	quotedName, remainder := consumeQuotedValue(source)
	name, err := unquoteName(quotedName)
	if err != nil {
		return "", "", err
	}

	return name, remainder, nil
}

// parseIndentGroup parses the optional indent-group marker attribute.
//
// Parameters:
// source - provides marker text beginning with indent-group.
//
// Returns:
// string - unquoted indentation group name.
// string - unconsumed marker text.
// error - when the attribute is malformed or empty.
func parseIndentGroup(source string) (string, string, error) {
	remainder := strings.TrimLeft(strings.TrimPrefix(source, "indent-group"), "\t ")
	if !strings.HasPrefix(remainder, "=") {
		return "", "", fmt.Errorf("indent-group must use the form indent-group=\"name\"")
	}
	remainder = strings.TrimLeft(strings.TrimPrefix(remainder, "="), "\t ")
	if !strings.HasPrefix(remainder, "\"") {
		value := strings.Fields(remainder)
		unquotedValue := ""
		if len(value) > 0 {
			unquotedValue = value[0]
		}

		return "", "", fmt.Errorf("indent-group value `%s` must be quoted", unquotedValue)
	}

	quotedGroup, rest := consumeQuotedValue(remainder)
	indentGroup, err := strconv.Unquote(quotedGroup)
	if err != nil {
		return "", "", fmt.Errorf("failed to unquote indent-group `%s`: %w", quotedGroup, err)
	}
	if indentGroup == "" {
		return "", "", fmt.Errorf("indent-group must not be empty")
	}

	return indentGroup, rest, nil
}

// validateDeclarationRemainder rejects unconsumed marker attributes.
//
// Parameters:
// source - provides text after the parsed declaration.
// prefix - identifies the fragment marker for diagnostics.
//
// Returns an error when the remainder begins with an attribute-shaped token.
func validateDeclarationRemainder(source string, prefix string) error {
	remainder := strings.TrimSpace(source)
	if !startsWithAttribute(remainder) {
		return nil
	}

	return fmt.Errorf("unexpected attribute after `%s` declaration: `%s`", prefix, remainder)
}

// startsWithAttribute reports whether source begins with an attribute-shaped token.
//
// Attribute names start with an ASCII letter, continue with ASCII letters, hyphens,
// or underscores, and may have horizontal whitespace before the equals sign.
//
// Parameters:
// source - provides unconsumed marker text.
//
// Returns true when source begins with a name followed by an equals sign.
func startsWithAttribute(source string) bool {
	return attributePrefixPattern.MatchString(source)
}

// consumeQuotedValue separates the first quoted string from the remaining marker text.
//
// Parameters:
// source - provides marker text beginning with a quoted string.
//
// Returns:
// string - quoted value, or the first unquoted token when no quoted value is present.
// string - unconsumed marker text.
func consumeQuotedValue(source string) (string, string) {
	if !strings.HasPrefix(source, "\"") {
		valueEnd := strings.IndexAny(source, ",\t \n")
		if valueEnd < 0 {
			return source, ""
		}

		return source[:valueEnd], source[valueEnd:]
	}

	escaped := false
	for index := 1; index < len(source); index++ {
		character := source[index]
		if character == '"' && !escaped {
			return source[:index+1], source[index+1:]
		}
		if character == '\\' {
			escaped = !escaped
		} else {
			escaped = false
		}
	}

	return source, ""
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
	nameCleaned, err := strconv.Unquote(quotedName)
	if err != nil {
		return "", fmt.Errorf("failed to unquote name `%s`: %w", quotedName, err)
	}
	if nameCleaned == "" {
		return "", fmt.Errorf("fragment name must not be empty")
	}

	return nameCleaned, nil
}
