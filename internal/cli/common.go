package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/lorentzforces/selfman/internal/git"
	"github.com/lorentzforces/selfman/internal/ops"
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
	self.cobraCmd.RunE = self.RunSelfmanCmd
}

func (self *SelfmanCommand) RunSelfmanCmd(cmd *cobra.Command, args []string) error {
	// TODO: support a dry-run
	actions, err := self.opsCmd(cmd, args)
	if err != nil {
		return err
	}

	opErrs := make([]error, 0)
	for _, action := range actions {
		fmt.Fprintf(os.Stderr, "Performing op: %s\n", action.Describe())
		msg, err := action.Execute()
		if len(msg) > 0 {
			fmt.Println(msg)
		}
		if err != nil {
			opErrs = append(opErrs, err)
			break
		}
	}

	if len(opErrs) > 0 {
		return errors.Join(opErrs...)
	} else {
		return nil
	}
}
