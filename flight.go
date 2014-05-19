package main

import (
  "./lib"
  // "./lib/config"
  "fmt"
  "os"
  // "github.com/fsouza/go-dockerclient"
  // "flag"
  "github.com/jessevdk/go-flags"
)

type CommandOptions struct{}

var (
  app           lib.TestFlight
  options       CommandOptions
  checkCommand  = lib.CheckCommand{AppState: &app.AppState}
  launchCommand lib.LaunchCommand
)

// == App ==
func init() {
  err := app.Init()
  if (err != nil) {
    fmt.Println(err)
    os.Exit(1)
  }

  app.Parser = flags.NewParser(&options, flags.Default)

  app.Parser.AddCommand("check",
    "pre-flight check",
    "Used for pre-flight check of the ansible playbook.",
    &checkCommand)

  app.Parser.AddCommand("launch",
    "flight launch",
    "Launch an ansible playbook test.",
    &launchCommand)
}

func main() {
  fmt.Println("Started...")
  app.Parse()
  fmt.Println(*app.AppState.BuildFile)

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
