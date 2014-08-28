package lib

import (
    "reflect"
    "testing"
)

// take buildfile, construct build matrix out of it
// for each matrix have a pointer to the original build
// file
// TODO: last here
func TestCreateBuildMatrixFromBuildFile(t *testing.T) {
    var testData = []struct {
        sampleBuildFile    BuildFile
        expectedMatrixKeys []string
    }{
        {   // n1 as-is buildfile
            BuildFile{
                From: "c++",
                Env: []DockerEnv{
                    DockerEnv{
                        Variable: "TMPDIR",
                        Value:    "/tmp/trash",
                    },
                },
            },
            []string{"(c++,TMPDIR=/tmp/trash)"},
        },
        {   // n2
            BuildFile{
                OS:              []string{"ubuntu"},
                Language:        "c++",
                LanguageVersion: "4.6.1",
            },
            []string{
                "(c++,ubuntu,4.6.1)",
            },
        },
        {   // n3
            BuildFile{
                OS:              []string{"centos", "ubuntu"},
                Language:        "c++",
                LanguageVersion: "4.6.1",
            },
            []string{
                "(c++,centos,4.6.1)", "(c++,ubuntu,4.6.1)",
            },
        },
    }

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
