package cli

import (
	"fmt"

	"github.com/lorentzforces/selfman/internal/config"
	"github.com/spf13/cobra"
)

func CreateListCmd() *cobra.Command {
	return &cobra.Command{
		Use: "list",
		Short: "List all applications managed by selfman (currently debug placeholder)",
		RunE: runListCmd,
	}
}

func runListCmd(cmd *cobra.Command, args []string) error {
	config, err := config.Produce()
	fmt.Printf("==DEBUG== Resolved config: %#v\n", config)
	if err != nil {
		return err
	}

	return nil
}
