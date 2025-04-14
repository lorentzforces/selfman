package cli

import (
	"path"
	"testing"

	"github.com/lorentzforces/selfman/internal/data"
	"github.com/lorentzforces/selfman/internal/data/mocks"
	"github.com/lorentzforces/selfman/internal/ops"
	"github.com/lorentzforces/selfman/internal/run"
	"github.com/stretchr/testify/assert"
)

func TestUpdateCommandValidatesNameExists(t *testing.T) {
	systemConfig := data.DefaultTestConfig()
	mockStorage := mocks.MockManagedFiles{}
	selfmanData := data.Selfman{
		SystemConfig: systemConfig,
		AppConfigs: map[string]data.AppConfig{},
		Storage: &mockStorage,
	}

	_, err := installApp("not-available-name", selfmanData)
	assert.Error(
		t, err,
		"An error must be thrown if the user attempts to install an application that is not " +
			"configured",
	)
}

func TestUpdateCommandProducesSaneOperations(t *testing.T) {
	systemConfig := data.DefaultTestConfig()
	appToUpdate := data.AppConfig{
		SystemConfig: systemConfig,
		Name: "updatable-app",
		Type: "git",
		RemoteRepo: run.StrPtr("doesn't matter"),
		BuildAction: "script",
		BuildCmd: run.StrPtr("make deez"),
	}

	mockStorage := mocks.MockManagedFiles{}
	mockStorage.MockAllFilesPresent()
	selfmanData, err := data.SelfmanFromValues(
		systemConfig,
		[]data.AppConfig{appToUpdate},
		&mockStorage,
	)
	assert.NoError(t, err)
	run.BailIfFailed(t)

	actions, err := updateApp(appToUpdate.Name, selfmanData)
	assert.NoError(t, err)
	run.BailIfFailed(t)
	expectedActions := []ops.Operation{
		ops.GitPull{
			RepoPath: path.Join(
				selfmanData.SystemConfig.SourcesPath(),
				appToUpdate.Name,
			),
		},
		ops.BuildWithScript{
			SourcePath: path.Join(
				selfmanData.SystemConfig.SourcesPath(),
				appToUpdate.Name,
			),
			ScriptShell: "/bin/sh",
			ScriptCmd: "make deez",
		},
		ops.MoveTarget{
			SourcePath: path.Join(
				selfmanData.SystemConfig.SourcesPath(),
				appToUpdate.Name,
				appToUpdate.Name,
			),
			DestinationPath: path.Join(
				selfmanData.SystemConfig.ArtifactsPath(), appToUpdate.Name,
			),
		},
		ops.LinkArtifact{
			SourcePath: path.Join(selfmanData.SystemConfig.ArtifactsPath(), appToUpdate.Name),
			DestinationPath: path.Join(*selfmanData.SystemConfig.BinaryDir, appToUpdate.Name),
		},
	}
	assert.Equal(t, expectedActions, actions)
}
