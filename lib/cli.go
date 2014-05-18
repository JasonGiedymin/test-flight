/*
 * [x] Moved all cli here
 * [ ] Might be helpful to create custom struct and move everything within
 */

package lib

import (
  "./config"
  "fmt"
  // "github.com/jessevdk/go-flags"
  // "os"
)

// == Commands ==
// var checkCommand CheckCommand
// var launchCommand LaunchCommand

// == Check Command ==
type CheckCommand struct {
  Dir      string `short:"d" long:"dir" description:"directory to run in"`
  AppState *ApplicationState
}

func (cmd *CheckCommand) Execute(args []string) error {
  fmt.Printf("Running Pre-Flight Check... in dir: [%v]\n", cmd.Dir)

  cmd.AppState.CurrentMode = "CHECK_FILES"
  cmd.AppState.State()

  _, err := HasRequiredFiles(&cmd.Dir, RequiredFiles)
  if err != nil {
    return err
  }

  buildFile, err := config.ReadBuildFile(cmd.Dir)
  if err != nil {
    return err
  }

  cmd.AppState.BuildFile = buildFile

  fmt.Println("Done.")
  return nil
}

// == Launch Command ==
type LaunchCommand struct {
  Dir string `short:"d" long:"dir" description:"directory to run in"`
}

func (cmd *LaunchCommand) Execute(args []string) error {
  fmt.Printf("Launching Tests... in dir: [%v]\n", cmd.Dir)
  _, err := HasRequiredFiles(&cmd.Dir, RequiredFiles)

  if err != nil {
    return err
  }

  fmt.Println("Done.")
  return nil
}
