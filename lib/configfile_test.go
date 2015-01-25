package lib

import (
    "reflect"
    "testing"
)

// TODO: Finish this
func TestConfigPaths(t *testing.T) {
    testData := []struct {
        dirFlag     string
        userHomeDir string
        expected    []string
    }{
        {
            "/Users/me/projects/somedir",
            "/Users/me",
            []string{
                "/Users/me/projects/somedir/test-flight-config.json", // dirFlag
                "./test-flight-config.json",                          // current working dir
                "/Users/me/test-flight-config.json",                  // user home
            },
        },
    }

    for _, test := range testData {
        actual := configPaths(test.dirFlag, test.userHomeDir)
        if !reflect.DeepEqual(test.expected, actual) {
            t.Errorf("Config path generation for lookup failed, \nactual: \n%s\n\nexpected: \n%s", actual, test.expected)
            t.Logf("buildfile:\n%s\n", buildFile)
        }
    }
}
