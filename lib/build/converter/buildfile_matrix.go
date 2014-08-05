package converter

import (
    "github.com/JasonGiedymin/test-flight/lib"
    "github.com/JasonGiedymin/test-flight/lib/build"
)

// 1. Take buildfile
// 1. create vectors
// 1. create build matrix
func ConvertBuildFileToMatrix(buildfile lib.BuildFile) build.BuildMatrix {
    vectors := build.BuildMatrixVectors{
        // Language: // would have to take from buildfile
        // Version:
        From: []string{buildfile.From},
        Env:  buildfile.Env,
    }

    return vectors.Product()
}
