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

func (self *AppManagedFiles) AppStatus(appName string) AppStatus {
	foundApp, present := self.AppConfigs[appName]
	if !present { return AppStatus{} }

	statusReport := AppStatus{}
	statusReport.IsConfigured = true
	statusReport.DesiredVersion = foundApp.Version

	if foundApp.Flavor == FlavorGit {
		statusReport.SourcePresent = isGitRepoPresent(foundApp.SourcePath())
		// result value will be nil if there's an error
		// TODO: right now we're just munching the error... log it?
		statusReport.AvailableVersions, _ = git.GetAllNamedRevs(foundApp.SourcePath())
	} else {
		statusReport.SourcePresent = dirExistsNotEmpty(foundApp.SourcePath())
		statusReport.AvailableVersions =
			getSourceVersions(foundApp.SystemConfig.SourcesPath(), foundApp.Name)
	}

	statusReport.TargetPresent = executableExists(foundApp.ArtifactPath())
	statusReport.LinkPresent = linkExists(foundApp.BinaryPath())
	statusReport.LibLinkPresent = linkExists(foundApp.LibPath())

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

func getSourceVersions(sourcesTopPath string, appName string) []string {
	appVersionsDirPath := path.Join(sourcesTopPath, appName)
	entries, err := os.ReadDir(appVersionsDirPath)
	if err != nil { return nil }

	results := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		// we only consider non-empty dirs just in case a directory got created for an operation
		// which ended up failing
		versionDirPath := path.Join(appVersionsDirPath, entry.Name())
		versionDirContents, err := os.ReadDir(versionDirPath)
		if err != nil { continue }
		if len(versionDirContents) > 0 {
			results = append(results, entry.Name())
		}
	}
	return results
}
