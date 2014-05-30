package lib

import (
  "io/ioutil"
  "os"
)

// var Info = factorlog.New(os.Stdout, factorlog.NewStdFormatter(`%{Color "green"}%{Date} %{Time} %{File}:%{Line} %{Message}%{Color "reset"}`))
// var Error = factorlog.New(os.Stdout, factorlog.NewStdFormatter(`%{Color "red"}%{Date} %{Time} %{File}:%{Line} %{Message}%{Color "reset"}`))

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
      if len(requiredFile.requiredFiles) > 0 && requiredFile.FileType == "d" {
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

func CreateFile(dir *string, requiredFile RequiredFile) (*os.File, error) {
  var fileName = *dir + "/" + requiredFile.fileName
  var err error
  var file *os.File

  if (requiredFile.FileType == "d") {
    if err = os.Mkdir(fileName, 0755); err != nil {
      return nil, err
    }
  } else if (requiredFile.FileType == "f") {
    if _, err = os.Create(fileName); err != nil {
      return nil, nil
    }
  }

  return file, nil
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

func HasRequiredFile(dir *string, requiredFile RequiredFile) (bool, error) {
  return HasRequiredFiles(dir, []RequiredFile{requiredFile})
}
