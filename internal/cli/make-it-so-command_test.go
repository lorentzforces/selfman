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

func TestMakeItSoFailsWithNonConfiguredApp(t *testing.T) {
	systemConfig := data.DefaultTestConfig()
	selfmanData := data.Selfman{
		SystemConfig: systemConfig,
		AppConfigs: map[string]data.AppConfig{},
		Storage: nil,
	}

	_, err := makeItSo("not-available-name", selfmanData)
	assert.Error(
		t, err,
		"An error must be thrown if the user attempts to update an application which is not " +
			"configured",
	)
}

func TestMakeItSoInstallGitAppFromScratch(t *testing.T) {
	systemConfig := data.DefaultTestConfig()

	appToInstall := data.AppConfig{
		SystemConfig: systemConfig,
		Name: "git-repo-app",
		Flavor: data.FlavorGit,
		RemoteRepo: run.StrPtr("git@github.com:github/gitignore.git"),
		BuildAction: data.BuildActionScript,
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

	actions, err := makeItSo(appToInstall.Name, selfmanData)
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

func TestMakeItSoInstallWebFetchAppFromScratch(t *testing.T) {
	systemConfig := data.DefaultTestConfig()

	appToInstall := data.AppConfig{
		SystemConfig: systemConfig,
		Name: "web-fetch-app",
		Flavor: data.FlavorWebFetch,
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

	actions, err := makeItSo(appToInstall.Name, selfmanData)
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

func TestMakeItSoGitWithPresentVersionStillFetches(t *testing.T) {
	systemConfig := data.DefaultTestConfig()
	gitApp := data.AppConfig{
		SystemConfig: systemConfig,
		Name: "git-app-with-version",
		Flavor: data.FlavorGit,
		Version: "origin/main",
		RemoteRepo: run.StrPtr("git@github.com:github/gitignore.git"),
		BuildAction: data.ActionNone,
	}

	mockStorage := mocks.MockManagedFiles{}
	mockStorage.On("AppStatus", gitApp.Name).Return(data.AppStatus{
		IsConfigured: true,
		SourcePresent: true,
		VersionPresent: true,
		TargetPresent: false,
		LinkPresent: false,
	})

	selfmanData, err := data.SelfmanFromValues(
		systemConfig,
		[]data.AppConfig{ gitApp },
		&mockStorage,
	)
	assert.NoError(t, err)
	run.BailIfFailed(t)

	actions, err := makeItSo(gitApp.Name, selfmanData)
	assert.NoError(t, err)
	run.BailIfFailed(t)

	expectedActions := []ops.Operation{
		ops.GitFetch{
			RepoPath: gitApp.SourcePath(),
		},
		ops.GitCheckoutRef{
			RepoPath: gitApp.SourcePath(),
			RefName: gitApp.Version,
		},
		ops.NoBuildOp,
		ops.MoveTarget{
			SourcePath: path.Join(gitApp.SourcePath(), gitApp.Name),
			DestinationPath: gitApp.ArtifactPath(),
		},
		ops.LinkArtifact{
			SourcePath: gitApp.ArtifactPath(),
			DestinationPath: path.Join(*selfmanData.SystemConfig.BinaryDir, gitApp.Name),
		},
	}
	assert.Equal(t, expectedActions, actions)
}

func TestMakeItSoWebFetchWithNoChangesDoesNothing(t *testing.T) {
	systemConfig := data.DefaultTestConfig()

	unchangedApp := data.AppConfig{
		SystemConfig: systemConfig,
		Name: "unchanged-web-fetch-app",
		Flavor: data.FlavorWebFetch,
		Version: "1.0.0",
		WebUrl: run.StrPtr("https://example.com/%VERSION%/web-fetch-app-%VERSION%.zip"),
		BuildAction: data.BuildActionScript,
		BuildCmd: run.StrPtr("exit 0"),
	}

	mockStorage := mocks.MockManagedFiles{}
	mockStorage.On("AppStatus", unchangedApp.Name).Return(data.AppStatus{
		IsConfigured: true,
		SourcePresent: true,
		VersionPresent: true,
		TargetPresent: true,
		LinkPresent: true,
	})

	selfmanData, err := data.SelfmanFromValues(
		systemConfig,
		[]data.AppConfig{ unchangedApp },
		&mockStorage,
	)
	assert.NoError(t, err)
	run.BailIfFailed(t)

	actions, err := makeItSo(unchangedApp.Name, selfmanData)
	assert.NoError(t, err)
	run.BailIfFailed(t)

	expectedActions := []ops.Operation{
		// linking the artifact is an unconditional easy operation
		ops.LinkArtifact{
			SourcePath: path.Join(unchangedApp.ArtifactPath()),
			DestinationPath: path.Join(*selfmanData.SystemConfig.BinaryDir, unchangedApp.Name),
		},
	}
	assert.Equal(t, expectedActions, actions)
}

func TestMakeItSoWithSourceAndNothingElseDoesntFetch(t *testing.T) {
	systemConfig := data.DefaultTestConfig()

	appToInstall := data.AppConfig{
		SystemConfig: systemConfig,
		Name: "web-fetch-app",
		Flavor: data.FlavorWebFetch,
		Version: "1.69.500",
		WebUrl: run.StrPtr("https://example.com/%VERSION%/web-fetch-app-%VERSION%.zip"),
		BuildAction: "script",
		BuildCmd: run.StrPtr("tar -xzf web-fetch-app-*.zip"),
	}

	mockStorage := mocks.MockManagedFiles{}
	mockStorage.On("AppStatus", appToInstall.Name).Return(data.AppStatus{
		IsConfigured: true,
		SourcePresent: true,
		VersionPresent: true,
	})

	selfmanData, err := data.SelfmanFromValues(
		systemConfig,
		[]data.AppConfig{ appToInstall },
		&mockStorage,
	)
	assert.NoError(t, err)
	run.BailIfFailed(t)

	actions, err := makeItSo(appToInstall.Name, selfmanData)
	assert.NoError(t, err)
	run.BailIfFailed(t)

	expectedActions := []ops.Operation{
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

func TestMakeItSoKeepingbinInPlace(t *testing.T) {
	systemConfig := data.DefaultTestConfig()

	inPlaceApp := data.AppConfig{
		SystemConfig: systemConfig,
		Name: "no-move-bin-app",
		Flavor: data.FlavorWebFetch,
		Version: "1.0.0",
		WebUrl: run.StrPtr("https://example.com/%VERSION%/app.zip"),
		BuildAction: data.BuildActionScript,
		BuildCmd: run.StrPtr("exit 0"),
		KeepBinWithSource: true,
	}

	mockStorage := mocks.MockManagedFiles{}
	mockStorage.On("AppStatus", inPlaceApp.Name).Return(data.AppStatus{
		IsConfigured: true,
	})

	selfmanData, err := data.SelfmanFromValues(
		systemConfig,
		[]data.AppConfig{ inPlaceApp },
		&mockStorage,
	)
	assert.NoError(t, err)
	run.BailIfFailed(t)

	actions, err := makeItSo(inPlaceApp.Name, selfmanData)
	assert.NoError(t, err)
	run.BailIfFailed(t)

	expectedActions := []ops.Operation{
		ops.FetchFromWeb{
			SourceUrl: *inPlaceApp.WebUrl,
			Version: inPlaceApp.Version,
			DestinationDir: inPlaceApp.SourcePath(),
		},
		ops.BuildWithScript{
			SourcePath: inPlaceApp.SourcePath(),
			ScriptShell: "/bin/sh",
			ScriptCmd: *inPlaceApp.BuildCmd,
		},
		ops.LinkArtifact{
			SourcePath: path.Join(inPlaceApp.SourcePath(), inPlaceApp.Name),
			DestinationPath: path.Join(*selfmanData.SystemConfig.BinaryDir, inPlaceApp.Name),
		},
	}
	assert.Equal(t, expectedActions, actions)
}
