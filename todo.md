## Dev Todo

Dev todos:

    [x] read files
    [x] cli params {check, launch}
    [x] wire in check files to check & launch command
    [x] wire in param to specify dir `flight check` or `flight check -d ./thisDir`.
        when specifying nothing, imply `./`
    [-] wire in dir to all commands, shared struct? => Nope. Bad go practice.
    [-] look at golang code, find examples of composition + errors. => if err tracks are the norm... oh well.
    [x] read build.json file
    [ ] look for `main.yml` file under each dir
    [ ] run ansible lint/check within a docker, call this `runup` (an actual aircraft term)
    [ ] better logging, visit spacemonkey libs

## Production Todo

Before heading to production.

Prod todos:

    [ ] silence full stack trace? Only allow in dev?
    [ ] docs
    [ ] more tests :-)
