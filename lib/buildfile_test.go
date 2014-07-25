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
    Tag:       "BeachTag",
    From:      "FromTheBeach",
    Requires:  []string{"ABeach", "BBeach"},
    Version:   "WhiteSands",
    Env: map[string]string{
        "BeachName": "LovelySands",
        "BeachSize": "AwesomeEnough",
    },
    Expose: []int{3000, 3001},
    Ignore: []string{".git"},
    Add: DockerAdd{
        Simple: []string{"simpleAdd1", "simpleAdd2"},
        Complex: []DockerAddComplexEntry{
            DockerAddComplexEntry{Name: "ComplexName1", Location: "ComplexLocation1"},
        },
    },
    Cmd:           "run to the beach",
    LaunchCmd:     []string{"run faster", "to the beach"},
    WorkDir:       "/tmp/beach",
    RunTests:      true,
    ResourceShare: ResourceShare{Mem: 256, Cpu: 1},
}

// This slice is a test data table.
// Each entry contains an input file in the form of a template and
// the expected result. The expected result in turn is used to
// supply the template vars by carefully matching the vars.
// The expected result will also be deep equaled to against actual output in
// a future test. If all is well, the tested parser's output should match
// the expected data once the realized template is fed to it.
// Why? I Found naming differences to require this test.
var YamlTestData = []struct {
    fileDataTemplate  string
    expectedBuildFile BuildFile
}{
    {`
      location: "{{.Location}}"
      owner: "{{.Owner}}"
      imagename: "{{.ImageName}}"
      tag: {{.Tag}}
      from: {{.From}}
      requires: {{range $key, $value := .Requires}}
        - {{$value}}{{end}}
      version: {{.Version}}
      env: {{range $key, $value := .Env}}
        {{$key}}: {{$value}}{{end}}
      expose: {{range $key, $value := .Expose}}
        - {{$value}}{{end}}
      ignore: {{range $key, $value := .Ignore}}
        - {{$value}}{{end}}
      add:
        simple: {{range $key, $value := .Add.Simple}}
        - {{$value}}{{end}}
        complex: {{range $key, $value := .Add.Complex}}
          - name: {{$value.Name}}
            location: {{$value.Location}}{{end}}
      cmd: "{{.Cmd}}"
      launchcmd: {{range $key, $value := .LaunchCmd}}
        - {{$value}}{{end}}
      workdir: "{{.WorkDir}}"
      runTests: {{.RunTests}}
      resourceshare:
        mem : {{.ResourceShare.Mem}}
        cpu: {{.ResourceShare.Cpu}}
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

        mockFile.Flush()

        if err := buildFile.ParseYaml(mockFileBuffer.Bytes()); err != nil {
            t.Log(buildFile)
            t.Error("Failed yaml parsing!")
        } else {
            if !reflect.DeepEqual(buildFile, testData.expectedBuildFile) {
                t.Error("Failed yaml parsing!")
                t.Log("*********buildfile*********")
                t.Log(buildFile)
                t.Log("*********expected*********")
                t.Log(testData.expectedBuildFile)
            } else { // if equals, output the templated file that was parsed
                t.Log("******************")
                t.Log(mockFileBuffer.String())
                t.Log("******************")
            }
        }
    }
}

// func TestParseJson(t *testing.T) {
//     file := `{
//       "runTests": true,
//       "owner": "123456789"
//     }`

//     var buildFile BuildFile
//     buildFile.ParseJson([]byte(file))

//     t.Log(buildFile)
// }
