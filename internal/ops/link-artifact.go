package ops

import (
	"errors"
	"fmt"
	"os"
)

type LinkArtifact struct {
	SourcePath string
	DestinationPath string
}

func (self LinkArtifact) Execute() (string, error) {
	err := os.Symlink(self.SourcePath, self.DestinationPath)

	if err != nil && !errors.Is(err, os.ErrExist) {
		return "", fmt.Errorf("Linking artifact failed: %w", err)
	}
	return "Linked artifact", nil
}

func (self LinkArtifact) Describe() OpDescription {
	topLine := "Link app artifact binary"
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
