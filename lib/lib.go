// ## Deps
//     go get github.com/SpaceMonkeyGo/errors
//

package lib

import (
  Logger "./logging"
  "./types"
  "bitbucket.org/kardianos/osext"
  "github.com/SpaceMonkeyGo/errors"
  // "github.com/jessevdk/go-flags"
  "os"
)

// == App Related ==

type TestFlight struct {
  AppState      types.ApplicationState
  requiredFiles []types.RequiredFile
}

func (app *TestFlight) SetState(state string) {
  app.AppState.SetState(state)
}

func (app *TestFlight) SetConfigFile(file *types.ConfigFile) {
  app.AppState.ConfigFile = file
}

func (app *TestFlight) SetBuildFile(file *types.BuildFile) {
  app.AppState.BuildFile = file
}

func (app *TestFlight) SetDir(dir string) {
  app.AppState.Meta.Dir = dir
}

func (app *TestFlight) Init() error {
  app.AppState.Meta = &meta
  app.SetState("INIT")

  execPath, error := osext.Executable()
  if error != nil {
    Logger.Error("Could not find executable path.")
    os.Exit(ExitCodes["init_fail"])
  }
  // app.AppState.Meta.ExecPath = execPath
  meta.ExecPath = execPath

  pwd, error := os.Getwd()
  if error != nil {
    Logger.Error("Could not find working directory.")
    os.Exit(ExitCodes["init_fail"])
  }
  meta.Pwd = pwd

  return nil
}

// == Default vars ==
var (
  defaultDir    = "./"
  mainYaml      = types.RequiredFile{Name: "main yaml", FileName: "main.yml", FileType: "f"}
  BadDir        = errors.NewClass("Can't read the current directory")
  FileCheckFail = errors.NewClass("File Check Failed")
  ansibleFiles  = []types.RequiredFile{mainYaml}
)

var meta = types.ApplicationMeta{
  Version: "0.9.4",
}

var RequiredFiles = []types.RequiredFile{
  {Name: "Test-Flight json build file", FileName: "build.json", FileType: "f"},
  {Name: "Test-Flight dir", FileName: ".test-flight", FileType: "d",
    RequiredFiles: []types.RequiredFile{
      {Name: "Ansible inventory file used for Test-Flight", FileName: "inventory", FileType: "f"},
      {Name: "Ansible playbook file used for Test-Flight", FileName: "playbook.yml", FileType: "f"},
    },
  },
  {Name: "Ansible handlers dir", FileName: "handlers", FileType: "d", RequiredFiles: ansibleFiles},
  {Name: "Ansible meta dir", FileName: "meta", FileType: "d", RequiredFiles: ansibleFiles},
  {Name: "Ansible tasks dir", FileName: "tasks", FileType: "d", RequiredFiles: ansibleFiles},
  {Name: "Ansible templates dir", FileName: "templates", FileType: "d"},
  {Name: "Ansible test dir", FileName: "tests", FileType: "d", RequiredFiles: ansibleFiles},
  {Name: "Ansible var dir for variables", FileName: "vars", FileType: "d", RequiredFiles: ansibleFiles},
  {Name: "Ansible vault dir for encrypted files", FileName: "vault", FileType: "d"},
}
