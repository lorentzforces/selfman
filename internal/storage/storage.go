package storage

import (
	"os"
	"path"

	"github.com/lorentzforces/selfman/internal/config"
)

// TODO/WIP: I need to figure out the proper interface for stuff that is going to live on-disk.
// Just using the file system is all well and good, but I want to actually write _tests_ for the
// damn thing.

// for a given app:
// - the state


type Apps struct {
	config *config.Config
}

func InitFileSystem()

type AppStatus string
const (
	AppStatusNotConfigured AppStatus = "NOT_CONFIGURED"
	AppStatusNotPresent AppStatus = "NOT_PRESENT"
	AppStatusPresent AppStatus = "PRESENT"
)

func (self Apps) Status(name string) AppStatus {
	var foundApp *config.AppConfig
	for _, app := range self.config.AppConfigs {
		if app.Name == name {
			foundApp = &app
		}
	}

	_ = config.MustBeAppType(foundApp.Type)
	// TODO: figure out whether git app is present
	return AppStatusNotConfigured
}

func isGitAppPresent(appPath string) bool {
	gitFilePath := path.Join(appPath, ".git")
	stat, err := os.Stat(gitFilePath)
	if stat.IsDir() && err != nil {
		return true
	}
	return false
}
