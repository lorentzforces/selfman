package ops

import (
	"fmt"
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
	os.Rename(tmpFile, path.Join(self.DestinationDir, path.Base(tmpFile)))

	return "Fetched app from the web", nil
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
