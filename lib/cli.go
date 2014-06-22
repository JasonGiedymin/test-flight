package lib

import (
  "./config"
  Logger "./logging"
  "./types"
  "os"
  // "time"
  // "fmt"
)

type FlightControls struct {}

func (fc *FlightControls) Init(app *TestFlight) {}

func (fc *FlightControls) CheckConfig() (*types.ConfigFile, error) {
  // Prereqs
  configFile, err := config.ReadConfigFile()
  if config.ReadFileError.Contains(err) {
    os.Exit(ExitCodes["config_missing"])
  }

  return configFile, nil
}

func (fc *FlightControls) CheckBuild(dir string, requiredFiles []types.RequiredFile) (*types.BuildFile, error) {
  if _, err := HasRequiredFiles(dir, requiredFiles); err != nil {
    return nil, err
  }

  if buildFile, err := getBuildFile(dir); err != nil {
    return nil, err
  } else {
    return buildFile, nil
  }
}

// TODO: Make as a member of Parser later...
func getBuildFile(dir string) (*types.BuildFile, error) {
  buildFile, err := config.ReadBuildFile(dir)
  if err != nil {
    Logger.Error("Error reading build file:", err)
    return nil, err
  }

  Logger.Debug("Buildfile found, contents:", *buildFile)
  return buildFile, nil
}

func commandPreReq(app *TestFlight) {
  // Prereqs
  configFile, err := config.ReadConfigFile()
  if config.ReadFileError.Contains(err) {
    os.Exit(ExitCodes["config_missing"])
  }
  app.SetConfigFile(configFile)
}

// == Version Command ==
type VersionCommand struct {
  Controls *FlightControls
  App *TestFlight
}

func (cmd *VersionCommand) Execute(args []string) error {
  cmd.App.SetState("VERSION_QUERY")
  Logger.Info("Test-Flight Version:", cmd.App.AppState.Meta.Version)
  return nil
}

// == Check Command ==
type CheckCommand struct {
  Controls *FlightControls
  App *TestFlight
  Dir string `short:"d" long:"dir" description:"directory to run in"`
}

func (cmd *CheckCommand) Execute(args []string) error {
  commandPreReq(cmd.App)

  Logger.Info("Running Pre-Flight Check... in dir:", cmd.Dir)
  cmd.App.AppState.Meta.Dir = cmd.Dir

  _, err := HasRequiredFiles(cmd.Dir, RequiredFiles)
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

// == Ground Command ==
// Should stop running containers
type GroundCommand struct {
  Controls *FlightControls
  App *TestFlight
  Dir string `short:"d" long:"dir" description:"directory to run in"`
}

func (cmd *GroundCommand) Execute(args []string) error {
  configFile, _ := cmd.Controls.CheckConfig()
  cmd.App.SetConfigFile(configFile)

  Logger.Info("Launching Tests... in dir:", cmd.Dir)
  cmd.App.SetDir(cmd.Dir)

  buildFile, _ := cmd.Controls.CheckBuild(cmd.Dir, RequiredFiles)
  cmd.App.SetBuildFile(buildFile)

  var dc = NewDockerApi(cmd.App.AppState.Meta, configFile, buildFile)
  dc.ShowInfo()
  // dc.ShowImages()

  if err := testFlightTemplates(dc, configFile); err != nil {
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

  // Stop
  // Delete container
  // Delete Image

  // if image, err := dc.CreateDockerImage(fqImageName); err != nil {
  //   return err
  // } else {
  //   if resp, err := dc.CreateContainer(image); err != nil {
  //     return err
  //   } else {
  //     Logger.Trace("Docker Container to start:", resp.Id)
  //     // dc.StopContainer(resp.Id)
  //   }
  // }

  return nil
}

// == Launch Command ==
type LaunchCommand struct {
  Controls *FlightControls
  App *TestFlight
  Dir string `short:"d" long:"dir" description:"directory to run in"`
}

func watchForEventsOn(channel ApiChannel) {
  for msg := range channel {
    Logger.Trace("DOCKER EVENT:", *msg)
  }
}

func (cmd *LaunchCommand) Execute(args []string) error {
  configFile, _ := cmd.Controls.CheckConfig()
  cmd.App.SetConfigFile(configFile)

  Logger.Info("Launching Tests... in dir:", cmd.Dir)
  cmd.App.SetDir(cmd.Dir)

  buildFile, _ := cmd.Controls.CheckBuild(cmd.Dir, RequiredFiles)
  cmd.App.SetBuildFile(buildFile)

  var dc = NewDockerApi(cmd.App.AppState.Meta, configFile, buildFile)
  dc.ShowInfo()
  // dc.ShowImages()

  if err := testFlightTemplates(dc, configFile); err != nil {
    return err
  }

  // Register channel so we can watch for events as they happen
  eventsChannel := make(ApiChannel)
  go watchForEventsOn(eventsChannel)
  dc.RegisterChannel(eventsChannel)

  fqImageName := cmd.App.AppState.BuildFile.ImageName + ":" + buildFile.Tag

  if image, err := dc.CreateDockerImage(fqImageName); err != nil {
    return err
  } else {
    if resp, err := dc.CreateContainer(image); err != nil {
      return err
    } else {
      Logger.Trace("Docker Container to start:", resp.Id)
      dc.StartContainer(resp.Id)
    }
  }

  return nil
}

// == Images Command
type ImagesCommand struct {
  Controls *FlightControls
  App *TestFlight
  Dir string `short:"d" long:"dir" description:"directory to run in"`
}

func (cmd *ImagesCommand) Execute(args []string) error {
  commandPreReq(cmd.App)

  cmd.App.SetState("IMAGES")
  Logger.Info("Listing images... using config from dir:", cmd.Dir)
  cmd.App.AppState.Meta.Dir = cmd.Dir

  dc := NewDockerApi(cmd.App.AppState.Meta, cmd.App.AppState.ConfigFile, cmd.App.AppState.BuildFile)
  return dc.ShowImages()
}

// == Template Command ==
type TemplateCommand struct {
  Controls *FlightControls
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
