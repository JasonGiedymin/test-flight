package lib

import (
    Logger "github.com/JasonGiedymin/test-flight/lib/logging"
)

type BuildCommand struct {
    Controls *FlightControls
    Options  *CommandOptions
    App      *TestFlight
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

    // TODO: create build matrix from buildfile
    // TODO: feed matrix entry to all needing buildfile

    // Api interaction here
    dc := NewDockerApi(cmd.App.AppState.Meta, configFile, buildFile)
    dc.ShowInfo()

    // Generate Templates
    // TODO: fails here with filemode
    if err := dc.createTestTemplates(*cmd.Options); err != nil {
        return err
    }

    // Register channel so we can watch for events as they happen
    eventsChannel := make(ApiChannel)
    go watchForEventsOn(eventsChannel)
    dc.RegisterChannel(eventsChannel)

    fqImageName := buildFile.ImageName + ":" + buildFile.Tag

    image, err := dc.CreateDockerImage(fqImageName, cmd.Options)
    if err != nil {
        return err
    }

    msg := "Created Docker Image: " + image
    Logger.Console(msg)
    return nil
}
