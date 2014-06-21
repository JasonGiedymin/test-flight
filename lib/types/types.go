package types

import (
  Logger "../logging"
  "encoding/json"
)

func bytes(data interface{}) ([]byte, error) {
  bytes, err := json.Marshal(data)
  if err != nil {
    return nil, err
  }
  return bytes, nil
}

type ApiPostRequest struct {
  Image       string
  OpenStdin   bool
  AttachStdin bool
}

func (post *ApiPostRequest) Bytes() ([]byte, error) {
  bytes, err := json.Marshal(post)
  if err != nil {
    return nil, err
  }
  return bytes, nil
}

type ApiPostResponse struct {
  Id       string
  Warnings []string
}

type DockerHostConfig struct {
  Binds           []string
  Links           []string
  LxcConf         []string
  PortBindings    map[string][]map[string]string
  PublishAllPorts bool
  Privileged      bool
  Dns             []string
  VolumesFrom     []string
}

func (post *DockerHostConfig) Bytes() ([]byte, error) {
  return bytes(post)
}

type ApiContainer struct {
  Command string
  Created int64
  Id      string
  Image   string
  Names   []string
  Ports   []int
  Status  string
}

type ApiDockerConfig struct {
  CpuShares    int
  ExposedPorts map[string]interface{} // empty interface, for future use
  Hostname     string
  Image        string
  Memory       int
  MemorySwap   int
}

type ApiDockerImage struct {
  Architecture    string
  Author          string
  Comment         string
  Config          ApiDockerConfig
  Container       string
  ContainerConfig ApiDockerConfig
  DockerVersion   string
  Id              string
  Os              string
  Parent          string
}

type ResourceShare struct {
  Mem int
  Cpu int
}

type DockerAddComplexEntry struct {
  Name     string
  Location string
}

type ConfigFileDockerAdd struct {
  Simple []string
  // User   []map[string]string
  Complex []DockerAddComplexEntry
}

type ConfigFile struct {
  TemplateDir        string
  DockerEndpoint     string
  WorkDir            string
  DockerAdd          ConfigFileDockerAdd
  OverwriteTemplates bool
}

type BuildFile struct {
  Owner             string
  ImageName         string
  Tag               string
  From              string
  Version           string
  RequiresDocker    string
  RequiresDockerUrl string
  Env               map[string]string
  Expose            []int
  Add               []DockerAddComplexEntry
  Cmd               string
  RunTests          bool
  ResourceShare     ResourceShare
}

type ApplicationMeta struct {
  Version     string
  LastCommand string
  CurrentMode string
  ExecPath    string
  Pwd         string
  Dir         string
}

type CommandOptions struct {
}

type ApplicationState struct {
  Meta       *ApplicationMeta
  Options    CommandOptions
  ConfigFile *ConfigFile
  BuildFile  *BuildFile
}

func (appState *ApplicationState) SetState(newState string) string {
  appState.Meta.CurrentMode = newState
  Logger.Trace("STATE changed to", appState.Meta.CurrentMode)
  return appState.Meta.CurrentMode
}

type RequiredFile struct {
  Name          string
  FileName      string
  FileType      string // [f]ile, [d]ir
  RequiredFiles []RequiredFile
}

// type MapLike interface {
//   Map(v interface{})
// }

type requiredFileMapFunc func(in RequiredFile) RequiredFile

// requiredFile.Map(MapLike)
func (req *RequiredFile) Map(fx requiredFileMapFunc) []RequiredFile {
  returnSlice := make([]RequiredFile, len(req.RequiredFiles))
  return returnSlice
}
