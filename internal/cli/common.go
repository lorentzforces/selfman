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
	runFunc func(*cobra.Command, []string) (*SelfmanResult, error)
}

type SelfmanResult struct {
	// Text output is always printed before any other messages
	textOutput fmt.Stringer
	// Any mutating operations to be executed as a result of running this command
	operations []ops.Operation
}

func (self *SelfmanCommand) RunSelfmanCommand(cmd *cobra.Command, args []string) error {
	cmdResult, err := self.runFunc(cmd, args)
	if err != nil { return err }

	dryRun, err := cmd.Flags().GetBool(globalOptionDryRun)
	run.AssertNoErr(err)

	if cmdResult.textOutput != nil {
		fmt.Println(cmdResult.textOutput)
	}

	if dryRun {
		dryRunOperations(cmdResult.operations)
		return nil
	} else {
		return executeOperations(cmdResult.operations)
	}
}

// Since the messages printed herein are progress updates, print to stderr
func executeOperations(actions []ops.Operation) error {
	// TODO: make this the verbose version, add non-verbose that only prints basic summary
	for _, action := range actions {
		fmt.Fprintln(os.Stderr, action.Describe())
		msg, err := action.Execute()
		if err != nil { return err }

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
