package converter

import (
    "github.com/JasonGiedymin/test-flight/lib"
    // "github.com/JasonGiedymin/test-flight/lib/build"
    "reflect"
    // "sort"
    "testing"
)

var testData = []struct {
    sampleBuildFile    lib.BuildFile
    expectedMatrixKeys []string
}{
    {   // n1 as-is buildfile
        lib.BuildFile{
            From: "c++",
            Env: []lib.DockerEnv{
                lib.DockerEnv{
                    Variable: "TMPDIR",
                    Value:    "/tmp/trash",
                },
            },
        },
        []string{"(c++,TMPDIR=/tmp/trash)"},
    },
    {   // n2
        lib.BuildFile{
            OS:       []string{"ubuntu"},
            Language: "c++",
            Version:  "4.6.1",
        },
        []string{
            "(c++,ubuntu,4.6.1)",
        },
    },
    {   // n3
        lib.BuildFile{
            OS:       []string{"centos", "ubuntu"},
            Language: "c++",
            Version:  "4.6.1",
        },
        []string{
            "(c++,centos,4.6.1)", "(c++,ubuntu,4.6.1)",
        },
    },
}

// take buildfile, construct build matrix out of it
// for each matrix have a pointer to the original build
// file
// TODO: last here
func TestCreateBuildMatrixFromBuildFile(t *testing.T) {
    for _, data := range testData {
        buildFile := data.sampleBuildFile
        expected := data.expectedMatrixKeys
        actual := ConvertBuildFileToMatrix(buildFile).Keys()

        // t.Logf("bf -> %v", buildFile)
        // t.Logf("conv -> %v", ConvertBuildFileToMatrix(buildFile))

        if !reflect.DeepEqual(expected, actual) {
            t.Errorf("Conversion of Buildfile to Matrix failed, \nactual: \n%s\n\nexpected: \n%s", actual, expected)
            t.Logf("buildfile:\n%s\n", buildFile)
        }
    }
}
