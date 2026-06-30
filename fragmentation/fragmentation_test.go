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
	"os"
	"path/filepath"
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

// TestFragmentation runs the fragmentation test suite.
func TestFragmentation(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Data Suite")
}

var _ = Describe("Fragmentation", func() {
	var config configuration.Configuration
	var resolver *fragmentation.Resolver

	BeforeEach(func() {
		resolver = newTestResolver(fragmentation.DefaultResolverCacheLimit)
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
		content := resolveTestFragment(resolver, correctFragmentsFileName, "main()", config)

		Expect(content).Should(Equal([]string{
			"public static void main(String[] args) {",
			indent + "System.out.println(\"Hello world\");",
			"}",
		}))
	})

	It("should resolve fragments without an end marker through the end of the file", func() {
		content := resolveTestFragment(
			resolver,
			unclosedFragmentFileName,
			"Fragment that never ends",
			config,
		)

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

	It("should skip a non-UTF-8 source and use the next code root", func() {
		invalidRoot := GinkgoT().TempDir()
		validRoot := GinkgoT().TempDir()
		fileName := "Example.java"
		Expect(os.WriteFile(filepath.Join(invalidRoot, fileName), []byte{0xff}, 0600)).To(Succeed())
		Expect(os.WriteFile(
			filepath.Join(validRoot, fileName),
			[]byte("class Example {}"),
			0600,
		)).To(Succeed())
		config.CodeRoots = _type.NamedPathList{
			_type.NamedPath{Path: invalidRoot},
			_type.NamedPath{Path: validRoot},
		}

		content, err := resolver.ResolveContent(
			fileName,
			fragmentation.DefaultFragmentName,
			config,
		)

		Expect(err).ShouldNot(HaveOccurred())
		Expect(content).Should(Equal([]string{"class Example {}"}))
	})

	It("should isolate cached source content between resolvers", func() {
		sourceRoot := GinkgoT().TempDir()
		fileName := "Example.java"
		sourcePath := filepath.Join(sourceRoot, fileName)
		config.CodeRoots = _type.NamedPathList{_type.NamedPath{Path: sourceRoot}}
		Expect(os.WriteFile(sourcePath, []byte("class First {}"), 0600)).To(Succeed())

		firstContent, err := resolver.ResolveContent(
			fileName,
			fragmentation.DefaultFragmentName,
			config,
		)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(os.WriteFile(sourcePath, []byte("class Second {}"), 0600)).To(Succeed())

		cachedContent, err := resolver.ResolveContent(
			fileName,
			fragmentation.DefaultFragmentName,
			config,
		)
		Expect(err).ShouldNot(HaveOccurred())
		freshResolver := newTestResolver(fragmentation.DefaultResolverCacheLimit)
		freshContent, err := freshResolver.ResolveContent(
			fileName,
			fragmentation.DefaultFragmentName,
			config,
		)
		Expect(err).ShouldNot(HaveOccurred())

		Expect(firstContent).Should(Equal([]string{"class First {}"}))
		Expect(cachedContent).Should(Equal(firstContent))
		Expect(freshContent).Should(Equal([]string{"class Second {}"}))
	})

	It("should evict the least recently used source from a limited resolver cache", func() {
		sourceRoot := GinkgoT().TempDir()
		config.CodeRoots = _type.NamedPathList{_type.NamedPath{Path: sourceRoot}}
		resolver = newTestResolver(2)
		writeSourceFile(sourceRoot, "A.java", "class A { String version = \"first\"; }")
		writeSourceFile(sourceRoot, "B.java", "class B { String version = \"first\"; }")
		writeSourceFile(sourceRoot, "C.java", "class C { String version = \"first\"; }")

		firstA := resolveSourceFile(resolver, "A.java", config)
		_ = resolveSourceFile(resolver, "B.java", config)
		cachedA := resolveSourceFile(resolver, "A.java", config)
		writeSourceFile(sourceRoot, "A.java", "class A { String version = \"second\"; }")
		writeSourceFile(sourceRoot, "B.java", "class B { String version = \"second\"; }")
		_ = resolveSourceFile(resolver, "C.java", config)

		stillCachedA := resolveSourceFile(resolver, "A.java", config)
		reloadedB := resolveSourceFile(resolver, "B.java", config)

		Expect(cachedA).Should(Equal(firstA))
		Expect(stillCachedA).Should(Equal(firstA))
		Expect(reloadedB).Should(Equal([]string{"class B { String version = \"second\"; }"}))
	})

	It("should reject cache limits below one", func() {
		resolver, err := fragmentation.NewResolver(0)

		Expect(resolver).Should(BeNil())
		Expect(err).Should(HaveOccurred())
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
		content := resolveTestFragment(resolver, complexFragmentsFileName, "Main", config)

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
		mainContent := resolveTestFragment(resolver, twoFragmentsFileName, "Main", config)
		helloContent := resolveTestFragment(resolver, twoFragmentsFileName, "Hello", config)

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
		mainContent := resolveTestFragment(resolver, overlappingFragmentsFileName, "Main", config)
		helloContent := resolveTestFragment(resolver, overlappingFragmentsFileName, "Hello", config)

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

// buildTestFragmentation creates Fragmentation for a source fixture.
func buildTestFragmentation(testFileName string,
	config configuration.Configuration) fragmentation.Fragmentation {
	codeRoot := config.CodeRoots[0]
	testFilePath := fmt.Sprintf("%s/org/example/%s", codeRoot.Path, testFileName)
	frag, err := fragmentation.NewFragmentation(testFilePath)

	Expect(err).ShouldNot(HaveOccurred())

	return frag
}

// doTestFragmentation fragments a source fixture and returns its rendered lines.
func doTestFragmentation(
	testFileName string,
	config configuration.Configuration,
) ([]string, map[string]fragmentation.Fragment) {
	frag := buildTestFragmentation(testFileName, config)

	lines, fragments, err := frag.DoFragmentation()

	Expect(err).ShouldNot(HaveOccurred())

	return lines, fragments
}

// resolveTestFragment returns one named fragment from a source fixture.
func resolveTestFragment(
	resolver *fragmentation.Resolver,
	testFileName string,
	fragmentName string,
	config configuration.Configuration,
) []string {
	content, err := resolver.ResolveContent(
		fmt.Sprintf("org/example/%s", testFileName),
		fragmentName,
		config,
	)

	Expect(err).ShouldNot(HaveOccurred())

	return content
}

// newTestResolver creates a resolver with a test-provided cache limit.
func newTestResolver(cacheLimit int) *fragmentation.Resolver {
	resolver, err := fragmentation.NewResolver(cacheLimit)
	Expect(err).ShouldNot(HaveOccurred())

	return resolver
}

// writeSourceFile writes one source fixture into a temporary source root.
func writeSourceFile(sourceRoot string, fileName string, content string) {
	Expect(os.WriteFile(filepath.Join(sourceRoot, fileName), []byte(content), 0600)).To(Succeed())
}

// resolveSourceFile returns the whole-file content resolved from a temporary source root.
func resolveSourceFile(
	resolver *fragmentation.Resolver,
	fileName string,
	config configuration.Configuration,
) []string {
	content, err := resolver.ResolveContent(
		fileName,
		fragmentation.DefaultFragmentName,
		config,
	)
	Expect(err).ShouldNot(HaveOccurred())

	return content
}
