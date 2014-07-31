package travis

import (
    "errors"
    "gopkg.in/yaml.v1"
)

// Travis Build Matrix
//  matrix:
//  exclude:
//      - rvm: 2.0.0
//        gemfile: Gemfile
//
type TravisBuildMatrix struct {
    Exclude []map[string]string
}

type TravisBranches struct {
    Except []string // blacklist
    Only   []string // whitelist
}

// multiple types :-(
type TravisEnvironment struct {
    Global []string
    Matrix []string
}

// Android components to be used by Travis
type TravisRuntimeAndroid struct {
    Components []string
    Licenses   []string
}

type TravisFile struct {
    // General Flow
    Branches      TravisBranches
    Env           []string // env vars, must also parse on '='
    Before_script []string // add to cmd
    After_script  []string // add to cmd
    Script        string   // add to cmd

    // Support
    Compiler []string // compiler list [unused] for us but used for c/c++
    Services []string // various services to use in tests (ruby)

    // Runtime List
    // android, c, c++, clojure, erlang, go, groovy, java,
    // nodejs, objective-c, perl, php, python, ruby, scala
    Language string

    // Java
    Jdk []string // java versions

    // Functional
    Ghc         []string // haskell
    Otp_release []string // erlang version (should also use language field)
    Lein        string   // clojure
    Scala       []string //scala versions

    // CLI
    Perl []string // perl versions
    Php  []string //php versions

    // Where the grass is greener - and it is lush
    Python []string //python versions
    Go     []string // go and go versions
    Nodejs []string // nodejs versions

    // Ruby
    Ruby         []string //ruby versions
    Rvm          []string //ruby versions via rvm
    Gemfile      []string // gemfile
    Bundler_args string   // bundler

    // Objective-c
    Xcode_project   string // instead of workspace
    Xcode_workspace string // instead of project
    Xcode_scheme    string
    Xcode_sdk       string // iphonesimulatorX.Y

    // Android
    Android TravisRuntimeAndroid // Specific for android
}

func (tf *TravisFile) ParseYaml(data []byte) error {
    if err := yaml.Unmarshal(data, &tf); err != nil {
        msg := "Could not parse yaml build file." + err.Error()
        return errors.New(msg)
    }

    return nil
}
