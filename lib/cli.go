package lib

import (
  Logger "./logging"
  // "os"
  // "time"
  // "fmt"
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

func (cmd *CheckCommand) Execute(args []string) error {
  Logger.Info("Running Pre-Flight Check... in dir:", cmd.Options.Dir)
  _, buildFile, err := cmd.Controls.CheckConfigs(cmd.App, cmd.Options)
  if err != nil {
    return err
  }

  cmd.App.AppState.Meta.Dir = cmd.Options.Dir

  _, err = HasRequiredFiles(cmd.Options.Dir, RequiredFiles)
  if err != nil {
    return err
  }

  buildFile, err = ReadBuildFile(FilePath(cmd.Options.Dir, "build.json"))
  if err != nil {
    return err
  }

  cmd.App.AppState.BuildFile = buildFile

  Logger.Debug("Buildfile found, contents:", *buildFile)
  return nil
}

func (cmd *GroundCommand) Execute(args []string) error {
  // Check Config and Buildfiles
  configFile, buildFile, err := cmd.Controls.CheckConfigs(cmd.App, cmd.Options)
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

func (cmd *DestroyCommand) Execute(args []string) error {
  Logger.Info("Destroying... using information from dir:", cmd.Dir)

  // Check Config and Buildfiles
  configFile, buildFile, err := cmd.Controls.CheckConfigs(cmd.App, cmd.Options)
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
func (cmd *BuildCommand) Execute(args []string) error {
  // Set vars
  Logger.Info("Building... using information from dir:", cmd.Options.Dir)

  // Check Config and Buildfiles
  configFile, buildFile, err := cmd.Controls.CheckConfigs(cmd.App, cmd.Options)
  if err != nil {
    return err
  }
  
  // Api interaction here
  var dc = NewDockerApi(cmd.App.AppState.Meta, configFile, buildFile)
  dc.ShowInfo()

  // Generate Templates
  // TODO: fails here with filemode
  if err := cmd.Controls.testFlightTemplates(dc, configFile, cmd.Options.SingleFileMode); err != nil {
    return err
  }

  // Register channel so we can watch for events as they happen
  eventsChannel := make(ApiChannel)
  go watchForEventsOn(eventsChannel)
  dc.RegisterChannel(eventsChannel)

  fqImageName := cmd.App.AppState.BuildFile.ImageName + ":" + cmd.App.AppState.BuildFile.Tag

  image, err := dc.CreateDockerImage(fqImageName, cmd.Options.SingleFileMode)
  if err != nil {
    return err
  }

  Logger.Trace("Created Docker Image:", image)
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

func (cmd *LaunchCommand) Execute(args []string) error {
  Logger.Info("Launching Tests... in dir:", cmd.Options.Dir)
  Logger.Debug("Force:", cmd.Options.Force)

  var wg sync.WaitGroup // used for channels

  // Check Config and Buildfiles
  configFile, buildFile, err := cmd.Controls.CheckConfigs(cmd.App, cmd.Options)
  if err != nil {
    return err
  }

  var dc = NewDockerApi(cmd.App.AppState.Meta, configFile, buildFile)
  dc.ShowInfo()

  if err := cmd.Controls.testFlightTemplates(dc, configFile, cmd.Options.SingleFileMode); err != nil {
    return err
  }

  // Register channel so we can watch for events as they happen
  eventsChannel := make(ApiChannel)
  go watchForEventsOn(eventsChannel)
  dc.RegisterChannel(eventsChannel)

  fqImageName := cmd.App.AppState.BuildFile.ImageName + ":" + buildFile.Tag

  if !cmd.Options.Force { // if not forcing, check to see if image exists.
    if image, err := dc.GetImageDetails(fqImageName); err != nil {
      return err
    } else if image != nil {
      Logger.Warn("Cannot launch new image, one with the same name already exists. User did not specify 'force' option.")
      return nil
    }
  }

  if image, err := dc.CreateDockerImage(fqImageName, cmd.Options.SingleFileMode); err != nil {
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

func (cmd *ImagesCommand) Execute(args []string) error {
  cmd.App.SetState("IMAGES")
  Logger.Info("Listing images... using config from dir:", cmd.Options.Dir)

  _, buildFile, err := cmd.Controls.CheckConfigs(cmd.App, cmd.Options)
  if err != nil {
    return err
  }  

  cmd.App.AppState.Meta.Dir = cmd.Options.Dir

  fqImageName := cmd.App.AppState.BuildFile.ImageName + ":" + buildFile.Tag

  dc := NewDockerApi(cmd.App.AppState.Meta, cmd.App.AppState.ConfigFile, cmd.App.AppState.BuildFile)

  dc.GetImageDetails(fqImageName)
  return nil
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
