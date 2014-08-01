package build

import (
    // "github.com/JasonGiedymin/test-flight/lib/build"
    "reflect"
    "sort"
    "testing"
)

var testData = []struct {
    Vectors        BuildMatrixVectors
    ExpectedMatrix []string
}{
    {
        BuildMatrixVectors{
            Language: "golang",
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
            "(golang,1,trash)",
            "(golang,1,can)",
            "(golang,1,car)",
            "(golang,2,trash)",
            "(golang,2,can)",
            "(golang,2,car)",
            "(golang,3,trash)",
            "(golang,3,can)",
            "(golang,3,car)",
        },
    },
    {
        BuildMatrixVectors{
            Language: "c++",
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
            "(c++,4.6.0,TMPDIR=/tmp/trash)",
            "(c++,4.6.0,LIBRARY_PATH=/usr/lib/my.lib.path/)",
            "(c++,4.6.1,TMPDIR=/tmp/trash)",
            "(c++,4.6.1,LIBRARY_PATH=/usr/lib/my.lib.path/)",
            "(c++,4.6.2,TMPDIR=/tmp/trash)",
            "(c++,4.6.2,LIBRARY_PATH=/usr/lib/my.lib.path/)",
        },
    },
}

func TestMatrixConstruction(t *testing.T) {
    for _, data := range testData {
        expected := data.ExpectedMatrix
        actual := func() []string {
            var out []string
            for key, _ := range data.Vectors.Product() {
                out = append(out, key)
            }
            sort.Strings(out)
            return out
        }()

        sort.Strings(expected)

        if !reflect.DeepEqual(expected, actual) {
            t.Errorf("Sets() failed, \nactual: \n%s\n\nexpected: \n%s", actual, expected)
        }
    }
}
