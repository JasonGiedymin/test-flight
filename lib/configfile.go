package lib

import (
    "encoding/json"
    Logger "github.com/JasonGiedymin/test-flight/lib/logging"
    // "github.com/SpaceMonkeyGo/errors"
    "errors"
    "io/ioutil"
    "os"
    "os/user"
    "strings"
)

var (
// ReadFileError = errors.NewClass("Could not read file.")
)

type ConfigFile struct {
    Location                 string // fq path of config file
    LocationDir              string // fq path where config file resides
    AnsibleTemplatesDir      string // playbook/.test-flight
    TestFlightAssets         string // $HOME/.test-flight
    UseSystemDockerTemplates bool   // to use DockerTemplatesDir/{system|user}
    DockerEndpoint           string
    WorkDir                  string
    DockerAdd                ConfigFileDockerAdd
    OverwriteTemplates       bool
}

type ConfigFileDockerAdd struct {
    Simple []string
    // User   []map[string]string
    Complex []DockerAddComplexEntry
}

// Used for defaults
func NewConfigFile() *ConfigFile {
    usr, _ := user.Current() // to get user home, get user first
    pwd, err := os.Getwd()   // use working dir
    if err != nil {
        pwd = usr.HomeDir // else use home dir by default
    }

    return &ConfigFile{ // optional values:
        DockerEndpoint:           "http://localhost:4243",
        AnsibleTemplatesDir:      FilePath(pwd, ".test-flight"),
        TestFlightAssets:         FilePath(usr.HomeDir, ".test-flight"),
        UseSystemDockerTemplates: true,
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
    Logger.Debug("Checking for config file in location:", file)

    var unmarshal = func(jsonBlob []byte) (*ConfigFile, error) {
        configFile := NewConfigFile()

        if err := json.Unmarshal(jsonBlob, &configFile); err != nil {
            msg := "Can't find or having trouble reading " + file +
                ". Please create the file or address syntax issues. System error:" + err.Error()
            // return nil, ReadFileError.New(msg)
            return nil, errors.New(msg)
        } else {
            // Augment configfile with it's location
            configFile.Location = file
            configFile.LocationDir = file[:strings.LastIndex(file, "/")]

            return configFile, nil
        }
    }

    if jsonBlob, err := ioutil.ReadFile(file); err != nil {
        return nil, errors.New("Can't find " + file)
    } else {
        if jsonBlob == nil {
            msg := file + " not found."
            return nil, errors.New(msg)
        } else {
            return unmarshal(jsonBlob)
        }
    }
}

// tries to find config file in user home, then if it cannot find one there
// will try to find a config file in the local running directory
func findConfig(dir string) (*ConfigFile, error) {
    configFileName := "test-flight-config.json"
    configFile := NewConfigFile()
    logConfigFile := func(configFile *ConfigFile) {
        Logger.Debug("Found config file.")
        Logger.Trace("Config file contents:", *configFile)
    }

    // try dir specified
    localConfigPath := FilePath(dir, configFileName)
    configFile, err := getConfig(localConfigPath)
    if err != nil {
        msg := "Config: " + localConfigPath + " may not exist or cannot be read. " + err.Error()
        Logger.Warn(msg)
    } else {
        logConfigFile(configFile)
        return configFile, nil
    }

    // try running directory next
    pwd, err := os.Getwd()
    if err != nil {
        return nil, err
    }
    Logger.Debug("Checking for config file in local pwd: " + pwd + "/" + configFileName)
    pwdConfigPath := FilePath(pwd, ".test-flight", configFileName)
    configFile, err = getConfig(pwdConfigPath)
    if err != nil {
        msg := "Config: " + localConfigPath + " may not exist or cannot be read. " + err.Error()
        Logger.Warn(msg)
    } else {
        logConfigFile(configFile)
        return configFile, nil
    }

    // try home
    usr, err := user.Current() // to get user home, get user first
    if err != nil {
        Logger.Error("Can't read user home.")
        // return nil, ReadFileError.New("Can't read user home.")
        return nil, errors.New("Can't read user home.")
    }

    homeConfigPath := FilePath(usr.HomeDir, ".test-flight", "test-flight-config.json")
    Logger.Debug("Checking for config file in user HOME: ", homeConfigPath)

    configFile, err = getConfig(homeConfigPath)
    if err != nil {
        msg := "Config: " + localConfigPath + " may not exist or cannot be read. " + err.Error()
        Logger.Warn(msg)
    } else {
        logConfigFile(configFile)
        return configFile, nil
    }

    return nil, errors.New("Cannot find config file, Please supply the config file.")
}

// can be called with default empty param which means user did not specify
// config file.
func ReadConfigFile(userSpecified string, dir string) (*ConfigFile, error) {
    if userSpecified != "" {
        return getConfig(userSpecified)
    }

    return findConfig(dir)
}
