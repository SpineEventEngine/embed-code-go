package cli

import (
	"embed-code/embed-code-go/configuration"
	"embed-code/embed-code-go/embedding"
	"embed-code/embed-code-go/fragmentation"
)

func CheckCodeSamples(config configuration.Configuration) {
	fragmentation.WriteFragmentFiles(config)
	embedding.CheckUpToDate(config)
}

func EmbedCodeSamples(config configuration.Configuration) {
	fragmentation.WriteFragmentFiles(config)
	embedding.EmbedAll(config)
}
