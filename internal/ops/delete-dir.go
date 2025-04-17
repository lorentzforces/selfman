package ops

import (
	"errors"
	"fmt"
	"os"
)

type DeleteDir struct {
	TypeOfDeletion string
	Path string
}

func (self DeleteDir) Execute() (string, error) {
	stat, err := os.Stat(self.Path)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		// TODO: should this be a success if the directory already doesn't exist?
		return "", fmt.Errorf("Dir to delete does not exist: %s", self.Path)
	} else if err != nil {
		return "", fmt.Errorf("Delete dir failed: %w", err)
	}

	if !stat.IsDir() {
		return "", fmt.Errorf("Went to delete dir, but found file: %s", self.Path)
	}

	err = os.Remove(self.Path)
	if err != nil {
		return "", fmt.Errorf("Delete dir failed: %w", err)
	}

	return "Deleted dir", nil
}

func (self DeleteDir) Describe() OpDescription {
	return OpDescription{
		TopLine: fmt.Sprintf("Dir deletion: %s", self.TypeOfDeletion),
		ContextLines: []string{
			fmt.Sprintf("path: %s", self.Path),
		},
	}
}
