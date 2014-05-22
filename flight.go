package main

import (
  "./lib"
  "os"
)

var (
  app lib.TestFlight
)

// == App ==
func init() {
  err := app.Init()
  if err != nil {
    os.Exit(lib.ExitCodes["init_fail"])
  }

}

// Runs Test-Flight
func main() {
  app.ProcessCommands() // parse command line options now

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
  app.AppState.SetState("END")
}
