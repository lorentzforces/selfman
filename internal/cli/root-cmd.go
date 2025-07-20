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

	addSelfmanCommands(
		rootCmd,
		[]SelfmanCommand{
			CreateListCmd(),
			CreateMakeItSoCmd(),
			CreateCheckCmd(),
			CreateRemoveCmd(),
		},
	)
	// TODO: intake binary for static-binary app
	// TODO(?): rollback?
	// TODO(?): list previous versions?
	// TODO: some kind of validation command for configuration? (roll into check?)

	return rootCmd
}

func addSelfmanCommands(rootCmd *cobra.Command, cmds []SelfmanCommand) {
	for _, cmd := range cmds {
		cmd.cobraCmd.RunE = cmd.RunSelfmanCommand
		rootCmd.AddCommand(cmd.cobraCmd)
	}
}
