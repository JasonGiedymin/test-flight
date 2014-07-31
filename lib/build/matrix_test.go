package build

import (
    // "github.com/JasonGiedymin/test-flight/lib/build"
    "reflect"
    "testing"
)

var testData = []struct {
    Vectors        BuildMatrixVectors
    ExpectedMatrix BuildMatrix
    ExpectedSet    BuildMatrix
}{
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
        BuildMatrix{
            []string{
                "c++", "c++", "c++", "c++", "c++", "c++",
            },
            []string{
                "4.6.0", "4.6.0", "4.6.1", "4.6.1", "4.6.2", "4.6.2",
            },
            []string{
                "TMPDIR=/tmp/trash",
                "LIBRARY_PATH=/usr/lib/my.lib.path/",
                "TMPDIR=/tmp/trash",
                "LIBRARY_PATH=/usr/lib/my.lib.path/",
                "TMPDIR=/tmp/trash",
                "LIBRARY_PATH=/usr/lib/my.lib.path/",
            },
        },
        BuildMatrix{
            []string{"c++"},
            []string{
                "4.6.0",
                "4.6.1",
                "4.6.2",
            },
            []string{
                "TMPDIR=/tmp/trash",
                "LIBRARY_PATH=/usr/lib/my.lib.path/",
            },
        },
    },
}

func TestMatrixSets(t *testing.T) {
    for _, data := range testData {
        expected := data.ExpectedSet
        actual := data.Vectors.Sets()
        if !reflect.DeepEqual(expected, actual) {
            t.Errorf("Sets() failed, \nactual: %s\nexpected: %s", actual, expected)
        }
    }
}

func TestMatrixConstruction(t *testing.T) {
    for _, data := range testData {
        expected := data.ExpectedMatrix
        actual := data.Vectors.Product()
        if !reflect.DeepEqual(expected, actual) {
            t.Errorf("Sets() failed, \nactual: %s\nexpected: %s", actual, expected)
        }
    }
}
