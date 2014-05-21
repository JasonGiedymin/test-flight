package lib

import (
  "./config"
  "fmt"
  Logger "./logging"
)

// == Check Command ==
type CheckCommand struct {
  Dir      string `short:"d" long:"dir" description:"directory to run in"`
  AppState *ApplicationState
}

func (cmd *CheckCommand) Execute(args []string) error {
  cmd.AppState.SetState("CHECK_FILES")
  Logger.Info("Running Pre-Flight Check... in dir:", cmd.Dir)

  _, err := HasRequiredFiles(&cmd.Dir, RequiredFiles)
  if err != nil {
    return err
  }

  buildFile, err := config.ReadBuildFile(cmd.Dir)
  if err != nil {
    return err
  }

  cmd.AppState.BuildFile = buildFile

  Logger.Debug("Buildfile found, contents:", *buildFile)
  return nil
}

// == Launch Command ==
type LaunchCommand struct {
  Dir string `short:"d" long:"dir" description:"directory to run in"`
}

func (cmd *LaunchCommand) Execute(args []string) error {
  fmt.Printf("Launching Tests... in dir: [%v]\n", cmd.Dir)
  _, err := HasRequiredFiles(&cmd.Dir, RequiredFiles)

  if err != nil {
    return err
  }

  fmt.Println("Done.")
  return nil
}
