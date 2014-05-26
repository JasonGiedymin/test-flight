package main

import (
  "./lib"
  // Logger "./lib/logging"
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

  app.AppState.SetState("END")
}
