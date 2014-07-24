package lib

import (
    // libtesting "github.com/JasonGiedymin/test-flight/lib/testing"
    // "os"
    "testing"
)

var YamlTestData = []struct {
    fileData  string
    buildFile BuildFile
}{
    {
        `
    runTests: true
    owner: "123456789"
    `,
        BuildFile{RunTests: true, Owner: "123456789"},
    },
}

func TestParseYaml(t *testing.T) {
    file := `
    runTests: true
    owner: "123456789"
    `

    for i, value := range YamlTestData {
        var buildFile BuildFile
        if err := buildFile.ParseYaml([]byte(value.fileData)); err != nil {
            t.Log(buildFile)
            t.Error("Failed yaml parsing!")
        } else {
            if buildFile != value.buildFile {
                t.Error("Failed yaml parsing!")
            }
        }
    }
}

func TestParseJson(t *testing.T) {
    file := `{
      "runTests": true,
      "owner": "123456789"
    }`

    var buildFile BuildFile
    buildFile.ParseJson([]byte(file))

    t.Log(buildFile)
}
