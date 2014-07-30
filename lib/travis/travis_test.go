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
    Compiler: []string{
        "clang",
        "gcc",
    },
    Services: []string{
        "riak",
        "rabbitmq",
    },
    Language: "cpp",
    Jdk: []string{
        "oraclejdk7",
        "openjdk6",
    },
    Ghc: []string{
        "7.6",
        "7.4",
    },
    Otp_release: []string{
        "17.0",
        "R16B03-1",
        "R16B03",
        "R16B02",
        "R16B01",
        "R15B03",
        "R15B02",
        "R15B01",
        "R14B04",
        "R14B03",
        "R14B02",
    },
    Lein: "lein2",
    Scala: []string{
        "2.9.3",
        "2.10.4",
        "2.11.0",
    },
    Perl: []string{
        "5.16",
        "5.14",
        "5.12",
        "5.10",
        "5.8",
    },
    Php: []string{
        "5.5",
        "5.4",
        "hhvm",
    },
    Python: []string{
        "2.6",
        "2.7",
        "3.2",
        "3.3",
    },
    Go: []string{
        "1.0",
        "1.1",
        "1.2",
        "1.3",
        "tip",
    },
    Nodejs: []string{
        "0.11",
        "0.10",
        "0.8",
        "0.6",
    },
    Rvm: []string{
        "2.1.0",
        "jruby-18mode",
        "jruby-19mode",
        "rbx-2",
        "ruby-head",
        "jruby-head",
        "ree",
    },
    Gemfile: []string{
        "Gemfile",
        "gemfiles/eventmachine-pre",
    },
    Bundler_args:    "--without production",
    Xcode_project:   "MyNewProject.xcodeproj",
    Xcode_workspace: "MyNewProject.workspace",
    Xcode_scheme:    "MyNewProjectTests",
    Xcode_sdk:       "iphonesimulatorX.Y",
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
      compiler: {{range $key, $value := .Compiler}}
        - {{$value}}{{end}}
      services: {{range $key, $value := .Services}}
        - {{$value}}{{end}}
      language: {{.Language}}
      jdk: {{range $key, $value := .Jdk}}
        - {{$value}}{{end}}
      ghc: {{range $key, $value := .Ghc}}
        - {{$value}}{{end}}
      otp_release: {{range $key, $value := .Otp_release}}
        - {{$value}}{{end}}
      lein: {{.Lein}}
      scala: {{range $key, $value := .Scala}}
        - {{$value}}{{end}}
      perl: {{range $key, $value := .Perl}}
        - {{$value}}{{end}}
      php: {{range $key, $value := .Php}}
        - {{$value}}{{end}}
      python: {{range $key, $value := .Python}}
        - {{$value}}{{end}}
      go: {{range $key, $value := .Go}}
        - {{$value}}{{end}}
      nodejs: {{range $key, $value := .Nodejs}}
        - {{$value}}{{end}}
      rvm: {{range $key, $value := .Rvm}}
        - {{$value}}{{end}}
      gemfile: {{range $key, $value := .Gemfile}}
        - {{$value}}{{end}}
      bundler_args: {{.Bundler_args}}
      xcode_project: {{.Xcode_project}}
      xcode_workspace: {{.Xcode_workspace}}
      xcode_scheme: {{.Xcode_scheme}}
      xcode_sdk: {{.Xcode_sdk}}
    `,
        standardTravisFile,
    },
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
