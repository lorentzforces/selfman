package data

import (
	"errors"
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/lorentzforces/selfman/internal/run"
	yaml "gopkg.in/yaml.v3"
)

type AppConfig struct {
	Name string
	Type string
	BuildTarget string
	RemoteRepo *string `yaml:"remote-repo,omitempty"`
}

type AppType string
var appTypes = map[string]AppType {
	"git": "git",
}

func GetAppType(label string) (AppType, error) {
	appType, present := appTypes[label]
	if !present {
		return "", fmt.Errorf("Invalid AppType label passed: %s", label)
	}
	return appType, nil
}

func MustBeAppType(label string) AppType {
	appType, present := appTypes[label]
	run.Assert(present, fmt.Sprintf("Invalid AppType label passed: %s", label))
	return appType
}

func (self *AppConfig) applyDefaults() {
	if len(self.BuildTarget) == 0 {
		self.BuildTarget = strings.ToLower(self.Name)
	}
}

// Validates an application config - error will be non-nil if validation failed.
func (self *AppConfig) validate() error {
	if len(self.Name) == 0 {
		return fmt.Errorf("Application name cannot be empty")
	}

	_, err := GetAppType(self.Type)
	if err != nil {
		return fmt.Errorf("(app %s) Invalid application type: %s", self.Name, self.Type)
	}

	if self.Type == "git" && self.RemoteRepo == nil {
		return fmt.Errorf("(app %s) Remote repo must be specified for apps of type git", self.Name)
	}

	return nil
}

func loadAppConfigs(appConfigPath string) ([]AppConfig, error) {
	stat, err := os.Stat(appConfigPath)
	if err != nil {
		// if the directory just doesn't exist, we say "okay" and return an empty list
		if errors.Is(err, os.ErrNotExist) {
			return make([]AppConfig, 0), nil
		}
		return nil, fmt.Errorf(
			"Could not load configured application config path at \"%s\": %w",
			appConfigPath, err,
		)
	}
	if !stat.IsDir() {
		return nil, fmt.Errorf(
			"Configured application config path was not a directory: \"%s\"",
			appConfigPath,
		)
	}

	appConfigPaths := make([]string, 0)
	entries, err := os.ReadDir(appConfigPath)
	for _, entry := range entries {
		if !entry.Type().IsRegular() {
			continue
		}
		if isAppConfigFileName(entry.Name()) {
			fullPath := path.Join(appConfigPath, entry.Name())
			appConfigPaths = append(appConfigPaths, fullPath)
		}
	}

	appConfigs := make([]AppConfig, 0, len(appConfigPaths))
	for _, path := range appConfigPaths {
		appConfig, err := parseAppConfig(path)
		if err != nil {
			return nil, err
		}
		appConfigs = append(appConfigs, appConfig)
	}

	return appConfigs, nil
}

var appConfigRegex = regexp.MustCompile(`.+\.config\.yaml\z`)

func isAppConfigFileName(fileName string) bool {
	return appConfigRegex.MatchString(fileName)
}

func parseAppConfig(appConfigPath string) (AppConfig, error) {
	configBytes, err := os.ReadFile(appConfigPath)
	run.AssertNoErr(err)

	appConfig := AppConfig{}
	err = yaml.Unmarshal(configBytes, &appConfig)
	if err != nil {
		newErr := fmt.Errorf(
			"Error parsing application config file \"%s\"",
			appConfigPath,
		)

		return AppConfig{}, errors.Join(newErr, err)
	}

	return appConfig, nil
}
