package analyzing

import (
	"embed-code/embed-code-go/configuration"
	"embed-code/embed-code-go/embedding"
	"embed-code/embed-code-go/fragmentation"
	"fmt"
	"os"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
)

const analyticsDir = "./build/analytics"
const embeddingsNotFoundFile = "embeddings-not-found-files.txt"
const embeddingChangedFiles = "embeddings-changed-files.txt"
const permission = 0644

// Analyzes all documentation files.
//
// If any error occurs during embedding, it is written to the analytics file with all the needed information.
//
// config — a configuration for embedding.
func AnalyzeAll(config configuration.Configuration) {
	var problemFiles []string
	var foundEmbeddingChangedFiles []string
	documentationRoot := config.DocumentationRoot
	docPatterns := config.DocIncludes
	for _, pattern := range docPatterns {
		globString := strings.Join([]string{documentationRoot, pattern}, "/")
		documentationFiles, _ := doublestar.FilepathGlob(globString)
		for _, documentationFile := range documentationFiles {
			processor := embedding.NewEmbeddingProcessor(documentationFile, config)
			changedEmbeddings, err := processor.FindChangedEmbeddings()
			if err != nil {
				problemFiles = append(problemFiles, err.Error())
			}
			if len(changedEmbeddings) > 0 {
				docRelPath := fragmentation.BuildDocRelativePath(documentationFile, config)
				for _, changedEmbedding := range changedEmbeddings {
					line := fmt.Sprintf("%s : %s", docRelPath, changedEmbedding.String())
					foundEmbeddingChangedFiles = append(foundEmbeddingChangedFiles, line)
				}
			}
		}
	}

	os.MkdirAll(analyticsDir, permission)
	fragmentation.WriteLinesToFile(fmt.Sprintf("%s/%s", analyticsDir, embeddingsNotFoundFile), problemFiles)
	fragmentation.WriteLinesToFile(fmt.Sprintf("%s/%s", analyticsDir, embeddingChangedFiles), foundEmbeddingChangedFiles)
}
