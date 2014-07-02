package lib

import (
  "os"
  Logger "./logging"
  "errors"
)

type FlightControls struct{}

func (fc *FlightControls) Init(app *TestFlight) {}

func (fc *FlightControls) CheckConfigs(app *TestFlight, options *CommandOptions) (*ConfigFile, *BuildFile, error) {
  if (options.Configfile != "") {
    Logger.Info("Using configfile:", options.Configfile)
  }

  // Prereqs
  app.SetDir(options.Dir)

  configFile, err := ReadConfigFile(options.Configfile)
  if ReadFileError.Contains(err) {
    os.Exit(ExitCodes["config_missing"])
  }
  app.SetConfigFile(configFile)

  requiredFiles := getRequiredFiles(options.SingleFileMode)

  // Get the buildfile
  // TODO: as more Control funcs get created refactor this below
  buildFile, err := fc.CheckBuild(options.Dir, requiredFiles)
  if err != nil {
    return nil, nil, err
  }
  app.SetBuildFile(buildFile)

  return configFile, buildFile, nil
}

func (fc *FlightControls) CheckBuild(dir string, requiredFiles []RequiredFile) (*BuildFile, error) {
  // Check for test-flight specific files first
  // These are common files
  if _, err := HasRequiredFiles(dir, AnsibleFiles); err != nil {
    msg := "Error reading required Ansible Files. Error: " + err.Error()
    return nil, errors.New(msg)
  }

  // Check for required files as specified by the user
  if _, err := HasRequiredFiles(dir, requiredFiles); err != nil {
    msg := "Error reading user specified required file. Error: " + err.Error()
    return nil, errors.New(msg)
  }

  buildFilePath := FilePath(dir, "build.json")
  if buildFile, err := ReadBuildFile(buildFilePath); err != nil {
    msg := "Error parsing buildfile: [" + buildFilePath + "]. Error: " + err.Error()
    return nil, errors.New(msg)
  } else {
    return buildFile, nil
  }
}

func (fc *FlightControls) testFlightTemplates(dc *DockerApi, 
  configFile *ConfigFile,
  singleFileMode bool) error {

  if configFile.OverwriteTemplates {
    return dc.createTestTemplates()
  }
  return nil
}