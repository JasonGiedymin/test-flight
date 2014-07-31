package build

import ()

// Basic combination
// Language x Version x Env
//
// | language |
// | c++      |
// | ruby     |

// Aim to create a string based matrix
type BuildMatrix [][]string

// Basic Build Matrix in Type form
type BuildMatrixVectors struct {
    // base
    Language string
    Version  []string
    Env      []string

    // custom

}

// Normalize the struct to a basic set of strings
// so that a generic cartesian product can be
// constructed. This is an implementation specific
// set, and is not generic. Thus it implies all
// data is unique :-O
func (v *BuildMatrixVectors) Sets() BuildMatrix {
    return BuildMatrix{
        []string{v.Language},
        v.Version,
        v.Env,
    }
}

// Cartesian Product
//
func (v *BuildMatrixVectors) Product() BuildMatrix {
    return BuildMatrix{}
}
