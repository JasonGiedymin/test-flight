package config

import (
  "encoding/json"
  "io/ioutil"
)

type BuildFile struct {
  Owner string
  ImageName string
  Version string
  RequiresDocker string
  RequiresDockerUrl string
}

func ReadBuildFile(path string) (*BuildFile, error) {
  jsonBlob, _ := ioutil.ReadFile(path + "/build.json")

  var buildFile BuildFile
  err := json.Unmarshal(jsonBlob, &buildFile)
  if err != nil {
    return nil, err
  }

  return &buildFile, nil
}
