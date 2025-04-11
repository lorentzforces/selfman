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
