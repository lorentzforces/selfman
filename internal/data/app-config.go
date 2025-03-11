package data

import (
	"errors"
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/lorentzforces/selfman/internal/ops"
	"github.com/lorentzforces/selfman/internal/run"
	yaml "gopkg.in/yaml.v3"
)

const (
	appTypeGit = "git"

	buildActionNone = "none"
	buildActionScript = "script"

	installActionGitClone = "git-clone"
)

// TODO: do we want to continue using the same struct for serialization and runtime usage?
type AppConfig struct {
	SystemConfig *SystemConfig `yaml:"-"`
	Name string
	Type string
	InstallAction string `yaml:"install-action"`
	BuildAction string `yaml:"build-action"`
	BuildTarget string `yaml:"build-target"`
	RemoteRepo *string `yaml:"remote-repo,omitempty"`
}

func (self *AppConfig) SourcePath() string {
	return path.Join(self.SystemConfig.SourcesPath(), self.Name)
}

func (self *AppConfig) ArtifactPath() string {
	return path.Join(self.SystemConfig.ArtifactsPath(), self.Name)
}

func (self *AppConfig) BuildTargetPath() string {
	return path.Join(self.SourcePath(), self.BuildTarget)
}

func (self *AppConfig) BinaryPath() string {
	return path.Join(*self.SystemConfig.BinaryDir, self.Name)
}

func (self *AppConfig) GetInstallOp() ops.Operation {
	switch self.InstallAction{
	case installActionGitClone: {
		return ops.GitClone{
			RepoUrl: *self.RemoteRepo,
			DestinationPath: self.BuildTarget,
		}
	}
	}

	run.FailOut("Unhandled install action -> operation mapping")
	panic("Unreachable in theory")
}

func (self *AppConfig) GetBuildOp() ops.Operation {
	switch self.BuildAction {
	case buildActionNone: {
		return ops.NoBuildOp
	}
	}

	run.FailOut("Unhandled build action -> operation mapping")
	panic("Unreachable in theory")
}

func (self *AppConfig) applyDefaults() {
	if len(self.BuildTarget) == 0 {
		self.BuildTarget = strings.ToLower(self.Name)
	}

	// TODO: create full structs of default configs for, say a git-based app, and then coalesce that
	// as a full default config (then apply name-dependent defaults)
	if self.Type == appTypeGit && self.InstallAction == "" {
		self.InstallAction = installActionGitClone
	}
}

// Validates an application config - error will be non-nil if validation failed.
func (self *AppConfig) validate() error {
	if len(self.Name) == 0 {
		return fmt.Errorf("Application name cannot be empty")
	}

	if !self.isValidAppType() {
		return fmt.Errorf("(app %s) Invalid application type: %s", self.Name, self.Type)
	}

	if !self.isValidBuildAction() {
		return fmt.Errorf("(app %s) Invalid build action: %s", self.Name, self.BuildAction)
	}

	if !self.isValidInstallAction() {
		return fmt.Errorf("(app %s) Invalid install action: %s", self.Name, self.InstallAction)
	}

	if self.Type == "git" && self.RemoteRepo == nil {
		return fmt.Errorf("(app %s) Remote repo must be specified for apps of type git", self.Name)
	}

	return nil
}

func (self *AppConfig) isValidAppType() bool {
	switch self.Type {
	case appTypeGit: return true
	default: return false
	}
}

func (self *AppConfig) isValidBuildAction() bool {
	switch self.BuildAction {
	case buildActionNone, buildActionScript: return true
	default: return false
	}
}

func (self *AppConfig) isValidInstallAction() bool {
	switch self.InstallAction {
	case installActionGitClone: return true
	default: return false
	}
}

func loadAppConfigs(systemConfig *SystemConfig) ([]AppConfig, error) {
	appConfigPath := *systemConfig.AppConfigDir
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
	entries, err := os.ReadDir(*systemConfig.AppConfigDir)
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
		appConfig.SystemConfig = systemConfig
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
