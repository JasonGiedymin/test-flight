package travis

import (
    "bufio"
    "bytes"
    libtesting "github.com/JasonGiedymin/test-flight/lib/testing"
    "gopkg.in/yaml.v1"
    "reflect"
    "testing"
    "text/template"
)

var standardTravisFile = TravisFile{
    Branches: TravisBranches{
        Except: []string{
            "legacy",
            "experimental",
        },
        Only: []string{
            "master",
            "stable",
        },
    },
    Env: []string{
        "FOO=foo BAR=bar",
        "FOO=bar BAR=foo",
    },
    Before_script: []string{
        "before_command_1",
        "before_command_2",
    },
    After_script: []string{
        "after_command_1",
        "after_command_2",
    },
    Script: "./script/ci/run_build.sh",
}

var YamlTestData = []struct {
    fileDataTemplate   string
    expectedTravisFile TravisFile
}{
    {`
      branches:
        except: {{range $key, $value := .Branches.Except}}
          - {{$value}}{{end}}
        only: {{range $key, $value := .Branches.Only}}
          - {{$value}}{{end}}
      env: {{range $key, $value := .Env}}
        - {{$value}}{{end}}
      before_script: {{range $key, $value := .Before_script}}
        - {{$value}}{{end}}
      after_script: {{range $key, $value := .After_script}}
        - {{$value}}{{end}}
      script: {{.Script}}
    `,
        standardTravisFile,
    },
    //   owner: {{.Owner}}
    //   imageName: {{.ImageName}}
    //   tag: {{.Tag}}
    //   from: {{.From}}
    //   requires: {{range $key, $value := .Requires}}
    //     - {{$value}}{{end}}
    //   version: {{.Version}}
    //   env: {{range $i, $entry := .Env}}
    //     - variable: {{$entry.Variable}}
    //       value: {{$entry.Value}}{{end}}
    //   expose: {{range $key, $value := .Expose}}
    //     - {{$value}}{{end}}
    //   ignore: {{range $key, $value := .Ignore}}
    //     - {{$value}}{{end}}
    //   add:
    //     simple: {{range $key, $value := .Add.Simple}}
    //     - {{$value}}{{end}}
    //     complex: {{range $key, $value := .Add.Complex}}
    //       - name: {{$value.Name}}
    //         location: {{$value.Location}}{{end}}
    //   cmd: {{.Cmd}}
    //   launchCmd: {{range $key, $value := .LaunchCmd}}
    //     - {{$value}}{{end}}
    //   workDir: {{.WorkDir}}
    //   runTests: {{.RunTests}}
    //   resources:
    //     cpu: {{.Resources.Cpu}}
    //     mem: {{.Resources.Mem}}
}

func TestTravisJson(t *testing.T) {
    // Loop through test data and generate templates.
    // Templates help write better tests.
    for i, testData := range YamlTestData {
        var travisFile TravisFile
        mockFileBuffer := bytes.NewBuffer(nil)      // backing buffer
        mockFile := bufio.NewWriter(mockFileBuffer) // soon to be realized template

        // Template the file data
        tmpl := template.Must(template.New("yaml-" + string(i)).
            Funcs(libtesting.FuncMap).
            Parse(testData.fileDataTemplate))
        err := tmpl.Execute(mockFile, testData.expectedTravisFile)
        if err != nil {
            t.Error("Test failed, could not execute template.", err)
        }

        mockFile.Flush()

        if err := yaml.Unmarshal(mockFileBuffer.Bytes(), &travisFile); err != nil {
            t.Error("Could not parse yaml file!")
            t.Log(mockFileBuffer.String())
        } else {
            if !reflect.DeepEqual(travisFile, testData.expectedTravisFile) {
                t.Error("Parsed yaml but values not what was expected!")
                t.Log("*********buildfile*********")
                t.Log(travisFile)
                t.Log("*********expected*********")
                t.Log(testData.expectedTravisFile)
            } else { // if equals, output the templated file that was parsed
                t.Log("*******Successfully Parsed Reference File***********")
                t.Log(mockFileBuffer.String())
                t.Log("************End proper parsed yaml file*************")
            }
        }
    }
}
