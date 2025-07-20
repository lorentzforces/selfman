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
)

const (
	FlavorGit = "git"
	FlavorWebFetch = "web-fetch"
	// TODO: binaryFile

	ActionNone = "none"

	BuildActionScript = "script"

	UpdateActionGitFetch = "git-fetch"
)

// TODO: do we want to continue using the same struct for serialization and runtime usage?
// TODO: make some form of config file with explainers (potentially accessible via command)
// TODO(jdtls): support some kind of "library" or "dir" target as well as the binary target, which
//              puts the full directory tree somewhere well-known
type AppConfig struct {
	SystemConfig *SystemConfig `yaml:"-"` // ignored in yaml
	Name string
	Flavor string
	Version string
	BuildAction string `yaml:"build-action"`
	BuildTarget string `yaml:"build-target"`
	RemoteRepo *string `yaml:"remote-repo,omitempty"`
	BuildCmd *string `yaml:"build-cmd,omitempty"`
	WebUrl *string `yaml:"web-url,omitempty"`
	KeepBinWithSource bool `yaml:"keep-bin-with-source"`
}

func (self *AppConfig) SourcePath() string {
	if self.Flavor == FlavorGit {
		return path.Join(self.SystemConfig.SourcesPath(), self.Name, "git")
	}
	return path.Join(self.SystemConfig.SourcesPath(), self.Name, self.Version)
}

// Will replace the path separator if it is found in the version (e.g. "origin/main")
func (self *AppConfig) ArtifactPath() string {
	if self.KeepBinWithSource {
		return self.BuildTargetPath()
	}

	rawFileName := self.Name + "---" + self.Version
	escapedFileName := strings.ReplaceAll(rawFileName, string(os.PathSeparator), "%SLASH%")
	return path.Join(self.SystemConfig.ArtifactsPath(), escapedFileName)
}

func (self *AppConfig) BuildTargetPath() string {
	return path.Join(self.SourcePath(), self.BuildTarget)
}

func (self *AppConfig) BinaryPath() string {
	return path.Join(*self.SystemConfig.BinaryDir, self.Name)
}

func (self *AppConfig) GetObtainSourceOp() ops.Operation {
	switch self.Flavor{
	case FlavorGit: {
		return ops.GitClone{
			RepoUrl: *self.RemoteRepo,
			DestinationPath: self.SourcePath(),
		}
	}
	case FlavorWebFetch: {
		return ops.FetchFromWeb{
			SourceUrl: *self.WebUrl,
			Version: self.Version,
			DestinationDir: self.SourcePath(),
		}
	}
	}

	run.FailOut(fmt.Sprintf(
		"Unhandled app flavor -> operation mapping: %s",
		self.Flavor,
	))
	panic("Unreachable in theory")
}

func (self *AppConfig) GetBuildOp() ops.Operation {
	switch self.BuildAction {
	case ActionNone: {
		return ops.NoBuildOp
	}
	case BuildActionScript: {
		return ops.BuildWithScript{
			SourcePath: self.SourcePath(),
			ScriptShell: *self.SystemConfig.ScriptShell,
			ScriptCmd: *self.BuildCmd,
		}
	}
	}

	run.FailOut(fmt.Sprintf("Unhandled build action -> operation mapping: %s", self.BuildAction))
	panic("Unreachable in theory")
}

// Returns the operation to select the appropriate version for the application, if the application
// flavor requires such an operation. If not, returns nil.
func (self *AppConfig) GetSelectVersionOp() ops.Operation {
	switch self.Flavor {
	case FlavorGit: {
		return ops.GitCheckoutRef{
			RepoPath: self.SourcePath(),
			RefName: self.Version,
		}
	}
	default: { return nil }
	}
}

func (self *AppConfig) GetFetchUpdatesOp() ops.Operation {
	switch self.Flavor {
	case FlavorGit: {
		return ops.GitFetch{
			RepoPath: self.SourcePath(),
		}
	}
	default: {
		return nil
	}
	}
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

	if !self.isValidAppFlavor() {
		return fmt.Errorf("(app %s) Invalid application flavor: %s", self.Name, self.Flavor)
	}

	if !self.isValidBuildAction() {
		return fmt.Errorf("(app %s) Invalid build action: %s", self.Name, self.BuildAction)
	}

	if self.Flavor == FlavorGit && self.RemoteRepo == nil {
		return fmt.Errorf(
			"(app %s) Remote repo must be specified for apps of flavor %s", self.Name, FlavorGit)
	}

	if self.Flavor == FlavorWebFetch && self.WebUrl == nil {
		return fmt.Errorf(
			"(app %s) Web URL must be specified for apps of flavor %s", self.Name, FlavorWebFetch)
	}

	return nil
}

func (self *AppConfig) isValidAppFlavor() bool {
	switch self.Flavor {
	case FlavorGit: return true
	case FlavorWebFetch: return true
	default: return false
	}
}

func (self *AppConfig) isValidBuildAction() bool {
	switch self.BuildAction {
	case ActionNone, BuildActionScript: return true
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
		if entry.Type().IsDir() {
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
		if err != nil { return nil, err }

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
	configFile, err := os.Open(appConfigPath)
	run.AssertNoErr(err)

	appConfig := AppConfig{}
	err = run.GetStrictDecoder(configFile).Decode(&appConfig)
	if err != nil {
		newErr := fmt.Errorf(
			"Error parsing application config file \"%s\"",
			appConfigPath,
		)

		return AppConfig{}, errors.Join(newErr, err)
	}

	return appConfig, nil
}
