package cli

import (
	"testing"

	"github.com/lorentzforces/selfman/internal/config"
	"github.com/lorentzforces/selfman/internal/ops"
	"github.com/lorentzforces/selfman/internal/run"
	"github.com/stretchr/testify/assert"
)

func TestInstallCommandValidatesNameExists(t *testing.T) {
	testConfig := config.Config{
		AppConfigs: []config.AppConfig{
			{
				Name: "test-app-config",
			},
		},
	}

	_, err := installApp("not-available-name", testConfig)
	assert.Error(
		t, err,
		"An error must be thrown if the user attempts to install an application that is not " +
			"configured",
	)
}

func TestInstallCommandProducesSaneOperations(t *testing.T) {
	testConfig := config.Config {
		AppConfigs: []config.AppConfig{
			{
				Name: "git-repo-app",
				Type: "git",
				RemoteRepo: "git@github.com:github/gitignore.git",
			},
		},
	}

	actions, err := installApp("git-repo-app", testConfig)
	assert.NoError(t, err)
	run.BailIfFailed(t)
	expectedActions := []ops.Operation{
		&ops.GitClone{
			RepoUrl: testConfig.AppConfigs[0].RemoteRepo,
			DestinationPath: "asdf",
		},
	}
	assert.Equal(t, expectedActions, actions)
}
