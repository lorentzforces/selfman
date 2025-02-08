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
	cmd := run.NewCmd("git", run.WithArgs("clone", url, destPath), run.WithTimeout(10))
	return cmd.Exec()
}
