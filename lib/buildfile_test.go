package lib

import (
    libtesting "github.com/JasonGiedymin/test-flight/lib/testing"
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
    // From:      []string{"FromTheBeach", "FromTheResort"},
    From:            "FromTheBeach",
    OS:              []string{"BeachOS", "SandyOS"},
    Language:        "BeachLang++",
    LanguageVersion: []string{"1.0.0", "2.0.0"},
    Requires:        []string{"ABeach", "BBeach"},
    Version:         "WhiteSands",
    Env: []DockerEnv{
        DockerEnv{Variable: "BeachName", Value: "LovelySands"},
        DockerEnv{Variable: "BeachSize", Value: "AwesomeEnough"},
    },
    Expose: []int{3000, 3001},
    Ignore: []string{".git"},
    Add: DockerAdd{
        Simple: []string{"simpleAdd1", "simpleAdd2"},
        Complex: []DockerAddComplexEntry{
            DockerAddComplexEntry{Name: "ComplexName1", Location: "ComplexLocation1"},
            DockerAddComplexEntry{Name: "ComplexName2", Location: "ComplexLocation2"},
        },
    },
    Cmd:       "run to the beach",
    LaunchCmd: []string{"run faster", "to the beach"},
    WorkDir:   "/tmp/beach",
    RunTests:  true,
    Resources: ResourceShare{Mem: 256, Cpu: 1},
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
    //from: {{range $key, $value := .From}}
    //      - {{$value}}{{end}}
    {`
      location: {{.Location}}
      owner: {{.Owner}}
      imageName: {{.ImageName}}
      tag: {{.Tag}}
      from: {{.From}}
      requires: {{range $key, $value := .Requires}}
        - {{$value}}{{end}}
      version: {{.Version}}
      env: {{range $i, $entry := .Env}}
        - variable: {{$entry.Variable}}
          value: {{$entry.Value}}{{end}}
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
      cmd: {{.Cmd}}
      launchCmd: {{range $key, $value := .LaunchCmd}}
        - {{$value}}{{end}}
      workDir: {{.WorkDir}}
      runTests: {{.RunTests}}
      resources:
        cpu: {{.Resources.Cpu}}
        mem: {{.Resources.Mem}}
      `,
        standardBuildFile,
    },
}

var JsonTestData = []struct {
    fileDataTemplate  string
    expectedBuildFile BuildFile
}{
    //"from": [ {{range $i, $value := .From}}
    //   "{{$value}}"{{if comma $i $.From | not}},{{end}}{{end}}
    //],
    {`{
      "location": "{{.Location}}",
      "owner": "{{.Owner}}",
      "imageName": "{{.ImageName}}",
      "tag": "{{.Tag}}",
      "from": "{{.From}}",
      "requires":[ {{range $i, $value := .Requires}}
        "{{$value}}"{{if comma $i $.Requires | not}},{{end}}{{end}}
      ],
      "version": "{{.Version}}",
      "env":[ {{range $i, $entry := .Env}}
        {"variable": "{{$entry.Variable}}", "value":"{{$entry.Value}}"}{{if comma $i $.Env | not}},{{end}}{{end}}
      ],
      "expose": [{{range $key, $value := .Expose}}
        {{$value}}{{if comma $key $.Expose | not}},{{end}}{{end}}
      ],
      "ignore": [{{range $key, $value := .Ignore}}
        "{{$value}}"{{if comma $key $.Ignore | not}},{{end}}{{end}}
      ],
      "add": {
        "simple": [{{range $key, $value := .Add.Simple}}
          "{{$value}}"{{if comma $key $.Add.Simple | not}},{{end}}{{end}}
        ],
        "complex": [ {{range $key, $value := .Add.Complex}}
          { "name": "{{$value.Name}}",
            "location": "{{$value.Location}}" }{{if comma $key $.Add.Complex | not}},{{end}}{{end}}
        ]
      },
      "cmd": "{{.Cmd}}",
      "launchCmd": [ {{range $key, $value := .LaunchCmd}}
        "{{$value}}"{{if comma $key $.Add.Simple | not}},{{end}}{{end}}
      ],
      "workDir": "{{.WorkDir}}",
      "runTests": {{.RunTests}},
      "resources": {
        "cpu": {{.Resources.Cpu}},
        "mem": {{.Resources.Mem}}
      }
    }`,
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
        tmpl := template.Must(template.New("yaml-" + string(i)).
            Funcs(libtesting.FuncMap).
            Parse(testData.fileDataTemplate))
        err := tmpl.Execute(mockFile, testData.expectedBuildFile)
        if err != nil {
            t.Error("Test failed, could not execute template.", err)
        }

        mockFile.Flush()

        if err := buildFile.ParseYaml(mockFileBuffer.Bytes()); err != nil {
            t.Error("Could not parse yaml file!")
            t.Log(mockFileBuffer.String())
        } else {
            if !reflect.DeepEqual(buildFile, testData.expectedBuildFile) {
                t.Error("Parsed yaml but values not what was expected!")
                t.Log("*********buildfile*********")
                t.Log(buildFile)
                t.Log("*********expected*********")
                t.Log(testData.expectedBuildFile)
            } else { // if equals, output the templated file that was parsed
                // t.Log("*******Successfully Parsed Reference File***********")
                // t.Log(mockFileBuffer.String())
                // t.Log("************End proper parsed yaml file*************")
            }
        }
    }
}

func TestParseJson(t *testing.T) {

    // Loop through test data and generate templates.
    // Templates help write better tests.
    for i, testData := range JsonTestData {
        var buildFile BuildFile
        mockFileBuffer := bytes.NewBuffer(nil)      // backing buffer
        mockFile := bufio.NewWriter(mockFileBuffer) // soon to be realized template

        // Template the file data
        tmpl := template.Must(template.New("json-" + string(i)).
            Funcs(libtesting.FuncMap).
            Parse(testData.fileDataTemplate))
        err := tmpl.Execute(mockFile, testData.expectedBuildFile)
        if err != nil {
            t.Error("Test failed, could not execute template.", err)
        }

        mockFile.Flush()

        if err := buildFile.ParseJson(mockFileBuffer.Bytes()); err != nil {
            t.Error("Could not parse json file!")
            t.Log(mockFileBuffer.String())
        } else {
            if !reflect.DeepEqual(buildFile, testData.expectedBuildFile) {
                t.Error("Parsed json but values not what was expected!")
                t.Log("*********buildfile*********")
                t.Log(buildFile)
                t.Log("*********expected*********")
                t.Log(testData.expectedBuildFile)
            } else { // if equals, output the templated file that was parsed
                // t.Log("*******Successfully Parsed Reference File***********")
                // t.Log(mockFileBuffer.String())
                // t.Log("************End proper parsed yaml file*************")
            }
        }
    }
}

// Test the default build file values
func TestDefaultBuildFile(t *testing.T) {
    actual := NewBuildFile()
    expected := &BuildFile{
        RunTests:  false,
        Owner:     "Test-Flight-User",
        ImageName: "Test-Flight-Test-Image",
        Tag:       "latest",
        Version:   "0.0.1",
        WorkDir:   "/tmp/build",
        Add:       DockerAdd{},
        Ignore:    []string{".git"},
    }

    if !reflect.DeepEqual(actual, expected) {
        t.Error("Default buildfile incorrect, values not expected!")
        t.Log("*********buildfile*********")
        t.Log(actual)
        t.Log("*********expected*********")
        t.Log(expected)
    }
}
