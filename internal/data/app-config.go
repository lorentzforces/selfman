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
// TODO(lombok): support some kind of "library" or "dir" target as well as the binary target, which
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
	MiscVars map[string]string `yaml:"misc-vars"`
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

	if self.MiscVars == nil {
		self.MiscVars = make(map[string]string, 1)
	}
	self.MiscVars["VERSION"] = self.Version
}

// Will apply misc vars to replace appropriate placeholders in these fields:
//   - BuildAction
//   - BuildTarget
//   - BuildCmd
//   - WebUrl
func (self *AppConfig) applyMiscVarsToPlaceholders() error {
	var err error
	self.BuildAction, err = replacePlaceholders(self.BuildAction, self.MiscVars)
	if err != nil {
		return errors.Join(fmt.Errorf("Error filling placeholders in BuildAction"), err)
	}
	self.BuildTarget, err = replacePlaceholders(self.BuildTarget, self.MiscVars)
	if err != nil {
		return errors.Join(fmt.Errorf("Error filling placeholders in BuildTarget"), err)
	}
	if self.BuildCmd != nil {
		*self.BuildCmd, err = replacePlaceholders(*self.BuildCmd, self.MiscVars)
		if err != nil {
			return errors.Join(fmt.Errorf("Error filling placeholders in BuildCmd"), err)
		}
	}
	if self.WebUrl != nil {
		*self.WebUrl, err = replacePlaceholders(*self.WebUrl, self.MiscVars)
		if err != nil {
			return errors.Join(fmt.Errorf("Error filling placeholders in WebUrl"), err)
		}
	}

	return nil
}

var placeholderPattern = regexp.MustCompile(`%(?<label>[-A-Za-z_]{3,})%`)
// CAPTURE GROUPS (submatches) label: 1
func replacePlaceholders(original string, keyVals map[string]string) (string, error) {
	matches := placeholderPattern.FindAllStringSubmatch(original, 32)
	if matches == nil { return original, nil }

	finalString := original
	badLabels := make([]string, 0)
	for _, match := range matches {
		label := match[1]
		labelVal, labelPresent := keyVals[label]
		if labelPresent {
			finalString = strings.ReplaceAll(finalString, "%" + label + "%", labelVal)
		} else {
			badLabels = append(badLabels, label)
		}
	}

	if len(badLabels) > 0 {
		return finalString, fmt.Errorf(
			"Placeholder labels referenced but no value found: %s",
			strings.Join(badLabels, ", "),
		)
	}

	return finalString, nil
}

var validLabelPattern = regexp.MustCompile(`\A[-A-Za-z_]*\z`)
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

	for label, _ := range self.MiscVars {
		if nil == validLabelPattern.FindStringIndex(label) {
			return fmt.Errorf(
				"(app %s) Label \"%s\" must be only upper & lower case letters, hyphens, and " +
					"underscores",
				self.Name, label,
			)
		}

		// all our valid characters are 1-byte in utf-8, so this is reasonable
		if len(label) < 3 {
			return fmt.Errorf(
				"(app %s) Label \"%s\" is less than the required three characters",
				self.Name, label,
			)
		}
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
