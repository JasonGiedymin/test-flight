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
  defaultDir    = "./"
  BadDir        = errors.NewClass("Can't read the current directory")
  FileCheckFail = errors.NewClass("File Check Failed")
)

// todo: break out into separate file
type RequiredFile struct {
  name          string
  fileName      string
  fileType      string // [f]ile, [d]ir
  requiredFiles []RequiredFile
}

// const defaultDir = "."

// var mainYaml RequiredFile

// mainYaml = new(RequiredFile)
var mainYaml = RequiredFile{name: "main yaml", fileName: "main.yml", fileType: "f"}
var ansibleFiles = []RequiredFile{mainYaml}

var requiredFiles = [...]RequiredFile{
  {name: "voom json build file", fileName: "build.json", fileType: "f"},
  {name: "voom docker dir", fileName: "docker", fileType: "d"},
  {name: "ansible meta dir", fileName: "meta", fileType: "d"},
  {name: "ansible tasks dir", fileName: "tasks", fileType: "d"},
  {name: "ansible test dir", fileName: "tests", fileType: "d"},
  {name: "ansible var dir", fileName: "vars", fileType: "d", requiredFiles: ansibleFiles},
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
