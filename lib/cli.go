package lib

import (
  "./config"
  Logger "./logging"
  "./types"
  "os"
)

// TODO: Make as a member of Parser later...
func setConfigFiles(dir string, appState *types.ApplicationState) error {
  buildFile, err := config.ReadBuildFile(dir)
  if err != nil {
    Logger.Error(err)
    return err
  }

  appState.BuildFile = buildFile
  Logger.Debug("Buildfile found, contents:", *buildFile)
  return nil
}

func commandPreReq() *types.ConfigFile {
  configFile, err := config.ReadConfigFile()
  if config.ReadFileError.Contains(err) {
    os.Exit(ExitCodes["config_missing"])
  }

  return configFile
}

// == Version Command ==
type VersionCommand struct {
  App *TestFlight
}

func (cmd *VersionCommand) Execute(args []string) error {
  cmd.App.AppState.SetState("VERSION_QUERY")
  Logger.Info("Test-Flight Version:", cmd.App.AppState.Meta.Version)
  return nil
}

// == Check Command ==
type CheckCommand struct {
  App *TestFlight
  Dir string `short:"d" long:"dir" description:"directory to run in"`
}

func (cmd *CheckCommand) Execute(args []string) error {
  cmd.App.AppState.SetState("CHECK_PREREQS")
  cmd.App.AppState.ConfigFile = commandPreReq()

  cmd.App.AppState.SetState("CHECK_FILES")
  Logger.Info("Running Pre-Flight Check... in dir:", cmd.Dir)
  cmd.App.AppState.Meta.Dir = cmd.Dir

  _, err := HasRequiredFiles(&cmd.Dir, RequiredFiles)
  if err != nil {
    return err
  }

  buildFile, err := config.ReadBuildFile(cmd.Dir)
  if err != nil {
    return err
  }

  cmd.App.AppState.BuildFile = buildFile

  Logger.Debug("Buildfile found, contents:", *buildFile)
  return nil
}

// == Launch Command ==
type LaunchCommand struct {
  App *TestFlight
  Dir string `short:"d" long:"dir" description:"directory to run in"`
}

func (cmd *LaunchCommand) Execute(args []string) error {
  cmd.App.SetState("CHECK_PREREQS")
  cmd.App.AppState.ConfigFile = commandPreReq()

  cmd.App.AppState.SetState("LAUNCH")
  Logger.Info("Launching Tests... in dir:", cmd.Dir)
  cmd.App.AppState.Meta.Dir = cmd.Dir

  if _, err := HasRequiredFiles(&cmd.Dir, RequiredFiles); err != nil {
    Logger.Error(err)
    return err
  }

  if err := setConfigFiles(cmd.Dir, &cmd.App.AppState); err != nil {
    return err
  }

  var dc = NewDockerApi(cmd.App.AppState.Meta, cmd.App.AppState.ConfigFile, cmd.App.AppState.BuildFile)
  dc.ShowInfo()
  dc.ShowImages()
  // dc.CreateDocker()
  dc.createTestTemplates()
  // dc.CreateTemplate()

  return nil
}
