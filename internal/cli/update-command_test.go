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

	_, err := updateApp("not-available-name", selfmanData)
	assert.Error(
		t, err,
		"An error must be thrown if the user attempts to update an application that is not " +
			"configured",
	)
}

func TestUpdateCommandProducesSaneOperations(t *testing.T) {
	systemConfig := data.DefaultTestConfig()
	appToUpdate := data.AppConfig{
		SystemConfig: systemConfig,
		Name: "updatable-app",
		Flavor: "git",
		RemoteRepo: run.StrPtr("doesn't matter"),
		BuildAction: "script",
		BuildCmd: run.StrPtr("make deez"),
		Version: "main",
	}

	mockStorage := mocks.MockManagedFiles{}
	mockStorage.On("AppStatus", appToUpdate.Name).Return(data.AppStatus{
		IsConfigured: true,
		SourcePresent: true,
		TargetPresent: true,
		LinkPresent: true,
	})

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
			RepoPath: appToUpdate.SourcePath(),
		},
		ops.BuildWithScript{
			SourcePath: appToUpdate.SourcePath(),
			ScriptShell: "/bin/sh",
			ScriptCmd: "make deez",
		},
		ops.MoveTarget{
			SourcePath: path.Join(appToUpdate.SourcePath(), appToUpdate.Name),
			DestinationPath: appToUpdate.ArtifactPath(),
		},
		ops.LinkArtifact{
			SourcePath: appToUpdate.ArtifactPath(),
			DestinationPath: path.Join(*selfmanData.SystemConfig.BinaryDir, appToUpdate.Name),
		},
	}
	assert.Equal(t, expectedActions, actions)
}

func TestUpdateCommandErrorsWithNonPresentApp(t *testing.T) {
	systemConfig := data.DefaultTestConfig()
	notPresentApp := data.AppConfig{
		SystemConfig: systemConfig,
		Name: "not-present-app",
		Flavor: "git",
		RemoteRepo: run.StrPtr("doesn't matter"),
		BuildAction: "script",
		BuildCmd: run.StrPtr("make deez"),
		Version: "main",
	}

	mockStorage := mocks.MockManagedFiles{}
	mockStorage.On("AppStatus", notPresentApp.Name).Return(data.AppStatus{
		IsConfigured: true,
	})


	selfmanData, err := data.SelfmanFromValues(
		systemConfig,
		[]data.AppConfig{ notPresentApp },
		&mockStorage,
	)
	assert.NoError(t, err)
	run.BailIfFailed(t)

	_, err = updateApp(notPresentApp.Name, selfmanData)
	assert.Error(
		t, err,
		"An error must be thrown if the user attempts to update an application which is configured " +
			"but not present",
	)
}
