package lib

import (
  "encoding/json"
  "io/ioutil"
  // "strings"
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

func ReadBuildFile(filePath string) (*BuildFile, error) {
  jsonBlob, _ := ioutil.ReadFile(filePath)

  var buildFile = NewBuildFile()
  err := json.Unmarshal(jsonBlob, buildFile)
  if err != nil {
    return nil, err
  }

  return buildFile, nil
}

