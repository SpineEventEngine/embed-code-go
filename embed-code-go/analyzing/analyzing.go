package analyzing

import (
	"embed-code/embed-code-go/configuration"
	"embed-code/embed-code-go/embedding"
	"embed-code/embed-code-go/embedding/parsing"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
)

type ProblemFile struct {
	Context parsing.ParsingContext
	Error   error
}

func (p ProblemFile) String() string {
	relativeMarkdownPath, err := filepath.Rel(p.Context.Embedding.Configuration.DocumentationRoot, p.Context.MarkdownFile)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%s | %s — %s | %s", relativeMarkdownPath, p.Context.Embedding.CodeFile, p.Context.Embedding.Fragment, p.Error.Error())
}

func AnalyzeAll(config configuration.Configuration) {
	var problemFiles []string
	documentationRoot := config.DocumentationRoot
	docPatterns := config.DocIncludes
	for _, pattern := range docPatterns {
		globString := strings.Join([]string{documentationRoot, pattern}, "/")
		documentationFiles, _ := doublestar.FilepathGlob(globString)
		for _, documentationFile := range documentationFiles {
			processor := embedding.NewEmbeddingProcessor(documentationFile, config)
			context, err := processor.EmbedWithError()
			if err != nil {
				problemFiles = append(problemFiles, ProblemFile{Context: context, Error: err}.String())
			}
		}
	}
	WriteStringsToFile("problem-files.txt", problemFiles)
}

func WriteStringsToFile(filepath string, strings []string) {
	file, err := os.Create(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	for _, s := range strings {
		_, err := file.WriteString(s + "\n")
		if err != nil {
			panic(err)
		}
	}
}
