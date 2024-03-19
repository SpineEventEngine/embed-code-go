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

// Splits the given file into fragments.
//
// The fragments are named parts of the file that are surrounded by "fragment brackets":
// ```
//
//	class HelloWorld {
//	    // #docfragment main_method
//	    public static void main(String[] argv) {
//	        // #docfragment printing
//	        System.out.println("Hello World");
//	        // #enddocfragment printing
//	    }
//	    // #enddocfragment main_method
//	}
//
// ```
//
// Fragments with the same name may appear multiple times in the same document.
//
// Even if no fragments are defined explicitly, the whole file is always a fragment on its own.
package fragmentation

import (
	"embed-code/embed-code-go/configuration"
	"fmt"
	"path/filepath"
)

const (
	FragmentStart = "#docfragment"
	FragmentEnd   = "#enddocfragment"
)

type Fragmentation struct {
	Configuration configuration.Configuration
	SourcesRoot   string
	CodeFile      string
}

// TODO: handle the errors
func NewFragmentation(
	config configuration.Configuration,
	sourcesRootRelative string,
	codeFileRelative string,
) Fragmentation {

	fragmentation := Fragmentation{}

	absoluteSourcesRoot, err := filepath.Abs(sourcesRootRelative)
	fragmentation.SourcesRoot = absoluteSourcesRoot
	if err != nil {
		fmt.Println(err)
	}

	absoluteCodeFile, err := filepath.Abs(codeFileRelative)
	fragmentation.CodeFile = absoluteCodeFile
	if err != nil {
		fmt.Println(err)
	}

	fragmentation.Configuration = config

	return fragmentation
}

// TODO: Implement
// @return (content, fragments) a refined content of the file to be cut into fragments, and the Fragments
func (fragmentation Fragmentation) fragmentize([]string, []Fragmentation) {}
