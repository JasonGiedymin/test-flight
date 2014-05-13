package main

import (
	"./lib"
	"fmt"
	// "github.com/fsouza/go-dockerclient"
)

/*
- read build.json
- require:
  - build.json
  - vars/
  - tests/
  - meta/
  - tasks/
  - docker/
- require:
  - base image (based on build.json runtime selection)
*/

func main() {
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

	files, filesError := lib.CheckForFiles()
	fmt.Printf("Found these files %v, and these errors: [%v].\n", files, filesError)
}
