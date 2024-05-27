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
const embeddingChangedFile = "embeddings-changed-files.txt"
const permission = 0644

// Analyzes all documentation files.
//
// If any error occurs during embedding, it is written to the analytics file with all the needed information.
//
// config — a configuration for embedding.
func AnalyzeAll(config configuration.Configuration) {
	docFiles := findDocumentationFiles(config)
	changedFiles, problemFiles := findChangedAndProblematicFiles(config, docFiles)

	os.MkdirAll(analyticsDir, permission)
	fragmentation.WriteLinesToFile(fmt.Sprintf("%s/%s", analyticsDir, embeddingChangedFile), changedFiles)
	fragmentation.WriteLinesToFile(fmt.Sprintf("%s/%s", analyticsDir, embeddingsNotFoundFile), problemFiles)

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

// Returns a list of documentation files that are not up-to-date with their code files.
// Also returns a list of files which cause an error.
func findChangedAndProblematicFiles(
	config configuration.Configuration,
	docFiles []string) (
	changedFiles []string,
	problemFiles []string) {

	for _, docFile := range docFiles {
		processor := embedding.NewEmbeddingProcessor(docFile, config)
		changedEmbeddings, err := processor.FindChangedEmbeddings()
		if err != nil {
			problemFiles = append(problemFiles, err.Error())
		}
		if len(changedEmbeddings) > 0 {
			docRelPath := fragmentation.BuildDocRelativePath(docFile, config)
			for _, changedEmbedding := range changedEmbeddings {
				line := fmt.Sprintf("%s : %s", docRelPath, changedEmbedding.String())
				changedFiles = append(changedFiles, line)
			}
		}
	}
	return
}
