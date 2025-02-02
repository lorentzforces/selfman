package git

import (
	"context"
	"os/exec"
	"time"
)

func ExecExists() bool {
	_, err := exec.LookPath("git")
	return err == nil
}

func Clone(url string, destPath string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10) * time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "git", "clone", url, destPath)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
