package build

import (
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

type BuildMatrixEntry struct {
    Language string
    Version  string
    Env      string
}

func (e BuildMatrixEntry) Key() string {
    return "(" + strings.Join([]string{
        e.Language,
        e.Version,
        e.Env,
    }, ",") + ")"
}

// Basic Build Matrix in Type form
type BuildMatrixVectors struct {
    // base
    Language string
    Version  []string
    Env      []string

    // custom
}

func (v *BuildMatrixVectors) Product() BuildMatrix {
    var matrix = make(BuildMatrix)

    for _, version := range v.Version {
        entry := BuildMatrixEntry{
            Language: v.Language,
        }

        entry.Version = version

        for _, env := range v.Env {
            entry.Env = env
            matrix[entry.Key()] = entry
        }

    }

    return matrix
}
