package build

import (
    "github.com/JasonGiedymin/test-flight/lib/build"
    "testing"
)

var testData = struct {
    Vectors        build.BuildMatrixVectors
    ExpectedMatrix BuildMatrix
}{
    Vectors: build.BuildMatrixVectors{
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
}

func testMatrixSets(t *testing.T) {
    t.Error("Pending")
}

func testMatrixConstruction(t *testing.T) {
    t.Error("Pending")
}
