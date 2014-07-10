package lib

import (
  Logger "github.com/JasonGiedymin/test-flight/lib/logging"
)

// == Template Command ==
type TemplateCommand struct {
  Controls *FlightControls
  Options  *CommandOptions
  App      *TestFlight
}

func (cmd *TemplateCommand) Execute(args []string) error {
  Logger.Info("Creating Templates... in dir:", cmd.Options.Dir)

  _, _, err := cmd.Controls.CheckConfigs(cmd.App, cmd.Options)
  if err != nil {
    return err
  }

  cmd.App.AppState.Meta.Dir = cmd.Options.Dir

  dc := NewDockerApi(cmd.App.AppState.Meta, cmd.App.AppState.ConfigFile, cmd.App.AppState.BuildFile)
  return dc.createTestTemplates(*cmd.Options)
}