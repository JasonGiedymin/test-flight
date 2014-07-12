package lib

/*
* Exit Codes:
*   0. Everything worked!
*   1. Misc undocumented error, yeah this shouldn't be.
*   2. Init failed
*   3. `test-flight-config.json` could not be found, required
*   4. Processing a command failed
 */
var ExitCodes = map[string]int{
    "success":        0,
    "misc":           1,
    "init_fail":      2,
    "config_missing": 3,
    "command_fail":   4,
    "docker_error":   5,
    "user_exit":      6,
    "system_error":   7,
}
