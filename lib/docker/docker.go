package docker

import (
  // "github.com/fsouza/go-dockerclient"
  Logger "../logging"
  "../config"
)

// Proxy Client
type DockerClient struct{
  ConfigFile  *config.ConfigFile
}

func (client *DockerClient) ShowInfo() {
  Logger.Debug("---------- Test-Flight Docker Info ----------")
  Logger.Debug("Docker Endpoint:", client.ConfigFile.DockerEndpoint)
  Logger.Debug("---------------------------------------------")
}

func (client *DockerClient) ShowDockers() {
  // endpoint := "http://localhost:4243"
  // client, _ := docker.NewClient(endpoint)
  // imgs, _ := client.ListImages(true)

  // for _, img := range imgs {
  // 	fmt.Println("ID: ", img.ID)
  // 	fmt.Println("RepoTags: ", img.RepoTags)
  // 	fmt.Println("Created: ", img.Created)
  // 	fmt.Println("Size: ", img.Size)
  // 	fmt.Println("VirtualSize: ", img.VirtualSize)
  // 	fmt.Println("ParentId: ", img.ParentId)
  // 	fmt.Println("Repository: ", img.Repository)
  // }
}
