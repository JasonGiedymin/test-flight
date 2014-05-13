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
	BadDir = errors.NewClass("Can't read the current directory")
)

// todo: break out into separate file
type RequiredFile struct {
	name     string
	fileName string
	fileType string // [f]ile, [d]ir
}

const currDir = "."

var requiredFiles = [...]RequiredFile{
	RequiredFile{name: "buildFile", fileName: "build.json", fileType: "f"},
	RequiredFile{name: "vars", fileName: "docker", fileType: "d"},
	RequiredFile{name: "vars", fileName: "meta", fileType: "d"},
	RequiredFile{name: "vars", fileName: "tasks", fileType: "d"},
	RequiredFile{name: "vars", fileName: "tests", fileType: "d"},
	RequiredFile{name: "vars", fileName: "vars", fileType: "d"},
}

// Take a FileInfo slice and convert to string slice
func ConvertFiles(files []os.FileInfo) []string {
	convertedFiles := []string{}
	for _, value := range files {
		convertedFiles = append(convertedFiles, value.Name())
	}

	return convertedFiles
}

// Check for files or return an error
func CheckForFiles() ([]string, error) {
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
