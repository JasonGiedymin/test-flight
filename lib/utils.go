package lib

import (
  "./types"
  "archive/tar"
  "io/ioutil"
  "os"
  Logger "./logging"
  "strings"
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

func findFile(filesFound []string, requiredFile types.RequiredFile, currDir string) (bool, error) {
  for _, file := range filesFound {
    if file == requiredFile.FileName {
      if len(requiredFile.RequiredFiles) > 0 && requiredFile.FileType == "d" {
        nextDir := currDir + "/" + requiredFile.FileName
        _, err := HasRequiredFiles(&nextDir, requiredFile.RequiredFiles)
        if err != nil {
          return false, err
        }
      }
      return true, nil
    }
  }

  return false, FileCheckFail.New("Required file/dir not found: [%v/%v]", currDir, requiredFile.FileName)
}

func CreateFile(dir *string, requiredFile types.RequiredFile) (*os.File, error) {
  var fileName = *dir + "/" + requiredFile.FileName
  var err error
  var file *os.File

  if requiredFile.FileType == "d" {
    if err = os.Mkdir(fileName, 0755); err != nil {
      return nil, err
    }
  } else if requiredFile.FileType == "f" {
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
func HasRequiredFiles(dir *string, requiredFiles []types.RequiredFile) (bool, error) {
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

func HasRequiredFile(dir *string, requiredFile types.RequiredFile) (bool, error) {
  return HasRequiredFiles(dir, []types.RequiredFile{requiredFile})
}

func TarDirectory(tw *tar.Writer, dir string) error {
  // var filesFromDisk []os.FileInfo
  // var err error
  Logger.Trace("Taring: ", dir)

  var archive = func(files []os.FileInfo) error {
    Logger.Trace(files)

    for _, file := range files {
  		hdr := &tar.Header{
  			Name: file.Name(),
  			Size: file.Size(),
  		}
  		if err := tw.WriteHeader(hdr); err != nil {
  			Logger.Error("Could not write context archive header", err)
        return err
  		}

      fullFilePath := strings.Join([]string{dir, file.Name()}, "/")
      if bytes, err := ioutil.ReadFile(fullFilePath); err != nil {
        Logger.Error("Could not read context file: [" + fullFilePath + "]", err)
        return err
      } else {
        if _, err := tw.Write(bytes); err != nil {
          Logger.Error("Could not archive context file: [" + fullFilePath + "]", err)
          return err
        }
      }
	  }

    Logger.Trace("Successfully archived context", dir)
    return nil
  }

  if filesFromDisk, err := ioutil.ReadDir(dir); err != nil {
    Logger.Error("Error while trying to tar [" + dir + "]", err)
    return err
  } else {
    return archive(filesFromDisk)
  }
}
