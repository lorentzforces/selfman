package cli

import (
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

	appWithStatus := data.AppConfig{
		SystemConfig: systemConfig,
		Name: "high-status-app",
		Flavor: "git",
		Version: "1.0.0",
		RemoteRepo: run.StrPtr("doesn't matter"),
		BuildAction: "none",
	}

	mockStorage := mocks.MockManagedFiles{}
	mockStorage.On("AppStatus", appWithStatus.Name).Return(data.AppStatus{
		IsConfigured: true,
		SourcePresent: false,
		LinkPresent: false,
		TargetPresent: true,
		DesiredVersion: appWithStatus.Version,
		AvailableVersions: []string{ appWithStatus.Version, "origin/main", "origin/next-version" },
	})

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
	assert.Equal(t,
		[]string{appWithStatus.Version, "origin/main", "origin/next-version"},
		result.status.AvailableVersions,
	)
}

func TestCheckShowsLibLinkStatusIfAppIsLib(t *testing.T) {
	systemConfig := data.DefaultTestConfig()

	libApp := data.AppConfig{
		SystemConfig: systemConfig,
		Name: "lib-app",
		Flavor: "web-fetch",
		Version: "1.0.0",
		WebUrl: run.StrPtr("https://example.com/%VERSION%/app.zip"),
		BuildAction: "none",
		LinkSourceAsLib: true,
	}

	mockStorage := mocks.MockManagedFiles{}
	mockStorage.On("AppStatus", libApp.Name).Return(data.AppStatus{
		IsConfigured: true,
		SourcePresent: true,
		LinkPresent: true,
		TargetPresent: true,
		LibLinkPresent: true,
		DesiredVersion: libApp.Version,
		AvailableVersions: []string{},
	})

	selfmanData, err := data.SelfmanFromValues(
		systemConfig,
		[]data.AppConfig{ libApp },
		&mockStorage,
	)
	assert.NoError(t, err)
	run.BailIfFailed(t)

	result, err := checkApp(libApp.Name, selfmanData)

	assert.Equal(t, true, result.status.LibLinkPresent)
}
