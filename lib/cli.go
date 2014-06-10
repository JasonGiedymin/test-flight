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
    Logger.Error("Error reading build file:", err)
    return err
  }

  appState.BuildFile = buildFile
  Logger.Debug("Buildfile found, contents:", *buildFile)
  return nil
}

// func commandPreReq() *types.ConfigFile {
//   configFile, err := config.ReadConfigFile()
//   if config.ReadFileError.Contains(err) {
//     os.Exit(ExitCodes["config_missing"])
//   }
//
//   return configFile
// }

func commandPreReq(app *TestFlight) {
  // Prereqs
  app.SetState("CHECK_PREREQS")
  configFile, err := config.ReadConfigFile()
  if config.ReadFileError.Contains(err) {
    os.Exit(ExitCodes["config_missing"])
  }
  app.SetConfigFile(configFile)
}

// == Version Command ==
type VersionCommand struct {
  App *TestFlight
}

func (cmd *VersionCommand) Execute(args []string) error {
  cmd.App.SetState("VERSION_QUERY")
  Logger.Info("Test-Flight Version:", cmd.App.AppState.Meta.Version)
  return nil
}

// == Check Command ==
type CheckCommand struct {
  App *TestFlight
  Dir string `short:"d" long:"dir" description:"directory to run in"`
}

func (cmd *CheckCommand) Execute(args []string) error {
  commandPreReq(cmd.App)

  cmd.App.SetState("CHECK_FILES")
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

func watchForEventsOn(channel ApiChannel) {
  for msg := range channel {
    Logger.Info("DOCKER EVENT:", *msg)
  }
}

func (cmd *LaunchCommand) Execute(args []string) error {
  commandPreReq(cmd.App)

  cmd.App.SetState("LAUNCH")
  Logger.Info("Launching Tests... in dir:", cmd.Dir)
  cmd.App.AppState.Meta.Dir = cmd.Dir

  if _, err := HasRequiredFiles(&cmd.Dir, RequiredFiles); err != nil {
    return err
  }

  if err := setConfigFiles(cmd.Dir, &cmd.App.AppState); err != nil {
    return err
  }

  var dc = NewDockerApi(cmd.App.AppState.Meta, cmd.App.AppState.ConfigFile, cmd.App.AppState.BuildFile)
  dc.ShowInfo()
  // dc.ShowImages()

  if err := testFlightTemplates(dc, cmd.App.AppState.ConfigFile); err != nil {
    return err
  }

  // Register channel so we can watch for events as they happen
  eventsChannel := make(ApiChannel)
  go watchForEventsOn(eventsChannel)
  dc.RegisterChannel(eventsChannel)

  dc.CreateDockerImage()
  // dc.CreateContainer()
  // dc.ShowImages()
  // dc.DeleteImage()

  return nil
}

// == Template Command ==
type TemplateCommand struct {
  App *TestFlight
  Dir string `short:"d" long:"dir" description:"directory to run in"`
}

func testFlightTemplates(dc *DockerApi, configFile *types.ConfigFile) error {
  if configFile.OverwriteTemplates {
    return dc.createTestTemplates()
  }

  return nil
}

func (cmd *TemplateCommand) Execute(args []string) error {
  commandPreReq(cmd.App)

  cmd.App.SetState("TEMPLATE")
  Logger.Info("Creating Templates... in dir:", cmd.Dir)
  cmd.App.AppState.Meta.Dir = cmd.Dir

  dc := NewDockerApi(cmd.App.AppState.Meta, cmd.App.AppState.ConfigFile, cmd.App.AppState.BuildFile)
  return testFlightTemplates(dc, cmd.App.AppState.ConfigFile)
}
