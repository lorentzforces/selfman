package cli

import (
	"testing"

	"github.com/lorentzforces/selfman/internal/data"
	"github.com/lorentzforces/selfman/internal/data/mocks"
	"github.com/lorentzforces/selfman/internal/ops"
	"github.com/lorentzforces/selfman/internal/run"
	"github.com/stretchr/testify/assert"
)

func TestRemoveCommandErrorsWithNonPresentApp(t *testing.T) {
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

	_, err = removeApp(notPresentApp.Name, false, selfmanData)
	assert.Error(
		t, err,
		"An error must be thrown if the user attempts to remove an application which is configured " +
			"but not present",
	)
}

func TestRemoveCommandProducesSaneOperations(t *testing.T) {
	systemConfig := data.DefaultTestConfig()

	appToRemove := data.AppConfig{
		SystemConfig: systemConfig,
		Name: "git-repo-app",
		Flavor: "git",
		RemoteRepo: run.StrPtr("git@github.com:github/gitignore.git"),
		BuildAction: data.ActionNone,
		Version: "main",
	}

	mockStorage := mocks.MockManagedFiles{}
	mockStorage.On("AppStatus", appToRemove.Name).Return(data.AppStatus{
		IsConfigured: true,
		SourcePresent: true,
		TargetPresent: true,
		LinkPresent: true,
	})

	selfmanData, err := data.SelfmanFromValues(
		systemConfig,
		[]data.AppConfig{appToRemove},
		&mockStorage,
	)
	assert.NoError(t, err)
	run.BailIfFailed(t)

	actions, err := removeApp(appToRemove.Name, false, selfmanData)
	assert.NoError(t, err)
	run.BailIfFailed(t)
	expectedActions := []ops.Operation{
		ops.DeleteFile{
			TypeOfDeletion: "Delete binary symlink",
			Path: appToRemove.BinaryPath(),
		},
		ops.DeleteFile{
			TypeOfDeletion: "Delete library link",
			Path: appToRemove.LibPath(),
		},
		ops.DeleteFilesWithPrefix{
			TypeOfDeletion: "Delete built artifacts",
			DirPath: systemConfig.ArtifactsPath(),
			FilePrefix: appToRemove.Name + "---",
		},
	}
	assert.Equal(t, expectedActions, actions)
}

func TestRemoveCommandRemovesSourceWhenAsked(t *testing.T) {
	systemConfig := data.DefaultTestConfig()

	appToRemove := data.AppConfig{
		SystemConfig: systemConfig,
		Name: "git-repo-app",
		Flavor: "git",
		RemoteRepo: run.StrPtr("git@github.com:github/gitignore.git"),
		BuildAction: data.ActionNone,
		Version: "main",
	}

	mockStorage := mocks.MockManagedFiles{}
	mockStorage.On("AppStatus", appToRemove.Name).Return(data.AppStatus{
		IsConfigured: true,
		SourcePresent: true,
		TargetPresent: true,
		LinkPresent: true,
	})

	selfmanData, err := data.SelfmanFromValues(
		systemConfig,
		[]data.AppConfig{appToRemove},
		&mockStorage,
	)
	assert.NoError(t, err)
	run.BailIfFailed(t)

	actions, err := removeApp(appToRemove.Name, true, selfmanData)
	assert.NoError(t, err)
	run.BailIfFailed(t)

	expectedActions := []ops.Operation{
		ops.DeleteFile{
			TypeOfDeletion: "Delete binary symlink",
			Path: appToRemove.BinaryPath(),
		},
		ops.DeleteFile{
			TypeOfDeletion: "Delete library link",
			Path: appToRemove.LibPath(),
		},
		ops.DeleteFilesWithPrefix{
			TypeOfDeletion: "Delete built artifacts",
			DirPath: systemConfig.ArtifactsPath(),
			FilePrefix: appToRemove.Name + "---",
		},
		ops.DeleteDir{
			TypeOfDeletion: "Delete source directory",
			Path: appToRemove.SourcePath(),
		},
	}
	assert.Equal(t, expectedActions, actions)
}
