package main

import (
  "./lib"
  "./lib/docker"
  "os"
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

  var dc = docker.NewApi(app.AppState.ConfigFile)
  dc.ShowInfo()
  dc.ShowImages()
  dc.CreateDocker()

  app.AppState.SetState("END")
}
