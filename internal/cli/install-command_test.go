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
		SystemConfig: systemConfig,
		Name: "git-repo-app",
		Flavor: "git",
		RemoteRepo: run.StrPtr("git@github.com:github/gitignore.git"),
		BuildAction: "none",
		Version: "main",
	}

	mockStorage := mocks.MockManagedFiles{}
	mockStorage.On("AppStatus", appToInstall.Name).Return(data.AppStatus{
		IsConfigured: true,
	})

	selfmanData, err := data.SelfmanFromValues(
		systemConfig,
		[]data.AppConfig{appToInstall},
		&mockStorage,
	)
	assert.NoError(t, err)
	run.BailIfFailed(t)

	actions, err := installApp(appToInstall.Name, selfmanData)
	assert.NoError(t, err)
	run.BailIfFailed(t)

	expectedActions := []ops.Operation{
		ops.GitClone{
			RepoUrl: *selfmanData.AppConfigs[appToInstall.Name].RemoteRepo,
			DestinationPath: path.Join(appToInstall.SourcePath()),
		},
		// TODO: change this to have an actual build step in it
		ops.NoBuildOp,
		ops.MoveTarget{
			SourcePath: path.Join(appToInstall.SourcePath(), appToInstall.Name),
			DestinationPath: appToInstall.ArtifactPath(),
		},
		ops.LinkArtifact{
			SourcePath: appToInstall.ArtifactPath(),
			DestinationPath: path.Join(*selfmanData.SystemConfig.BinaryDir, appToInstall.Name),
		},
	}
	assert.Equal(t, expectedActions, actions)
}

func TestInstallGitDoesNotCloneIfSourceAlreadyPresent(t *testing.T) {
	systemConfig := data.DefaultTestConfig()

	appToInstall := data.AppConfig{
		SystemConfig: systemConfig,
		Name: "git-repo-app",
		Flavor: "git",
		RemoteRepo: run.StrPtr("git@github.com:github/gitignore.git"),
		BuildAction: "none",
		Version: "main",
	}

	mockStorage := mocks.MockManagedFiles{}
	mockStorage.On("AppStatus", appToInstall.Name).Return(data.AppStatus{
		IsConfigured: true,
		SourcePresent: true,
	})

	selfmanData, err := data.SelfmanFromValues(
		systemConfig,
		[]data.AppConfig{ appToInstall },
		&mockStorage,
	)
	assert.NoError(t, err)
	run.BailIfFailed(t)

	actions, err := installApp(appToInstall.Name, selfmanData)
	assert.NoError(t, err)
	run.BailIfFailed(t)

	expectedActions := []ops.Operation{
		ops.SkipCloneOp,
		ops.NoBuildOp,
		ops.MoveTarget{
			SourcePath: path.Join(appToInstall.SourcePath(), appToInstall.Name),
			DestinationPath: appToInstall.ArtifactPath(),
		},
		ops.LinkArtifact{
			SourcePath: appToInstall.ArtifactPath(),
			DestinationPath: path.Join(*selfmanData.SystemConfig.BinaryDir, appToInstall.Name),
		},
	}
	assert.Equal(t, expectedActions, actions)
}
