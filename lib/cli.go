package lib

import (
  Logger "github.com/JasonGiedymin/test-flight/lib/logging"
  "sync"
  "runtime"
)

func getRequiredFiles(filemode bool) []RequiredFile {
  if filemode {
    return AnsibleFiles
  } else {
    return RequiredFiles
  }
}

func (cmd *VersionCommand) Execute(args []string) error {
  cmd.App.SetState("VERSION_QUERY")
  Logger.Info("Test-Flight Version:", cmd.App.AppState.Meta.Version)
  return nil
}

func (cmd *GroundCommand) Execute(args []string) error {
  // Check Config and Buildfiles
  configFile, buildFile, err := cmd.Controls.CheckConfigs(cmd.App, cmd.Options)
  if err != nil {
    return err
  }
  
  Logger.Info("Grounding Tests... in dir:", cmd.Dir)

  dc := NewDockerApi(cmd.App.AppState.Meta, configFile, buildFile)
  dc.ShowInfo()

  if err := cmd.Controls.testFlightTemplates(dc, configFile, cmd.SingleFileMode); err != nil {
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

func (cmd *DestroyCommand) Execute(args []string) error {
  Logger.Info("Destroying... using information from dir:", cmd.Dir)

  // Check Config and Buildfiles
  configFile, buildFile, err := cmd.Controls.CheckConfigs(cmd.App, cmd.Options)
  if err != nil {
    return err
  }

  dc := NewDockerApi(cmd.App.AppState.Meta, configFile, buildFile)
  dc.ShowInfo()

  // Register channel so we can watch for events as they happen
  eventsChannel := make(ApiChannel)
  go watchForEventsOn(eventsChannel)
  dc.RegisterChannel(eventsChannel)

  fqImageName := cmd.App.AppState.BuildFile.ImageName + ":" + cmd.App.AppState.BuildFile.Tag

  if _, err := dc.DeleteContainer(fqImageName); err != nil {
    Logger.Error("Could not delete container,", err)
    return err
  }

  if _, err := dc.DeleteImage(fqImageName); err != nil {
    Logger.Error("Could not delete image,", err)
    return err
  }

  // Nothing to do
  return nil
}

func watchForEventsOn(channel ApiChannel) {
  for msg := range channel {
    Logger.Trace("DOCKER EVENT:", *msg)
  }
}

func watchContainerOn(channel ContainerChannel, wg *sync.WaitGroup) {
  for msg := range channel {
    runtime.Gosched()
    Logger.Console(msg)
  }
  
  wg.Done()
}

func (cmd *TemplateCommand) Execute(args []string) error {
  cmd.App.SetState("TEMPLATE")
  Logger.Info("Creating Templates... in dir:", cmd.Options.Dir)

  _, _, err := cmd.Controls.CheckConfigs(cmd.App, cmd.Options)
  if err != nil {
    return err
  }

  cmd.App.AppState.Meta.Dir = cmd.Options.Dir

  dc := NewDockerApi(cmd.App.AppState.Meta, cmd.App.AppState.ConfigFile, cmd.App.AppState.BuildFile)
  return cmd.Controls.testFlightTemplates(dc, cmd.App.AppState.ConfigFile, cmd.Options.SingleFileMode)
}
