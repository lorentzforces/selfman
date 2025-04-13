package cli

import (
	"fmt"

	"github.com/lorentzforces/selfman/internal/data"
	"github.com/lorentzforces/selfman/internal/ops"
	"github.com/spf13/cobra"
)

func CreateInstallCmd() SelfmanCommand {
	cmd := SelfmanCommand{
		cobraCmd: &cobra.Command{
			Use: "install [flags] app-name",
			Short: "Install an application with a pre-existing configuration file",
		},
		opsCmd: runInstallCmd,
	}
	cmd.InitCobraFunctions()
	return cmd
}

func runInstallCmd(cmd *cobra.Command, args []string) ([]ops.Operation, error) {
	if err := validatePrereqs(); err != nil {
		return nil, err
	}
	selfmanData, err := data.Produce()
	if err != nil {
		return nil, err
	}

	if len(args) < 1 {
		return nil,
			fmt.Errorf("Install command expects an application name, but one was not provided")
	}

	ops, err := installApp(args[0], selfmanData)
	if err != nil {
		return nil, err
	}

	return ops, nil
}

// TODO: check if app is already installed
func installApp(name string, selfmanData data.Selfman) ([]ops.Operation, error) {
	app, configured := selfmanData.AppConfigs[name]
	if !configured {
		return nil, fmt.Errorf("Could not find a configured application with name \"%s\"", name)
	}

	// TODO: right now we only support apps of type git
	repoPath := app.SourcePath()
	buildTargetPath := app.BuildTargetPath()
	artifactPath := app.ArtifactPath()
	binPath := app.BinaryPath()
	buildOp := app.GetBuildOp()
	actions := []ops.Operation{
		ops.GitClone{
			RepoUrl: *app.RemoteRepo,
			DestinationPath: repoPath,
		},
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
