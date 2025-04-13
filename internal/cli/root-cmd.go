package cli

import (
	"github.com/spf13/cobra"
)

const (
	globalOptionDryRun = "dry-run"
)

func CreateRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use: "selfman",
		Short: "A tool for managing self-managed and self-build applications & tools",
		SilenceUsage: true,
		SilenceErrors: true,
	}

	rootCmd.InitDefaultHelpFlag()
	rootCmd.PersistentFlags().Bool(
		globalOptionDryRun,
		false,
		"For commands which would make changes, print operations that would be taken instead of " +
			"executing them",
	)

	rootCmd.AddCommand(CreateListCmd())
	rootCmd.AddCommand(CreateInstallCmd().cobraCmd)
	rootCmd.AddCommand(CreateUpdateCmd().cobraCmd)
	// TODO: intake binary for static-binary app
	// TODO: uninstall
	// TODO(?): rollback?
	// TODO(?): list previous versions?
	// TODO: some kind of validation/check command

	return rootCmd
}
