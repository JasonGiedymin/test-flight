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

// todo: break out into separate file
type RequiredFile struct {
  name          string
  fileName      string
  fileType      string // [f]ile, [d]ir
  requiredFiles []RequiredFile
}

var (
  defaultDir    = "./"
  mainYaml      = RequiredFile{name: "main yaml", fileName: "main.yml", fileType: "f"}
  BadDir        = errors.NewClass("Can't read the current directory")
  FileCheckFail = errors.NewClass("File Check Failed")
  ansibleFiles  = []RequiredFile{mainYaml}
)

var RequiredFiles = []RequiredFile{
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

func findFile(filesFound []string, requiredFile RequiredFile, currDir string) (bool, error) {
  for _, file := range filesFound {
    if file == requiredFile.fileName {
      if len(requiredFile.requiredFiles) > 0 {
        nextDir := currDir + "/" + requiredFile.fileName
        _, err := HasRequiredFiles(&nextDir, requiredFile.requiredFiles)
        if err != nil {
          return false, err
        }
      }
      return true, nil
    }
  }

  return false, FileCheckFail.New("Required file/dir not found: [%v/%v]", currDir, requiredFile.fileName)
}

// TODO: goroutine + channels to further optimize
/*
 * Reads the current directory and returns
 *   - bool: if all the required files were found
 */
func HasRequiredFiles(dir *string, requiredFiles []RequiredFile) (bool, error) {
  var currDir = *dir

  if currDir == "" {
    currDir = defaultDir
  }

  filesFromDisk, err := ioutil.ReadDir(currDir)

  if err != nil {
    return false, err
  }

  for _, requiredFile := range requiredFiles {
    found, err := findFile(ConvertFiles(filesFromDisk), requiredFile, currDir)

    if !found {
      return false, err
    }
  }

  return true, nil
}
