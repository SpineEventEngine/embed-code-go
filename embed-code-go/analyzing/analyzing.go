package analyzing

import (
	"fmt"
	"os"
	"strings"

	"embed-code/embed-code-go/configuration"
	"embed-code/embed-code-go/embedding"
	"embed-code/embed-code-go/fragmentation"
	"github.com/bmatcuk/doublestar/v4"
)

const analyticsDir = "./build/analytics"
const embeddingsNotFoundFile = "embeddings-not-found-files.txt"
const embeddingChangedFile = "embeddings-changed-files.txt"

// Represents read and write permissions for the owner of the file, and read-only permissions for group and others.
const permission = 0755

// Analyzes all documentation files.
//
// If any error occurs during embedding, it is written to the analytics file with all the needed information.
//
// config — a configuration for embedding.
func AnalyzeAll(config configuration.Configuration) {
	docFiles := findDocumentationFiles(config)
	changedEmbeddings, problemEmbeddings := extractAnalyticsForDocs(config, docFiles)

	os.MkdirAll(analyticsDir, permission)
	fragmentation.WriteLinesToFile(fmt.Sprintf("%s/%s", analyticsDir, embeddingChangedFile), changedEmbeddings)
	fragmentation.WriteLinesToFile(fmt.Sprintf("%s/%s", analyticsDir, embeddingsNotFoundFile), problemEmbeddings)
}

// Finds all documentation files for given config.
func findDocumentationFiles(config configuration.Configuration) []string {
	documentationRoot := config.DocumentationRoot
	docPatterns := config.DocIncludes
	var documentationFiles []string
	for _, pattern := range docPatterns {
		globString := strings.Join([]string{documentationRoot, pattern}, "/")
		matches, err := doublestar.FilepathGlob(globString)
		if err != nil {
			panic(err)
		}
		documentationFiles = append(documentationFiles, matches...)
	}
	return documentationFiles
}

// Returns a list of embeddings that are not up-to-date with their code files.
// Also returns a list of embeddings which cause an error.
func extractAnalyticsForDocs(
	config configuration.Configuration,
	docFiles []string,
) (changedEmbeddingsLines []string, problemEmbeddingsLines []string) {

	for _, docFile := range docFiles {
		processor := embedding.NewEmbeddingProcessor(docFile, config)
		changedEmbeddings, err := processor.FindChangedEmbeddings()

		// If there is an error during embedding, it is written to the analytics file.
		if err != nil {
			problemEmbeddingsLines = append(problemEmbeddingsLines, err.Error())
		}
		// Even if error occurs, there might be embeddings that are changed.
		if len(changedEmbeddings) > 0 {
			docRelPath := fragmentation.BuildDocRelativePath(docFile, config)
			for _, changedEmbedding := range changedEmbeddings {
				line := fmt.Sprintf("%s : %s", docRelPath, changedEmbedding.String())
				changedEmbeddingsLines = append(changedEmbeddingsLines, line)
			}
		}
	}
	return
}
