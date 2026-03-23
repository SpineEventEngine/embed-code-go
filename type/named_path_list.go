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
