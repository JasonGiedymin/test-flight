package lib

import (
    Logger "github.com/JasonGiedymin/test-flight/lib/logging"
)

// == Destroy Command ==
// Should destroy running containers
type DestroyCommand struct {
    Controls *FlightControls
    App      *TestFlight
    Options  *CommandOptions
}

func (cmd *DestroyCommand) Execute(args []string) error {
    Logger.Info("Destroying... using information from dir:", cmd.Options.Dir)

    // Check Config and Buildfiles
    configFile, buildFile, err := cmd.Controls.CheckConfigs(cmd.App, cmd.Options)
    if err != nil {
        return err
    }

    dc := NewDockerApi(cmd.App.AppState.Meta, configFile, buildFile)
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

    Logger.Info("Complete.")

    // Nothing to do
    return nil
}
