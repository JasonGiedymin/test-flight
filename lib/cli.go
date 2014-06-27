package lib

import (
  "./config"
  Logger "./logging"
  "./types"
  "os"
  // "time"
  // "fmt"
  "sync"
  "runtime"
)

func (fc *types.FlightControls) Init(app *types.TestFlight) {}

func (fc *types.FlightControls) CheckConfigs(app *types.TestFlight, singleFileMode bool, dir string) (*types.ConfigFile, *types.BuildFile, error) {
  // Prereqs
  app.SetDir(dir)

  configFile, err := config.ReadConfigFile()
  if config.ReadFileError.Contains(err) {
    os.Exit(ExitCodes["config_missing"])
  }
  app.SetConfigFile(configFile)

  requiredFiles := getRequiredFiles(singleFileMode)

  // Get the buildfile
  // TODO: as more Control funcs get created refactor this below
  buildFile, err := fc.CheckBuild(dir, requiredFiles)
  if err != nil {
    Logger.Error(err)
    return nil, nil, err
  }
  app.SetBuildFile(buildFile)

  return configFile, buildFile, nil
}

func (fc *types.FlightControls) CheckBuild(dir string, requiredFiles []types.RequiredFile) (*types.BuildFile, error) {
  // Check for test-flight specific files first
  // These are common files
  if _, err := HasRequiredFiles(dir, AnsibleFiles); err != nil {
    return nil, err
  }

  // Check for required files as specified by the user
  if _, err := HasRequiredFiles(dir, requiredFiles); err != nil {
    return nil, err
  }

  if buildFile, err := getBuildFile(dir); err != nil {
    return nil, err
  } else {
    return buildFile, nil
  }
}

func (fc *types.FlightControls) testFlightTemplates(dc *DockerApi, 
  configFile *types.ConfigFile,
  singleFileMode bool) error {

  if configFile.OverwriteTemplates {
    return dc.createTestTemplates()
  }
  return nil
}

func getRequiredFiles(filemode bool) []types.RequiredFile {
  if filemode {
    return AnsibleFiles
  } else {
    return RequiredFiles
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
  Controls *types.FlightControls
  App      *TestFlight
}

func (cmd *VersionCommand) Execute(args []string) error {
  cmd.App.SetState("VERSION_QUERY")
  Logger.Info("Test-Flight Version:", cmd.App.AppState.Meta.Version)
  return nil
}

// == Check Command ==
type CheckCommand struct {
  Controls *types.FlightControls
  App      *TestFlight
  Dir      string `short:"d" long:"dir" description:"directory to run in"`
  SingleFileMode bool `short:"s" long:"singlefile" description:"single ansible file to use"`
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
  Controls *types.FlightControls
  App      *TestFlight
  Dir      string `short:"d" long:"dir" description:"directory to run in"`
  SingleFileMode bool `short:"s" long:"singlefile" description:"single ansible file to use"`
}

func (cmd *GroundCommand) Execute(args []string) error {
  // Check Config and Buildfiles
  configFile, buildFile, err := cmd.Controls.CheckConfigs(cmd.App, cmd.SingleFileMode, cmd.Dir)
  if err != nil {
    return err
  }
  
  Logger.Info("Grounding Tests... in dir:", cmd.Dir)

  var dc = NewDockerApi(cmd.App.AppState.Meta, configFile, buildFile)
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

// == Destroy Command ==
// Should destroy running containers
type DestroyCommand struct {
  Controls *types.FlightControls
  App      *TestFlight
  Dir      string `short:"d" long:"dir" description:"directory to run in"`
  SingleFileMode bool `short:"s" long:"singlefile" description:"single ansible file to use"`
}

func (cmd *DestroyCommand) Execute(args []string) error {
  Logger.Info("Destroying... using information from dir:", cmd.Dir)

  // Check Config and Buildfiles
  configFile, buildFile, err := cmd.Controls.CheckConfigs(cmd.App, cmd.SingleFileMode, cmd.Dir)
  if err != nil {
    return err
  }

  var dc = NewDockerApi(cmd.App.AppState.Meta, configFile, buildFile)
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

// == Build Command ==
// Should build a docker image
func (cmd *types.BuildCommand) Execute(args []string) error {
  // Set vars
  Logger.Info("Building... using information from dir:", cmd.Dir)

  // Check Config and Buildfiles
  configFile, buildFile, err := cmd.Controls.CheckConfigs(cmd.App, cmd.SingleFileMode, cmd.Dir)
  if err != nil {
    return err
  }
  
  // Api interaction here
  var dc = NewDockerApi(cmd.App.AppState.Meta, configFile, buildFile)
  dc.ShowInfo()

  // Generate Templates
  // TODO: fails here with filemode
  if err := cmd.Controls.testFlightTemplates(dc, configFile, cmd.SingleFileMode); err != nil {
    return err
  }

  // Register channel so we can watch for events as they happen
  eventsChannel := make(ApiChannel)
  go watchForEventsOn(eventsChannel)
  dc.RegisterChannel(eventsChannel)

  fqImageName := cmd.App.AppState.BuildFile.ImageName + ":" + cmd.App.AppState.BuildFile.Tag

  image, err := dc.CreateDockerImage(fqImageName)
  if err != nil {
    return err
  }

  Logger.Trace("Created Docker Image:", image)
  return nil
}

// == Launch Command ==
type LaunchCommand struct {
  Controls *types.FlightControls
  App      *TestFlight
  Dir      string `short:"d" long:"dir" description:"directory to run in"`
  Force    bool   `short:"f" long:"force" description:"force new image"`
  SingleFileMode bool `short:"s" long:"singlefile" description:"single ansible file to use"`
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

func (cmd *LaunchCommand) Execute(args []string) error {
  Logger.Info("Launching Tests... in dir:", cmd.Dir)
  Logger.Debug("Force:", cmd.Force)

  var wg sync.WaitGroup // used for channels

  // Check Config and Buildfiles
  configFile, buildFile, err := cmd.Controls.CheckConfigs(cmd.App, cmd.SingleFileMode, cmd.Dir)
  if err != nil {
    return err
  }

  var dc = NewDockerApi(cmd.App.AppState.Meta, configFile, buildFile)
  dc.ShowInfo()

  if err := cmd.Controls.testFlightTemplates(dc, configFile, cmd.SingleFileMode); err != nil {
    return err
  }

  // Register channel so we can watch for events as they happen
  eventsChannel := make(ApiChannel)
  go watchForEventsOn(eventsChannel)
  dc.RegisterChannel(eventsChannel)

  fqImageName := cmd.App.AppState.BuildFile.ImageName + ":" + buildFile.Tag

  if !cmd.Force { // if not forcing, check to see if image exists.
    if image, err := dc.GetImageDetails(fqImageName); err != nil {
      return err
    } else if image != nil {
      Logger.Warn("Cannot launch new image, one with the same name already exists. User did not specify 'force' option.")
      return nil
    }
  }

  if image, err := dc.CreateDockerImage(fqImageName); err != nil {
    return err
  } else {
    if resp, err := dc.CreateContainer(image); err != nil {
      return err
    } else {
      Logger.Trace("Docker Container to start:", resp.Id)
      if _, err := dc.StartContainer(resp.Id); err != nil {
        return err
      } else {
        wg.Add(1)
        containerChannel := dc.Attach(resp.Id)
        go watchContainerOn(containerChannel, &wg)
        wg.Wait()
      }
    }
  }

  return nil
}

// == Images Command
type ImagesCommand struct {
  Controls *types.FlightControls
  App      *TestFlight
  Dir      string `short:"d" long:"dir" description:"directory to run in"`
}

func (cmd *ImagesCommand) Execute(args []string) error {
  commandPreReq(cmd.App)

  cmd.App.SetState("IMAGES")
  Logger.Info("Listing images... using config from dir:", cmd.Dir)
  cmd.App.AppState.Meta.Dir = cmd.Dir

  buildFile, _ := cmd.Controls.CheckBuild(cmd.Dir, RequiredFiles)
  fqImageName := cmd.App.AppState.BuildFile.ImageName + ":" + buildFile.Tag

  dc := NewDockerApi(cmd.App.AppState.Meta, cmd.App.AppState.ConfigFile, cmd.App.AppState.BuildFile)

  dc.GetImageDetails(fqImageName)
  return nil
}

// == Template Command ==
type TemplateCommand struct {
  Controls *types.FlightControls
  App      *TestFlight
  Dir      string `short:"d" long:"dir" description:"directory to run in"`
  SingleFileMode bool `short:"s" long:"singlefile" description:"single ansible file to use"`
}

func (cmd *TemplateCommand) Execute(args []string) error {
  commandPreReq(cmd.App)

  cmd.App.SetState("TEMPLATE")
  Logger.Info("Creating Templates... in dir:", cmd.Dir)
  cmd.App.AppState.Meta.Dir = cmd.Dir

  dc := NewDockerApi(cmd.App.AppState.Meta, cmd.App.AppState.ConfigFile, cmd.App.AppState.BuildFile)
  return cmd.Controls.testFlightTemplates(dc, cmd.App.AppState.ConfigFile, cmd.SingleFileMode)
}
