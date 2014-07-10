package lib

import (
  Logger "github.com/JasonGiedymin/test-flight/lib/logging"
)


// == Version Command ==
type VersionCommand struct {
  Controls *FlightControls
  Options  *CommandOptions
  App      *TestFlight
}

func (cmd *VersionCommand) Execute(args []string) error {
  Logger.Info("Test-Flight Version:", cmd.App.AppState.Meta.Version)

  return nil
}