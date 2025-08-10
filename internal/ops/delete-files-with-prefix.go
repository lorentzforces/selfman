package ops

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
)

type DeleteFilesWithPrefix struct {
	TypeOfDeletion string
	DirPath string
	FilePrefix string
}

func (self DeleteFilesWithPrefix) Execute() (string, error) {
	stat, err := os.Lstat(self.DirPath)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		return "", fmt.Errorf("Provided directory path does not exist: %w", err)
	} else if err != nil {
		return "", fmt.Errorf("Delete multiple files failed: %w", err)
	}

	if !stat.IsDir() {
		return "", fmt.Errorf("Went to delete files in dir, but found non-dir: %s", self.DirPath)
	}

	files, err := os.ReadDir(self.DirPath)
	if err != nil {
		return "", fmt.Errorf("Delete multiple files failed: %w", err)
	}

	deleteErrors := make([]error, 0)

	for _, file := range files {
		if !strings.HasPrefix(file.Name(), self.FilePrefix) {
			continue
		}
		err = os.Remove(path.Join(self.DirPath, file.Name()))
		if err != nil {
			deleteErrors = append(deleteErrors, err)
		}
	}

	if len(deleteErrors) > 0 {
		return "Deleted multiple files with errors", errors.Join(deleteErrors...)
	}

	return "Deleted multiple files", nil
}

func (self DeleteFilesWithPrefix) Describe() OpDescription {
	return OpDescription{
		TopLine: fmt.Sprintf("File deletion for prefix: %s", self.TypeOfDeletion),
		ContextLines: []string{
			fmt.Sprintf("dir path: %s", self.DirPath),
			fmt.Sprintf("file prefix for bulk delete: %s", self.FilePrefix),
		},
	}
}
