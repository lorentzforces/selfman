package data

import (
	"io/fs"
	"os"
	"path"

	"golang.org/x/sys/unix"
)

type ManagedFiles interface {
	IsGitAppPresent(repoPath string) bool
	DirExistsNotEmpty(path string) bool
	ExecutableExists(path string) bool
	LinkExists(path string) bool
}

type OnDiskManagedFiles struct { }

func (self *OnDiskManagedFiles) IsGitAppPresent(repoPath string) bool {
	gitFilePath := path.Join(repoPath, ".git")
	stat, err := os.Stat(gitFilePath)
	if err == nil && stat.IsDir() {
		return true
	}
	return false
}

func (self *OnDiskManagedFiles) DirExistsNotEmpty(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}
	if !stat.IsDir() {
		return false
	}

	dirContents, err := os.ReadDir(path)
	return len(dirContents) > 0
}

func (self *OnDiskManagedFiles) ExecutableExists(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}
	const anyExecBitmask fs.FileMode = 0111
	if stat.Mode() & anyExecBitmask == 0 {
		return false
	}
	if unix.Access(path, unix.X_OK) != nil {
		return false
	}
	return true
}

func (self *OnDiskManagedFiles) LinkExists(path string) bool {
	stat, err := os.Lstat(path)
	if err != nil {
		return false
	}
	if stat.Mode() & os.ModeSymlink == 0 {
		return false
	}
	return true
}
