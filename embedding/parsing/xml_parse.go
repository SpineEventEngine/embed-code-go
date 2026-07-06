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

package parsing

import (
	"embed-code/embed-code-go/configuration"
	"encoding/xml"
	"fmt"
	"strings"
)

// Item contains the XML name and attributes decoded from an embedding instruction.
type Item struct {
	// XMLName is the instruction tag name.
	XMLName xml.Name

	// Attrs contains instruction attributes with their names and values.
	Attrs []xml.Attr `xml:",any,attr"`
}

// FromXML parses an XML-like `<embed-code>` tag into an Instruction.
//
// The line can be self-closing:
// `<embed-code file="org/example/Hello.java" fragment="Hello class"/>`.
// It can also use a closing tag:
// `<embed-code file="org/example/Hello.java" fragment="Hello class"></embed-code>`.
//
// Supported instruction attributes:
//   - file - mandatory relative path to the source file;
//   - fragment - optional source fragment name. When omitted, the whole file is embedded;
//   - start - optional glob-like pattern. Matching lines before it are excluded;
//   - end - optional glob-like pattern. Matching lines after it are excluded;
//   - line - optional glob-like pattern. Only the matching line is embedded;
//   - comments - optional comment filtering mode. When omitted, all comments are retained.
//
// Parameters:
// line - provides raw instruction text.
// config - provides embedding configuration.
//
// Returns:
// Instruction - parsed embedding instruction.
// error - when XML or instruction attributes are invalid.
func FromXML(line string, config configuration.Configuration) (Instruction, error) {
	fields, err := ParseXMLLine(line)
	if err != nil {
		return Instruction{}, err
	}

	return NewInstruction(fields, config)
}

// ParseXMLLine parses an XML-like `<embed-code>` tag into attribute key-value pairs.
//
// Parameters:
// xmlLine - provides raw instruction text.
//
// Returns:
// map[string]string - instruction attributes by name.
// error - when the line is not a valid embed-code XML element.
func ParseXMLLine(xmlLine string) (map[string]string, error) {
	var root Item
	err := xml.Unmarshal([]byte(quoteEscapedXMLLine(xmlLine)), &root)
	if err != nil {
		return map[string]string{}, err
	}

	if root.XMLName.Local != EmbeddingTag {
		return map[string]string{},
			fmt.Errorf("the provided line's header is not `%s`:\n%s", EmbeddingTag, xmlLine)
	}

	attributes := make(map[string]string)
	for _, subItem := range root.Attrs {
		attributes[subItem.Name.Local] = subItem.Value
	}

	return attributes, nil
}

// quoteEscapedXMLLine converts backslash-escaped quotes into XML entities.
func quoteEscapedXMLLine(xmlLine string) string {
	return strings.ReplaceAll(xmlLine, `\"`, "&quot;")
}
