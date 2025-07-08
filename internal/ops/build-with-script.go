package ops

import (
	"fmt"

	"github.com/lorentzforces/selfman/internal/run"
)

type BuildWithScript struct {
	SourcePath string
	ScriptShell string
	ScriptCmd string
}

func (self BuildWithScript) Execute() (string, error) {
	err := run.NewCmd(self.ScriptShell, run.WithArgs("-c", self.ScriptCmd), run.WithWorkingDir(self.SourcePath)).Exec()
	if err != nil {
		return "", fmt.Errorf("Error while running build script: %w", err)
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
