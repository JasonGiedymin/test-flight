package lib

import (
  // "fmt"
  // "io/ioutil"
  "testing"
)

// == Setup ==
var goodFiles = []string{
  "build.json",
  "docker",
  "meta",
  "tasks",
  "tests",
  "vars",
}

var badFiles = []string{
  "build.json",
}

func MockFileIO(good bool) []string {
  if good {
    return goodFiles
  }

  return badFiles
}

// ====

func TestRequiredFiles(t *testing.T) {
  expectedSize := 6

  if len(MockFileIO(true)) != expectedSize {
    t.Error("Expected a array slice of length", expectedSize)
  }

  if len(MockFileIO(false)) != 1 {
    t.Error("Expected single file in array.")
  }
}

func TestCheckFiles(t *testing.T) {
  result, err := CheckFiles(MockFileIO(true))
  if result == false { // only error if not expected
    t.Error("Expected the MockFileIO good set of files to be valid. Error was", err)
  }

  result, err = CheckFiles(MockFileIO(false))
  if result == true { // only error if not expected
    t.Error("Expected the MockFileIO good set of files to be valid. Error was", err)
  }

  // } else {
  // 	fmt.Println(err)
  // }
}
