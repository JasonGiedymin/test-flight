package lib

import (
    Logger "github.com/JasonGiedymin/test-flight/lib/logging"
)

// == Check Command ==
type CheckCommand struct {
    Controls *FlightControls
    Options  *CommandOptions
    App      *TestFlight
    // Dir      *string //`short:"d" long:"dir" description:"directory to run in"`
    // SingleFileMode *bool //`short:"s" long:"singlefile" description:"single ansible file to use"`
}

func (cmd *CheckCommand) Execute(args []string) error {
    Logger.Info("Running Pre-Flight Check... in dir:", cmd.Options.Dir)

    // Check Config and Buildfiles
    _, _, err := cmd.Controls.CheckConfigs(cmd.App, cmd.Options)
    if err != nil {
        Logger.Error("Could not verify config files. " + err.Error())
    } else {
        Logger.Info("All checks passed! Files found!")
    }

    return nil
}
