// ## Deps
//     go get github.com/SpaceMonkeyGo/errors
//

package lib

import (
	"fmt"
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

const currDir = "."

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

// Returns files as []string
func GetFiles() ([]string, error) {
	var error error
	convertedFiles := []string{}

	// Initial read of files
	files, err := ioutil.ReadDir(currDir)

	if err != nil {
		error = BadDir.New("Can't read the directory: [%#v]", currDir)
	} else {
		convertedFiles = ConvertFiles(files)
		error = nil
	}

	return convertedFiles, error
}

func CheckFiles(filesFound []string) (bool, error) {
	found := []string{}

	for _, file := range requiredFiles {
		for _, fFile := range filesFound {
			if fFile == file.fileName {
				found = append(found, fFile)
			}
		}
	}

	if len(found) == len(requiredFiles) {
		return true, nil
	} else {
		return false, FileCheckFail.New("Length was: ", len(found))
	}

}
