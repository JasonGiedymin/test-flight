package lib

import (
  Logger "github.com/JasonGiedymin/test-flight/lib/logging"
)

// == Images Command
type ImagesCommand struct {
  Controls *FlightControls
  Options  *CommandOptions
  App      *TestFlight
}

func (cmd *ImagesCommand) Execute(args []string) error {
  Logger.Info("Listing images...")

  // Check Config and Buildfiles
  configFile, buildFile, err := cmd.Controls.CheckConfigs(cmd.App, cmd.Options)
  if err != nil {
    return err
  }
  
  // Api interaction here
  dc := NewDockerApi(cmd.App.AppState.Meta, configFile, buildFile)
  dc.ShowInfo()

  fqImageName := buildFile.ImageName + ":" + buildFile.Tag

  if imageDetails, err := dc.GetImageDetails(fqImageName); err != nil {
    return err
  } else {
    imageDetails.Print();
  }


  return nil
}