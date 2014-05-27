package config

import (
  "encoding/json"
  "io/ioutil"
)

type BuildFile struct {
  Owner             string
  ImageName         string
  Version           string
  RequiresDocker    string
  RequiresDockerUrl string
  Env               map[string]string
  Expose            []int
  Add               []ConfigFileUserAdd
}

// For specific defaults
func NewBuildFile() *BuildFile {
  return &BuildFile{
    Owner:             "Test-Flight-User",
    ImageName:         "Test-Flight-Test-Image",
    Version:           "0.0.1",
    RequiresDocker:    "",
    RequiresDockerUrl: "",
  }
}

func ReadBuildFile(path string) (*BuildFile, error) {
  jsonBlob, _ := ioutil.ReadFile(path + "/build.json")

  var buildFile = NewBuildFile()
  err := json.Unmarshal(jsonBlob, buildFile)
  if err != nil {
    return nil, err
  }

  return buildFile, nil
}
