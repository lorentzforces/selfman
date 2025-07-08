package cli

import (
	"fmt"

	"github.com/lorentzforces/selfman/internal/data"
	"github.com/lorentzforces/selfman/internal/ops"
	"github.com/spf13/cobra"
)

func CreateUpdateCmd() SelfmanCommand {
	return SelfmanCommand{
		cobraCmd: &cobra.Command{
			Use: "update [flags] app-name",
			Short: "Update the given application",
		},
		runFunc: runUpdateCmd,
	}
}

func runUpdateCmd(cmd *cobra.Command, args []string) (*SelfmanResult, error) {
	if err := validatePrereqs(); err != nil { return nil, err
	}
	selfmanData, err := data.Produce()
	if err != nil { return nil, err }

	if len(args) < 1 {
		return nil,
		fmt.Errorf("Update command expects an application name, but one was not provided")
	}

	ops, err := updateApp(args[0], selfmanData)
	if err != nil { return nil, err }

	return &SelfmanResult{
		textOutput: nil,
		operations: ops,
	}, nil
}

func updateApp(name string, selfmanData data.Selfman) ([]ops.Operation, error) {
	app, status := selfmanData.AppStatus(name)
	if !status.IsConfigured {
		return nil, fmt.Errorf("Could not find a configured application with name \"%s\"", name)
	}
	if !status.SourcePresent {
		return nil, fmt.Errorf("Application \"%s\" has not been installed, no source present", name)
	}

	buildTargetPath := app.BuildTargetPath()
	artifactPath := app.ArtifactPath()
	binPath := app.BinaryPath()
	selectVersionOp := app.GetSelectVersionOp()
	buildOp := app.GetBuildOp()
	updateOp := app.GetFetchUpdatesOp()

	actions := []ops.Operation{
		updateOp,
		selectVersionOp,
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
