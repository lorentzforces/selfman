package ops

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/lorentzforces/selfman/internal/run"
)

type FetchFromWeb struct {
	SourceUrl string
	Version string
	DestinationDir string
}

func (self FetchFromWeb) Execute() (string, error) {
	fullUrl := strings.ReplaceAll(self.SourceUrl, "%VERSION%", self.Version)

	tmpFile, err := run.GetFileFromUrl(fullUrl)
	if err != nil { return "", fmt.Errorf("Fetch from web failed: %w", err) }

	err = run.VerifyDirExists(self.DestinationDir)
	if err != nil {
		return "", fmt.Errorf("Error creating destination dir (%s): %w", self.DestinationDir, err)
	}

	destPath := path.Join(self.DestinationDir, path.Base(tmpFile))
	err = os.Rename(tmpFile, destPath)
	// os.Rename notably does not work across filesystem boundaries on Unix (and some other cases),
	// so attempt to do a copy-and-remove operation if that doesn't work (but renames are way
	// faster so we try that first)
	if err != nil && isIncompatibleRenameError(err) {
		err = moveFileWithCopy(tmpFile, destPath)
	}
	if err != nil {
		return "", fmt.Errorf("Error moving fetched file: %w", err)
	}

	return "Fetched app from the web", nil
}

// Attempt to determine if the given error is from a rename operation which is not supported on the
// current platform for whatever reason (e.g. moving across a filesystme boundary on Unix). These
// may present very differently on different platforms, so be ready to update this.
func isIncompatibleRenameError(err error) bool {
	return strings.Contains(err.Error(), "invalid cross-device link")
}

// this function closely copied from this StackOverflow answer:
// https://stackoverflow.com/questions/50740902/move-a-file-to-a-different-drive-with-go
func moveFileWithCopy(sourcePath, destPath string) error {
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("move file with copy: couldn't open source file: %w", err)
	}
	defer inputFile.Close()

	outputFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("move file with copy: couldn't open dest file: %w", err)
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, inputFile)
	if err != nil {
		return fmt.Errorf("move file with copy: couldn't copy to dest from source: %w", err)
	}

	inputFile.Close() // for Windows, close before trying to remove: https://stackoverflow.com/a/64943554/246801
	err = os.Remove(sourcePath)
	if err != nil {
		return fmt.Errorf("move file with copy: couldn't remove source file: %w", err)
	}

	return nil
}

func (self FetchFromWeb) Describe() OpDescription {
	topLine := "Fetch app version from web"
	sourceUrl := fmt.Sprintf("web source URL: %s", self.SourceUrl)
	versionString := fmt.Sprintf("version label: %s", self.Version)
	destination := fmt.Sprintf("destination dir: %s", self.DestinationDir)

	return OpDescription {
		TopLine: topLine,
		ContextLines: []string{
			sourceUrl,
			versionString,
			destination,
		},
	}
}
