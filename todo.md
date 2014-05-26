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
    [x] better logging, visit ~~spacemonkey libs~~ factorlog
    [x] find better golang flow (perf functional) => custom functional every time
    [x] apply logging levels => Error, Warn, Info, Debug, Trace
    [x] create debug.log
    [x] define usage of logging levels
    [x] use logging where there is fmt.Println
    [x] start linting (though the atom goplus plugin which calls golint seems to fail often)
    [-] create test-flight file in lib/
    [-] move parsing into test-flight within lib
    [x] add version flag
    [x] cleanup code from flight to lib
    [x] add exit codes
    [x] add more state to the config file reading
    [x] remove logging line and class (since I proxy it's useless - for now)
    [x] set state to trace logging level

#### v0.9.1 - Alpha

    [x] change version
    [x] fix bug - forgot to set configfile when reading
    [x] docker endpoint specified in test-flight config & build.json (build overrides config)
    [-] check that docker socket/ip exists => client lib should do it
    [x] cleanup docker api struct
    [x] show docker images
    [x] create dockerfile methods and placeholder
    [x] placeholder for Dockerfile template
    [x] placeholder for inventory template
    [x] placeholder for playbook template
    [x] move templates out
    [x] documentation about building
    [x] template var struct
    [x] populate template var with test-flight-config.json + build.json
    [x] add additional template fields required by template
    [x] add better defaults via New() methods to conf/build files
    [x] add build/config files directly into TemplateVar for complete access (helpful for user)
    [x] start of `wiki/buildfile.md` which describes the makeup of the buildfile
    [x] create template constructs for use with creating dockerfile
    [ ] finalize template
    [x] extra error handling in cli & parse methods
    [x] consolidate docker package pack into lib, for more consolidation in future

#### v0.9.2 - Alpha

    [ ] new repo
    [ ] version
    [ ] create template constructs for use with creating inventory and playbook
    [ ] change docker dir name to test-flight (in reqs)
    [ ] create templates for inventory and playbook
    [ ] create templates in .test-flight/cache
    [ ] create dockerfile and pass to docker api client
    [ ] run ansible lint/check within a docker, call this `runup` (an actual aircraft term)
    [ ] add help flag
    [ ] cleanup code (lib is a mess, need more {})
    [ ] constrain test-flight-config info to sections, pass appropriate sections to libs
    [ ] create TestFlight.New()
    [ ] create/show test-flight specific images?
    [ ] test-flight cleanup

### v0.9.5 - Beta release

    [ ] build.json - add entries to specify list of tests to run and their order
    [ ] test-flight config.json - add entry to specify run timeout
    [ ] get rid of commandPreReq plz
    [ ] template file newline cleanup

### v1.0.0

    [ ] user feedback...

### v1.1.0

    [ ] user defined template vars
    [ ] test-flight config - keep track of dockerfiles on filesystem that it built
    [ ] test-flight config - entries to enable/disable dependency building (all)
    [ ] test-flight config - entries to enable/disable dependency building (parent)
    [ ] test-flight config - entries to enable/disable dependency building (children)
    [ ] test-flight launch command sub-command to disable dependency building of children
    [ ] build.json - add entries for next docker to build
    [ ] build.json - add triggers to build dependent dockers
    [ ] nested docker template sections for when creating Dockerfile would be nice

### v1.5.0

    [ ] ansible npm repo - oh yeah...
    [ ] atom apm plugin
    [ ] test-flight UI to report status/state
    [ ] sub-templates in each template, to start hooking into (for plugins later)
    [ ] plugins for templates, config, build...

## Production Todo

Before heading to production.

Prod todos:

    [ ] silence full stack trace? Only allow in dev?
    [ ] docs
    [ ] more tests :-)
