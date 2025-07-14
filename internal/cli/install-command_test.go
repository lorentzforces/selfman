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

func TestInstallGitProducesSaneOperations(t *testing.T) {
	systemConfig := data.DefaultTestConfig()

	appToInstall := data.AppConfig{
		SystemConfig: systemConfig,
		Name: "git-repo-app",
		Flavor: "git",
		RemoteRepo: run.StrPtr("git@github.com:github/gitignore.git"),
		BuildAction: "script",
		BuildCmd: run.StrPtr("make build"),
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
			DestinationPath: appToInstall.SourcePath(),
		},
		ops.GitCheckoutRef{
			RepoPath: appToInstall.SourcePath(),
			RefName: appToInstall.Version,
		},
		ops.BuildWithScript{
			SourcePath: appToInstall.SourcePath(),
			ScriptShell: "/bin/sh",
			ScriptCmd: "make build",
		},
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
		ops.GitCheckoutRef{
			RepoPath: appToInstall.SourcePath(),
			RefName: appToInstall.Version,
		},
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

func TestInstallWebFetchProducesSaneOperations(t *testing.T) {
	systemConfig := data.DefaultTestConfig()

	appToInstall := data.AppConfig{
		SystemConfig: systemConfig,
		Name: "web-fetch-app",
		Flavor: "web-fetch",
		Version: "1.69.500",
		WebUrl: run.StrPtr("https://example.com/%VERSION%/web-fetch-app-%VERSION%.zip"),
		BuildAction: "script",
		BuildCmd: run.StrPtr("tar -xzf web-fetch-app-*.zip"),
	}

	mockStorage := mocks.MockManagedFiles{}
	mockStorage.On("AppStatus", appToInstall.Name).Return(data.AppStatus{
		IsConfigured: true,
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
		ops.FetchFromWeb{
			SourceUrl: *appToInstall.WebUrl,
			Version: appToInstall.Version,
			DestinationDir: appToInstall.SourcePath(),
		},
		ops.BuildWithScript{
			SourcePath: appToInstall.SourcePath(),
			ScriptShell: "/bin/sh",
			ScriptCmd: "tar -xzf web-fetch-app-*.zip",
		},
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

func TestInstallAppKeepingBinInPlace(t *testing.T) {
	systemConfig := data.DefaultTestConfig()

	appToInstall := data.AppConfig{
		SystemConfig: systemConfig,
		Name: "no-move-bin-app",
		Flavor: "web-fetch",
		Version: "1.0.0",
		WebUrl: run.StrPtr("https://example.com/%VERSION%/app.zip"),
		BuildAction: "script",
		BuildCmd: run.StrPtr("exit 0"),
		KeepBinWithSource: true,
	}

	mockStorage := mocks.MockManagedFiles{}
	mockStorage.On("AppStatus", appToInstall.Name).Return(data.AppStatus{
		IsConfigured: true,
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
		ops.FetchFromWeb{
			SourceUrl: *appToInstall.WebUrl,
			Version: appToInstall.Version,
			DestinationDir: appToInstall.SourcePath(),
		},
		ops.BuildWithScript{
			SourcePath: appToInstall.SourcePath(),
			ScriptShell: "/bin/sh",
			ScriptCmd: *appToInstall.BuildCmd,
		},
		ops.LinkArtifact{
			SourcePath: path.Join(appToInstall.SourcePath(), appToInstall.Name),
			DestinationPath: path.Join(*selfmanData.SystemConfig.BinaryDir, appToInstall.Name),
		},
	}
	assert.Equal(t, expectedActions, actions)
}
