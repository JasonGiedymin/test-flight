package main

import (
  "./lib"
  // Logger "./lib/logging"
  "./lib/types"
  "github.com/jessevdk/go-flags"
  "os"
)

/* Limited to the app, parser, and commands */
var (
  app            lib.TestFlight
  parser         *flags.Parser
  checkCommand   lib.CheckCommand
  launchCommand  lib.LaunchCommand
  versionCommand lib.VersionCommand
  options        types.CommandOptions
)

// Func to parse the app commands
func ProcessCommands() {
  app.SetState("PARSE_COMMAND_LINE")
  if _, err := parser.Parse(); err != nil {
    os.Exit(lib.ExitCodes["command_fail"])
  }
}

// == App ==
func init() {
  err := app.Init()
  if err != nil {
    os.Exit(lib.ExitCodes["init_fail"])
  }

  flightControls := lib.FlightControls{}
  checkCommand := lib.CheckCommand{Controls: &flightControls, App: &app}
  imagesCommand := lib.ImagesCommand{Controls: &flightControls, App: &app}
  launchCommand := lib.LaunchCommand{Controls: &flightControls, App: &app}
  groundCommand := lib.GroundCommand{Controls: &flightControls, App: &app}
  versionCommand := lib.VersionCommand{Controls: &flightControls, App: &app}
  templateCommand := lib.TemplateCommand{Controls: &flightControls, App: &app}

  parser = flags.NewParser(&options, flags.Default)

  parser.AddCommand("check",
    "pre-flight check",
    "Used for pre-flight check of the ansible playbook.",
    &checkCommand)

  parser.AddCommand("images",
    "shows all docker images",
    "Used for to show all docker images",
    &imagesCommand)

  parser.AddCommand("launch",
    "flight launch",
    "Launch an ansible playbook test.",
    &launchCommand)

  parser.AddCommand("ground",
    "flight ground",
    "Grounds an ansible playbook test.",
    &groundCommand)

  parser.AddCommand("template",
    "flight template",
    "Create templates required for test-flight.",
    &templateCommand)

  parser.AddCommand("version",
    "shows version",
    "Show Test-Flight version number.",
    &versionCommand)
}

// Runs Test-Flight
func main() {
  ProcessCommands() // parse command line options now
  app.SetState("END")
}
