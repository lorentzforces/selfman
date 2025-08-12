package cli

import (
	"fmt"

	"github.com/lorentzforces/selfman/internal/data"
	"github.com/lorentzforces/selfman/internal/ops"
	"github.com/lorentzforces/selfman/internal/run"
	"github.com/spf13/cobra"
)

const removeCmdOptionRemoveSource = "remove-source"

func CreateRemoveCmd() SelfmanCommand {
	selfmanCmd := SelfmanCommand{
		cobraCmd: &cobra.Command{
			Use: "remove [flags] app-name",
			Short: "Remove an application's files managed by selfman",
			Aliases: []string{ "rm" },
		},
		runFunc: runRemoveCmd,
	}

	selfmanCmd.cobraCmd.Flags().Bool(
		removeCmdOptionRemoveSource,
		false,
		"Fully remove application source in addition to artifacts, links, & products",
	)

	return selfmanCmd
}

func runRemoveCmd(cmd *cobra.Command, args []string) (*SelfmanResult, error) {
	if err := validatePrereqs(); err != nil { return nil, err }
	selfmanData, err := data.Produce()
	if err != nil { return nil, err }

	if len(args) < 1 {
		return nil,
			fmt.Errorf("Remove command expects an application name, but one was not provided")
	}
	removeSource, err := cmd.Flags().GetBool(removeCmdOptionRemoveSource)
	run.AssertNoErr(err)
	ops, err := removeApp(args[0], removeSource, selfmanData)
	if err != nil { return nil, err }

	return &SelfmanResult{
		textOutput: nil,
		operations: ops,
	}, nil
}

func removeApp(name string, removeSource bool, selfmanData data.Selfman) ([]ops.Operation, error) {
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
		ops.DeleteFile{
			TypeOfDeletion: "Delete library link",
			Path: app.LibPath(),
		},
		ops.DeleteFilesWithPrefix{
			TypeOfDeletion: "Delete built artifacts",
			DirPath: app.SystemConfig.ArtifactsPath(),
			FilePrefix: app.Name + "---",
		},
	}

	if removeSource {
		actions = append(actions, ops.DeleteDir{
			TypeOfDeletion: "Delete source directory",
			Path: app.SourcePath(),
		})
	}

	return actions, nil
}
