package lib

import (
    // "errors"
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

    getImageId := func() (string, error) {

        if image, err := dc.GetImageDetails(fqImageName); err != nil {
            return "", err
        } else { // response from endpoint (but could be 404)
            if image == nil { // doesn't exist
                if newImage, err := dc.CreateDockerImage(fqImageName, cmd.Options); err != nil {
                    return "", err
                } else {
                    return newImage, nil
                }
            } else { // exists
                if cmd.Options.Force { // recreate anyway
                    if newImage, err := dc.CreateDockerImage(fqImageName, cmd.Options); err != nil {
                        return "", err
                    } else {
                        return newImage, nil
                    }
                } else { // warn user
                    Logger.Warn("Will be starting a container that already exists.")
                    return image.Id, nil
                }
            }
            return image.Id, nil
        }
    }

    createAndStartContainerFrom := func(imageId string) error {
        if container, err := dc.CreateContainer(imageId); err != nil {
            return err
        } else {
            Logger.Trace("Docker Container to start:", container.Id)
            if _, err := dc.StartContainer(container.Id); err != nil {
                return err
            } else {
                wg.Add(1)
                containerChannel := dc.Attach(container.Id)
                go watchContainerOn(containerChannel, &wg)
                wg.Wait()
                return nil
            }
        }
    }

    if imageId, err := getImageId(); err != nil {
        return err
    } else {
        msg := "Launching docker container: " + imageId
        Logger.Console(msg)

        return createAndStartContainerFrom(imageId)
    }

    return nil
}
