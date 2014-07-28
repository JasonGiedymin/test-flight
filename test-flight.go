package main

import (
    "github.com/JasonGiedymin/test-flight/lib"
    Logger "github.com/JasonGiedymin/test-flight/lib/logging"
    "github.com/jessevdk/go-flags"
    "os"
)

/* Limited to the app, parser, and commands */
var (
    app     lib.TestFlight
    parser  *flags.Parser
    options lib.CommandOptions
)

var meta = lib.ApplicationMeta{
    Version: "0.9.8.2",
}

// == App ==
func init() {
    err := app.Init(&meta)
    if err != nil {
        os.Exit(lib.ExitCodes["init_fail"])
    }

    flightControls := lib.FlightControls{}

    checkCommand := lib.CheckCommand{Controls: &flightControls, App: &app, Options: &options}
    imagesCommand := lib.ImagesCommand{Controls: &flightControls, App: &app, Options: &options}
    buildCommand := lib.BuildCommand{Controls: &flightControls, App: &app, Options: &options}
    launchCommand := lib.LaunchCommand{Controls: &flightControls, App: &app, Options: &options}
    groundCommand := lib.GroundCommand{Controls: &flightControls, App: &app, Options: &options}
    destroyCommand := lib.DestroyCommand{Controls: &flightControls, App: &app, Options: &options}
    versionCommand := lib.VersionCommand{Controls: &flightControls, App: &app, Options: &options}
    templateCommand := lib.TemplateCommand{Controls: &flightControls, App: &app, Options: &options}

    options = lib.CommandOptions{}

    parser = flags.NewParser(&options, flags.Default)

    parser.AddCommand("check",
        "flight check",
        "Checks if pre-reqs are satisfied to launch this ansible playbook.",
        &checkCommand)

    parser.AddCommand("images",
        "flight images",
        "Shows all images",
        &imagesCommand)

    parser.AddCommand("build",
        "flight build",
        "Build will build an ansible playbook docker image.",
        &buildCommand)

    parser.AddCommand("launch",
        "flight launch",
        "Launch builds, runs, and tests an ansible playbook.",
        &launchCommand)

    parser.AddCommand("ground",
        "flight ground",
        "Ground stops all containers found running this ansible playbook.",
        &groundCommand)

    parser.AddCommand("destroy",
        "flight destroy",
        "Destroy destroys all containers and images found running this ansible playbook.",
        &destroyCommand)

    parser.AddCommand("version",
        "shows version",
        "Show Test-Flight version number.",
        &versionCommand)

    parser.AddCommand("template",
        "flight template",
        "Creates templates required for test-flight.",
        &templateCommand)
}

// Runs Test-Flight
func main() {
    if _, err := parser.Parse(); err != nil {
        // Exclude flags.Error case because `--help` will return an error
        // https://github.com/jessevdk/go-flags/issues/45
        if _, ok := err.(*flags.Error); !ok { // type assertion returns (var,bool)
            Logger.Error(err)
            os.Exit(lib.ExitCodes["command_fail"])
        }
    }
}
