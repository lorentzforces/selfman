package cli

import (
	"fmt"

	"github.com/lorentzforces/selfman/internal/data"
	"github.com/lorentzforces/selfman/internal/ops"
	"github.com/spf13/cobra"
)

func CreateRemoveCommand() SelfmanCommand {
	cmd := SelfmanCommand{
		cobraCmd: &cobra.Command{
			Use: "remove [flags] app-name",
			Short: "Remove an application's files managed by selfman",
		},
		opsCmd: runRemoveCmd,
	}
	cmd.InitCobraFunctions()
	return cmd
}

func runRemoveCmd(cmd *cobra.Command, args []string) ([]ops.Operation, error) {
	if err := validatePrereqs(); err != nil {
		return nil, err
	}
	selfmanData, err := data.Produce()
	if err != nil {
		return nil, err
	}

	if len(args) < 1 {
		return nil,
			fmt.Errorf("Remove command expects an application name, but one was not provided")
	}
	ops, err := removeApp(args[0], selfmanData)
	if err != nil { return nil, err }

	return ops, nil
}

func removeApp(name string, selfmanData data.Selfman) ([]ops.Operation, error) {
	app, appStatus := selfmanData.AppStatus(name)
	if !appStatus.IsConfigured {
		return nil, fmt.Errorf("Could not find a configured application with name \"%s\"", name)
	}

	// by default, do not delete the source path
	actions := []ops.Operation{
		ops.DeleteFile{
			TypeOfDeletion: "Delete binary symlink",
			Path: app.BinaryPath(),
		},
		ops.DeleteFile{
			TypeOfDeletion: "Delete built artifact",
			Path: app.ArtifactPath(),
		},
	}

	return actions, nil
}
