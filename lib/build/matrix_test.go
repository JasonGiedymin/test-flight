package build

import (
    "github.com/JasonGiedymin/test-flight/lib"
    "reflect"
    "sort"
    "testing"
)

var testData = []struct {
    Vectors            BuildMatrixVectors
    ExpectedMatrixKeys []string
    SampleBuildFile    lib.BuildFile
}{
    {
        BuildMatrixVectors{
            Language: "golang",
            From:     []string{"ubuntu"},
            Version: []string{
                "1",
                "2",
                "3",
            },
            Env: []string{
                "trash",
                "can",
                "car",
            },
        },
        []string{
            "(golang,ubuntu,1,trash)",
            "(golang,ubuntu,1,can)",
            "(golang,ubuntu,1,car)",
            "(golang,ubuntu,2,trash)",
            "(golang,ubuntu,2,can)",
            "(golang,ubuntu,2,car)",
            "(golang,ubuntu,3,trash)",
            "(golang,ubuntu,3,can)",
            "(golang,ubuntu,3,car)",
        },
        *lib.NewBuildFile(),
    },
    {
        BuildMatrixVectors{
            Language: "c++",
            From:     []string{"ubuntu"},
            Version: []string{
                "4.6.0",
                "4.6.1",
                "4.6.2",
            },
            Env: []string{
                "TMPDIR=/tmp/trash",
                "LIBRARY_PATH=/usr/lib/my.lib.path/",
            },
        },
        []string{
            "(c++,ubuntu,4.6.0,TMPDIR=/tmp/trash)",
            "(c++,ubuntu,4.6.0,LIBRARY_PATH=/usr/lib/my.lib.path/)",
            "(c++,ubuntu,4.6.1,TMPDIR=/tmp/trash)",
            "(c++,ubuntu,4.6.1,LIBRARY_PATH=/usr/lib/my.lib.path/)",
            "(c++,ubuntu,4.6.2,TMPDIR=/tmp/trash)",
            "(c++,ubuntu,4.6.2,LIBRARY_PATH=/usr/lib/my.lib.path/)",
        },
        *lib.NewBuildFile(),
    },
}

func TestMatrixConstruction(t *testing.T) {
    for _, data := range testData {
        expected := data.ExpectedMatrixKeys
        actual := data.Vectors.Product().Keys()

        sort.Strings(expected)

        if !reflect.DeepEqual(expected, actual) {
            t.Errorf("Sets() failed, \nactual: \n%s\n\nexpected: \n%s", actual, expected)
        }
    }
}
