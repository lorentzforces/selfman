package ops

import (
	"fmt"

	"github.com/lorentzforces/selfman/internal/git"
)

type GitCheckoutRef struct {
	RepoPath string
	RefName string
}

func (self GitCheckoutRef) Execute() (string, error) {
	err := git.Checkout(self.RepoPath, self.RefName)
	// TODO: consider figuring out some more graceful way of handling the case where the requested
	// ref just doesn't exist
	if err != nil { return "", fmt.Errorf("Git checkout failed: %w", err) }
	return "Executed git checkout", nil
}

func (self GitCheckoutRef) Describe() OpDescription {
	topLine := "Git checkout ref"
	repoPath := fmt.Sprintf("local repository path: %s", self.RepoPath)
	refName := fmt.Sprintf("ref name: %s", self.RefName)

	return OpDescription {
		TopLine: topLine,
		ContextLines: []string{
			repoPath,
			refName,
		},
	}
}
