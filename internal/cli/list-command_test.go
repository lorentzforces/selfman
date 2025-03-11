package cli

import (
	"path"
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
		Type: "git",
		RemoteRepo: run.StrPtr("doesn't matter"),
		BuildAction: "none",
	}
	notPresentApp := data.AppConfig{
		SystemConfig: systemConfig,
		Name: "NotPresentApp",
		Type: "git",
		RemoteRepo: run.StrPtr("doesn't matter"),
		BuildAction: "none",
	}

	mockStorage := mocks.MockManagedFiles{}
	mockStorage.On(
		"IsGitAppPresent",
		path.Join(systemConfig.SourcesPath(), presentApp.Name),
	).Return(true)
	mockStorage.On(
		"IsGitAppPresent",
		path.Join(systemConfig.SourcesPath(), notPresentApp.Name),
	).Return(false)

	selfmanData, err := data.SelfmanFromValues(
		systemConfig,
		[]data.AppConfig{presentApp, notPresentApp},
		&mockStorage,
	)
	assert.NoError(t, err)
	run.BailIfFailed(t)

	results := listApplications(selfmanData)

	expected := []listResult{
		{ name: presentApp.Name, status: data.AppStatusPresent },
		{ name: notPresentApp.Name, status: data.AppStatusNotPresent },
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
		Type: "git",
		RemoteRepo: run.StrPtr("test"),
		BuildAction: "none",
	}
	bravoApp := data.AppConfig{
		SystemConfig: systemConfig,
		Name: "bravo",
		Type: "git",
		RemoteRepo: run.StrPtr("test"),
		BuildAction: "none",
	}
	charlieApp := data.AppConfig{
		SystemConfig: systemConfig,
		Name: "charlie",
		Type: "git",
		RemoteRepo: run.StrPtr("test"),
		BuildAction: "none",
	}
	deltaApp := data.AppConfig{
		SystemConfig: systemConfig,
		Name: "delta",
		Type: "git",
		RemoteRepo: run.StrPtr("test"),
		BuildAction: "none",
	}
	foxtrotApp := data.AppConfig{
		SystemConfig: systemConfig,
		Name: "foxtrot",
		Type: "git",
		RemoteRepo: run.StrPtr("test"),
		BuildAction: "none",
	}

	mockStorage := mocks.MockManagedFiles{}
	mockStorage.On("IsGitAppPresent", mock.Anything).Return(false)

	selfmanData, err := data.SelfmanFromValues(
		systemConfig,
		[]data.AppConfig{deltaApp, bravoApp, charlieApp, alphaApp, foxtrotApp},
		&mockStorage,
	)
	assert.NoError(t, err)
	run.BailIfFailed(t)

	results := listApplications(selfmanData)

	expected := []listResult{
		{ name: alphaApp.Name, status: data.AppStatusNotPresent },
		{ name: bravoApp.Name, status: data.AppStatusNotPresent },
		{ name: charlieApp.Name, status: data.AppStatusNotPresent },
		{ name: deltaApp.Name, status: data.AppStatusNotPresent },
		{ name: foxtrotApp.Name, status: data.AppStatusNotPresent },
	}
	assert.Equal(t, expected, results)
}
