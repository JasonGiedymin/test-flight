package config

import (
  "encoding/json"
  "fmt"
  "github.com/SpaceMonkeyGo/errors"
  "io/ioutil"
  "os"
  "os/user"
)

var (
  ReadFileError = errors.NewClass("Could not read file.")
)

type ConfigFile struct {
  DockerEndpoint string
}

func ReadConfigFile() (*ConfigFile, error) {
  var configFile ConfigFile

  fmt.Println("Looking for test-flight-config.json...")

  usr, err := user.Current()
  if err != nil {
    return nil, ReadFileError.New("Can't read test-flight-config.json file in user home.")
  }

  jsonBlob, _ := ioutil.ReadFile(usr.HomeDir + "/test-flight-config.json")
  err = json.Unmarshal(jsonBlob, &configFile)

  if err != nil {
    //TODO: log noting prog is trying local now
    fmt.Println("Checking for local pwd config file...")

    // Get user home first
    pwd, err := os.Getwd()
    if err != nil {
      return nil, err
    }

    // with user home find config file
    jsonBlob, err = ioutil.ReadFile(pwd + "/test-flight-config.json")
    err = json.Unmarshal(jsonBlob, &configFile)
    if err != nil {
      return nil, ReadFileError.New("Can't find test-flight-config.json file in local pwd.")
    }
  }

  fmt.Printf("Got config file: %v\n", configFile)
  return &configFile, nil
}
