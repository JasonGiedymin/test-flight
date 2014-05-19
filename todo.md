## Dev Todo

Dev todos:

### v0.9.0 - Alpha

    [x] read files
    [x] cli params {check, launch}
    [x] wire in check files to check & launch command
    [x] wire in param to specify dir `flight check` or `flight check -d ./thisDir`.
        when specifying nothing, imply `./`
    [-] wire in dir to all commands, shared struct? => Nope. Bad go practice.
    [-] look at golang code, find examples of composition + errors. => if err tracks are the norm... oh well.
    [x] read build.json file
    [x] look for `main.yml` file under each dir
    [x] look for all dirs and all files required (vault, templates, docker, inventory file...)
    [x] modify for better oop design
    [x] add app state obj
    [x] modify app state to be populated build file
    [x] read user home config file then local to bin
    [ ] cleanup `ReadConfigFile`, my gosh is it ugly
    [ ] better logging, visit spacemonkey libs
    [ ] find better golang flow (perf functional)
    [ ] create test-flight file in lib/
    [ ] move parsing into test-flight within lib
    [ ] docker endpoint specified in test-flight config & build.json (build overrides config)
    [ ] check that docker socket/ip exists
    [ ] create dockerfile and pass to docker
    [ ] run ansible lint/check within a docker, call this `runup` (an actual aircraft term)

    [ ] log file

### v0.9.5 - Beta release

    [ ] build.json - add entries to specify list of tests to run and their order
    [ ] test-flight config.json - add entry to specify run timeout

### v1.0.0

    [ ] user feedback...

### v1.1.0

    [ ] test-flight config - keep track of dockerfiles on filesystem that it built
    [ ] test-flight config - entries to enable/disable dependency building (all)
    [ ] test-flight config - entries to enable/disable dependency building (parent)
    [ ] test-flight config - entries to enable/disable dependency building (children)
    [ ] test-flight launch command sub-command to disable dependency building of children
    [ ] build.json - add entries for next docker to build
    [ ] build.json - add triggers to build dependent dockers

### v1.5.0

    [ ] atom apm plugin
    [ ] test-flight UI to report status/state

## Production Todo

Before heading to production.

Prod todos:

    [ ] silence full stack trace? Only allow in dev?
    [ ] docs
    [ ] more tests :-)
