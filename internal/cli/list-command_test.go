package cli

import (
	"testing"

	"github.com/lorentzforces/selfman/internal/data"
	"github.com/lorentzforces/selfman/internal/data/mocks"
	"github.com/lorentzforces/selfman/internal/run"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAppStatusesAreReflected(t *testing.T) {
	systemConfig := data.DefaultTestConfig()

	presentApp := data.AppConfig{
		SystemConfig: systemConfig,
		Name: "PresentApp",
		Flavor: "git",
		RemoteRepo: run.StrPtr("doesn't matter"),
		BuildAction: "none",
	}

	notPresentApp := data.AppConfig{
		SystemConfig: systemConfig,
		Name: "NotPresentApp",
		Flavor: "git",
		RemoteRepo: run.StrPtr("doesn't matter"),
		BuildAction: "none",
	}

	inconsistentApp := data.AppConfig{
		SystemConfig: systemConfig,
		Name: "InconsistentApp",
		Flavor: "git",
		RemoteRepo: run.StrPtr("doesn't matter"),
		BuildAction: "none",
	}

	storageData := mocks.StartMockingManagedFiles(systemConfig)
	mockStorage := storageData.SetMocks(
		storageData.GitAppPresent(presentApp.Name, true),
		storageData.ExecutableExists(presentApp.Name, true),
		storageData.LinkExists(presentApp.Name, true),
		storageData.MetaData(presentApp.Name, &data.Meta{}),

		storageData.GitAppPresent(notPresentApp.Name, false),
		storageData.ExecutableExists(notPresentApp.Name, false),
		storageData.LinkExists(notPresentApp.Name, false),
		storageData.MetaData(notPresentApp.Name, &data.Meta{}),

		storageData.GitAppPresent(inconsistentApp.Name, true),
		storageData.ExecutableExists(inconsistentApp.Name, true),
		storageData.LinkExists(inconsistentApp.Name, false),
		storageData.MetaData(inconsistentApp.Name, &data.Meta{}),
	)

	selfmanData, err := data.SelfmanFromValues(
		systemConfig,
		[]data.AppConfig{ presentApp, notPresentApp, inconsistentApp },
		mockStorage,
	)
	assert.NoError(t, err)
	run.BailIfFailed(t)

	results := listApplications(selfmanData)

	expected := []listResult{
		{ name: presentApp.Name, status: data.AppStatusLinkPresent },
		{ name: notPresentApp.Name, status: data.AppStatusIsConfigured },
		{ name: inconsistentApp.Name, status: data.AppStatusInconsistent },
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
	}
	bravoApp := data.AppConfig{
		SystemConfig: systemConfig,
		Name: "bravo",
		Flavor: "git",
		RemoteRepo: run.StrPtr("test"),
		BuildAction: "none",
	}
	charlieApp := data.AppConfig{
		SystemConfig: systemConfig,
		Name: "charlie",
		Flavor: "git",
		RemoteRepo: run.StrPtr("test"),
		BuildAction: "none",
	}
	deltaApp := data.AppConfig{
		SystemConfig: systemConfig,
		Name: "delta",
		Flavor: "git",
		RemoteRepo: run.StrPtr("test"),
		BuildAction: "none",
	}
	foxtrotApp := data.AppConfig{
		SystemConfig: systemConfig,
		Name: "foxtrot",
		Flavor: "git",
		RemoteRepo: run.StrPtr("test"),
		BuildAction: "none",
	}

	mockStorage := mocks.MockManagedFiles{}
	mockStorage.On("IsGitAppPresent", mock.Anything).Return(false)
	mockStorage.On("ExecutableExists", mock.Anything).Return(false)
	mockStorage.On("LinkExists", mock.Anything).Return(false)
	mockStorage.On("GetMetaData", mock.Anything).Return(data.Meta{})

	selfmanData, err := data.SelfmanFromValues(
		systemConfig,
		[]data.AppConfig{deltaApp, bravoApp, charlieApp, alphaApp, foxtrotApp},
		&mockStorage,
	)
	assert.NoError(t, err)
	run.BailIfFailed(t)

	results := listApplications(selfmanData)

	expected := []listResult{
		{ name: alphaApp.Name, status: data.AppStatusIsConfigured },
		{ name: bravoApp.Name, status: data.AppStatusIsConfigured },
		{ name: charlieApp.Name, status: data.AppStatusIsConfigured },
		{ name: deltaApp.Name, status: data.AppStatusIsConfigured },
		{ name: foxtrotApp.Name, status: data.AppStatusIsConfigured },
	}
	assert.Equal(t, expected, results)
}
