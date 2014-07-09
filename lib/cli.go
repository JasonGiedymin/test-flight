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
  return cmd.Controls.testFlightTemplates(dc, cmd.App.AppState.ConfigFile, *cmd.Options)
}
