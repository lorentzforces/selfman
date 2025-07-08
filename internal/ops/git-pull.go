package ops

import (
	"fmt"

	"github.com/lorentzforces/selfman/internal/git"
)

type GitPull struct {
	RepoPath string
}

func (self GitPull) Execute() (string, error) {
	err := git.Pull(self.RepoPath)
	if err != nil {
		return "", fmt.Errorf("Git pull failed: %w", err)
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
