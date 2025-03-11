package data

import (
	"errors"
	"fmt"

	"github.com/lorentzforces/selfman/internal/run"
)

type Selfman struct {
	SystemConfig *SystemConfig
	AppConfigs map[string]AppConfig
	Storage ManagedFiles
}

func Produce() (Selfman, error) {
	systemConfig, err := loadSystemConfig()
	if err != nil {
		return Selfman{}, err
	}

	appConfigs, err := loadAppConfigs(&systemConfig)
	if err != nil {
		return Selfman{}, err
	}

	selfman, err := SelfmanFromValues(&systemConfig, appConfigs, &OnDiskManagedFiles{})
	if err != nil {
		return Selfman{}, err
	}

	selfman.VerifyAllDirectoriesExist()
	return selfman, nil
}

func SelfmanFromValues(
	system *SystemConfig,
	apps []AppConfig,
	storage ManagedFiles,
) (Selfman, error) {
	appConfigMap := make(map[string]AppConfig, len(apps))
	for _, app := range apps {
		app.applyDefaults()
		err := app.validate()
		if err != nil {
			newErr := fmt.Errorf(
				"Invalid app config in app directory \"%s\"",
				*system.AppConfigDir,
			)
			return Selfman{}, errors.Join(newErr, err)
		}
		appConfigMap[app.Name] = app
	}

	return Selfman{
		SystemConfig: system,
		AppConfigs: appConfigMap,
		Storage: storage,
	}, nil
}

type AppStatus string
const (
	AppStatusNotConfigured AppStatus = "not configured"
	AppStatusNotPresent AppStatus = "not present"
	AppStatusPresent AppStatus = "present"
)

func (self Selfman) AppStatus(appName string) AppStatus {
	foundApp, present := self.AppConfigs[appName]
	if !present { return AppStatusNotConfigured }

	appSourcePath := foundApp.SourcePath()
	switch foundApp.Type {
	case appTypeGit: {
		appPresent := self.Storage.IsGitAppPresent(appSourcePath)
		if appPresent {
			return AppStatusPresent
		} else {
			return AppStatusNotPresent
		}
	}
	}

	run.FailOut(fmt.Sprintf("Undetermined case for app name: %s", appName))
	panic("unreachable in theory")
}

func (self Selfman) VerifyAllDirectoriesExist() {
	run.VerifyDirExists(*self.SystemConfig.DataDir)
	run.VerifyDirExists(*self.SystemConfig.AppConfigDir)
	run.VerifyDirExists(*self.SystemConfig.BinaryDir)
	run.VerifyDirExists(self.SystemConfig.ArtifactsPath())
}
