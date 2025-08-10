package git

import (
	"os/exec"
	"strings"

	"github.com/lorentzforces/selfman/internal/run"
)

func ExecExists() bool {
	_, err := exec.LookPath("git")
	return err == nil
}

func Clone(url string, destPath string) error {
	_, err := run.NewCmd("git", run.WithArgs("clone", url, destPath), run.WithTimeout(30)).Exec()
	return err
}

func Fetch(repoPath string) error {
	_, err := run.NewCmd(
		"git",
		run.WithArgs("fetch", "--tags"),
		run.WithTimeout(30),
		run.WithWorkingDir(repoPath),
	).Exec()

	return err
}

func Checkout(repoPath string, ref string) error {
	_, err := run.NewCmd(
		"git",
		run.WithArgs("checkout", ref),
		run.WithTimeout(5),
		run.WithWorkingDir(repoPath),
	).Exec()
	return err
}

// Returns the list of human-named revs (i.e. branches and tags)
func GetAllNamedRevs(repoPath string) ([]string, error) {
	tagOutput, err := run.NewCmd("git", run.WithArgs("tag")).Exec()
	if err != nil { return nil, err }

	branchOutput, err := run.NewCmd(
		"git",
		run.WithArgs("branch", "--all", "--format=%(refname:lstrip=2)"),
		run.WithWorkingDir(repoPath),
	).Exec()
	if err != nil { return nil, err }

	branches := removeFinalEmptyString(strings.Split(branchOutput, "\n"))

	refNames := removeFinalEmptyString(strings.Split(tagOutput, "\n"))
	for _, branchName := range branches {
		isIgnoredBranch :=
			strings.HasSuffix(branchName, "/HEAD") ||
			strings.HasPrefix(branchName, "(HEAD detached")
		if isIgnoredBranch {
			continue
		}
		refNames = append(refNames, branchName)
	}
	return refNames, nil
}

func removeFinalEmptyString(strs []string) []string {
	if strs[len(strs) - 1] == "" {
		strs = strs[0:len(strs) - 1]
	}
	return strs
}

func RevExists(repoPath string, revName string) bool {
	// appending "^{commit}" should help make sure we don't false-positive on potential other
	// types of revs that git might track
	_, err := run.NewCmd(
		"git",
		run.WithArgs(
			"rev-parse",
			"-c",
			repoPath,
			"--verify",
			"--end-of-options",
			revName + "^{commit}",
		),
	).Exec()

	return err == nil
}
