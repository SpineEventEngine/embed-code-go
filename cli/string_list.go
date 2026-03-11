package cli

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"strings"
)

// StringList is a list of strings.
type StringList []string

// UnmarshalYAML implements yaml.Unmarshaler.
//
// Supported formats:
//
//	A comma-separated string:
//
//	  list: "a,b,c"
//
//	A YAML sequence:
//
//	  list:
//	    - a
//	    - b
//	    - c
func (s *StringList) UnmarshalYAML(value *yaml.Node) error {
	switch value.Kind {

	case yaml.ScalarNode:
		parts := strings.Split(value.Value, ",")
		res := make([]string, 0, len(parts))
		for _, p := range parts {
			if trimmed := strings.TrimSpace(p); trimmed != "" {
				res = append(res, trimmed)
			}
		}
		*s = res
		return nil

	case yaml.SequenceNode:
		var res []string
		for _, n := range value.Content {
			res = append(res, strings.TrimSpace(n.Value))
		}
		*s = res
		return nil
	default:
		return fmt.Errorf("invalid format for string list")
	}
}
