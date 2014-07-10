package lib

import (
  "strings"
  Logger "github.com/JasonGiedymin/test-flight/lib/logging"
)

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

func (api *ApiDockerImage) Print() {
  info := []string {
    "",
    "Architecture: " + api.Architecture,
    "Author: " + api.Author,
    "Comment: " + api.Comment,
    "Container: " + api.Container,
    "DockerVersion: " + api.DockerVersion,
    "Container ID: " + api.Id,
    "OS: " + api.Os,
    "Parent: " + api.Parent,
    "",
  }

  Logger.Info(strings.Join(info, "\n"))
}