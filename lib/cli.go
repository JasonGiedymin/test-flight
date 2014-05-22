package lib

import (
  "./config"
  Logger "./logging"
  "os"
)

func commandPreReq(state func(string) string) {
  state("CHECK_PREREQS")
  _, err := config.ReadConfigFile()
  if config.ReadFileError.Contains(err) {
    os.Exit(ExitCodes["config_missing"])
  }
}

// == Version Command ==
type VersionCommand struct {
  AppState *ApplicationState
}

func (cmd *VersionCommand) Execute(args []string) error {
  cmd.AppState.SetState("VERSION_QUERY")
  Logger.Info("Test-Flight Version:", cmd.AppState.Meta.Version)
  return nil
}

// == Check Command ==
type CheckCommand struct {
  Dir      string `short:"d" long:"dir" description:"directory to run in"`
  AppState *ApplicationState
}

func (cmd *CheckCommand) Execute(args []string) error {
  commandPreReq(cmd.AppState.SetState) // I'm lazy

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
  AppState *ApplicationState
}

func (cmd *LaunchCommand) Execute(args []string) error {
  commandPreReq(cmd.AppState.SetState)

  Logger.Info("Launching Tests... in dir:", cmd.Dir)
  _, err := HasRequiredFiles(&cmd.Dir, RequiredFiles)

  if err != nil {
    return err
  }
  return nil
}
