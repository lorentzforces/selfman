package data

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/lorentzforces/selfman/internal/run"
	yaml "gopkg.in/yaml.v3"
)

const ConfigurationEnvVar = "SELFMAN_CONFIG"

func DefaultTestConfig() *SystemConfig {
	return &SystemConfig{
		AppConfigDir: run.StrPtr("/tmp/selfman-test/apps"),
		DataDir: run.StrPtr("/tmp/selfman-test/data"),
		BinaryDir: run.StrPtr("/tmp/selfman-test/bin"),
	}
}

func loadSystemConfig() (SystemConfig, error) {
	path, err := resolveConfigPath(ConfigurationEnvVar)
	if err != nil {
		return SystemConfig{}, fmt.Errorf("Could not resolve config file: %w", err)
	}

	defaultConfig := defaultConfig()

	if len(path) == 0 {
		return defaultConfig, nil
	}

	configBytes, err := os.ReadFile(path)
	run.AssertNoErrReason(err, "Config file was resolved but later reading failed")
	configData := SystemConfig{}
	err = yaml.Unmarshal(configBytes, &configData)
	if err != nil {
		return SystemConfig{}, fmt.Errorf("Error parsing config file: %w", err)
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

	configXdgPath := path.Join(resolveXdgConfigDir(), "selfman", "config.yaml")
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
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return true, err
	}

	// in theory if the above succeeded, this cannot fail
	fileStat, err := os.Stat(path)
	run.AssertNoErrReason(err, "file stat failed after file successfully opened")
	if fileStat.IsDir() {
		return true, fmt.Errorf("Resolved config file is a directory: \"%s\"", path)
	}

	return true, nil
}

func coalesceConfigs(a, b SystemConfig) SystemConfig {
	result := SystemConfig{}
	result.AppConfigDir = run.Coalesce(b.AppConfigDir, a.AppConfigDir)
	result.DataDir = run.Coalesce(b.DataDir, a.DataDir)
	result.BinaryDir = run.Coalesce(b.BinaryDir, a.BinaryDir)
	return result
}

type SystemConfig struct {
	AppConfigDir *string `yaml:"app-config-dir,omitempty"`
	DataDir *string `yaml:"data-dir,omitempty"`
	BinaryDir *string `yaml:"binary-dir,omitempty"`
}

func (self *SystemConfig) expandPaths() {
	*self.AppConfigDir = os.ExpandEnv(*self.AppConfigDir)
	*self.DataDir = os.ExpandEnv(*self.DataDir)
	*self.BinaryDir = os.ExpandEnv(*self.BinaryDir)
}

func (self *SystemConfig) SourcesPath() string {
	return path.Join(*self.DataDir, "sources")
}

func (self *SystemConfig) ArtifactsPath() string {
	return path.Join(*self.DataDir, "artifacts")
}

func defaultConfig() SystemConfig {
	return SystemConfig{
		AppConfigDir: run.StrPtr(path.Join(resolveXdgConfigDir(), "selfman", "apps")),
		DataDir: run.StrPtr(path.Join(resolveXdgDataDir(), "selfman")),
		BinaryDir: run.StrPtr(resolveXdgBinDir()),
	}
}

