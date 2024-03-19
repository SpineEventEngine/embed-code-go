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
//
// TODO: handle the errors
// Composes a FragmentFile for the given fragment in the given code file.
//
// @param [string] code_file an absolute path to a code file
// @param [string] fragment the fragment
func (fragmentFile FragmentFile) FromAbsoluteFile(
	codeFile string,
	fragmentName string,
	configuration configuration.Configuration,
) FragmentFile {

	relativeCodeFile, err := filepath.Rel(configuration.CodeRoot, codeFile)
	if err != nil {
		fmt.Println("Error:", err)
	}

	return FragmentFile{
		CodeFile:      relativeCodeFile,
		FragmentName:  fragmentName,
		Configuration: configuration,
	}
}

// Private methods

// TODO: Handle the errors
func (fragmentFile FragmentFile) absolutePath() string {

	fileExtension := filepath.Ext(fragmentFile.CodeFile)
	fragmentsAbsDir, err := filepath.Abs(fragmentFile.Configuration.FragmentsDir)
	if err != nil {
		fmt.Println("Error:", err)
	}

	if fragmentFile.FragmentName == DefaultFragment {
		return filepath.Join(fragmentsAbsDir, fragmentFile.CodeFile)
	} else {
		baseName := strings.TrimSuffix(fragmentFile.CodeFile, filepath.Ext(fragmentFile.CodeFile))
		withoutExtension := filepath.Join(filepath.Dir(fragmentFile.CodeFile), baseName)
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
	os.WriteFile(fragmentFile.absolutePath(), byteStr, 0777)
}

// TODO: Handle the errors
func (fragmentFile FragmentFile) Content() []string {
	path := fragmentFile.absolutePath()
	isPathFileExits, err := isFileExists(path)
	if isPathFileExits {
		return readLines(path)
	} else {
		fmt.Println("Error:", err)
	}
	return []string{}
}

func (fragmentFile FragmentFile) String() string {
	return fragmentFile.absolutePath()
}
