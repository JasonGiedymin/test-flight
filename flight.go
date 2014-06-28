package main

import (
  "./lib"
  Logger "./lib/logging"
  "github.com/jessevdk/go-flags"
  "os"
)

type CommandOptions struct {
  Configfile      string `short:"c" long:"config" description:"test-flight config file to use"`
  Build           *lib.BuildCommand
}

/* Limited to the app, parser, and commands */
var (
  app            lib.TestFlight
  parser         *flags.Parser
  checkCommand   lib.CheckCommand
  launchCommand  lib.LaunchCommand
  versionCommand lib.VersionCommand
  options        CommandOptions
)

var meta = lib.ApplicationMeta{
  Version: "0.9.5",
}

// Func to parse the app commands
func ProcessCommands() {
  app.SetState("PARSE_COMMAND_LINE")
  if _, err := parser.Parse(); err != nil {
    Logger.Error(err)
    os.Exit(lib.ExitCodes["command_fail"])
  } else {
    Logger.Info("--> Config file to use:", *options.Build)
  }
}

// == App ==
func init() {
  err := app.Init(&meta)
  if err != nil {
    os.Exit(lib.ExitCodes["init_fail"])
  }

  flightControls := lib.FlightControls{}
  // checkCommand := lib.CheckCommand{Controls: &flightControls, App: &app}
  // imagesCommand := lib.ImagesCommand{Controls: &flightControls, App: &app}
  buildCommand := lib.BuildCommand{Controls: &flightControls, App: &app}
  // launchCommand := lib.LaunchCommand{Controls: &flightControls, App: &app}
  // groundCommand := lib.GroundCommand{Controls: &flightControls, App: &app}
  // destroyCommand := lib.DestroyCommand{Controls: &flightControls, App: &app}
  // versionCommand := lib.VersionCommand{Controls: &flightControls, App: &app}
  // templateCommand := lib.TemplateCommand{Controls: &flightControls, App: &app}
  options = CommandOptions{Build: &buildCommand}

  parser = flags.NewParser(&options, flags.Default)

  // parser.AddCommand("check",
  //   "flight check",
  //   "Checks if pre-reqs are satisfied to launch this ansible playbook.",
  //   &checkCommand)

  // parser.AddCommand("images",
  //   "flight images",
  //   "Shows all images",
  //   &imagesCommand)

  // parser.AddCommand("build",
  //   "flight build",
  //   "Build will build an ansible playbook docker image.",
  //   &buildCommand)

  // parser.AddCommand("launch",
  //   "flight launch",
  //   "Launch builds, runs, and tests an ansible playbook.",
  //   &launchCommand)

  // parser.AddCommand("ground",
  //   "flight ground",
  //   "Ground stops all containers found running this ansible playbook.",
  //   &groundCommand)

  // parser.AddCommand("destroy",
  //   "flight destroy",
  //   "Destroy destroys all containers and images found running this ansible playbook.",
  //   &destroyCommand)

  // parser.AddCommand("template",
  //   "flight template",
  //   "Creates templates required for test-flight.",
  //   &templateCommand)

  // parser.AddCommand("version",
  //   "shows version",
  //   "Show Test-Flight version number.",
  //   &versionCommand)
}

// Runs Test-Flight
func main() {
  ProcessCommands() // parse command line options now
  app.SetState("END")
}
