package lib

import (
    // libtesting "github.com/JasonGiedymin/test-flight/lib/testing"
    // "os"
    "bufio"
    "bytes"
    "reflect"
    "testing"
    "text/template"
)

var standardBuildFile = BuildFile{
    Location:  "/at/the/beach",
    Owner:     "123456789",
    ImageName: "AtTheBeachOS",
    RunTests:  true,
}

// This slice is a test data table.
// Each entry contains an input file in the form of a template and
// the expected result. The expected result in turn is used to
// supply the template vars by carefully matching the vars.
// The expected result will also be deep equaled to against actual output in
// a future test. If all is well, the tested parser's output should match
// the expected data once the realized template is fed to it.
var YamlTestData = []struct {
    fileDataTemplate  string
    expectedBuildFile BuildFile
}{
    {
        `location: "{{.Location}}"
         owner: "{{.Owner}}"
         imageName: "{{.ImageName}}"
         runTests: {{.RunTests}}
        `,
        standardBuildFile,
    },
}

func TestParseYaml(t *testing.T) {

    // Loop through test data and generate templates.
    // Templates help write better tests.
    for i, testData := range YamlTestData {
        var buildFile BuildFile
        mockFileBuffer := bytes.NewBuffer(nil)      // backing buffer
        mockFile := bufio.NewWriter(mockFileBuffer) // soon to be realized template

        // Template the file data
        tmpl := template.Must(template.New("yaml-" + string(i)).Parse(testData.fileDataTemplate))
        err := tmpl.Execute(mockFile, testData.expectedBuildFile)
        if err != nil {
            t.Error("Test failed, could not execute template.", err)
        }

        if err := buildFile.ParseYaml(mockFileBuffer.Bytes()); err != nil {
            t.Log(buildFile)
            t.Error("Failed yaml parsing!")
        } else {
            if reflect.DeepEqual(buildFile, testData.expectedBuildFile) {
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
