package cli

import (
	"fmt"

	"github.com/lorentzforces/selfman/internal/data"
	"github.com/spf13/cobra"
)

func CreateCheckCmd() *cobra.Command {
	return &cobra.Command{
		Use: "check",
		Short: "Get detailed information about an application",
		RunE: runCheckCmd,
	}
}

func runCheckCmd(cmd *cobra.Command, args []string) error {
	if err := validatePrereqs(); err != nil {
		return fmt.Errorf("Well, there's your problem: %w", err)
	}
	selfmanData, err := data.Produce()
	if err != nil {
		return err
	}

	// TODO: return a result object with detailed info
	err = checkApp(args[0], selfmanData)

	return err
}

func checkApp(name string, selfmanData data.Selfman) error {
	_, status := selfmanData.AppStatus(name)
	if !status.IsConfigured {
		return fmt.Errorf(
			"Well, there's your problem: no configuration for an app named %s was found",
			name,
		)
	}

	return nil
}
