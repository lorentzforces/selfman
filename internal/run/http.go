package run

import (
	"fmt"
	"io"
	"net/http"
	urlPkg "net/url"
	"os"
	"path"
	"time"
)

var httpClient = &http.Client{
	Timeout: 15 * time.Second,
}

// Fetch a file from the given URL using an http GET request. If no error is encountered, returns
// the path of the resulting file (which will be created in a temp directory). File name is determined from the path component of the given
// URL.
func GetFileFromUrl(url string) (string, error) {
	parsedUrl, err := urlPkg.Parse(url)
	if err != nil { return "", fmt.Errorf("Invalid URL: %s", url) }

	response, err := httpClient.Get(url)
	if err != nil { return "", fmt.Errorf("Failed to fetch from URL (%s): %w", url, err) }
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf(
			"Fetch responded with non-200 status code (%d from %s)",
			response.StatusCode, url,
		)
	}

	destPath := path.Join(os.TempDir(), path.Base(parsedUrl.Path))
	destFile, err := os.Create(destPath)
	if err != nil {
		return "", fmt.Errorf("Failed to create destination file for download: %w", err)
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, response.Body)
	if err != nil { return "", fmt.Errorf("Error while copying response buffer to file: %w", err) }

	return destPath, nil
}
