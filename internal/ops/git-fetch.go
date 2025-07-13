package ops

import (
	"fmt"

	"github.com/lorentzforces/selfman/internal/git"
)

type GitFetch struct {
	RepoPath string
}

func (self GitFetch) Execute() (string, error) {
	err := git.Fetch(self.RepoPath)
	if err != nil { return "", fmt.Errorf("Git fetch failed: %w", err) }
	return "Executed git fetch", nil
}

func (self GitFetch) Describe() OpDescription {
	topLine := "Git fetch latest source code"
	repoPath := fmt.Sprintf("local repository path: %s", self.RepoPath)

	return OpDescription{
		TopLine: topLine,
		ContextLines: []string{
			repoPath,
		},
	}
}
