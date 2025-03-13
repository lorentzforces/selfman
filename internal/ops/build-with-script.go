package ops

import (
	"fmt"
	"os"

	"github.com/lorentzforces/selfman/internal/run"
)

type BuildWithScript struct {
	SourcePath string
	ScriptShell string
	ScriptCmd string
}

func (self BuildWithScript) Execute() (string, error) {
	oldWorkingDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("Determining working dir failed: %w", err)
	}

	err = os.Chdir(self.SourcePath)
	if err != nil {
		return "", fmt.Errorf("Changing to build dir failed: %w", err)
	}

	err = run.NewCmd(self.ScriptShell, run.WithArgs("-c", self.ScriptCmd)).Exec()
	if err != nil {
		return "", fmt.Errorf("Error while running build script: %w", err)
	}

	err = os.Chdir(oldWorkingDir)
	if err != nil {
		return "", fmt.Errorf("Failed to reset working dir after running build script: %w", err)
	}

	return "Executed build script", nil
}

func (self BuildWithScript) Describe() OpDescription {
	topLine := "Build app with script"
	sourcePath := fmt.Sprintf("Source path: %s", self.SourcePath)
	scriptShell := fmt.Sprintf("Shell: %s -c", self.ScriptShell)
	scriptCmd := fmt.Sprintf("Script command: %s", self.ScriptCmd)

	return OpDescription{
		TopLine: topLine,
		ContextLines: []string{
			sourcePath,
			scriptShell,
			scriptCmd,
		},
	}
}
