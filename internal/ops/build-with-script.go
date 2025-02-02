package ops

import "fmt"

type BuildWithScript struct {
	SourcePath string
}

func (self BuildWithScript) Execute() (string, error) {
	return "", fmt.Errorf("")
}

func (self BuildWithScript) Describe() OpDescription {
	topLine := "Build app with script"
	sourcePath := fmt.Sprintf("Source path: %s", self.SourcePath)

	return OpDescription{
		TopLine: topLine,
		ContextLines: []string{
			sourcePath,
		},
	}
}
