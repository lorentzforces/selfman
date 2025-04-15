package ops

import (
	"errors"
	"fmt"
	"os"
)

type DeleteFile struct {
	TypeOfDeletion string
	Path string
}

func (self DeleteFile) Execute() (string, error) {
	stat, err := os.Stat(self.Path)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		return "", fmt.Errorf("File to delete does not exist: %s", self.Path)
	} else if err != nil {
		return "", fmt.Errorf("Delete file failed: %w", err)
	}

	if stat.IsDir() {
		return "", fmt.Errorf("Went to delete file, but found dir: %s", self.Path)
	}

	err = os.Remove(self.Path)
	if err != nil {
		return "", fmt.Errorf("Delete file failed: %w", err)
	}

	return "Deleted file", nil
}

func (self DeleteFile) Describe() OpDescription {
	return OpDescription{
		TopLine: fmt.Sprintf("File deletion: %s", self.TypeOfDeletion),
		ContextLines: []string{
			fmt.Sprintf("path: %s", self.Path),
		},
	}
}
