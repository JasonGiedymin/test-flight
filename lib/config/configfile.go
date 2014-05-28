package config

import (
  Logger "../logging"
  "encoding/json"
  "github.com/SpaceMonkeyGo/errors"
  "io/ioutil"
  "os"
  "os/user"
)

var (
  ReadFileError = errors.NewClass("Could not read file.")
)

type DockerAddComplexEntry struct {
  Name     string
  Location string
}

type ConfigFileDockerAdd struct {
  Simple    []string
  // User   []map[string]string
  Complex   []DockerAddComplexEntry
}

type ConfigFile struct {
  DockerEndpoint string
  WorkDir string
  DockerAdd ConfigFileDockerAdd
}

func NewConfigFile() *ConfigFile {
  return &ConfigFile{
    DockerEndpoint: "http://localhost:4243",
  }
}

// type ConfigProcessor func(inStr string) (string, error)
//
// // If err, next; otherwise return err
// func ConfigCompose(fx []ConfigProcessor, path string) (*ConfigFile, error) {
//   var configFile *ConfigFile
//   var err error
//
//   for i := range fx {
//     configFile, err := fx[i](path)
//     if configFile != nil && err != nil {
//       return configFile, nil
//     }
//   }
//
//   return nil, err
// }

func ReadConfigFile() (*ConfigFile, error) {
  configFileName := "test-flight-config.json"
  configFile := NewConfigFile()

  usr, err := user.Current()
  if err != nil {
    Logger.Error("Can't read user home.")
    return nil, ReadFileError.New("Can't read user home.")
  }

  Logger.Debug("Checking for config file in user HOME: " + usr.HomeDir + "/test-flight-config.json")
  jsonBlob, _ := ioutil.ReadFile(usr.HomeDir + "/test-flight-config.json")
  err = json.Unmarshal(jsonBlob, &configFile)

  if err != nil {
    Logger.Warn(configFileName + " not found in user HOME: " + usr.HomeDir)

    pwd, err := os.Getwd()
    if err != nil {
      return nil, err
    }

    // with user home find config file
    Logger.Debug("Checking for config file in local pwd: " + pwd + "/" + configFileName)
    jsonBlob, err = ioutil.ReadFile(pwd + "/" + configFileName)
    err = json.Unmarshal(jsonBlob, &configFile)
    if err != nil {
      Logger.Error("Can't find or having trouble reading " + configFileName +
        " in local pwd or user home. Please create the file or address syntax issues.")
      Logger.Error(err)
      return nil, ReadFileError.New("Can't find test-flight-config.json file in local pwd.")
    }
  }

  Logger.Debug("Found config file, contents:", configFile)
  return configFile, nil
}
