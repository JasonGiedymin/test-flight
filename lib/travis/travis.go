package travis

type TravisBranches struct {
    Except []string // blacklist
    Only   []string // whitelist
}

//
//  matrix:
//  exclude:
//      - rvm: 2.0.0
//        gemfile: Gemfile
//
type TravisBuildMatrix struct {
    Exclude []map[string]string
}

// multiple types :-(
type TravisEnvironment struct {
    Global []string
    Matrix []string
}

type TravisFile struct {
    // general
    Branches      TravisBranches
    Env           []string // env vars, must also parse on '='
    Before_script []string // add to cmd
    After_script  []string // add to cmd
    Script        string   // add to cmd

    // runtimes
    Otp_release []string // erlang
    Language    string   // groovy, java
    Nodejs      []string // nodejs
    Perl        []string // perl
    Php         []string //php
    Python      []string //python
    Scala       []string //scala

    // Ruby
    Ruby    []string //ruby
    Rvm     []string //ruby via rvm
    Gemfile []string // gemfile

}
