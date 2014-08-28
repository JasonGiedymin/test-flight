package lib

import (
    "errors"
    // "fmt"
    "sort"
    "strings"
)

// Basic combination
// Language x Version x Env
//
// | language |
// | c++      |
// | ruby     |

// Aim to create a string based matrix
type BuildMatrix map[string]BuildMatrixEntry

func (m BuildMatrix) Keys() []string {
    var out []string
    for key, _ := range m { // interestingly enough one could also call value.Key()
        out = append(out, key)
    }
    sort.Strings(out)
    return out
}

type BuildMatrixEntryResult struct {
    Entry  BuildMatrixEntry
    Stdout *string
    Err    error
}

type BuildMatrixEntryResults struct {
    Results []BuildMatrixEntryResult
}

func (r BuildMatrixEntryResults) Errors() []string {
    var resultErrors []string
    for _, v := range r.Results {
        if v.Err != nil {
            resultErrors = append(resultErrors, v.Err.Error())
        }
    }
    return resultErrors
}

type BuildMatrixFunc func(buildMatrixEntry BuildMatrixEntry, result BuildMatrixEntryResult) BuildMatrixEntryResult

func (m BuildMatrix) Map(fx BuildMatrixFunc) BuildMatrixEntryResults {
    var results []BuildMatrixEntryResult
    for _, e := range m {
        result := BuildMatrixEntryResult{Entry: e}
        results = append(results, fx(e, result))
    }

    return BuildMatrixEntryResults{results}
}

type BuildMatrixEntry struct {
    // base ----------------
    Language string
    OS       string
    Version  string
    Env      DockerEnv
    //----------------------
    From string // always ubuntu for travis
}

func (e BuildMatrixEntry) Key() string {
    var keys []string

    if e.From != "" { // `From` overrides
        keys = append(keys, e.From)
    } else { // use `Lan`, `OS`, `Ver`, to construct `From`
        keys = append(keys,
            e.Language,
            e.OS,
            e.Version,
        )
    }

    if e.Env.String() != "" {
        keys = append(keys, e.Env.String())
    }

    return "(" + strings.Join(keys, ",") + ")"
}

// Basic Build Matrix in Type form
type BuildMatrixVectors struct {
    // base ----------------
    Language string
    OS       []string
    Version  []string
    Env      []DockerEnv
    //----------------------
    From []string // can take precedence over Lang + Version, Legacy override

    // custom
}

// Will generate a product
// Uses vectors of input unless `From` is given. `From` is used
// to completely override a matrix. This is because the vectors
// used to build a matrix will in turn construct the final
// `From` automatically. So in this means you can supply the
// vectors to have it be generated, or you can directly supply
// it and override any behavior necessary.
func (v *BuildMatrixVectors) Product() BuildMatrix {
    var matrix = make(BuildMatrix)

    generate := func() {
        for _, version := range v.Version {
            entry := BuildMatrixEntry{
                Language: v.Language,
            }

            entry.Version = version

            for _, os := range v.OS {
                entry.OS = os
                if len(v.Env) > 0 {
                    for _, env := range v.Env {
                        entry.Env = env
                        matrix[entry.Key()] = entry // add to matrix
                    }
                } else { // no env!
                    matrix[entry.Key()] = entry
                }
            }
        }
    }

    override := func() {
        for _, from := range v.From {
            entry := BuildMatrixEntry{From: from}

            for _, env := range v.Env {
                entry.Env = env
                matrix[entry.Key()] = entry // add to matrix
            }
        }
    }

    if len(v.From) > 0 {
        override()
    } else {
        generate()
    }

    return matrix
}

func RunEntry(buildFile *BuildFile, handler BuildMatrixFunc, errorMessage string) error {
    matrix := ConvertBuildFileToMatrix(*buildFile)

    results := matrix.Map(handler)
    resultErrors := results.Errors()
    if len(resultErrors) > 0 {
        return errors.New(errorMessage + " " + strings.Join(resultErrors, ","))
    } else {
        return nil
    }
}
