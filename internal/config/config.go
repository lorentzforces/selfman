package config

import (
	"fmt"
	"os"
	"os/user"
	"path"

	"github.com/lorentzforces/selfman/internal/run"
	yaml "gopkg.in/yaml.v3"
)

const ConfigurationEnvVar = "SELFMAN_CONFIG"

type Config struct {
	Test string
	OptionalTest string `yaml:"optional-test,omitempty"`
}

func defaultConfig() Config {
	return Config{
		Test: "default test value",
		OptionalTest: "default optional test value",
	}
}

func Produce() (Config, error) {
	path, err := resolveConfigPath(ConfigurationEnvVar)
	if err != nil {
		return Config{}, fmt.Errorf("Could not resolve config file: %w", err)
	}
	fmt.Printf("==DEBUG== Config file at: %s\n", path)

	configData := defaultConfig()

	if len(path) == 0 {
		return configData, nil
	}

	configBytes, err := os.ReadFile(path)
	run.AssertNoErrReason(err, "Config file was resolved but later reading failed")
	err = yaml.Unmarshal(configBytes, &configData)
	if err != nil {
		return Config{}, fmt.Errorf("Error parsing config file: %w", err)
	}

	return configData, nil
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

	configDir := os.Getenv("XDG_CONFIG_HOME")
	if len(configDir) == 0 {
		usr, err := user.Current()
		run.AssertNoErr(err)
		configDir = usr.HomeDir
	}
	configXdgPath := path.Join(configDir, "selfman", "config.yaml")
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
