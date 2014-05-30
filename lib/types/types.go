package types

import Logger "../logging"

type DockerAddComplexEntry struct {
  Name     string
  Location string
}

type ConfigFileDockerAdd struct {
  Simple    []string
  // User   []map[string]string
  Complex   []DockerAddComplexEntry
}

type ConfigFile struct {
  TemplateDir string
  DockerEndpoint string
  WorkDir string
  DockerAdd ConfigFileDockerAdd
}

type BuildFile struct {
  Owner             string
  ImageName         string
  Version           string
  RequiresDocker    string
  RequiresDockerUrl string
  Env               map[string]string
  Expose            []int
  Add               []DockerAddComplexEntry
  Cmd               string
  RunTests          bool
}

type ApplicationMeta struct {
  Version string
  LastCommand string
  CurrentMode string
  ExecPath    string
  Pwd         string
  Dir         string
}

type CommandOptions struct {
}

type ApplicationState struct {
  Meta        *ApplicationMeta
  Options     CommandOptions
  ConfigFile  *ConfigFile
  BuildFile   *BuildFile
}

func (appState *ApplicationState) SetState(newState string) string {
  appState.Meta.CurrentMode = newState
  Logger.Trace("STATE changed to", appState.Meta.CurrentMode)
  return appState.Meta.CurrentMode
}
