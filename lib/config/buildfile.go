package config

import (
  "../types"
  "encoding/json"
  "io/ioutil"
  "strings"
)

// For specific defaults
func NewBuildFile() *types.BuildFile {
  return &types.BuildFile{
    Owner:             "Test-Flight-User",
    ImageName:         "Test-Flight-Test-Image",
    From:              "",
    Version:           "0.0.1",
    RequiresDocker:    "",
    RequiresDockerUrl: "",
    RunTests:          false,
  }
}

func ReadBuildFile(dir string) (*types.BuildFile, error) {
  path := []string{dir, "build.json"}
  jsonBlob, _ := ioutil.ReadFile(strings.Join(path, "/"))

  var buildFile = NewBuildFile()
  err := json.Unmarshal(jsonBlob, buildFile)
  if err != nil {
    return nil, err
  }

  return buildFile, nil
}
