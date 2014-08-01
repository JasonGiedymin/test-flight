package lib

import (
    "errors"
    Logger "github.com/JasonGiedymin/test-flight/lib/logging"
    // "os"
)

type FlightControls struct{}

func (fc *FlightControls) Init(app *TestFlight) {}

// Checks if files are proper and can be parsed
func (fc *FlightControls) CheckConfigs(app *TestFlight, options *CommandOptions) (*ConfigFile, *BuildFile, error) {
    // Setup Logging
    // TODO: Replace with only Load once app.Init() is gone
    Logger.Load(len(options.Verbose))
    Logger.Setup()

    // Set vars
    if options.Dir == "" {
        options.Dir = "./" // set to local working dir as default
    }
    Logger.Info("Using directory:", options.Dir)

    // Prereqs
    app.SetDir(options.Dir)

    configFile, err := ReadConfigFile(options.Configfile, options.Dir)
    if err != nil {
        return nil, nil, err
    }

    app.SetConfigFile(configFile)
    Logger.Info("Using configfile:", configFile.Location)

    // Get the buildfile
    // TODO: as more Control funcs get created refactor this below
    buildFile, err := fc.CheckBuild(options.Dir)
    if err != nil {
        return nil, nil, err
    }
    app.SetBuildFile(buildFile)
    Logger.Info("Using buildfile:", buildFile.Location)
    Logger.Debug("buildfile - contents:", *buildFile)

    if err := fc.CheckRequiredFiles(options.Dir, options.SingleFileMode, buildFile); err != nil {
        return nil, nil, err
    }

    return configFile, buildFile, nil
}

func (fc *FlightControls) CheckRequiredFiles(
    dir string,
    singleFileMode bool,
    buildFile *BuildFile,
) error {
    var requiredFiles []RequiredFile

    if singleFileMode {
        requiredFiles = AnsibleFiles
    } else {
        if files, err := generateRequiredFilesFrom(buildFile); err != nil {
            return err
        } else {
            requiredFiles = files
        }
    }

    if _, err := HasRequiredFiles(dir, requiredFiles); err != nil {
        msg := "Error reading user specified required file. Error: " + err.Error()
        return errors.New(msg)
    }

    return nil
}

func (fc *FlightControls) CheckBuild(dir string) (*BuildFile, error) {
    requiredFiles := []RequiredFile{buildFile}
    // Check for required files as specified by the user
    if _, err := HasRequiredFiles(dir, requiredFiles); err != nil {
        msg := "Error reading user specified required file. Error: " + err.Error()
        return nil, errors.New(msg)
    }

    buildFilePath := FilePath(dir, Constants().buildFileName)
    if buildFile, err := ReadBuildFile(buildFilePath); err != nil {
        msg := "Error parsing buildfile: [" + buildFilePath + "]. Error: " + err.Error()
        return nil, errors.New(msg)
    } else {
        return buildFile, nil
    }
}
