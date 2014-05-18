package lib

import (
  "os"
  "io/ioutil"
)

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
      if len(requiredFile.requiredFiles) > 0 && requiredFile.fileType == "d" {
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
