package lib

import (
  "encoding/json"
  "io/ioutil"
  "strings"
  // "errors"
)

// For specific defaults
func NewBuildFile() *BuildFile {
  return &BuildFile{
    Owner:             "Test-Flight-User",
    ImageName:         "Test-Flight-Test-Image",
    Tag:               "latest",
    From:              "",
    Version:           "0.0.1",
    RequiresDocker:    "",
    RequiresDockerUrl: "",
    RunTests:          false,
  }
}

func ReadBuildFile(dir string) (*BuildFile, error) {
  path := []string{dir, "build.json"}
  jsonBlob, _ := ioutil.ReadFile(strings.Join(path, "/"))

  var buildFile = NewBuildFile()
  err := json.Unmarshal(jsonBlob, buildFile)
  if err != nil {
    return nil, err
  }

  return buildFile, nil
}

// TODO: Make as a member of Parser later...
func GetBuildFile(dir string) (*BuildFile, error) {
  buildFile, err := ReadBuildFile(dir)
  if err != nil {
    // errors.New("Error reading build file: " + err.Error())
    return nil, err
  }

  return buildFile, nil
}
