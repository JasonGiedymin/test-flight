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

type ConfigFile struct {
  DockerEndpoint string
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
  var configFile ConfigFile

  usr, err := user.Current()
  if err != nil {
    Logger.Error("Can't read user home.")
    return nil, ReadFileError.New("Can't read user home.")
  }

  Logger.Debug("Looking for config file in user HOME: " + usr.HomeDir + "/test-flight-config.json")
  jsonBlob, _ := ioutil.ReadFile(usr.HomeDir + "/test-flight-config.json")
  err = json.Unmarshal(jsonBlob, &configFile)

  if err != nil {
    Logger.Warn(configFileName + " not found in user HOME: " + usr.HomeDir)

    // Get user home first
    pwd, err := os.Getwd()
    if err != nil {
      return nil, err
    }

    // with user home find config file
    Logger.Debug("Now checking for config file in local running pwd: " + pwd + "/" + configFileName)
    jsonBlob, err = ioutil.ReadFile(pwd + "/" + configFileName)
    err = json.Unmarshal(jsonBlob, &configFile)
    if err != nil {
      Logger.Error("Can't find " + configFileName + "in local pwd or user home. Please create the file.")
      return nil, ReadFileError.New("Can't find test-flight-config.json file in local pwd.")
    }
  }

  Logger.Debug("Found config file, contents:", configFile)
  return &configFile, nil
}
