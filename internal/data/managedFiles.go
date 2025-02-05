package data

import (
	"os"
	"path"
)

type ManagedFiles interface {
	isGitAppPresent(appPath string) bool
}

type OnDiskManagedFiles struct { }

func (self OnDiskManagedFiles) isGitAppPresent(appPath string) bool {
	gitFilePath := path.Join(appPath, ".git")
	stat, err := os.Stat(gitFilePath)
	if err == nil && stat.IsDir() {
		return true
	}
	return false
}
