package embedding_instruction

import (
	"encoding/xml"
	"fmt"
)

const xmlStringHeader string = "embed-code"

type Item struct {
	XMLName xml.Name
	Attrs   []xml.Attr `xml:",any,attr"`
}

func ParseXmlLine(xmlLine string) map[string]string {
	var root Item
	err := xml.Unmarshal([]byte(xmlLine), &root)
	if err != nil {
		panic(err)
	}

	if root.XMLName.Local != xmlStringHeader {
		panic(fmt.Sprintf("The provided line's header is not 'embed-code':\n%s", xmlLine))
	}

	attributes := make(map[string]string)
	for _, subItem := range root.Attrs {
		attributes[subItem.Name.Local] = subItem.Value
	}

	return attributes
}
