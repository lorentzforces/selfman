package cli

import (
	"path"
	"testing"

	"github.com/lorentzforces/selfman/internal/data"
	"github.com/lorentzforces/selfman/internal/ops"
	"github.com/lorentzforces/selfman/internal/run"
	"github.com/stretchr/testify/assert"
)

func TestInstallCommandValidatesNameExists(t *testing.T) {
	systemConfig := data.DefaultTestConfig()
	selfmanData := data.Selfman{
		SystemConfig: systemConfig,
		AppConfigs: map[string]data.AppConfig{},
		Storage: nil,
	}

	_, err := installApp("not-available-name", selfmanData)
	assert.Error(
		t, err,
		"An error must be thrown if the user attempts to install an application that is not " +
			"configured",
	)
}

func TestInstallCommandProducesSaneOperations(t *testing.T) {
	systemConfig := data.DefaultTestConfig()
	appToInstall := data.AppConfig{
		Name: "git-repo-app",
		Type: "git",
		RemoteRepo: run.StrPtr("git@github.com:github/gitignore.git"),
	}

	selfmanData, err := data.SelfmanFromValues(
		systemConfig,
		[]data.AppConfig{appToInstall},
		nil,
	)
	assert.NoError(t, err)
	run.BailIfFailed(t)

	actions, err := installApp(appToInstall.Name, selfmanData)
	assert.NoError(t, err)
	run.BailIfFailed(t)
	expectedActions := []ops.Operation{
		&ops.GitClone{
			RepoUrl: *selfmanData.AppConfigs[appToInstall.Name].RemoteRepo,
			DestinationPath: path.Join(selfmanData.SystemConfig.SourcesPath(), appToInstall.Name),
		},
		// TODO: this will need to take into account app-specific build target paths
		&ops.MoveTarget{
			SourcePath: path.Join(
				selfmanData.SystemConfig.SourcesPath(),
				appToInstall.Name,
				appToInstall.Name,
			),
			DestinationPath: path.Join(selfmanData.SystemConfig.ArtifactsPath(), appToInstall.Name),
		},
	}
	assert.Equal(t, expectedActions, actions)
}
