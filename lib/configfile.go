package lib

import (
  Logger "./logging"
  "encoding/json"
  "github.com/SpaceMonkeyGo/errors"
  "io/ioutil"
  "os"
  "os/user"
)

var (
  ReadFileError = errors.NewClass("Could not read file.")
)

type ConfigFileDockerAdd struct {
  Simple []string
  // User   []map[string]string
  Complex []DockerAddComplexEntry
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

func getConfig(file string) (*ConfigFile, error) {
  configFile := NewConfigFile()

  Logger.Debug("Checking for config file in location:", file)
  jsonBlob, err := ioutil.ReadFile(file)
  err = json.Unmarshal(jsonBlob, &configFile)
  if err != nil {
    msg := "Can't find or having trouble reading " + file +
      ". Please create the file or address syntax issues. System error:" + err.Error()
    return nil, ReadFileError.New(msg)
  }

  return configFile, nil
}

// tries to find config file in user home, then if it cannot find one there
// will try to find a config file in the local running directory
func findConfig() (*ConfigFile, error) {
  configFileName := "test-flight-config.json"
  configFile := NewConfigFile()

  // get home
  usr, err := user.Current()
  if err != nil {
    Logger.Error("Can't read user home.")
    return nil, ReadFileError.New("Can't read user home.")
  }  
  
  Logger.Debug("Checking for config file in user HOME: " + usr.HomeDir + "/test-flight-config.json")

  // try home first
  configFile, err = getConfig(usr.HomeDir + "/test-flight-config.json")
  if err != nil {
    Logger.Warn(configFileName + " not found in user HOME: " + usr.HomeDir)

    pwd, err := os.Getwd()
    if err != nil {
      return nil, err
    }

    // try running directory next
    Logger.Debug("Checking for config file in local pwd: " + pwd + "/" + configFileName)
    configFile, err = getConfig(pwd + "/" + configFileName)
    if err != nil {
      return nil, errors.New("Could not find configfile in user home or local running" +
        " directory. Please supply the config file")
    }
  }

  Logger.Debug("Found config file, contents:", configFile)
  return configFile, nil
}

// can be called with default empty param which means user did not specify
// config file.
func ReadConfigFile(userSpecified string) (*ConfigFile, error) {
  if (userSpecified != "") {
    return getConfig(userSpecified)
  }
  
  return findConfig()  
}
