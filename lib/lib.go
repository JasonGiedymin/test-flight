// ## Deps
//     go get github.com/SpaceMonkeyGo/errors
//

package lib

import (
	// "fmt"
	"github.com/SpaceMonkeyGo/errors"
	"io/ioutil"
	"os"
)

var (
	BadDir        = errors.NewClass("Can't read the current directory")
	FileCheckFail = errors.NewClass("File Check Failed")
)

// todo: break out into separate file
type RequiredFile struct {
	name     string
	fileName string
	fileType string // [f]ile, [d]ir
}

const defaultDir = "."

var requiredFiles = [...]RequiredFile{
	RequiredFile{name: "voom json build file", fileName: "build.json", fileType: "f"},
	RequiredFile{name: "voom docker dir", fileName: "docker", fileType: "d"},
	RequiredFile{name: "ansible meta dir", fileName: "meta", fileType: "d"},
	RequiredFile{name: "ansible tasks dir", fileName: "tasks", fileType: "d"},
	RequiredFile{name: "ansible test dir", fileName: "tests", fileType: "d"},
	RequiredFile{name: "ansible var dir", fileName: "vars", fileType: "d"},
}

// Converts []FileInfo => []string
func ConvertFiles(files []os.FileInfo) []string {
	convertedFiles := []string{}
	for _, value := range files {
		convertedFiles = append(convertedFiles, value.Name())
	}

	return convertedFiles
}

func findFile(filesFound []string, requiredFile RequiredFile) (bool, error) {
	for _, file := range filesFound {
		if file == requiredFile.fileName {
			return true, nil
		}
	}

	return false, FileCheckFail.New("Required file/dir not found: [%#v]", requiredFile.fileName)
}

// TODO: goroutine + channels to further optimize
// Required files are necessary before starting.
// Check that each are found within the local dir.
func CheckFiles(filesFound []string) (bool, error) {
	for _, file := range requiredFiles {
		found, err := findFile(filesFound, file)

		if !found {
			return false, err
		}
	}

	return true, nil
}

/*
 * Reads the current directory and returns
 *   - bool: if all the required files were found
 */
func HasRequiredFiles(dir *string) (bool, error) {
	var currDir = *dir

	if currDir == "" {
		currDir = defaultDir
	}

	filesFromDisk, err := ioutil.ReadDir(currDir)

	if err != nil {
		return false, err
	}

	return CheckFiles(ConvertFiles(filesFromDisk))
}
