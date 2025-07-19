package cli

import (
	"testing"

	"github.com/lorentzforces/selfman/internal/data"
	"github.com/stretchr/testify/assert"
)

// TODO: test install git app from scratch

// TODO: test install web fetch app from scratch

// TODO: test git app even with same version still fetches and updates (update branch)

// TODO: test web fetch app with no changes does nothing

// TODO: test install app with source but nothing else puts everything in its place

// TODO: test an app with a different version available pulls and does all the install stuff

func TestMakeItSoFailsWithNonConfiguredapp(t *testing.T) {

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
