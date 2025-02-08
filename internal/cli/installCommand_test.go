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
		RemoteRepo: "git@github.com:github/gitignore.git",
	}

	selfmanData := data.Selfman{
		SystemConfig: systemConfig,
		AppConfigs: map[string]data.AppConfig{
			appToInstall.Name: appToInstall,
		},
		Storage: nil,
	}

	actions, err := installApp("git-repo-app", selfmanData)
	assert.NoError(t, err)
	run.BailIfFailed(t)
	expectedActions := []ops.Operation{
		&ops.GitClone{
			RepoUrl: selfmanData.AppConfigs["git-repo-app"].RemoteRepo,
			DestinationPath: path.Join(*selfmanData.SystemConfig.AppSourceDir, "git-repo-app"),
		},
	}
	assert.Equal(t, expectedActions, actions)
}
