package main

import (
	// "./lib"
	"./lib/cli"
	"fmt"
	// "os"
	// "github.com/fsouza/go-dockerclient"
	// "flag"
	// "github.com/jessevdk/go-flags"
)

// == App ==
func init() {
	cli.Init()
}

func main() {
	fmt.Println("Started...")

	cli.Parse()
	// Verify required commands exist
	// if _, err := cli.Parser.Parse(); err != nil {
	// 	os.Exit(1)
	// }

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
