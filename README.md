Ansible Shipyard Test-Flight
----------------------------

Test + Build your Ansible Playbooks with Docker.

## Index
- [Production](#production)
- [Develeopment](#development)
- [Testing](#testing)
- [Notes](#notes)
- [Building](#building)

## Usage

### Production Build <a name="production"></a>

Once built, test-flight can be used like so:

    # To run test-flight in the current directory
    flight launch

    # Or
    flight launch .

    # To run test-flight in the another directory named 'test'
    flight launch -d ./test


### Development Build <a name="development"></a>

During development these commands can be used:

    # Run Check
    clear && go run flight.go check

    # Or In another directory
    clear && go run flight.go check -d ./test

    # Run Launch
    clear && go run flight.go launch

    # Or In another directory
    clear && go run flight.go launch -d ./test


## Building

Note that the binary which is created is a complete executable which includes
not only the application but a built in GC, scheduler, etc... This is why the
file size is large. Once running, the go app will NOT take up much memory.

One can build the application like so

    go build

For production only build, without debug

    go build -ldflags "-s"

Further optimizations can be done using UPX, for go one needs goupx. This
can be obtained by obtaining it like so:

    go get github.com/pwaller/goupx/


## Testing<a name="testing"></a>

Run the following:

    go test ./...


## Notes<a name="notes"></a>

  - deps managed by [godev](https://github.com/tools/godep).
    Installed via `go get github.com/tools/godep`
  - uses:
    - [golint](go get github.com/golang/lint/golint)
    - [go-dockerclient](https://github.com/fsouza/go-dockerclient)
    - [SpaceMonkeyGo Errors](https://github.com/SpaceMonkeyGo/errors)
    - [Go Flags](https://github.com/jessevdk/go-flags)
    - [FactorLog](https://github.com/kdar/factorlog)
    - [osext](https://bitbucket.org/kardianos/osext/src)
  - go get:
    - go get github.com/golang/lint/golint
    - go get github.com/fsouza/go-dockerclient
    - go get github.com/SpaceMonkeyGo/errors
    - go get github.com/jessevdk/go-flags
    - go get github.com/kdar/factorlog
    - go get bitbucket.org/kardianos/osext
  - update libs:
    - go get -u all

## Building<a name="building"></a>

- gofmt settings: `-w -s -tabs=false -tabwidth=2`
