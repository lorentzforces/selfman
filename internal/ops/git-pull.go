package ops

import (
	"fmt"
	"os"

	"github.com/lorentzforces/selfman/internal/git"
)

type GitPull struct {
	RepoPath string
}

func (self GitPull) Execute() (string, error) {
	oldWorkingDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("Determining working dir failed: %w", err)
	}

	err = os.Chdir(self.RepoPath)
	if err != nil {
		return "", fmt.Errorf("Changing to repo dir for pull failed: %w", err)
	}

	err = git.Pull()
	if err != nil {
		return "", fmt.Errorf("Git pull failed: %w", err)
	}

	err = os.Chdir(oldWorkingDir)
	if err != nil {
		return "", fmt.Errorf("Failed to reset working dir after running build script: %w", err)
	}

	return "Executed git pull", nil
}

func (self GitPull) Describe() OpDescription {
	topLine := "Git pull latest source code"
	repoPath := fmt.Sprintf("local repository path: %s", self.RepoPath)

	return OpDescription{
		TopLine: topLine,
		ContextLines: []string{
			repoPath,
		},
	}
}
