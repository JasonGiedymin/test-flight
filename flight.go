package main

import (
  "./lib"
  // Logger "./lib/logging"
  "./lib/types"
  "github.com/jessevdk/go-flags"
  "os"
)

var (
  app            lib.TestFlight
  checkCommand   lib.CheckCommand2
  launchCommand  lib.LaunchCommand
  versionCommand lib.VersionCommand
  options        types.CommandOptions
)

// == App ==
func init() {
  err := app.Init()
  if err != nil {
    os.Exit(lib.ExitCodes["init_fail"])
  }

  checkCommand := lib.CheckCommand2{App: &app}
  // launchCommand := lib.LaunchCommand{App: &app, AppState: &app.AppState}
  // versionCommand := lib.VersionCommand{App: &app, AppState: &app.AppState}

  parser := flags.NewParser(&options, flags.Default)

  parser.AddCommand("check",
    "pre-flight check",
    "Used for pre-flight check of the ansible playbook.",
    &checkCommand)

  // parser.AddCommand("launch",
  //   "flight launch",
  //   "Launch an ansible playbook test.",
  //   &launchCommand)
  //
  // parser.AddCommand("version",
  //   "shows version",
  //   "Show Test-Flight version number.",
  //   &versionCommand)
}

// Runs Test-Flight
func main() {
  // app.ProcessCommands() // parse command line options now
  app.AppState.SetState("END")
}
