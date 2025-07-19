package cli

import (
	"fmt"

	"github.com/lorentzforces/selfman/internal/data"
	"github.com/lorentzforces/selfman/internal/ops"
	"github.com/spf13/cobra"
)

func CreateMakeItSoCmd() SelfmanCommand {
	return SelfmanCommand{
		cobraCmd: &cobra.Command{
			Use: "make-it-so",
			Short: "Update, install, or otherwise make an application up-to-date with its " +
				"configuration",
			Aliases: []string{ "mis" },
		},
		runFunc: runMakeItSoCmd,
	}
}

func runMakeItSoCmd(cmd *cobra.Command, args []string) (*SelfmanResult, error) {
	if err := validatePrereqs(); err != nil { return nil, err }
	selfmanData, err := data.Produce()
	if err != nil { return nil, err }

	if len(args) < 1 {
		if len(args) < 1 {
			return nil, fmt.Errorf(
				"make-it-so command expects an application name, but one was not provided")
		}
	}

	ops, err := makeItSo(args[0], selfmanData)
	if err != nil { return nil, err }

	return &SelfmanResult{
		textOutput: nil,
		operations: ops,
	}, nil
}

func makeItSo(name string, selfmanData data.Selfman) ([]ops.Operation, error) {
	_, appStatus := selfmanData.AppStatus(name)
	if !appStatus.IsConfigured {
		return nil, fmt.Errorf("Could not find a configured application with name \"%s\"", name)
	}

	return []ops.Operation{ }, nil
}
