package config

import (
  "encoding/json"
  "fmt"
  "io/ioutil"
)

type ConfigFile struct {
  DockerEndpoint    string
}

func ReadConfigFile(path string) {
  jsonBlob, _ := ioutil.ReadFile(path + "/config.json")

  var configFile ConfigFile
  err := json.Unmarshal(jsonBlob, &configFile)
  if err != nil {
    fmt.Println("error:", err)
  }
  fmt.Printf("%+v\n", configFile)
}
