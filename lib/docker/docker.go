package docker

import (
  "../"
  "../config"
  Logger "../logging"
  "github.com/fsouza/go-dockerclient"
  "os"
  "path/filepath"
  "text/template"
)

type TemplateVar struct {
  Owner string
  ImageName string
  Version string
  RequiresDocker string
  RequiresDockerUrl string
  WorkDir string
  Env map[string]string
  Expose map[string]string
  Cmd string
}

// Proxy Client
type DockerApi struct {
  configFile *config.ConfigFile
  buildFile  *config.BuildFile
  client     *docker.Client
}

func NewApi(configFile *config.ConfigFile, buildFile *config.BuildFile) *DockerApi {
  api := DockerApi{configFile: configFile, buildFile: buildFile}
  client, err := docker.NewClient(configFile.DockerEndpoint)
  if err != nil {
    Logger.Error("Docker API Client Error:", err)
    os.Exit(lib.ExitCodes["docker_error"])
  }

  api.client = client
  return &api
}

func (api *DockerApi) getTemplateVar() *TemplateVar {
  return &TemplateVar{
    Owner: api.buildFile.Owner,
    ImageName: api.buildFile.ImageName,
    Version: api.buildFile.Version,
    RequiresDocker: api.buildFile.RequiresDocker,
    RequiresDockerUrl: api.buildFile.RequiresDockerUrl,
  }
}

func (api *DockerApi) CreateTemplate() {
  pwd, err := os.Getwd()
  if err != nil {
    // return nil, err
  }

  pattern := filepath.Join(pwd + "/templates/", "*.tmpl")
  tmpl := template.Must(template.ParseGlob(pattern))

  // var x = TemplateVar{
  //   Owner: api.buildFile.Owner,
  // }

  Logger.Trace("-->", *api.getTemplateVar())

  err = tmpl.Execute(os.Stdout, *api.getTemplateVar())
  if err != nil {
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
  Logger.Trace( api.createDockerFile() )
}
