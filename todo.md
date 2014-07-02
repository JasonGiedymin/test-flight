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
  - [x] Rename `createDocker()` to `createDockerImage()` which is more accurate
  - [x] Log info about container created with slice of image hash
  - [x] Add tag info in dockerfile for when building images
  - [x] Start mixing in rest api directly into client
  - [x] Add tag name from buildfile to image name in url for rest call
  - [x] Add GetImage() func to get specific image (uses 3rd party)
  - [x] Rename GetImage() to GetImageDetails()
  - [x] Add ListContainers()
  - [x] Start on DeleteImage
  - [x] Add ApiContainer struct
  - [x] Add DeleteContainer()
  - [x] try http package instead of 3rd party http client: http://golang.org/pkg/encoding/json/#example_Unmarshal
        => need to replace more, but good for now
  - [x] showImage use image:tag param
  - [x] Check if container already exists, delete if so then create image
  - [~] __Work/Design Api-builder to dry api out (moved up in priority)__
    - [x] Base Api
    - [x] Base Api using Url
    - [x] Mock Server
    - [x] Add methods
    - [x] Mock Docker Test
    - [x] RC Release: [see here](https://github.com/JasonGiedymin/go-apibuilder)
  - [x] Makefile
  - [x] Finish on DeleteImage
  - [x] Redo CreateContainer with go lib (not go-apibuilder, yet)
  - [x] when image creation successful, pass name to container creation
  - [x] Add `make help`
  - [x] Add `make run-images`
  - [x] Add `make run-launch`
  - [x] Add startContainer
  - [x] remove some state commands [refactor-prime]
  - [-] ~~Add `flight delete image`~~ This is effectively a `destroy`.
  - [x] Start Refactor CLI, dry it up
  - [x] Modify ListContainers(ImageName) - to only get containers running "imageName"
  - [x] Add `flight ground` which will stop container` (stop container)
  - [x] Add slightly better messages during deletion
  - [x] Prevent creation if cannot delete
  - [x] Add `flight destroy ~~delete~~` which will do both container and image (del all)
    - [x] Add `destroy` command
    - [x] Wire in deletion
  - [x] Bug where Delete Container untags but doesn't delete =>
        call delete on container and image
  - [x] Remove sleeps
  - [x] Api out of sync with run-destroy => no return (tdd)
  - [x] Rename make file commands from `run-command` to `test-command`
  - [x] Bug when 'launch, ground, then destroy', cannot find to destroy => found
        that always returned running containers when calling ListContainers()
  - [x] Add `flight build` which will build the image (new image)
  - [x] Modify makefile
  - [x] Modify `flight launch` to only launch a new container when supplying the
        `-f` parameter.
  - [-] ~~pass container name to runner to run docker container~~ Seems to work
        fine as-is with Id. Would rather use monads.
  - [~] Modify config for stdin by adding attach api call
    - [x] Add `Attach` method
    - [x] Complete `Attach` so that it reads from the stream
    - [x] channel the console output => still waits on exit, for now
    - [x] wait on container channel => WaitGroup
  - [x] Fix interfact wrapping in logging messages => temporary fix with 
        v[0] (index 0), and only for console logging
  - [x] control-c watch => use channel and signal notify
  - [x] Add file mode command line option for building/launching `-f`, however
        don't code for it, only the flag

#### v0.9.5 - Alpha

  - [x] bug, filemode was a string, should have been a bool
  - [x] Add template command to makefile
  - [x] Add filemode to template command (though doesn't use it yet)
  - [x] change 'filemode' to 'single file mode' and `-s` , and leave force `-f`
  - [x] bug, destroy command is creating templates
  - [x] bug, `.test-flight` dir not being created
  - [x] Massive Refactor
    - [x] test-flight struct into separate file [refactor-types]
    - [x] flatten all dirs
    - [x] extract classes to single out functionality and testing,
          future task will be to break apart again
    - [x] make to work after crazy refactor => must be done one command at a
          time
    - [x] refactor CheckConfigs by passing CommandOptions instead
    - [x] build command, stop looking for config when it is specified
      - [x] make CheckConfigs return config file if `-c` specified
  - [x] Sub commands should be within Test-Flight Options?
  - [x] Remove exit code in Test-Flight
  - [x] Add config file command param `-c`
  - [~] but, filemode does not skip dir mode required files => the template
        is calling for a directory. Require a new Dockerfile template when
        in single file mode. Test command: `clear && make test-build-s`
     - [x] Add FilePath() to generate file paths
     - [x] Move templates to sub dir `dirmode` and `filemode`
     - [~] Modify CreateDockerImage() to generate dockerfile based on template
           which changes if set to filemode.
  - [ ] Change singlefilemode to just filemode
  - [ ] finish filemode on all commands (build, launch, template, abstract it)
  - [ ] Sync up & Refactor all commands (build should be latest)
  - [ ] Add RunDocker Tests and code using go-apibuilder
  - [ ] wire in buildfile resource share specs to container options
  - [ ] Run container from docker image just created
  - [ ] Finish refactor of CLI, dry it up
  - [ ] Tests and refactor
  - [ ] Rebuild DockerApi calls that use golang http with go-apibuilder
  - [ ] Refactor and check returns of APIs
  - [ ] Remove napping
  - [ ] Replace waits with channels
  - [ ] Update README with make commands
  - [ ] Add cleanup command (Removes images/tags with None)
  - [ ] endpoint timeout

#### v0.9.6 - Alpha

  - [ ] tests
  - [ ] api docs
  - [ ] refactor params for consistency
  - [ ] revisit all error codes, refactor
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

### v0.9.7 - Beta release

  - [ ] add check for minimal docker version allowed
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
  - [ ] refactor-dockerapi - refactor docker Api so that it no longer relies on go lib dockerclient
