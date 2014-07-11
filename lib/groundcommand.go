package lib

import (
  Logger "github.com/JasonGiedymin/test-flight/lib/logging"
)

// == Ground Command ==
// Should stop running containers
type GroundCommand struct {
  Controls *FlightControls
  Options  *CommandOptions
  App      *TestFlight
}

func (cmd *GroundCommand) Execute(args []string) error {
  Logger.Info("Grounding Tests... in dir:", cmd.Options.Dir)

  // Check Config and Buildfiles
  configFile, buildFile, err := cmd.Controls.CheckConfigs(cmd.App, cmd.Options)
  if err != nil {
    return err
  }

  dc := NewDockerApi(cmd.App.AppState.Meta, configFile, buildFile)
  dc.ShowInfo()

  if err := dc.createTestTemplates(*cmd.Options); err != nil {
    return err
  }

  // Register channel so we can watch for events as they happen
  eventsChannel := make(ApiChannel)
  go watchForEventsOn(eventsChannel)
  dc.RegisterChannel(eventsChannel)

  fqImageName := cmd.App.AppState.BuildFile.ImageName + ":" + cmd.App.AppState.BuildFile.Tag
  if running, err := dc.ListContainers(fqImageName); err != nil {
    Logger.Trace("Error while trying to get a list of containers for ", fqImageName)
    return err
  } else {
    for _, container := range running {
      dc.StopContainer(container)
    }
  }

  return nil
}
