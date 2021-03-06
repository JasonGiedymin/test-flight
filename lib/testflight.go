package lib

import (
    "bitbucket.org/kardianos/osext"
    "errors"
    Logger "github.com/JasonGiedymin/test-flight/lib/logging"
    "os"
)

type TestFlight struct {
    AppState      ApplicationState
    requiredFiles []RequiredFile
}

func (app *TestFlight) SetState(state string) {
    app.AppState.SetState(state)
}

func (app *TestFlight) SetConfigFile(file *ConfigFile) {
    app.AppState.ConfigFile = file
}

func (app *TestFlight) SetBuildFile(file *BuildFile) {
    app.AppState.BuildFile = file
}

func (app *TestFlight) SetDir(dir string) {
    app.AppState.Meta.Dir = dir
}

func (app *TestFlight) Init(meta *ApplicationMeta) error {
    app.AppState.Meta = meta

    execPath, error := osext.Executable()
    if error != nil {
        msg := "Could not find executable path."
        Logger.Error(msg)
        return errors.New(msg)
    }
    // app.AppState.Meta.ExecPath = execPath
    meta.ExecPath = execPath

    pwd, error := os.Getwd()
    if error != nil {
        msg := "Could not find working directory."
        Logger.Error(msg)
        return errors.New(msg)
    }
    meta.Pwd = pwd

    return nil
}
