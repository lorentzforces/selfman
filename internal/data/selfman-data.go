package data

import (
	"errors"
	"fmt"
	"path"

	"github.com/lorentzforces/selfman/internal/run"
)

// TODO: most external things don't actually care about app configs, so provide a "is app
// configured" method
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

	appConfigs, err := loadAppConfigs(*systemConfig.AppConfigDir)
	if err != nil {
		return Selfman{}, err
	}

	return SelfmanFromValues(&systemConfig, appConfigs, &OnDiskManagedFiles{})
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

	appType := MustBeAppType(foundApp.Type)
	appSourcePath := self.AppSourcePath(appName)

	switch appType {
	case "git": {
		appPresent := self.Storage.IsGitAppPresent(appSourcePath)
		if appPresent {
			return AppStatusPresent
		} else {
			return AppStatusNotPresent
		}
	}
	}

	run.FailOut(fmt.Sprintf("Undetermined case for app name: %s", appName))
	panic("unreachable")
}

func (self Selfman) AppSourcePath(appName string) string {
	_, present := self.AppConfigs[appName]
	run.Assert(present, fmt.Sprintf("Invalid app name: %s", appName))
	return path.Join(self.SystemConfig.SourcesPath(), appName)
}

func (self Selfman) AppArtifactPath(appName string) string {
	_, present := self.AppConfigs[appName]
	run.Assert(present, fmt.Sprintf("Invalid app name: %s", appName))
	return path.Join(self.SystemConfig.ArtifactsPath(), appName)
}

func (self Selfman) AppBuildTargetPath(appName string) string {
	app, present := self.AppConfigs[appName]
	run.Assert(present, fmt.Sprintf("Invalid app name: %s", appName))
	return path.Join(self.AppSourcePath(appName), app.BuildTarget)
}
