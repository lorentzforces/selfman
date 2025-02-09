package cli

import (
	"path"
	"testing"

	"github.com/lorentzforces/selfman/internal/data"
	"github.com/lorentzforces/selfman/internal/data/mocks"
	"github.com/stretchr/testify/assert"
)

func TestAppStatusesAreReflected(t *testing.T) {
	systemConfig := data.DefaultTestConfig()
	presentApp := data.AppConfig{ Name: "PresentApp", Type: "git" }
	notPresentApp := data.AppConfig{ Name: "NotPresentApp", Type: "git" }

	mockStorage := mocks.MockManagedFiles{}
	mockStorage.On(
		"IsGitAppPresent",
		path.Join(systemConfig.SourcesPath(), presentApp.Name),
	).Return(true)
	mockStorage.On(
		"IsGitAppPresent",
		path.Join(systemConfig.SourcesPath(), notPresentApp.Name),
	).Return(false)

	selfmanData := data.Selfman{
		SystemConfig: systemConfig,
		AppConfigs: map[string]data.AppConfig{
			presentApp.Name: presentApp,
			notPresentApp.Name: notPresentApp,
		},
		Storage: &mockStorage,
	}

	results := listApplications(selfmanData)

	expected := []listResult{
		{ name: presentApp.Name, status: data.AppStatusPresent },
		{ name: notPresentApp.Name, status: data.AppStatusNotPresent },
	}
	// TODO: iterating over a map is nondeterministic, put a deliberate sort in place on the list command
	assert.ElementsMatch(t, expected, results)
}
