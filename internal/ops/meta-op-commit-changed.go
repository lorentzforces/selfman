package ops

import (
	"fmt"
	"strings"

	"github.com/lorentzforces/selfman/internal/git"
)

type MetaOpCommitChanged struct {
	RepoPath string
	OrigCommitHash string
	IfChangedOps []Operation
}

func (self MetaOpCommitChanged) Execute() (string, error) {
	hash, err := git.CurrentHeadCommit(self.RepoPath)
	if err != nil { return "", fmt.Errorf("Determining current HEAD commit failed: %w", err) }

	if hash == self.OrigCommitHash {
		return "Current and original commit hashes match, successfully did nothing", nil
	}

	var output strings.Builder
	output.WriteString("New commit hash detected, executing conditional operations...")
	for _, op := range self.IfChangedOps {
		opOutput, err := op.Execute()
		if err != nil {
			output.WriteString("\nStep failed")
			if len(opOutput) > 0 {
				output.WriteString(": " + opOutput)
			}
			return output.String(), err
		}

		output.WriteString("\n" + opOutput)
	}

	return output.String(), nil
}

func (self MetaOpCommitChanged) Describe() OpDescription {
	return OpDescription{
		TopLine: "If the head commit changes, execute operations",
		ContextLines: []string{
			fmt.Sprintf("local repository path: %s", self.RepoPath),
			fmt.Sprintf("starting commit: %s", self.OrigCommitHash),
		},
	}
}

func (self MetaOpCommitChanged) InnerOps() []Operation {
	return self.IfChangedOps
}
