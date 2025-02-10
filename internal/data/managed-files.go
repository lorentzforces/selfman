package data

import (
	"os"
	"path"
)

type ManagedFiles interface {
	IsGitAppPresent(appPath string) bool
}

type OnDiskManagedFiles struct { }

func (self *OnDiskManagedFiles) IsGitAppPresent(appPath string) bool {
	gitFilePath := path.Join(appPath, ".git")
	stat, err := os.Stat(gitFilePath)
	if err == nil && stat.IsDir() {
		return true
	}
	return false
}
