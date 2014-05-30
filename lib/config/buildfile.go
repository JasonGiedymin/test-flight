package config

import (
  "encoding/json"
  "io/ioutil"
  "../types"
)

// For specific defaults
func NewBuildFile() *types.BuildFile {
  return &types.BuildFile{
    Owner:             "Test-Flight-User",
    ImageName:         "Test-Flight-Test-Image",
    Version:           "0.0.1",
    RequiresDocker:    "",
    RequiresDockerUrl: "",
    RunTests:          false,
  }
}

func ReadBuildFile(path string) (*types.BuildFile, error) {
  jsonBlob, _ := ioutil.ReadFile(path + "/build.json")

  var buildFile = NewBuildFile()
  err := json.Unmarshal(jsonBlob, buildFile)
  if err != nil {
    return nil, err
  }

  return buildFile, nil
}
