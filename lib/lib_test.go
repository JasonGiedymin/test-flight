package lib

import (
	"fmt"
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
	} else {
		return badFiles
	}
}

// ====

func TestRequiredFiles(t *testing.T) {
	expectedSize := 6
	fmt.Println("Testing required files...")
	if len(MockFileIO(true)) != expectedSize {
		t.Error("Expected a array slice of length", expectedSize)
	}
}

func TestCheckFiles(t *testing.T) {
	result, err := CheckFiles(MockFileIO(true))
	if result != true {
		t.Error("Expected the MockFileIO good set of files to be valid. Error was", err)
	}
}
