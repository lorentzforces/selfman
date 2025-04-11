package cli

import (
	"github.com/spf13/cobra"
)

func CreateRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use: "selfman",
		Short: "A tool for managing self-managed and self-build applications & tools",
		SilenceUsage: true,
		SilenceErrors: true,
	}

	rootCmd.InitDefaultHelpFlag()

	rootCmd.AddCommand(CreateListCmd())
	rootCmd.AddCommand(CreateInstallCmd().cobraCmd)
	rootCmd.AddCommand(CreateUpdateCmd().cobraCmd)

	return rootCmd
}
