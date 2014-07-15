package lib

import (
    Logger "github.com/JasonGiedymin/test-flight/lib/logging"
    "sync"
)

// == Launch Command ==
type LaunchCommand struct {
    Controls *FlightControls
    Options  *CommandOptions
    App      *TestFlight
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

    dc := NewDockerApi(cmd.App.AppState.Meta, configFile, buildFile)
    dc.ShowInfo()

    if err := dc.createTestTemplates(*cmd.Options); err != nil {
        return err
    }
    Logger.Trace("Created test-flight templates.")

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

    if image, err := dc.CreateDockerImage(fqImageName, cmd.Options); err != nil {
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

    Logger.Info("Complete.")

    return nil
}
