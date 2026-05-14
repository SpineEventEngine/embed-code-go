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

package fragmentation

import (
	"fmt"
	"path/filepath"
	"strings"

	config "embed-code/embed-code-go/configuration"
	_type "embed-code/embed-code-go/type"
)

// resolverCacheLimit is the maximum number of source files retained in the resolver cache.
const resolverCacheLimit = 100

// cachedFragmentation stores cleaned source lines and parsed fragments for one source file.
type cachedFragmentation struct {
	lines     []string
	fragments map[string]Fragment
}

// resolverCache stores source fragmentations already resolved during the current run.
var resolverCache = newCache[resolvedSource, cachedFragmentation](
	resolverCacheLimit,
	loadSourceFragments,
)

// ResolveContent returns source lines for the requested code file fragment.
//
// Named fragments are extracted directly from the source file on demand and cached by source file.
func ResolveContent(codePath string, fragmentName string, config config.Configuration) ([]string, error) {
	if fragmentName == "" {
		fragmentName = DefaultFragmentName
	}

	source, found, err := resolveSource(codePath, config)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, unresolvedSourceError(codePath, fragmentName, config)
	}

	content, err := cachedSourceFragments(source)
	if err != nil {
		return nil, err
	}

	fragment, found := content.fragments[fragmentName]
	if !found {
		codeFileReference := "file://" + source.absolutePath
		return nil, fmt.Errorf("fragment `%s` from code file `%s` not found",
			fragmentName, codeFileReference)
	}

	return fragmentLines(fragment, content.lines, config.Separator), nil
}

// ClearResolverCache removes cached source fragmentations.
func ClearResolverCache() {
	resolverCache.clear()
}

// resolvedSource describes a source file resolved from a user-facing embedding path.
type resolvedSource struct {
	root         _type.NamedPath
	relativePath string
	absolutePath string
}

// resolveSource resolves the user-facing code path to an included textual source file.
func resolveSource(codePath string, config config.Configuration) (resolvedSource, bool, error) {
	codeRootName, relativePath, named := splitNamedPath(codePath)
	for _, root := range config.CodeRoots {
		if named && strings.TrimSpace(root.Name) != codeRootName {
			continue
		}

		source, err := sourceFromRoot(root, relativePath)
		if err != nil {
			return resolvedSource{}, false, err
		}
		if !shouldDoFragmentation(source.absolutePath) {
			continue
		}

		return source, true, nil
	}

	return resolvedSource{}, false, nil
}

// splitNamedPath separates a named-code-root prefix from a code path.
func splitNamedPath(codePath string) (string, string, bool) {
	normalizedPath := filepath.ToSlash(filepath.Clean(codePath))
	if !strings.HasPrefix(normalizedPath, NamedPathPrefix) {
		return "", normalizedPath, false
	}

	withoutPrefix := strings.TrimPrefix(normalizedPath, NamedPathPrefix)
	rootName, relativePath, _ := strings.Cut(withoutPrefix, "/")

	return rootName, relativePath, true
}

// sourceFromRoot builds a source path from a code root and a relative path.
func sourceFromRoot(root _type.NamedPath, relativePath string) (resolvedSource, error) {
	rootAbs, err := filepath.Abs(root.Path)
	if err != nil {
		return resolvedSource{}, err
	}

	return resolvedSource{
		root:         root,
		relativePath: filepath.FromSlash(relativePath),
		absolutePath: filepath.Join(rootAbs, filepath.FromSlash(relativePath)),
	}, nil
}

// cachedSourceFragments returns cached source fragmentation for a resolved source file.
func cachedSourceFragments(source resolvedSource) (cachedFragmentation, error) {
	return resolverCache.get(source)
}

// loadSourceFragments reads and fragments the source file when it is not already cached.
func loadSourceFragments(source resolvedSource) (cachedFragmentation, error) {
	fragmentation := NewFragmentation(source.absolutePath, source.root, config.Configuration{})
	lines, fragments, err := fragmentation.DoFragmentation()
	if err != nil {
		return cachedFragmentation{}, err
	}
	return cachedFragmentation{
		lines:     lines,
		fragments: fragments,
	}, nil
}

// fragmentLines renders a fragment into lines.
func fragmentLines(fragment Fragment, lines []string, separator string) []string {
	text := fragment.text(lines, separator)
	if text == "" {
		return []string{}
	}

	return strings.Split(strings.TrimSuffix(text, "\n"), "\n")
}

// unresolvedSourceError builds an error for a code path that cannot be resolved from sources.
func unresolvedSourceError(codePath string, fragmentName string, config config.Configuration) error {
	codeFileReference, err := codeFileReference(codePath, config)
	if err != nil {
		return err
	}
	if fragmentName == DefaultFragmentName {
		return fmt.Errorf("code file `%s` not found", codeFileReference)
	}

	return fmt.Errorf("fragment `%s` from code file `%s` not found",
		fragmentName, codeFileReference)
}

// codeFileReference builds a user-facing source reference for unresolved code paths.
func codeFileReference(codePath string, config config.Configuration) (string, error) {
	codeRootName, relativePath, named := splitNamedPath(codePath)
	for _, root := range config.CodeRoots {
		if named && strings.TrimSpace(root.Name) != codeRootName {
			continue
		}

		source, err := sourceFromRoot(root, relativePath)
		if err != nil {
			return "", err
		}
		if named {
			return fmt.Sprintf("%s (%s)", codePath, source.absolutePath), nil
		}
		if len(config.CodeRoots) == 1 {
			return source.absolutePath, nil
		}
	}

	if named {
		return "", fmt.Errorf("code root with name `%s` not found for path `%s`",
			codeRootName, codePath)
	}

	return codePath, nil
}
