package data

import (
	"io/fs"
	"os"
	"path"

	"github.com/lorentzforces/selfman/internal/git"
	"golang.org/x/sys/unix"
)

type ManagedFiles interface {
	AppStatus(appName string) AppStatus
}

type AppManagedFiles struct {
	AppConfigs map[string]AppConfig
}

// TODO(jdtls): consider if there is a binary (and later: library) present, to do some
// back-checking on it and maybe doing no work if the requested version is already fully present
func (self *AppManagedFiles) AppStatus(appName string) AppStatus {
	foundApp, present := self.AppConfigs[appName]
	if !present { return AppStatus{} }

	statusReport := AppStatus{}
	statusReport.IsConfigured = true
	statusReport.DesiredVersion = foundApp.Version

	if foundApp.Flavor == FlavorGit {
		statusReport.SourcePresent = isGitRepoPresent(foundApp.SourcePath())
		if statusReport.SourcePresent {
			statusReport.VersionPresent =
				isGitRevPresent(foundApp.SourcePath(), statusReport.DesiredVersion)
		}
	} else {
		statusReport.SourcePresent = dirExistsNotEmpty(foundApp.SourcePath())
		statusReport.VersionPresent = statusReport.SourcePresent
	}

	statusReport.TargetPresent = executableExists(foundApp.ArtifactPath())
	statusReport.LinkPresent = linkExists(foundApp.BinaryPath())

	return statusReport
}

func isGitRepoPresent(repoPath string) bool {
	gitFilePath := path.Join(repoPath, ".git")
	return dirExistsNotEmpty(gitFilePath)
}

func dirExistsNotEmpty(path string) bool {
	stat, err := os.Stat(path)
	if err != nil { return false }
	if !stat.IsDir() { return false }

	dirContents, err := os.ReadDir(path)
	return len(dirContents) > 0
}

func executableExists(path string) bool {
	stat, err := os.Stat(path)
	if err != nil { return false }
	// rely on the OS mode mask to enforce sanity
	const anyExecBitmask fs.FileMode = 0111
	if stat.Mode() & anyExecBitmask == 0 {
		return false
	}
	if unix.Access(path, unix.X_OK) != nil {
		return false
	}
	return true
}

func linkExists(path string) bool {
	stat, err := os.Lstat(path)
	if err != nil { return false }
	if stat.Mode() & os.ModeSymlink == 0 {
		return false
	}
	return true
}

func isGitRevPresent(repoPath string, rev string) bool {
	present := git.RevExists(repoPath, rev)
	return present
}
