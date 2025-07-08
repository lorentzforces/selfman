package ops

import (
	"fmt"
)

type GitCheckoutRef struct {
	RepoPath string
	RefName string
}

func (self GitCheckoutRef) Execute() (string, error) {
	return "", fmt.Errorf("Not yet implemented: GitCheckoutRef::Execute")
}

func (self GitCheckoutRef) Describe() OpDescription {
	topLine := "Git checkout ref"
	repoPath := fmt.Sprintf("Local repository path: %s", self.RepoPath)
	refName := fmt.Sprintf("Ref name: %s", self.RefName)

	return OpDescription {
		TopLine: topLine,
		ContextLines: []string{
			repoPath,
			refName,
		},
	}
}
