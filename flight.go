package main

import (
  "./lib"
  "./lib/docker"
  "os"
  Logger "./lib/logging"
)

var (
  app lib.TestFlight
)

// == App ==
func init() {
  err := app.Init()
  if err != nil {
    os.Exit(lib.ExitCodes["init_fail"])
  }
}

// Runs Test-Flight
func main() {
  app.ProcessCommands() // parse command line options now

  Logger.Trace(*app.AppState.ConfigFile)
  Logger.Trace(*app.AppState.BuildFile)

  var dc = docker.NewApi(app.AppState.ConfigFile, app.AppState.BuildFile)
  // dc.ShowInfo()
  // dc.ShowImages()
  // dc.CreateDocker()
  dc.CreateTemplate()

  app.AppState.SetState("END")
}
