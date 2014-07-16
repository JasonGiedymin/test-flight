package lib

import (
    "encoding/json"
    Logger "github.com/JasonGiedymin/test-flight/lib/logging"
    "io/ioutil"
    // "strings"
    // "errors"
)

type BuildFile struct {
    Location          string
    Owner             string
    ImageName         string
    Tag               string
    From              string
    Version           string
    RequiresDocker    string
    RequiresDockerUrl string
    Env               map[string]string
    Expose            []int
    Ignore            []string
    Add               []DockerAddComplexEntry
    Cmd               string
    RunTests          bool
    ResourceShare     ResourceShare
}

// For specific defaults
func NewBuildFile() *BuildFile {
    return &BuildFile{
        Owner:             "Test-Flight-User",
        ImageName:         "Test-Flight-Test-Image",
        Tag:               "latest",
        From:              "",
        Version:           "0.0.1",
        RequiresDocker:    "",
        RequiresDockerUrl: "",
        RunTests:          false,
    }
}

func ReadBuildFile(filePath string) (*BuildFile, error) {
    jsonBlob, _ := ioutil.ReadFile(filePath)

    var buildFile = NewBuildFile()
    err := json.Unmarshal(jsonBlob, buildFile)
    if err != nil {
        Logger.Error("Error while trying to parse buildfile,", filePath, err)
        return nil, err
    }

    buildFile.Location = filePath

    return buildFile, nil
}
