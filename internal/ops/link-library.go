package ops

import (
	"errors"
	"fmt"
	"os"
)

type LinkLibrary struct {
	SourcePath string
	DestinationPath string
}

func (self LinkLibrary) Execute() (string, error) {
	err := os.Remove(self.DestinationPath)
	if err != nil && ! errors.Is(err, os.ErrNotExist) {
		return "", fmt.Errorf("Linking library failed while deleting existing link: %w", err)
	}

	err = os.Symlink(self.SourcePath, self.DestinationPath)
	if err != nil { return "", fmt.Errorf("Linking source as library failed: %w", err) }
	return "Linked app source as library", nil
}

func (self LinkLibrary) Describe() OpDescription {
	topLine := "Link app source as library"
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
