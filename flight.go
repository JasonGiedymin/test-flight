package main

import (
  "./lib"
  "./lib/config"
  "fmt"
  "os"
  // "github.com/fsouza/go-dockerclient"
  // "flag"
  "github.com/jessevdk/go-flags"
)

type TestFlight struct {
  appState      lib.ApplicationState
  configFile    config.ConfigFile
  buildFile     config.BuildFile
  requiredFiles []lib.RequiredFile
  parser        *flags.Parser
}

type CommandOptions struct{}

func (app *TestFlight) Parse() {
  if _, err := app.parser.Parse(); err != nil {
    os.Exit(1)
  }
}

func (app *TestFlight) State() {
  app.appState.State()
}

var (
  app           TestFlight
  options       CommandOptions
  checkCommand  = lib.CheckCommand{AppState: &app.appState}
  launchCommand lib.LaunchCommand
)

// == App ==
func init() {
  app.appState.CurrentMode = "INIT"
  app.State()

  app.parser = flags.NewParser(&options, flags.Default)

  app.parser.AddCommand("check",
    "pre-flight check",
    "Used for pre-flight check of the ansible playbook.",
    &checkCommand)

  app.parser.AddCommand("launch",
    "flight launch",
    "Launch an ansible playbook test.",
    &launchCommand)
}

func main() {
  fmt.Println("Started...")
  app.Parse()
  fmt.Println(*app.appState.BuildFile)

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
