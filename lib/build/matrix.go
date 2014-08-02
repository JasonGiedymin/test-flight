package build

import (
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

type BuildMatrixEntry struct {
    Language string
    From     string // always ubuntu for travis
    Version  string
    Env      string
}

func (e BuildMatrixEntry) Key() string {
    return "(" + strings.Join([]string{
        e.Language,
        e.From,
        e.Version,
        e.Env,
    }, ",") + ")"
}

// Basic Build Matrix in Type form
type BuildMatrixVectors struct {
    // base
    Language string
    From     []string
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

        for _, from := range v.From {
            entry.From = from
            for _, env := range v.Env {
                entry.Env = env
                matrix[entry.Key()] = entry
            }
        }
    }

    return matrix
}
