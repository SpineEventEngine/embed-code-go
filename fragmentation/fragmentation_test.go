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

package fragmentation_test

import (
	"embed-code/embed-code-go/configuration"
	"embed-code/embed-code-go/fragmentation"
	_type "embed-code/embed-code-go/type"
	"fmt"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const (
	correctFragmentsFileName     = "Hello.java"
	unclosedFragmentFileName     = "Unclosed.java"
	unopenedFragmentFileName     = "Unopen.java"
	complexFragmentsFileName     = "Complex.java"
	twoFragmentsFileName         = "TwoFragments.java"
	overlappingFragmentsFileName = "OverlappingFragments.java"
	emptyFileName                = "Empty.java"
	indent                       = "    "
)

func TestFragmentation(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Data Suite")
}

var _ = Describe("Fragmentation", func() {
	var config configuration.Configuration

	BeforeEach(func() {
		fragmentation.ClearResolverCache()
		config = configuration.NewConfiguration()
		config.DocumentationRoot = "../test/resources/docs"
		config.CodeRoots = _type.NamedPathList{_type.NamedPath{Path: "../test/resources/code/java"}}
	})

	It("should do file fragmentation successfully", func() {
		lines, fragments := doTestFragmentation(correctFragmentsFileName, config)

		Expect(lines).ShouldNot(ContainElement(ContainSubstring("#docfragment")))
		Expect(lines).ShouldNot(ContainElement(ContainSubstring("#enddocfragment")))
		Expect(fragments).Should(HaveKey(fragmentation.DefaultFragmentName))
		Expect(fragments).Should(HaveKey("Without License"))
		Expect(fragments).Should(HaveKey("Hello class"))
		Expect(fragments).Should(HaveKey("main()"))
	})

	It("should resolve named fragments", func() {
		content := resolveTestFragment(correctFragmentsFileName, "main()", config)

		Expect(content).Should(Equal([]string{
			"public static void main(String[] args) {",
			indent + "System.out.println(\"Hello world\");",
			"}",
		}))
	})

	It("should resolve fragments without an end marker through the end of the file", func() {
		content := resolveTestFragment(unclosedFragmentFileName, "Fragment that never ends", config)

		Expect(content).Should(Equal([]string{
			indent + indent + "System.out.println(\"Hello world\");",
			indent + "}",
			"}",
		}))
	})

	It("should fragment an empty file", func() {
		lines, fragments := doTestFragmentation(emptyFileName, config)

		Expect(lines).Should(BeEmpty())
		Expect(fragments).Should(HaveLen(1))
		Expect(fragments).Should(HaveKey(fragmentation.DefaultFragmentName))
	})

	It("should fail on an unopened fragment", func() {
		frag := buildTestFragmentation(unopenedFragmentFileName, config)

		_, _, err := frag.DoFragmentation()

		Expect(err).Should(HaveOccurred())
	})

	Context("fragments parsing", func() {
		mainFragment := "main"
		subMainFragment := "sub-main"

		It("should correctly find fragment openings", func() {
			docFragment := fmt.Sprintf(
				"// #docfragment \"%s\",\"%s\"", mainFragment, subMainFragment)

			openings, _ := fragmentation.FindDocFragments(docFragment)
			Expect(openings).Should(HaveLen(2))
			Expect(openings[0]).Should(Equal(mainFragment))
			Expect(openings[1]).Should(Equal(subMainFragment))
		})

		It("should correctly find fragment endings", func() {
			endDocFragment := fmt.Sprintf(
				"// #enddocfragment \"%s\",\"%s\"", mainFragment, subMainFragment)

			endings, _ := fragmentation.FindEndDocFragments(endDocFragment)
			Expect(endings).Should(HaveLen(2))
			Expect(endings[0]).Should(Equal(mainFragment))
			Expect(endings[1]).Should(Equal(subMainFragment))
		})

		It("should not find fragment endings as there are openings", func() {
			docFragment := fmt.Sprintf(
				"// #docfragment \"%s\",\"%s\"", mainFragment, subMainFragment)

			openings, _ := fragmentation.FindEndDocFragments(docFragment)
			Expect(openings).Should(BeEmpty())
		})

		It("should not find fragment openings as there are endings", func() {
			endDocFragment := fmt.Sprintf(
				"// #enddocfragment \"%s\",\"%s\"", mainFragment, subMainFragment)

			endings, _ := fragmentation.FindDocFragments(endDocFragment)
			Expect(endings).Should(BeEmpty())
		})
	})

	It("should correctly parse file into many partitions", func() {
		content := resolveTestFragment(complexFragmentsFileName, "Main", config)

		expected := []string{
			"public class Main {",
			indent + config.Separator,
			indent + "public static void main(String[] args) {",
			indent + indent + config.Separator,
			indent + indent + "System.out.println(helperMethod());",
			"",
			indent + "}",
			config.Separator,
			"}",
		}
		Expect(content).Should(Equal(expected))
	})

	It("should correctly parse file with several different fragments", func() {
		mainContent := resolveTestFragment(twoFragmentsFileName, "Main", config)
		helloContent := resolveTestFragment(twoFragmentsFileName, "Hello", config)

		Expect([][]string{mainContent, helloContent}).Should(ConsistOf([][]string{
			{
				"public class TwoFragments {",
				indent + config.Separator,
				indent + "public static void main(String[] args) {",
				indent + indent + config.Separator,
				indent + indent + "System.out.println(helperMethod());",
				"",
				indent + "}",
				config.Separator,
				"}",
			},
			{
				"public static void hello(String[] args) {",
				indent + config.Separator,
				indent + "var coolText = \"Cool Text\";",
				indent + "System.out.println(coolText);",
				"}",
			},
		}))
	})

	It("should correctly parse file with several overlapping fragments", func() {
		mainContent := resolveTestFragment(overlappingFragmentsFileName, "Main", config)
		helloContent := resolveTestFragment(overlappingFragmentsFileName, "Hello", config)

		Expect([][]string{mainContent, helloContent}).Should(ConsistOf([][]string{
			{
				"public class OverlappingFragments {",
				indent + config.Separator,
				indent + "public static void main(String[] args) {",
				indent + indent + config.Separator,
				indent + indent + "System.out.println(helperMethod());",
				"",
				indent + "}",
				config.Separator,
				"}",
			},
			{
				"public class OverlappingFragments {",
				indent + config.Separator,
				indent + "public static void hello(String[] args) {",
				indent + indent + config.Separator,
				indent + indent + "var coolText = \"Cool Text\";",
				indent + indent + "System.out.println(coolText);",
				indent + "}",
				config.Separator,
				"}",
			},
		}))
	})
})

func buildTestFragmentation(testFileName string,
	config configuration.Configuration) fragmentation.Fragmentation {
	codeRoot := config.CodeRoots[0]
	testFilePath := fmt.Sprintf("%s/org/example/%s", codeRoot.Path, testFileName)

	return fragmentation.NewFragmentation(testFilePath, codeRoot, config)
}

func doTestFragmentation(
	testFileName string,
	config configuration.Configuration,
) ([]string, map[string]fragmentation.Fragment) {
	frag := buildTestFragmentation(testFileName, config)

	lines, fragments, err := frag.DoFragmentation()

	Expect(err).ShouldNot(HaveOccurred())
	return lines, fragments
}

func resolveTestFragment(
	testFileName string,
	fragmentName string,
	config configuration.Configuration,
) []string {
	content, err := fragmentation.ResolveContent(
		fmt.Sprintf("org/example/%s", testFileName),
		fragmentName,
		config,
	)

	Expect(err).ShouldNot(HaveOccurred())
	return content
}
