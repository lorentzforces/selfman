package git

import (
	"os/exec"

	"github.com/lorentzforces/selfman/internal/run"
)

func ExecExists() bool {
	_, err := exec.LookPath("git")
	return err == nil
}

func Clone(url string, destPath string) error {
	return run.NewCmd("git", run.WithArgs("clone", url, destPath), run.WithTimeout(30)).Exec()
}

func Fetch(repoPath string) error {
	return run.NewCmd(
		"git",
		run.WithArgs("fetch", "--tags"),
		run.WithTimeout(30),
		run.WithWorkingDir(repoPath),
	).Exec()
}

func Checkout(repoPath string, ref string) error {
	return run.NewCmd(
		"git",
		run.WithArgs("checkout", ref),
		run.WithTimeout(5),
		run.WithWorkingDir(repoPath),
	).Exec()
}

func RevExists(repoPath string, revName string) bool {
	// appending "^{commit}" should help make sure we don't false-positive on potential other
	// types of revs that git might track
	err := run.NewCmd(
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
