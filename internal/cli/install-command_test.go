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
	mockStorage := mocks.MockManagedFiles{}

	appToInstall := data.AppConfig{
		SystemConfig: systemConfig,
		Name: "git-repo-app",
		BaseType: "git",
		RemoteRepo: run.StrPtr("git@github.com:github/gitignore.git"),
		BuildAction: "none",
	}
	mockStorage.On(
		"IsGitAppPresent",
		path.Join(systemConfig.SourcesPath(), appToInstall.Name),
	).Return(false)
	mockStorage.On(
		"ExecutableExists",
		path.Join(systemConfig.ArtifactsPath(), appToInstall.Name),
	).Return(false)
	mockStorage.On(
		"LinkExists",
		path.Join(*systemConfig.BinaryDir, appToInstall.Name),
	).Return(false)

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
			DestinationPath: path.Join(selfmanData.SystemConfig.SourcesPath(), appToInstall.Name),
		},
		// TODO: change this to have an actual build step in it
		ops.NoBuildOp,
		ops.MoveTarget{
			SourcePath: path.Join(
				selfmanData.SystemConfig.SourcesPath(),
				appToInstall.Name,
				appToInstall.Name,
			),
			DestinationPath: path.Join(selfmanData.SystemConfig.ArtifactsPath(), appToInstall.Name),
		},
		ops.LinkArtifact{
			SourcePath: path.Join(selfmanData.SystemConfig.ArtifactsPath(), appToInstall.Name),
			DestinationPath: path.Join(*selfmanData.SystemConfig.BinaryDir, appToInstall.Name),
		},
	}
	assert.Equal(t, expectedActions, actions)
}

func TestInstallGitDoesNotCloneIfSourceAlreadyPresent(t *testing.T) {
	systemConfig := data.DefaultTestConfig()
	mockStorage := mocks.MockManagedFiles{}

	appToInstall := data.AppConfig{
		SystemConfig: systemConfig,
		Name: "git-repo-app",
		BaseType: "git",
		RemoteRepo: run.StrPtr("git@github.com:github/gitignore.git"),
		BuildAction: "none",
	}
	mockStorage.On(
		"IsGitAppPresent",
		path.Join(systemConfig.SourcesPath(), appToInstall.Name),
	).Return(true)
	mockStorage.On(
		"ExecutableExists",
		path.Join(systemConfig.ArtifactsPath(), appToInstall.Name),
	).Return(false)
	mockStorage.On(
		"LinkExists",
		path.Join(*systemConfig.BinaryDir, appToInstall.Name),
	).Return(false)

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
			SourcePath: path.Join(
				selfmanData.SystemConfig.SourcesPath(),
				appToInstall.Name,
				appToInstall.Name,
			),
			DestinationPath: path.Join(selfmanData.SystemConfig.ArtifactsPath(), appToInstall.Name),
		},
		ops.LinkArtifact{
			SourcePath: path.Join(selfmanData.SystemConfig.ArtifactsPath(), appToInstall.Name),
			DestinationPath: path.Join(*selfmanData.SystemConfig.BinaryDir, appToInstall.Name),
		},
	}
	assert.Equal(t, expectedActions, actions)
}
