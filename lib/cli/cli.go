/*
 * [x] Moved all cli here
 * [ ] Might be helpful to create custom struct and move everything within
 */

package cli

import (
	"../"
	"fmt"
	"github.com/jessevdk/go-flags"
	"os"
)

type Options struct{}

var options Options

var Parser = flags.NewParser(&options, flags.Default)

// == Std Command Opts ==
type StandardOption struct {
	Dir string `short:"f" long:"force" description:"Force removal of files"`
}

// == Commands ==
var checkCommand CheckCommand
var launchCommand LaunchCommand

// == Check Command ==
type CheckCommand struct {
	Dir string `short:"d" long:"dir" description:"directory to run in"`
}

func (cmd *CheckCommand) Execute(args []string) error {
	fmt.Printf("Running Pre-Flight Check... in dir: [%v]\n", cmd.Dir)
	_, err := lib.HasRequiredFiles()

	return err
}

// == Launch Command ==
type LaunchCommand struct{}

func (cmd *LaunchCommand) Execute(args []string) error {
	fmt.Println("Launching tests...")
	_, err := lib.HasRequiredFiles()

	return err
}

// func getArgs() {
// 	args, _ := Parser.ParseArgs(os.Args)
// 	fmt.Println(args)
// 	// fmt.Printf("Found dir: %v ", strings.Join(args, " "))
// }

// == Init ==
func Init() {
	Parser.AddCommand("check",
		"pre-flight check",
		"Used for pre-flight check of the ansible playbook.",
		&checkCommand)

	Parser.AddCommand("launch",
		"flight launch",
		"Launch an ansible playbook test.",
		&launchCommand)
}

func Parse() {
	if _, err := Parser.Parse(); err != nil {
		os.Exit(1)
	}
}
