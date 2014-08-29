package lib

import (
    Logger "github.com/JasonGiedymin/test-flight/lib/logging"
    registry "github.com/JasonGiedymin/voom-registry/registry/client"

    "errors"
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

    // // TODO: create build matrix from buildfile
    // // TODO: feed matrix entry to all needing buildfile

    build := func(buildMatrixEntry BuildMatrixEntry, result BuildMatrixEntryResult) BuildMatrixEntryResult {

        if buildFile.From == "" { // Get the from if we don't have it
            // registry service interaction
            query := registry.SearchDataJSON{
                Os:       buildMatrixEntry.OS,
                Language: buildMatrixEntry.Language,
                Version:  buildMatrixEntry.Version,
            }
            registryApi := registry.RegistryApi{"http://registry.amuxbit.com"}
            fromDocker, err := registryApi.Search(query)
            if err != nil {
                result.Err = errors.New("Error connecting with registry. " + err.Error())
                return result
            } else {
                if len(fromDocker.Results) > 1 {
                    msg := `Combination of OS, Language, and Language Version yielded 
                    multiple results. Please modify.`

                    result.Err = errors.New(msg)
                }
                buildFile.From = fromDocker.Results[0]
            }
        }

        // Api interaction here
        dc := NewDockerApi(cmd.App.AppState.Meta, configFile, buildFile)
        dc.ShowInfo()

        // Generate Templates
        // TODO: fails here with filemode
        if err := dc.createTestTemplates(*cmd.Options); err != nil {
            result.Err = err
            return result
        }

        // Register channel so we can watch for events as they happen
        eventsChannel := make(ApiChannel)
        go watchForEventsOn(eventsChannel)
        dc.RegisterChannel(eventsChannel)

        fqImageName := buildFile.ImageName + ":" + buildFile.Tag

        image, err := dc.CreateDockerImage(fqImageName, cmd.Options)
        if err != nil {
            result.Err = err
            return result
        }

        msg := "Created Docker Image: " + image
        Logger.Console(msg)
        return result
    }

    return RunEntry(buildFile, build, "Error while building.")
}
