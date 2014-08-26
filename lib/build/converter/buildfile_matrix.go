package converter

import (
    "github.com/JasonGiedymin/test-flight/lib"
    "github.com/JasonGiedymin/test-flight/lib/build"
)

// 1. Take buildfile
// 1. create vectors
// 1. create build matrix
func ConvertBuildFileToMatrix(buildfile lib.BuildFile) build.BuildMatrix {
    var from, version []string

    if buildfile.From != "" {
        from = []string{buildfile.From}
    }

    if buildfile.Version != "" {
        version = []string{buildfile.Version}
    }

    vectors := build.BuildMatrixVectors{
        OS:       buildfile.OS,
        Language: buildfile.Language,
        Version:  version,

        From: from,

        Env: buildfile.Env,
    }

    return vectors.Product()
}
