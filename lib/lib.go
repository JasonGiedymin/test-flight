// ## Deps
//     go get github.com/SpaceMonkeyGo/errors
//

package lib

import (
  "./config"
  "fmt"
  "github.com/SpaceMonkeyGo/errors"
  "github.com/jessevdk/go-flags"
  "os"
)

// == App Related ==

type RequiredFile struct {
  name          string
  fileName      string
  fileType      string // [f]ile, [d]ir
  requiredFiles []RequiredFile
}

type ApplicationConfig struct {
  defaultDir string
}

type ApplicationState struct {
  ConfigFile  config.ConfigFile
  BuildFile   *config.BuildFile
  LastCommand string
  CurrentMode string
}

func (appState *ApplicationState) State() string {
  fmt.Printf("[STATE] - %v\n", appState.CurrentMode)
  return appState.CurrentMode
}

type TestFlight struct {
  AppState      ApplicationState
  requiredFiles []RequiredFile
  parser        *flags.Parser
}

func (app *TestFlight) Parse() {
  if _, err := app.parser.Parse(); err != nil {
    os.Exit(1)
  }
}

func (app *TestFlight) State() string {
  fmt.Printf("[STATE] - %v", app.AppState.CurrentMode)
  return app.AppState.CurrentMode
}

// == Default vars ==
var (
  defaultDir    = "./"
  mainYaml      = RequiredFile{name: "main yaml", fileName: "main.yml", fileType: "f"}
  BadDir        = errors.NewClass("Can't read the current directory")
  FileCheckFail = errors.NewClass("File Check Failed")
  ansibleFiles  = []RequiredFile{mainYaml}
)

var RequiredFiles = []RequiredFile{
  {name: "Test-Flight json build file", fileName: "build.json", fileType: "f"},
  {name: "Test-Flight docker dir", fileName: "docker", fileType: "d",
    requiredFiles: []RequiredFile{
      {name: "Ansible inventory file used for Test-Flight", fileName: "inventory", fileType: "f"},
      {name: "Ansible playbook file used for Test-Flight", fileName: "playbook.yml", fileType: "f"},
    },
  },
  {name: "Ansible handlers dir", fileName: "handlers", fileType: "d", requiredFiles: ansibleFiles},
  {name: "Ansible meta dir", fileName: "meta", fileType: "d", requiredFiles: ansibleFiles},
  {name: "Ansible tasks dir", fileName: "tasks", fileType: "d", requiredFiles: ansibleFiles},
  {name: "Ansible templates dir", fileName: "templates", fileType: "d"},
  {name: "Ansible test dir", fileName: "tests", fileType: "d", requiredFiles: ansibleFiles},
  {name: "Ansible var dir for variables", fileName: "vars", fileType: "d", requiredFiles: ansibleFiles},
  {name: "Ansible vault dir for encrypted files", fileName: "vault", fileType: "d"},
}
