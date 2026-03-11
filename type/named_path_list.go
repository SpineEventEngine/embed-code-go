package _type

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"strings"
)

// NamedPath represents a path that may optionally have a name.
type NamedPath struct {
	Name string `yaml:"name"`
	Path string `yaml:"path"`
}

// NamedPathList is a list of NamedPath values.
type NamedPathList []NamedPath

// UnmarshalYAML converts YAML nodes into NamedPathList objects.
//
// Supported formats:
//
//	Single string:
//
//	  paths: "../examples"
//
//	List of strings:
//
//	  paths:
//	    - "../examples"
//	    - "../runtime"
//
//	List of NamedPath objects:
//
//	  paths:
//	    - name: examples
//	      path: "../examples"
//	    - name: runtime
//	      path: "../runtime"
func (pathList *NamedPathList) UnmarshalYAML(value *yaml.Node) error {
	switch value.Kind {

	case yaml.ScalarNode:
		*pathList = []NamedPath{
			{Path: strings.TrimSpace(value.Value)},
		}
		return nil

	case yaml.SequenceNode:
		result := make([]NamedPath, 0, len(value.Content))

		for _, node := range value.Content {
			switch node.Kind {

			case yaml.ScalarNode:
				result = append(result, NamedPath{
					Path: strings.TrimSpace(node.Value),
				})

			case yaml.MappingNode:
				var p NamedPath
				if err := node.Decode(&p); err != nil {
					return err
				}
				result = append(result, p)

			default:
				return fmt.Errorf("invalid named path format")
			}
		}

		*pathList = result
		return nil
	default:
		return fmt.Errorf("invalid format for named paths")
	}
}
