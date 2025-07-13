package cli

import (
	"testing"

	"github.com/lorentzforces/selfman/internal/data"
	"github.com/lorentzforces/selfman/internal/data/mocks"
	"github.com/lorentzforces/selfman/internal/run"
	"github.com/stretchr/testify/assert"
)

func TestAppStatusesAreReflected(t *testing.T) {
	systemConfig := data.DefaultTestConfig()

	presentApp := data.AppConfig{
		SystemConfig: systemConfig,
		Name: "PresentApp",
		Flavor: "git",
		RemoteRepo: run.StrPtr("doesn't matter"),
		BuildAction: "none",
		Version: "main",
	}

	notPresentApp := data.AppConfig{
		SystemConfig: systemConfig,
		Name: "NotPresentApp",
		Flavor: "git",
		RemoteRepo: run.StrPtr("doesn't matter"),
		BuildAction: "none",
		Version: "main",
	}

	inconsistentApp := data.AppConfig{
		SystemConfig: systemConfig,
		Name: "InconsistentApp",
		Flavor: "git",
		RemoteRepo: run.StrPtr("doesn't matter"),
		BuildAction: "none",
		Version: "main",
	}

	mockStorage := mocks.MockManagedFiles{}
	mockStorage.On("AppStatus", presentApp.Name).Return(data.AppStatus{
		IsConfigured: true,
		SourcePresent: true,
		TargetPresent: true,
		LinkPresent: true,
	})
	mockStorage.On("AppStatus", notPresentApp.Name).Return(data.AppStatus{
		IsConfigured: true,
	})
	mockStorage.On("AppStatus", inconsistentApp.Name).Return(data.AppStatus{
		IsConfigured: true,
		SourcePresent: true,
		TargetPresent: true,
		LinkPresent: false,
	})

	selfmanData, err := data.SelfmanFromValues(
		systemConfig,
		[]data.AppConfig{ presentApp, notPresentApp, inconsistentApp },
		&mockStorage,
	)
	assert.NoError(t, err)
	run.BailIfFailed(t)

	results := listApplications(selfmanData)

	expected := []listResult{
		{ name: presentApp.Name, version: "main", status: data.AppStatusLinkPresent },
		{ name: notPresentApp.Name, version: "main", status: data.AppStatusIsConfigured },
		{ name: inconsistentApp.Name, version: "main", status: data.AppStatusInconsistent },
	}
	assert.ElementsMatch(t, expected, results)
}

func TestAppsAreSortedInLexicalOrder(t *testing.T) {
	systemConfig := data.DefaultTestConfig()
	// there are a lot of configs here, but it's so it (hopefully) never gets sorted randomly
	// in the case that tested code is just iterating over the map values
	alphaApp := data.AppConfig{
		SystemConfig: systemConfig,
		Name: "alpha",
		Flavor: "git",
		RemoteRepo: run.StrPtr("test"),
		BuildAction: "none",
		Version: "main",
	}
	bravoApp := data.AppConfig{
		SystemConfig: systemConfig,
		Name: "bravo",
		Flavor: "git",
		RemoteRepo: run.StrPtr("test"),
		BuildAction: "none",
		Version: "main",
	}
	charlieApp := data.AppConfig{
		SystemConfig: systemConfig,
		Name: "charlie",
		Flavor: "git",
		RemoteRepo: run.StrPtr("test"),
		BuildAction: "none",
		Version: "main",
	}
	deltaApp := data.AppConfig{
		SystemConfig: systemConfig,
		Name: "delta",
		Flavor: "git",
		RemoteRepo: run.StrPtr("test"),
		BuildAction: "none",
		Version: "main",
	}
	foxtrotApp := data.AppConfig{
		SystemConfig: systemConfig,
		Name: "foxtrot",
		Flavor: "git",
		RemoteRepo: run.StrPtr("test"),
		BuildAction: "none",
		Version: "main",
	}

	mockStorage := mocks.MockManagedFiles{}

	selfmanData, err := data.SelfmanFromValues(
		systemConfig,
		[]data.AppConfig{deltaApp, bravoApp, charlieApp, alphaApp, foxtrotApp},
		&mockStorage,
	)

	for appName, _ := range selfmanData.AppConfigs {
		mockStorage.On("AppStatus", appName).Return(data.AppStatus{
			IsConfigured: true,
		})
	}

	assert.NoError(t, err)
	run.BailIfFailed(t)

	results := listApplications(selfmanData)

	expected := []listResult{
		{ name: alphaApp.Name, version: "main", status: data.AppStatusIsConfigured },
		{ name: bravoApp.Name, version: "main", status: data.AppStatusIsConfigured },
		{ name: charlieApp.Name, version: "main", status: data.AppStatusIsConfigured },
		{ name: deltaApp.Name, version: "main", status: data.AppStatusIsConfigured },
		{ name: foxtrotApp.Name, version: "main", status: data.AppStatusIsConfigured },
	}
	assert.Equal(t, expected, results)
}
