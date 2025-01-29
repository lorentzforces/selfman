package config

import (
	"errors"
	"fmt"
	"os"
	"path"
	"regexp"

	"github.com/lorentzforces/selfman/internal/run"
	yaml "gopkg.in/yaml.v3"
)

const ConfigurationEnvVar = "SELFMAN_CONFIG"

func Produce() (Config, error) {
	selfmanConfig, err := loadSelfmanConfig()
	if err != nil {
		return Config{}, err
	}

	appConfigs, err := loadAppConfigs(*selfmanConfig.AppConfigDir)
	if err != nil {
		return Config{}, err
	}

	for _, app := range appConfigs {
		err := app.validate()
		if err != nil {
			newErr := fmt.Errorf(
				"Invalid app config in app directory \"%s\"",
				*selfmanConfig.AppConfigDir,
			)
			return Config{}, errors.Join(newErr, err)
		}
	}

	return Config{
		SelfmanConfig: selfmanConfig,
		AppConfigs: appConfigs,
	}, nil
}

type Config struct {
	SelfmanConfig Selfman
	AppConfigs []AppConfig
}

type Selfman struct {
	AppConfigDir *string `yaml:"app-config-dir,omitempty"`
}

func (self *Selfman) expandPaths() {
	*self.AppConfigDir = os.ExpandEnv(*self.AppConfigDir)
}

type AppConfig struct {
	Name string
}

// Validates an application config - error will be non-nil if validation failed.
func (self AppConfig) validate() error {
	if len(self.Name) == 0 {
		return fmt.Errorf("Application name cannot be empty")
	}

	return nil
}

func defaultConfig() Selfman {
	defaultAppConfigPath := path.Join(run.ResolveXdgConfigDir(), "selfman", "apps")
	return Selfman{
		AppConfigDir: &defaultAppConfigPath,
	}
}

func loadSelfmanConfig() (Selfman, error) {
	path, err := resolveConfigPath(ConfigurationEnvVar)
	if err != nil {
		return Selfman{}, fmt.Errorf("Could not resolve config file: %w", err)
	}

	defaultConfig := defaultConfig()

	if len(path) == 0 {
		return defaultConfig, nil
	}

	configBytes, err := os.ReadFile(path)
	run.AssertNoErrReason(err, "Config file was resolved but later reading failed")
	configData := Selfman{}
	err = yaml.Unmarshal(configBytes, &configData)
	if err != nil {
		return Selfman{}, fmt.Errorf("Error parsing config file: %w", err)
	}

	finalConfig := coalesceConfigs(defaultConfig, configData)

	finalConfig.expandPaths()
	return finalConfig, nil
}

// Returns the first-resolved configuration location. Resolves in this order of priority:
//   - $SELFMAN_CONFIG
//   - $XDG_CONFIG_HOME/selfman/selfman.config
//   - ~/.config/selfman/selfman.config
// If $SELFMAN_CONFIG is set but no readable file exists at that path, this function will return
// an error.
//
// If no readable file is resolved by the above process, an empty string is returned with no error.
func resolveConfigPath(configEnvName string) (string, error) {
	configEnvPath := os.Getenv(configEnvName)
	if len(configEnvPath) > 0 {
		_, err := checkFileAtPath(configEnvPath)
		if err != nil {
			return "", fmt.Errorf(
				"Configuration path was specified in env var %s but was not readable: %w",
				configEnvName, err,
			)
		}
		return configEnvPath, nil
	}

	configXdgPath := path.Join(run.ResolveXdgConfigDir(), "selfman", "config.yaml")
	found, err := checkFileAtPath(configXdgPath)

	switch {
	case found && err != nil:
		return "", fmt.Errorf(
			"Configuration found at path \"%s\" with error: %w",
			configXdgPath, err,
		)
	case found: return configXdgPath, nil
	default: return "", nil
	}
}

func checkFileAtPath(path string) (foundFile bool, err error) {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return false, err
	}

	// in theory if the above succeeded, this cannot fail
	fileStat, err := os.Stat(path)
	run.AssertNoErrReason(err, "file stat failed after file successfully opened")
	if fileStat.IsDir() {
		return true, fmt.Errorf("Resolved config file is a directory: \"%s\"", path)
	}

	return true, nil
}

func coalesceConfigs(a, b Selfman) Selfman {
	result := Selfman{}
	result.AppConfigDir = run.Coalesce(b.AppConfigDir, a.AppConfigDir)
	return result
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
