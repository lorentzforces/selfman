package ops

import "fmt"

type BuildWithScript struct {
	SourcePath string
	ScriptCmd string
}

func (self BuildWithScript) Execute() (string, error) {
	return "", fmt.Errorf("")
}

func (self BuildWithScript) Describe() OpDescription {
	topLine := "Build app with script"
	sourcePath := fmt.Sprintf("Source path: %s", self.SourcePath)
	scriptCmd := fmt.Sprintf("Script: %s", self.ScriptCmd)

	return OpDescription{
		TopLine: topLine,
		ContextLines: []string{
			sourcePath,
			scriptCmd,
		},
	}
}
