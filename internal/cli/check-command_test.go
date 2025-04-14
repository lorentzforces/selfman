package cli

import (
	"path"
	"testing"

	"github.com/lorentzforces/selfman/internal/data"
	"github.com/lorentzforces/selfman/internal/data/mocks"
	"github.com/lorentzforces/selfman/internal/run"
	"github.com/stretchr/testify/assert"
)

func TestCheckCommandValidatesNameExists(t *testing.T) {
	systemConfig := data.DefaultTestConfig()
	mockStorage := mocks.MockManagedFiles{}
	selfmanData := data.Selfman{
		SystemConfig: systemConfig,
		AppConfigs: map[string]data.AppConfig{},
		Storage: &mockStorage,
	}

	_, err := checkApp("not-available-name", selfmanData)
	assert.Error(
		t, err,
		"An error must be thrown if the user attempts to update an application that is not " +
			"configured",
	)
}

func TestCheckShowsDetailedStatusInformation(t *testing.T) {
	systemConfig := data.DefaultTestConfig()
	mockStorage := mocks.MockManagedFiles{}

	appWithStatus := data.AppConfig{
		SystemConfig: systemConfig,
		Name: "high-status-app",
		Type: "git",
		RemoteRepo: run.StrPtr("doesn't matter"),
		BuildAction: "none",
	}
	mockStorage.On(
		"IsGitAppPresent",
		path.Join(systemConfig.SourcesPath(), appWithStatus.Name),
	).Return(false)
	mockStorage.On(
		"ExecutableExists",
		path.Join(systemConfig.ArtifactsPath(), appWithStatus.Name),
	).Return(true)
	mockStorage.On(
		"LinkExists",
		path.Join(*systemConfig.BinaryDir, appWithStatus.Name),
	).Return(false)

	selfmanData, err := data.SelfmanFromValues(
		systemConfig,
		[]data.AppConfig{ appWithStatus },
		&mockStorage,
	)
	assert.NoError(t, err)
	run.BailIfFailed(t)

	result, err := checkApp(appWithStatus.Name, selfmanData)

	assert.Equal(t, true, result.status.IsConfigured)
	assert.Equal(t, false, result.status.SourcePresent)
	assert.Equal(t, true, result.status.TargetPresent)
	assert.Equal(t, false, result.status.LinkPresent)
}
