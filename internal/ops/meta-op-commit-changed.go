package ops

import "fmt"

type MetaOpCommitChanged struct {
	OrigCommitHash string
	IfChangedOp Operation
}

// TODO(commit-changed): execution implementation
func (self MetaOpCommitChanged) Execute() (string, error) {
	return "", fmt.Errorf("MetaOpCommitChanged::Execute - not yet implemented")
}

func (self MetaOpCommitChanged) Describe() OpDescription {
	return OpDescription{
		TopLine: "Execute if the current commit changes",
		ContextLines: []string{
			fmt.Sprintf("starting commit: %s", self.OrigCommitHash),
		},
	}
}

func (self MetaOpCommitChanged) InnerOps() []Operation {
	return []Operation{ self.IfChangedOp }
}
