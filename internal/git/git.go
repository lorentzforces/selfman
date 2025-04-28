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

func Pull() error {
	return run.NewCmd("git", run.WithArgs("pull"), run.WithTimeout(30)).Exec()
}

func RevExists(revName string) bool {
	// appending "^{commit}" should help make sure we don't false-positive on potential other
	// types of revs that git might track
	err := run.NewCmd(
		"git",
		run.WithArgs("rev-parse", "--verify", "--end-of-options", revName + "^{commit}"),
	).Exec()
	// really hacky - if we don't find a rev
	if err != nil {
		err = run.NewCmd(
			"git",
			run.WithArgs(
				"rev-parse", "--verify", "--end-of-options",
				"origin/" + revName + "^{commit}",
			),
		).Exec()
	}

	return err == nil
}
