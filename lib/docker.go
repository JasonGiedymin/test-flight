package lib

import (
  "./config"
  Logger "./logging"
  "github.com/fsouza/go-dockerclient"
  "os"
  "path/filepath"
  "text/template"
)

type TemplateVar struct {
  Owner             string
  ImageName         string
  Version           string
  RequiresDocker    string
  RequiresDockerUrl string
  WorkDir           string
  Env               map[string]string
  Expose            map[string]string
  Cmd               string
  AddSystem         []string
  AddUser           []config.ConfigFileUserAdd
  // AddUser            map[string]string
  ConfigFile        *config.ConfigFile
  BuildFile         *config.BuildFile
}

// Proxy Client
type DockerApi struct {
  configFile *config.ConfigFile
  buildFile  *config.BuildFile
  client     *docker.Client
}

func NewDockerApi(configFile *config.ConfigFile, buildFile *config.BuildFile) *DockerApi {
  api := DockerApi{configFile: configFile, buildFile: buildFile}
  client, err := docker.NewClient(configFile.DockerEndpoint)
  if err != nil {
    Logger.Error("Docker API Client Error:", err)
    os.Exit(ExitCodes["docker_error"])
  }

  api.client = client
  return &api
}

// One big proxy obj to help users. Slowly phase this out.
func (api *DockerApi) getTemplateVar() *TemplateVar {
  return &TemplateVar{
    // Direct:
    ConfigFile:        api.configFile,
    BuildFile:         api.buildFile,

    // Helpers for common accessors
    // Keep names simple!
    Owner:             api.buildFile.Owner,
    ImageName:         api.buildFile.ImageName,
    Version:           api.buildFile.Version,
    RequiresDocker:    api.buildFile.RequiresDocker,
    RequiresDockerUrl: api.buildFile.RequiresDockerUrl,
    WorkDir:           api.configFile.WorkDir,
    Env:               api.buildFile.Env,
    AddSystem:         api.configFile.DockerAdd.System,
    AddUser:           api.configFile.DockerAdd.User,
  }
}

func (api *DockerApi) CreateTemplate() {
  pwd, err := os.Getwd()
  if err != nil {
    // return nil, err
  }

  pattern := filepath.Join(pwd+"/templates/", "*.tmpl")
  tmpl := template.Must(template.ParseGlob(pattern))

  if err := tmpl.Execute(os.Stdout, *api.getTemplateVar()); err != nil {
    Logger.Error("template execution: %s", err)
  }
}

func (api *DockerApi) ShowInfo() {
  Logger.Debug("---------- Test-Flight Docker Info ----------")
  Logger.Debug("Docker Endpoint:", api.configFile.DockerEndpoint)
  Logger.Debug("---------------------------------------------")
}

func (api *DockerApi) ShowImages() {
  images, _ := api.client.ListImages(true)

  if len(images) <= 0 {
    Logger.Info("No docker images found.")
    return
  }

  for _, img := range images {
    Logger.Info("ID: ", img.ID)
    Logger.Info("RepoTags: ", img.RepoTags)
    Logger.Info("Created: ", img.Created)
    Logger.Info("Size: ", img.Size)
    Logger.Info("VirtualSize: ", img.VirtualSize)
    Logger.Info("ParentId: ", img.ParentId)
    Logger.Info("Repository: ", img.Repository)
  }
}

func (api *DockerApi) createDockerFile() string {
  dockerFile := `
  # Dockerfile
  # ----------
  #
  `

  return dockerFile
}

func (api *DockerApi) CreateDocker() {
}
