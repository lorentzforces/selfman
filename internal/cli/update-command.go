package cli

import (
	"fmt"

	"github.com/lorentzforces/selfman/internal/data"
	"github.com/lorentzforces/selfman/internal/ops"
	"github.com/spf13/cobra"
)

func CreateUpdateCmd() SelfmanCommand {
	cmd := SelfmanCommand{
		cobraCmd: &cobra.Command{
			Use: "update [flags] app-name",
			Short: "Update the given application",
		},
		opsCmd: runUpdateCmd,
	}
	cmd.InitCobraFunctions()
	return cmd
}

func runUpdateCmd(cmd *cobra.Command, args []string) ([]ops.Operation, error) {
	if err := validatePrereqs(); err != nil {
		return nil, err
	}
	selfmanData, err := data.Produce()
	if err != nil {
		return nil, err
	}

	if len(args) < 1 {
		return nil,
		fmt.Errorf("Update command expects an application name, but one was not provided")
	}

	ops, err := updateApp(args[0], selfmanData)
	if err != nil {
		return nil, err
	}

	return ops, nil
}

// TODO: check if app is even installed before trying to update it
func updateApp(name string, selfmanData data.Selfman) ([]ops.Operation, error) {
	app, configured := selfmanData.AppConfigs[name]
	if !configured {
		return nil, fmt.Errorf("Could not find a configured application with name \"%s\"", name)
	}

	buildTargetPath := app.BuildTargetPath()
	artifactPath := app.ArtifactPath()
	binPath := app.BinaryPath()
	buildOp := app.GetBuildOp()
	updateOp := app.GetFetchUpdatesOp()
	actions := []ops.Operation{
		updateOp,
		buildOp,
		ops.MoveTarget{
			SourcePath: buildTargetPath,
			DestinationPath: artifactPath,
		},
		ops.LinkArtifact{
			SourcePath: artifactPath,
			DestinationPath: binPath,
		},
	}
	return actions, nil
}
