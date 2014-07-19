package lib

import (
    "encoding/json"
    Logger "github.com/JasonGiedymin/test-flight/lib/logging"
    "io/ioutil"
    // "strings"
    // "errors"
)

type DockerAddComplexEntry struct {
    Name     string
    Location string
}

type DockerAdd struct {
    Simple  []string                // top level dir entries added by Docker template
    Complex []DockerAddComplexEntry // complex free form entries added by Docker template
}

type BuildFile struct {
    Location      string
    Owner         string
    ImageName     string
    Tag           string
    From          string
    Requires      []string
    Version       string
    Env           map[string]string
    Expose        []int
    Ignore        []string
    Add           DockerAdd
    Cmd           string
    WorkDir       string
    RunTests      bool
    ResourceShare ResourceShare
}

// most basic non prescriptive layout = nothing
// var basicAdd = DockerAdd{
//     Simple: []string{
//         "meta",
//         "tasks",
//         "tests",
//     },  // unless overridden
// }

// fully prescriptive layout
// TODO: implement wizard mode
var wizardAdd = DockerAdd{
    Simple: []string{
        "defaults",
        "handlers",
        "meta",
        "tasks",
        "templates",
        "tests",
        "vars",
        "vault",
    },  // unless overridden
}

// For specific defaults
func NewBuildFile() *BuildFile {
    return &BuildFile{
        RunTests:  false,
        Owner:     "Test-Flight-User", // must have something!
        ImageName: "Test-Flight-Test-Image", // must have an image name!
        Tag:       "latest", // implies latest, else you tell it
        From:      "", // from nothing
        Version:   "0.0.1", // must have a version!
        WorkDir:   "/tmp/build", // default working dir
        Add:       DockerAdd{}, // add nothing
        Ignore:    []string{".git"} // by default if user does not specify anything we will ignore git
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
