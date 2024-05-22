package analyzing

import (
	"embed-code/embed-code-go/configuration"
	"embed-code/embed-code-go/embedding"
	"embed-code/embed-code-go/fragmentation"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
)

func AnalyzeAll(config configuration.Configuration) {
	var problemFiles []string
	documentationRoot := config.DocumentationRoot
	docPatterns := config.DocIncludes
	for _, pattern := range docPatterns {
		globString := strings.Join([]string{documentationRoot, pattern}, "/")
		documentationFiles, _ := doublestar.FilepathGlob(globString)
		for _, documentationFile := range documentationFiles {
			processor := embedding.NewEmbeddingProcessor(documentationFile, config)
			err := processor.Embed()
			if err != nil {
				problemFiles = append(problemFiles, err.Error())
			}
		}
	}
	fragmentation.WriteLinesToFile("problem-files.txt", problemFiles)
}
