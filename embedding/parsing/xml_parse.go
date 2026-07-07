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

// quoteEscapedXMLLine escapes characters that XML forbids inside attribute
// values, so instruction patterns may contain raw source-code characters.
//
// Backslash-escaped quotes become the `&quot;` entity, and the XML
// metacharacters `&`, `<`, and `>` appearing inside a quoted attribute value
// are escaped to their entities. Source-line patterns such as
// `line="if (a < b)"` or `end="while (x && y)"` routinely contain these
// characters; without escaping, xml.Unmarshal rejects the whole instruction
// and the parser fails with an InstructionParseError. Markup outside quoted
// values, such as the `<embed-code` header and the `/>` terminator, is left
// untouched.
func quoteEscapedXMLLine(xmlLine string) string {
	var builder strings.Builder
	builder.Grow(len(xmlLine))
	insideValue := false
	for i := 0; i < len(xmlLine); i++ {
		char := xmlLine[i]
		if char == '\\' && i+1 < len(xmlLine) && xmlLine[i+1] == '"' {
			builder.WriteString("&quot;")
			i++

			continue
		}
		if char == '"' {
			insideValue = !insideValue
			builder.WriteByte(char)

			continue
		}
		if insideValue {
			switch char {
			case '&':
				// Preserve a pre-escaped entity such as `&quot;`; escape a
				// raw ampersand so xml.Unmarshal accepts it.
				if entity := xmlEntityPrefix(xmlLine[i:]); entity != "" {
					builder.WriteString(entity)
					i += len(entity) - 1
				} else {
					builder.WriteString("&amp;")
				}

				continue
			case '<':
				builder.WriteString("&lt;")

				continue
			case '>':
				builder.WriteString("&gt;")

				continue
			}
		}
		builder.WriteByte(char)
	}

	return builder.String()
}

// xmlEntityPrefix returns the leading XML character entity reference in text,
// or an empty string when text does not begin with one.
//
// It recognizes the predefined entities (`&amp;`, `&lt;`, `&gt;`, `&quot;`,
// `&apos;`) and numeric character references such as `&#39;` or `&#x1F;`, so a
// pattern author may include a pre-escaped entity without it being re-escaped.
func xmlEntityPrefix(text string) string {
	semicolon := strings.IndexByte(text, ';')
	if semicolon < 2 {
		return ""
	}
	name := text[1:semicolon]
	switch name {
	case "amp", "lt", "gt", "quot", "apos":
		return text[:semicolon+1]
	}
	if isNumericCharRef(name) {
		return text[:semicolon+1]
	}

	return ""
}

// isNumericCharRef reports whether name is the body of an XML numeric character
// reference, such as `#39` (decimal) or `#x1F` (hexadecimal).
func isNumericCharRef(name string) bool {
	if len(name) < 2 || name[0] != '#' {
		return false
	}
	digits := name[1:]
	hexadecimal := digits[0] == 'x' || digits[0] == 'X'
	if hexadecimal {
		digits = digits[1:]
	}
	if digits == "" {
		return false
	}
	for i := 0; i < len(digits); i++ {
		if !isReferenceDigit(digits[i], hexadecimal) {
			return false
		}
	}

	return true
}

// isReferenceDigit reports whether char is a valid digit for a numeric
// character reference, allowing hexadecimal digits when hexadecimal is set.
func isReferenceDigit(char byte, hexadecimal bool) bool {
	if char >= '0' && char <= '9' {
		return true
	}
	if !hexadecimal {
		return false
	}

	return (char >= 'a' && char <= 'f') || (char >= 'A' && char <= 'F')
}
