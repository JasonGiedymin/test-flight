package main

import (
	// "./lib"
	"./lib/cli"
	"fmt"
	"os"
	// "github.com/fsouza/go-dockerclient"
	// "flag"
	"github.com/jessevdk/go-flags"
)

/*
- require:
  - base image (based on build.json runtime selection)
*/

var checkCommand cli.CheckCommand

type LaunchCommand struct{}

type Options struct{}

var launchCommand LaunchCommand
var options Options
var parser = flags.NewParser(&options, flags.Default)

func init() {
	parser.AddCommand("check",
		"pre-flight check",
		"Used for pre-flight check of the ansible playbook.",
		&checkCommand)

	parser.AddCommand("launch",
		"flight launch",
		"Launch an ansible playbook test.",
		&launchCommand)
}

func main() {
	fmt.Println("Started...")

	// Verify required commands exist
	if _, err := parser.Parse(); err != nil {
		os.Exit(1)
	}

	// endpoint := "http://localhost:4243"
	// client, _ := docker.NewClient(endpoint)
	// imgs, _ := client.ListImages(true)

	// for _, img := range imgs {
	// 	fmt.Println("ID: ", img.ID)
	// 	fmt.Println("RepoTags: ", img.RepoTags)
	// 	fmt.Println("Created: ", img.Created)
	// 	fmt.Println("Size: ", img.Size)
	// 	fmt.Println("VirtualSize: ", img.VirtualSize)
	// 	fmt.Println("ParentId: ", img.ParentId)
	// 	fmt.Println("Repository: ", img.Repository)
	// }

	// files, filesError := lib.GetFiles()
	// fmt.Printf("Found these files %v, and these errors: [%v].\n", files, filesError)
}
