// ## Deps
//     go get github.com/SpaceMonkeyGo/errors
//

package lib

import (
  "./types"
  "github.com/SpaceMonkeyGo/errors"
  "github.com/jessevdk/go-flags"
  "os"
  "bitbucket.org/kardianos/osext"
  Logger "./logging"
)

// == App Related ==

type RequiredFile struct {
  name          string
  fileName      string
  FileType      string // [f]ile, [d]ir
  requiredFiles []RequiredFile
}

type TestFlight struct {
  AppState      types.ApplicationState
  requiredFiles []RequiredFile
  Parser        *flags.Parser
}

func (app *TestFlight) ProcessCommands() {
  app.AppState.SetState("PARSE_COMMAND_LINE")
  if _, err := app.Parser.Parse(); err != nil {
    os.Exit(ExitCodes["command_fail"])
  }
}

func (app *TestFlight) Init() error {
  app.AppState.Meta = &meta

  Logger.Trace("State", app.AppState)
  app.AppState.SetState("INIT")

  execPath, error := osext.Executable()
  if (error != nil) {
    Logger.Error("Could not find executable path.")
    os.Exit(ExitCodes["init_fail"])
  }
  app.AppState.Meta.ExecPath = execPath

  pwd, error := os.Getwd()
  if (error != nil) {
    Logger.Error("Could not find working directory.")
    os.Exit(ExitCodes["init_fail"])
  }
  app.AppState.Meta.Pwd = pwd

  checkCommand = CheckCommand{AppState: &app.AppState}
  launchCommand = LaunchCommand{AppState: &app.AppState}
  versionCommand = VersionCommand{AppState: &app.AppState}

  app.Parser = flags.NewParser(&app.AppState.Options, flags.Default)

  app.Parser.AddCommand("check",
    "pre-flight check",
    "Used for pre-flight check of the ansible playbook.",
    &checkCommand)

  app.Parser.AddCommand("launch",
    "flight launch",
    "Launch an ansible playbook test.",
    &launchCommand)

  app.Parser.AddCommand("version",
    "shows version",
    "Show Test-Flight version number.",
    &versionCommand)

  return nil
}

// == Default vars ==
var (
  defaultDir     = "./"
  mainYaml       = RequiredFile{name: "main yaml", fileName: "main.yml", FileType: "f"}
  BadDir         = errors.NewClass("Can't read the current directory")
  FileCheckFail  = errors.NewClass("File Check Failed")
  ansibleFiles   = []RequiredFile{mainYaml}
  checkCommand   CheckCommand
  launchCommand  LaunchCommand
  versionCommand VersionCommand
)

var meta = types.ApplicationMeta{
  Version: "0.9.2",
}

var RequiredFiles = []RequiredFile{
  {name: "Test-Flight json build file", fileName: "build.json", FileType: "f"},
  {name: "Test-Flight dir", fileName: ".test-flight", FileType: "d",
    requiredFiles: []RequiredFile{
      {name: "Ansible inventory file used for Test-Flight", fileName: "inventory", FileType: "f"},
      {name: "Ansible playbook file used for Test-Flight", fileName: "playbook.yml", FileType: "f"},
    },
  },
  {name: "Ansible handlers dir", fileName: "handlers", FileType: "d", requiredFiles: ansibleFiles},
  {name: "Ansible meta dir", fileName: "meta", FileType: "d", requiredFiles: ansibleFiles},
  {name: "Ansible tasks dir", fileName: "tasks", FileType: "d", requiredFiles: ansibleFiles},
  {name: "Ansible templates dir", fileName: "templates", FileType: "d"},
  {name: "Ansible test dir", fileName: "tests", FileType: "d", requiredFiles: ansibleFiles},
  {name: "Ansible var dir for variables", fileName: "vars", FileType: "d", requiredFiles: ansibleFiles},
  {name: "Ansible vault dir for encrypted files", fileName: "vault", FileType: "d"},
}
