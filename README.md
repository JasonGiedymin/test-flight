Ansible Shipyard Test-Flight
----------------------------

Test + Build your Ansible Playbooks with Docker.


## Usage

Once built, test-flight can be used like so:

    # To run test-flight in the current directory
    flight launch

    # Or
    flight launch .

    # To run test-flight in the another directory named 'test'
    clear && go run flight.go launch -d ./test

Other commands:

    # Run check
    clear && go run flight.go check

    # Or In another directory
    clear && go run flight.go check -d ./test


## Notes

  - deps managed by [godev](https://github.com/tools/godep).
    Installed via `go get github.com/tools/godep`
  - uses:
    - [golint](go get github.com/golang/lint/golint)
    - [go-dockerclient](https://github.com/fsouza/go-dockerclient)
    - [SpaceMonkeyGo Errors](https://github.com/SpaceMonkeyGo/errors)
    - [Go Flags](https://github.com/jessevdk/go-flags)
    - [FactorLog](github.com/kdar/factorlog)
  - go get:
    - go get github.com/fsouza/go-dockerclient
    - go get github.com/SpaceMonkeyGo/errors
    - go get github.com/jessevdk/go-flags


## Building

- gofmt settings: `-w -s -tabs=false -tabwidth=2`
