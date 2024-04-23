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

package embedding_instruction

import (
	"encoding/xml"
	"fmt"
)

const xmlStringHeader string = "embed-code"

// Needed for xml.Unmarshal parsing. The fields are filling up during the parsing.
//
// XMLName — a name of the tag in XML line.
//
// Attrs — a list of xml.Attr. The xml.Attr contains both names and values of attributes.
type Item struct {
	XMLName xml.Name
	Attrs   []xml.Attr `xml:",any,attr"`
}

// Parses given XML-encoded xmlLine and returns attributes data as key-value pairs.
//
// xmlLine — a XML-encoded line.
func ParseXMLLine(xmlLine string) (map[string]string, error) {
	var root Item
	err := xml.Unmarshal([]byte(xmlLine), &root)
	if err != nil {
		return map[string]string{}, err
	}

	if root.XMLName.Local != xmlStringHeader {
		return map[string]string{}, fmt.Errorf("The provided line's header is not 'embed-code':\n%s", xmlLine)
	}

	attributes := make(map[string]string)
	for _, subItem := range root.Attrs {
		attributes[subItem.Name.Local] = subItem.Value
	}

	return attributes, nil
}
