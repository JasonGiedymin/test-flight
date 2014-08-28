package lib

// 1. Take buildfile
// 1. create vectors
// 1. create build matrix
func ConvertBuildFileToMatrix(buildfile BuildFile) BuildMatrix {
    var from, languageVersion []string

    if buildfile.From != "" {
        from = []string{buildfile.From}
    }

    if len(buildfile.LanguageVersion) == 0 {
        languageVersion = []string{"latest"}
    }

    vectors := BuildMatrixVectors{
        OS:       buildfile.OS,
        Language: buildfile.Language,
        Version:  languageVersion,

        From: from,

        Env: buildfile.Env,
    }

    return vectors.Product()
}
