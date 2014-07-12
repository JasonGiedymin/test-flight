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
    pwd, err := os.Getwd() // use working dir
    if err != nil {
        pwd = "~" // else use home dir by default
    }

    return &ConfigFile{ // optional values:
        DockerEndpoint:           "http://localhost:4243",
        AnsibleTemplatesDir:      ".test-flight",
        TestFlightAssets:         FilePath(pwd, ".test-flight"),
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
func findConfig() (*ConfigFile, error) {
    configFileName := "test-flight-config.json"
    configFile := NewConfigFile()

    // get home
    usr, err := user.Current()
    if err != nil {
        Logger.Error("Can't read user home.")
        // return nil, ReadFileError.New("Can't read user home.")
        return nil, errors.New("Can't read user home.")
    }

    homeConfigPath := FilePath(usr.HomeDir, ".test-flight", "test-flight-config.json")

    Logger.Debug("Checking for config file in user HOME: ", homeConfigPath)

    // try home first
    configFile, err = getConfig(homeConfigPath)
    if err != nil {
        Logger.Warn(configFileName + " not found in user HOME: " + usr.HomeDir)

        pwd, err := os.Getwd()
        if err != nil {
            return nil, err
        }

        // try running directory next
        Logger.Debug("Checking for config file in local pwd: " + pwd + "/" + configFileName)
        pwdConfigPath := FilePath(pwd, ".test-flight", configFileName)
        configFile, err = getConfig(pwdConfigPath)
        if err != nil {
            return nil, errors.New("Could not find configfile in user home or local running" +
                " directory. Please supply the config file.")
        }
    }

    Logger.Debug("Found config file, contents:", configFile)
    return configFile, nil
}

// can be called with default empty param which means user did not specify
// config file.
func ReadConfigFile(userSpecified string) (*ConfigFile, error) {
    if userSpecified != "" {
        return getConfig(userSpecified)
    }

    return findConfig()
}
