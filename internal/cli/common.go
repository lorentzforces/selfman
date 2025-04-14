package cli

import (
	"fmt"
	"os"

	"github.com/lorentzforces/selfman/internal/git"
	"github.com/lorentzforces/selfman/internal/ops"
	"github.com/lorentzforces/selfman/internal/run"
	"github.com/spf13/cobra"
)

func validatePrereqs() error {
	if !git.ExecExists() {
		return fmt.Errorf("Cannot find a \"git\" executable on PATH")
	}
	return nil
}

type SelfmanCommand struct {
	cobraCmd *cobra.Command
	opsCmd func(*cobra.Command, []string) ([]ops.Operation, error)
}

func (self *SelfmanCommand) InitCobraFunctions() {
	self.cobraCmd.RunE = self.RunMutatingSelfmanCmd
}

func (self *SelfmanCommand) RunMutatingSelfmanCmd(cmd *cobra.Command, args []string) error {
	actions, err := self.opsCmd(cmd, args)
	if err != nil {
		return err
	}

	dryRun, err := cmd.Flags().GetBool(globalOptionDryRun)
	run.AssertNoErr(err)

	if dryRun {
		dryRunOperations(actions)
		return nil
	} else {
		return executeOperations(actions)
	}
}

// Since the messages printed herein are progress updates, print to stderr
func executeOperations(actions []ops.Operation) error {
	// TODO: make this the verbose version, add non-verbose that only prints basic summary
	for _, action := range actions {
		fmt.Fprintln(os.Stderr, action.Describe())
		msg, err := action.Execute()

		if err != nil {
			return err
		}

		fmt.Printf("✓")
		if len(msg) > 0 {
			fmt.Printf(" %s", msg)
		}
		fmt.Println()
	}

	return nil
}

// Since this is asked for as the main output, print to stdout
func dryRunOperations(actions []ops.Operation) {
	fmt.Printf("Would perform the following operations:\n\n")
	for _, action := range actions {
		fmt.Println(action.Describe())
	}
}
