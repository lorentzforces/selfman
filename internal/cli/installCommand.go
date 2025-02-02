package cli

import (
	"fmt"

	"github.com/lorentzforces/selfman/internal/config"
	"github.com/lorentzforces/selfman/internal/ops"
	"github.com/lorentzforces/selfman/internal/platform"
	"github.com/spf13/cobra"
)

func CreateInstallCmd() SelfmanCommand {
	cmd := SelfmanCommand{
		cobraCmd: &cobra.Command{
			Use: "install app-name",
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
	configData, err := config.Produce()
	if err != nil {
		return nil, err
	}

	if len(args) < 1 {
		return nil,
			fmt.Errorf("Install command expects an application name, but one was not provided")
	}

	ops, err := installApp(args[0], configData)
	if err != nil {
		return nil, err
	}

	return ops, nil
}

func installApp(name string, cfg config.Config) ([]ops.Operation, error) {
	var app *config.AppConfig
	for _, appCandidate := range cfg.AppConfigs {
		if appCandidate.Name == name {
			app = &appCandidate
		}
	}
	if app == nil {
		return nil, fmt.Errorf("Could not find a configured application with name \"%s\"", name)
	}

	repoPath := platform.ResolveRepoPathForApp(name)
	actions := []ops.Operation{
		&ops.GitClone{
			RepoUrl: app.RemoteRepo,
			DestinationPath: repoPath,
		},
	}
	return actions, nil
}
