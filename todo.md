## Dev Todo

Version todos:

### v0.9.0 - Alpha

  - [x] read files
  - [x] cli params {check, launch}
  - [x] wire in check files to check & launch command
  - [x] wire in param to specify dir `flight check` or `flight check -d ./thisDir`.
        when specifying nothing, imply `./`
  - [x] ~~wire in dir to all commands, shared struct? => Nope. Bad go practice.~~
  - [x] ~~look at golang code, find examples of composition + errors. => if err tracks are the norm... oh well.~~
  - [x] read build.json file
  - [x] look for `main.yml` file under each dir
  - [x] look for all dirs and all files required (vault, templates, docker, inventory file...)
  - [x] modify for better oop design
  - [x] add app state obj
  - [x] modify app state to be populated build file
  - [x] read user home config file then local to bin
  - [x] better logging, visit ~~spacemonkey libs~~ factorlog
  - [x] find better golang flow (perf functional) => custom functional every time
  - [x] apply logging levels => Error, Warn, Info, Debug, Trace
  - [x] create debug.log
  - [x] define usage of logging levels
  - [x] use logging where there is fmt.Println
  - [x] start linting (though the atom goplus plugin which calls golint seems to fail often)
  - [x] ~~create test-flight file in lib/~~
  - [x] ~~move parsing into test-flight within lib~~
  - [x] add version flag
  - [x] cleanup code from flight to lib
  - [x] add exit codes
  - [x] add more state to the config file reading
  - [x] remove logging line and class (since I proxy it's useless - for now)
  - [x] set state to trace logging level

#### v0.9.1 - Alpha

  - [x] change version
  - [x] fix bug - forgot to set configfile when reading
  - [x] docker endpoint specified in test-flight config & build.json (build overrides config)
  - [x] ~~check that docker socket/ip exists => client lib should do it~~
  - [x] cleanup docker api struct
  - [x] show docker images
  - [x] create dockerfile methods and placeholder
  - [x] placeholder for Dockerfile template
  - [x] placeholder for inventory template
  - [x] placeholder for playbook template
  - [x] move templates out
  - [x] documentation about building
  - [x] template var struct
  - [x] populate template var with test-flight-config.json + build.json
  - [x] add additional template fields required by template
  - [x] add better defaults via New() methods to conf/build files
  - [x] add build/config files directly into TemplateVar for complete access (helpful for user)
  - [x] start of `wiki/buildfile.md` which describes the makeup of the buildfile
  - [x] create template constructs for use with creating dockerfile
  - [x] extra error handling in cli & parse methods
  - [x] consolidate docker package pack into lib, for more consolidation in future

#### v0.9.2 - Alpha

  - [x] new repo
  - [x] version
  - [x] create template constructs for use with creating inventory and playbook
  - [x] First workable Dockerfile template
  - [x] create nested template for test-flight run command
  - [x] add .gitignore (better late than never :-))
  - [x] change config.user to config.complex, move user to build.json
  - [x] better test build.json (at least for what is available)
  - [x] clean up three add sections
  - [x] create nested template for other sections (add, expose, etc...)
  - [x] create friendlier template and nested template naming scheme
  - [x] wiki the naming concept
  - [x] add test-flight version in template
  - [x] fix meta in AppState
  - [x] create nested templates for inventory and playbook
  - [x] create nested templates in .test-flight/cache
  - [x] add test-flight version signature to all first level templates
  - [x] change docker dir name to test-flight (in reqs)
  - [x] fix `inventory`, and `playbook.yml` template not generating
  - [x] get dir of where test-flight binary is running from
  - [x] move more into types so refactor and pruning of structs can be done later
  - [x] generate test-flight files in `./.test-flight`
  - [x] invert test-flight struct, lib, and parser [refactor-parser]
  - [-] ~~slim down createTestTemplates()~~ Needs more intensive refactor of
        DockerApi which expects config and build files [refactor-dockerapi]
  - [x] bring things back into test-flight struct for state management [refactor-parser]
  - [x] move RequiredFiles into Types [refactor-types]

#### v0.9.3 - Alpha

  - [x] create flag to disable `./.test-flight` file generation
  - [x] command to just create templates `templates`
  - [x] slim down parser pre-reqs for parser commands, lots of repetition
  - [x] refactor cli template command as a function so they can be composed
  - [x] test byte buffer when creating Dockerfile
  - [x] create Dockerfile
  - [x] pass to docker api client
  - [x] start work on tar archiving the context dir
  - [x] if context dir files are sub dirs, recursively call archive func
  - [x] pass to docker api client successfully
  - [x] break out docker portion where archiving Dockerfile
  - [x] add channel event watcher for basic docker client events (start/die/etc...) ~~changes~~
  - [-] capture stdout from building => building is just building and I know when it fails.
        the next step should be running where there getting output is important (for now).
        However some build info would be great but that will have to wait.
  - [x] add logging messages describing build state

#### v0.9.4 - Alpha

  - [x] tag version
  - [x] create resource share type for buildfile
  - [x] start on creating container options
  - [x] Create container from docker image just created
  - [x] Rename `createDocker()`` to `createDockerImage()`` which is more accurate
  - [x] Log info about container created with slice of image hash
  - [ ] when image creation successful, pass name to container creation
  - [ ] pass container name to runner to run docker container
  - [ ] wire in buildfile resource share specs to container options
  - [ ] Run container from docker image just created

#### v0.9.5 - Alpha

  - [ ] run ansible lint/check within a docker, call this `runup` (an actual aircraft term)
  - [ ] allow docker diagnosis by preventing run and cmd commands, in buildfile
        give the flag: `debugContainer: {true|false}`
  - [ ] add help flag
  - [ ] remove AppState completely
  - [ ] cleanup code (lib is a mess, need more {})
  - [ ] test-flight verbose/debug via config and/or cli
  - [ ] create TestFlight.New()
  - [ ] dist/builds for various platforms
  - [ ] CI for builds and packaging
  - [ ] create/show test-flight specific images?
  - [ ] create simple map for docker hasFiles (I need functional programming!)
  - [ ] test-flight cleanup

### v0.9.5 - Beta release

  - [ ] build.json - add entries to specify list of tests to run and their order
  - [ ] test-flight config.json - add entry to specify run timeout
  - [ ] get rid of commandPreReq plz
  - [ ] template file newline cleanup

### v1.0.0

  - [ ] user feedback...
  - [ ] tests, tests, tests
  - [ ] some dirs should be optional, some mandatory
  - [ ] constrain test-flight-config info to sections, pass appropriate sections to libs
  - [ ] add more documentation along with sample build/config file usage (docs)
  - [ ] stop passing state everywhere [refactor-prime]
  - [ ] stop passing app everywhere, modularize [refactor-prime]

### v1.1.0

  - [ ] running list of containers created by test-flight for easy access to
        modify/delete etc...
  - [ ] user defined template vars
  - [ ] test-flight config - keep track of dockerfiles on filesystem that it built
  - [ ] test-flight config - entries to enable/disable dependency building (all)
  - [ ] test-flight config - entries to enable/disable dependency building (parent)
  - [ ] test-flight config - entries to enable/disable dependency building (children)
  - [ ] test-flight launch command sub-command to disable dependency building of children
  - [ ] build.json - add entries for next docker to build
  - [ ] build.json - add triggers to build dependent dockers
  - [ ] nested docker template sections for when creating Dockerfile would be nice

### v1.5.0

  - [ ] ansible npm repo - oh yeah...
  - [ ] atom apm plugin
  - [ ] test-flight UI to report status/state
  - [ ] sub-templates in each template, to start hooking into (for plugins later)
  - [ ] plugins for templates, config, build...


## Production Todo

Before heading to production.

Prod todos:

  - [ ] silence full stack trace? Only allow in dev?
  - [ ] docs
  - [ ] more tests :-)


## Major Refactor Tags

These tags represent major refactor intentions that will span multiple tasks. Each refactor
must be done within a defined section and each section must be done in sequence. There
is no particular order within the section.

Section 1:

  - [ ] refactor-types - Move all custom types to types package for future analysis and future refactors.
        But first must see what the lanscape is.


Section 2:

  - [x] refactor-parser - invert parser and test-flight so parser is not part of test-flight
  - [ ] refactor-prime - refactor so not passing so much app state around
  - [ ] refactor-dockerapi - refactor dockerApi so that it at least conforms to above and is slimmer
