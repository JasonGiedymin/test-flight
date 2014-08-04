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
    {   // n1
        lib.BuildFile{
            From: "c++",
            Env: []lib.DockerEnv{
                lib.DockerEnv{
                    Variable: "TMPDIR",
                    Value:    "/tmp/trash",
                },
            },
        },
        []string{
            "(c++,,TMPDIR=/tmp/trash)",
        },
    },
    // { // n2
    //     lib.BuildFile{
    //         From: "c++",
    //     },
    //     []string{
    //         "(c++,ubuntu,4.6.0,TMPDIR=/tmp/trash)",
    //     },
    // },
    // { // n3
}

// take buildfile, construct build matrix out of it
// for each matrix have a pointer to the original build
// file
func TestCreateBuildMatrixFromBuildFile(t *testing.T) {
    for _, data := range testData {
        buildFile := data.sampleBuildFile
        expected := data.expectedMatrixKeys
        actual := ConvertBuildFileToMatrix(buildFile)

        if !reflect.DeepEqual(expected, actual) {
            t.Errorf("Conversion of Buildfile to Matrix failed, \nactual: \n%s\n\nexpected: \n%s", actual, expected)
        }

        t.Log(buildFile)
    }

    t.Error("PENDING")
}
