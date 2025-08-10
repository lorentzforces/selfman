package ops

import (
	"fmt"

	"github.com/lorentzforces/selfman/internal/run"
)

type MoveTarget struct {
	SourcePath string
	DestinationPath string
}

func (self MoveTarget) Execute() (string, error) {
	err := run.MoveFile(self.SourcePath, self.DestinationPath)
	if err != nil {
		return "", fmt.Errorf("Target move failed: %w", err)
	}
	return "Moved target", nil
}

func (self MoveTarget) Describe() OpDescription {
	topLine := "Move app target"
	fromLine := fmt.Sprintf("from: %s", self.SourcePath)
	toLine := fmt.Sprintf("to: %s", self.DestinationPath)

	return OpDescription{
		TopLine: topLine,
		ContextLines: []string{
			fromLine,
			toLine,
		},
	}
}
