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
	if err != nil { return Selfman{}, err }

	appConfigs, err := loadAppConfigs(&systemConfig)
	if err != nil { return Selfman{}, err }

	selfman, err := SelfmanFromValues(&systemConfig, appConfigs, &OnDiskManagedFiles{})
	if err != nil { return Selfman{}, err }

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

type AppStatus struct {
	IsConfigured bool
	SourcePresent bool
	TargetPresent bool
	LinkPresent bool
}

func (self AppStatus) FullyPresent() bool {
	return self.IsConfigured && self.SourcePresent && self.TargetPresent && self.LinkPresent
}

func (self AppStatus) ConsistentState() bool {
	return self.FullyPresent() || (
		self.SourcePresent == self.TargetPresent &&
		self.TargetPresent == self.LinkPresent)
}

const (
	AppStatusLinkPresent = "installed & linked"
	AppStatusInconsistent = "partially present - inconsistent state"
	AppStatusIsConfigured = "not present"
	AppStatusNotConfigured = "unknown app - not configured"
)

func (self AppStatus) Label() string {
	switch {
	case self.FullyPresent(): return AppStatusLinkPresent
	case !self.ConsistentState(): return AppStatusInconsistent
	case self.IsConfigured: return AppStatusIsConfigured
	default: return AppStatusNotConfigured
	}
}

// Given an app name, returns a detailed status regarding that app. If the app has a valid
// configuration, will return
//
// NOTE: If AppStatus.IsConfigured is false, then AppConfig will be an invalid value.
func (self Selfman) AppStatus(appName string) (AppConfig, AppStatus) {
	foundApp, present := self.AppConfigs[appName]
	if !present { return AppConfig{}, AppStatus{} }

	statusReport := AppStatus{}
	statusReport.IsConfigured = true
	if foundApp.Type == appTypeGit {
		statusReport.SourcePresent = self.Storage.IsGitAppPresent(foundApp.SourcePath())
	} else {
		statusReport.SourcePresent = self.Storage.DirExistsNotEmpty(foundApp.SourcePath())
	}

	statusReport.TargetPresent = self.Storage.ExecutableExists(foundApp.ArtifactPath())
	statusReport.LinkPresent = self.Storage.LinkExists(foundApp.BinaryPath())

	return foundApp, statusReport
}

func (self Selfman) VerifyAllDirectoriesExist() {
	run.VerifyDirExists(*self.SystemConfig.DataDir)
	run.VerifyDirExists(*self.SystemConfig.AppConfigDir)
	run.VerifyDirExists(*self.SystemConfig.BinaryDir)
	run.VerifyDirExists(self.SystemConfig.ArtifactsPath())
}
