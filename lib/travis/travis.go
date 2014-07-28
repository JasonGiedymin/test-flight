package travis

type TravisFile struct {
    // general
    script        string   // add to cmd
    before_script []string // add to cmd
    after_script  []string // add to cmd
    // runtimes
    otp_release []string // erlang
    language    string   // groovy, java
    nodejs      []string // nodejs
    perl        []string // perl
    php         []string //php
    python      []string //python
    ruby        []string //ruby
    scala       []string //scala
}
