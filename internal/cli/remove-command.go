package cli

import (
	"fmt"

	"github.com/lorentzforces/selfman/internal/data"
	"github.com/lorentzforces/selfman/internal/ops"
	"github.com/spf13/cobra"
)

// TODO: implement a full removal option that also removes source files etc
// TODO: when library support is added, manage removing the library part as well

func CreateRemoveCmd() SelfmanCommand {
	return SelfmanCommand{
		cobraCmd: &cobra.Command{
			Use: "remove [flags] app-name",
			Short: "Remove an application's files managed by selfman",
			Aliases: []string{ "rm" },
		},
		runFunc: runRemoveCmd,
	}
}

func runRemoveCmd(cmd *cobra.Command, args []string) (*SelfmanResult, error) {
	if err := validatePrereqs(); err != nil { return nil, err }
	selfmanData, err := data.Produce()
	if err != nil { return nil, err }

	if len(args) < 1 {
		return nil,
			fmt.Errorf("Remove command expects an application name, but one was not provided")
	}
	ops, err := removeApp(args[0], selfmanData)
	if err != nil { return nil, err }

	return &SelfmanResult{
		textOutput: nil,
		operations: ops,
	}, nil
}

func removeApp(name string, selfmanData data.Selfman) ([]ops.Operation, error) {
	app, appStatus := selfmanData.AppStatus(name)
	if !appStatus.IsConfigured {
		return nil, fmt.Errorf("Could not find a configured application with name \"%s\"", name)
	}
	if !appStatus.SourcePresent {
		return nil, fmt.Errorf("Application \"%s\" has not been installed, no source present", name)
	}

	// by default, do not delete the source path
	actions := []ops.Operation{
		ops.DeleteFile{
			TypeOfDeletion: "Delete binary symlink",
			Path: app.BinaryPath(),
		},
		ops.DeleteFilesWithPrefix{
			TypeOfDeletion: "Delete built artifacts",
			DirPath: app.SystemConfig.ArtifactsPath(),
			FilePrefix: app.Name + "---",
		},
	}

	return actions, nil
}
