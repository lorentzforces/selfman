package cli

import (
	"fmt"

	"github.com/lorentzforces/selfman/internal/data"
	"github.com/lorentzforces/selfman/internal/ops"
	"github.com/spf13/cobra"
)

func CreateInstallCmd() SelfmanCommand {
	return SelfmanCommand{
		cobraCmd: &cobra.Command{
			Use: "install [flags] app-name",
			Short: "Install an application with a pre-existing configuration file",
		},
		runFunc: runInstallCmd,
	}
}

func runInstallCmd(cmd *cobra.Command, args []string) (*SelfmanResult, error) {
	if err := validatePrereqs(); err != nil { return nil, err }
	selfmanData, err := data.Produce()
	if err != nil { return nil, err }

	if len(args) < 1 {
		return nil,
			fmt.Errorf("Install command expects an application name, but one was not provided")
	}

	ops, err := installApp(args[0], selfmanData)
	if err != nil { return nil, err }

	return &SelfmanResult{
		textOutput: nil,
		operations: ops,
	}, nil
}

func installApp(name string, selfmanData data.Selfman) ([]ops.Operation, error) {
	app, appStatus := selfmanData.AppStatus(name)
	if !appStatus.IsConfigured {
		return nil, fmt.Errorf("Could not find a configured application with name \"%s\"", name)
	}

	buildTargetPath := app.BuildTargetPath()
	artifactPath := app.ArtifactPath()
	binPath := app.BinaryPath()

	actions := make([]ops.Operation, 0, 10)

	appSourceOp := app.GetObtainSourceOp()
	_, isGitApp := appSourceOp.(ops.GitClone)
	if isGitApp && appStatus.SourcePresent {
		actions = append(actions, ops.SkipCloneOp)
	} else {
		actions = append(actions, appSourceOp)
	}

	if versionOp := app.GetSelectVersionOp(); versionOp != nil {
		actions = append(actions, versionOp)
	}

	actions = append(actions, app.GetBuildOp())

	if !app.KeepBinWithSource {
		actions = append(
			actions,
			ops.MoveTarget{
				SourcePath: buildTargetPath,
				DestinationPath: artifactPath,
			},
		)
	}

	actions = append(
		actions,
		ops.LinkArtifact{
			SourcePath: artifactPath,
			DestinationPath: binPath,
		},
	)

	return actions, nil
}
