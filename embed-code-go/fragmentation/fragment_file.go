package fragmentation

import (
	"crypto/sha1"
	"embed-code/embed-code-go/configuration"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type FragmentFile struct {
	CodeFile      string
	FragmentName  string
	Configuration configuration.Configuration
}

// Iniitalizers

// Composes a FragmentFile for the given fragment in the given code file.
//
// @param [string] code_file an absolute path to a code file
// @param [string] fragment the fragment
func NewFragmentFileFromAbsolute(
	codeFile string,
	fragmentName string,
	configuration configuration.Configuration,
) FragmentFile {

	absoluteCodeRoot, err := filepath.Abs(configuration.CodeRoot)
	if err != nil {
		panic(err)
	}
	relativeCodeFile, err := filepath.Rel(absoluteCodeRoot, codeFile)
	if err != nil {
		panic(err)
	}

	return FragmentFile{
		CodeFile:      relativeCodeFile,
		FragmentName:  fragmentName,
		Configuration: configuration,
	}
}

// Private methods

func (fragmentFile FragmentFile) absolutePath() string {

	fileExtension := filepath.Ext(fragmentFile.CodeFile)
	fragmentsAbsDir, err := filepath.Abs(fragmentFile.Configuration.FragmentsDir)
	if err != nil {
		panic(err)
	}

	if fragmentFile.FragmentName == DefaultFragment {
		return filepath.Join(fragmentsAbsDir, fragmentFile.CodeFile)
	} else {
		withoutExtension := strings.TrimSuffix(fragmentFile.CodeFile, fileExtension)
		filename := fmt.Sprintf("%s-%s", withoutExtension, fragmentFile.getFragmentHash())
		return filepath.Join(fragmentsAbsDir, filename+fileExtension)
	}
}

// TODO: Investigate why does it use the hash of fragment name instead of the hash of fragment content
func (fragmentFile FragmentFile) getFragmentHash() string {
	hash := sha1.New()
	hash.Write([]byte(fragmentFile.FragmentName))
	sha1_hash := hex.EncodeToString(hash.Sum(nil))[:8]
	return sha1_hash
}

// Public methods

func (fragmentFile FragmentFile) Write(text string) {
	byteStr := []byte(text)
	filePath := fragmentFile.absolutePath()
	os.WriteFile(filePath, byteStr, 0777)
}

func (fragmentFile FragmentFile) Content() []string {
	path := fragmentFile.absolutePath()
	isPathFileExits, err := isFileExists(path)
	if isPathFileExits {
		return readLines(path)
	} else {
		panic(err)
	}
}

func (fragmentFile FragmentFile) String() string {
	return fragmentFile.absolutePath()
}
