package lib

import (
    "encoding/json"
    Logger "github.com/JasonGiedymin/test-flight/lib/logging"
    "strings"
)

type ApiDockerConfig struct {
    CpuShares    int
    ExposedPorts map[string]interface{} // empty interface, for future use
    Hostname     string
    Image        string
    Memory       int
    MemorySwap   int
}

// Docker Container API Post Body
type ApiPostRequest struct {
    Image        string
    OpenStdin    bool
    AttachStdin  bool
    AttachStdout bool
    AttachStderr bool
    Cmd          []string
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
    return toBytes(post)
}

type ApiContainerPortDetails struct {
    PrivatePort int
    PublicPort  int
    Type        string
}

type ApiContainer struct {
    Id         string
    Image      string
    Command    string
    Created    int64
    Status     string
    Ports      []ApiContainerPortDetails
    SizeRw     int
    SizeRootFs int
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

func (api *ApiDockerImage) Print() {
    info := []string{
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

    Logger.Console(strings.Join(info, "\n"))
}
