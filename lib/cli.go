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

func commandPreReq(appState *types.ApplicationState) {
  appState.SetState("CHECK_PREREQS")
  configFile, err := config.ReadConfigFile()
  if config.ReadFileError.Contains(err) {
    os.Exit(ExitCodes["config_missing"])
  }

  appState.ConfigFile = configFile
}

// == Version Command ==
type VersionCommand struct {
  AppState *types.ApplicationState
}

func (cmd *VersionCommand) Execute(args []string) error {
  cmd.AppState.SetState("VERSION_QUERY")
  Logger.Info("Test-Flight Version:", cmd.AppState.Meta.Version)
  return nil
}

// == Check Command ==
type CheckCommand struct {
  Dir      string `short:"d" long:"dir" description:"directory to run in"`
  AppState *types.ApplicationState
}

func (cmd *CheckCommand) Execute(args []string) error {
  commandPreReq(cmd.AppState) // I'm lazy

  cmd.AppState.SetState("CHECK_FILES")
  Logger.Info("Running Pre-Flight Check... in dir:", cmd.Dir)
  cmd.AppState.Meta.Dir = cmd.Dir

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
  Dir      string `short:"d" long:"dir" description:"directory to run in"`
  AppState *types.ApplicationState
}

func (cmd *LaunchCommand) Execute(args []string) error {
  commandPreReq(cmd.AppState)

  cmd.AppState.SetState("LAUNCH")
  Logger.Info("Launching Tests... in dir:", cmd.Dir)
  cmd.AppState.Meta.Dir = cmd.Dir

  if _, err := HasRequiredFiles(&cmd.Dir, RequiredFiles); err != nil {
    Logger.Error(err)
    return err
  }

  if err := setConfigFiles(cmd.Dir, cmd.AppState); err != nil {
    return err
  }

  var dc = NewDockerApi(cmd.AppState.Meta, cmd.AppState.ConfigFile, cmd.AppState.BuildFile)
  dc.ShowInfo()
  dc.ShowImages()
  // dc.CreateDocker()
  dc.createTestTemplates()
  // dc.CreateTemplate()

  return nil
}
