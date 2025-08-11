package cli

import (
	"fmt"

	"github.com/lorentzforces/selfman/internal/data"
	"github.com/lorentzforces/selfman/internal/ops"
	"github.com/spf13/cobra"
)

func CreateMakeItSoCmd() SelfmanCommand {
	return SelfmanCommand{
		cobraCmd: &cobra.Command{
			Use: "make-it-so",
			Short: "Update, install, or otherwise make an application up-to-date with its " +
				"configuration",
			Aliases: []string{ "mis" },
		},
		runFunc: runMakeItSoCmd,
	}
}

func runMakeItSoCmd(cmd *cobra.Command, args []string) (*SelfmanResult, error) {
	if err := validatePrereqs(); err != nil { return nil, err }
	selfmanData, err := data.Produce()
	if err != nil { return nil, err }

	if len(args) < 1 {
		if len(args) < 1 {
			return nil, fmt.Errorf(
				"make-it-so command expects an application name, but one was not provided")
		}
	}

	ops, err := makeItSo(args[0], selfmanData)
	if err != nil { return nil, err }

	return &SelfmanResult{
		textOutput: nil,
		operations: ops,
	}, nil
}

func makeItSo(name string, selfmanData data.Selfman) ([]ops.Operation, error) {
	app, appStatus := selfmanData.AppStatus(name)
	if !appStatus.IsConfigured {
		return nil, fmt.Errorf("Could not find a configured application with name \"%s\"", name)
	}

	buildTargetPath := app.BuildTargetPath()
	artifactPath := app.ArtifactPath()
	binPath := app.BinaryPath()

	actions := make([]ops.Operation, 0, 10)

	fetchUpdatesOp := app.GetFetchUpdatesOp()
	if !appStatus.SourcePresent {
		actions = append(actions, app.GetObtainSourceOp())
	} else if appStatus.SourcePresent && fetchUpdatesOp != nil {
		// don't need to fetch updates if we just obtained the source
		actions = append(actions, app.GetFetchUpdatesOp())
	}

	// TODO: consider deleting the default branch after cloning and checking out the desired branch
	//       IF AND ONLY IF we just cloned

	if versionOp := app.GetSelectVersionOp(); versionOp != nil {
		actions = append(actions, versionOp)
	}

	// TODO: git app which you're still on the same branch of and you want to build whatever is on the tip...
	// (right now this will not build)
	if !appStatus.TargetPresent {
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
	}

	actions = append(
		actions,
		ops.LinkArtifact{
			SourcePath: artifactPath,
			DestinationPath: binPath,
		},
	)

	if app.LinkSourceAsLib {
		actions = append(
			actions,
			ops.LinkLibrary{
				SourcePath: app.SourcePath(),
				DestinationPath: app.LibPath(),
			},
		)
	}

	return actions, nil
}
