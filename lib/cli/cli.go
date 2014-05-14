package cli

import (
	"fmt"
	// "github.com/jessevdk/go-flags"
)

type CheckCommand struct{}

func (cmd *CheckCommand) Execute(args []string) error {
	fmt.Println("Running Pre-Flight Check...")
	return nil
}
