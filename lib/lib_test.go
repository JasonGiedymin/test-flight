package lib

import (
	"fmt"
	"testing"
)

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

func TestRequiredFiles(t *testing.T) {
	fmt.Println("Testing required files...")
	fmt.Println(MockFileIO(true))
}
