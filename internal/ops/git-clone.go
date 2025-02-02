package ops

import (
	"fmt"

	"github.com/lorentzforces/selfman/internal/git"
)

type GitClone struct {
	RepoUrl string
	DestinationPath string
}

func (self GitClone) Execute() (string, error) {
	err := git.Clone(self.RepoUrl, self.DestinationPath)
	if err != nil {
		return "", fmt.Errorf("Git clone failed: %w", err)
	}
	return "Cloned git repo", nil
}

func (self GitClone) Describe() OpDescription {
	topLine := "Clone git repository"
	urlLine := fmt.Sprintf("repository URL: %s", self.RepoUrl)
	destLine := fmt.Sprintf("destination path: %s", self.DestinationPath)

	return OpDescription{
		TopLine: topLine,
		ContextLines: []string{
			urlLine,
			destLine,
		},
	}
}
