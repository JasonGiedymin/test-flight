package docker

import (
  "../"
  "../config"
  Logger "../logging"
  "github.com/fsouza/go-dockerclient"
  "os"
)

// Proxy Client
type DockerApi struct {
  configFile *config.ConfigFile
  client     *docker.Client
}

func NewApi(configFile *config.ConfigFile) *DockerApi {
  api := DockerApi{configFile: configFile}
  client, err := docker.NewClient(configFile.DockerEndpoint)
  if err != nil {
    Logger.Error("Docker API Client Error:", err)
    os.Exit(lib.ExitCodes["docker_error"])
  }

  api.client = client
  return &api
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
