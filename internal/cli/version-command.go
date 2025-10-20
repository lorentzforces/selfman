package cli

import (
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/lorentzforces/selfman/internal/data"
	"github.com/spf13/cobra"
)

func CreateVersionCmd() SelfmanCommand {
	return SelfmanCommand{
		cobraCmd: &cobra.Command{
			Use: "version",
			Short: "Print version & build information",
		},
		runFunc: runVersionCmd,
	}
}

func runVersionCmd(cmd *cobra.Command, args []string) (*SelfmanResult, error) {
	buildData, err := fetchBuildInfo()
	if err != nil {
		return nil, err
	}

	return &SelfmanResult {
		textOutput: &buildData,
		operations: nil,
	}, nil
}

func fetchBuildInfo() (selfmanBuildInfo, error) {
	selfBuildData := selfmanBuildInfo{}
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		return selfBuildData, fmt.Errorf("Could not read self build info")
	}

	selfBuildData.goVersion = buildInfo.GoVersion

	for _, setting := range buildInfo.Settings {
		switch setting.Key {
		case "vcs.revision":
			selfBuildData.gitRev = setting.Value
		case "vcs.modified":
			selfBuildData.vcsHadModifications = (setting.Value == "true")
		case "GOARCH":
			selfBuildData.buildArch = setting.Value
		case "GOOS":
			selfBuildData.buildOsTarget = setting.Value
		}
	}

	return selfBuildData, nil
}

type selfmanBuildInfo struct {
	goVersion string
	buildArch string
	buildOsTarget string
	gitRev string
	vcsHadModifications bool
}

func (self selfmanBuildInfo) String() string {
	var buf strings.Builder

	buf.WriteString(fmt.Sprintf("selfman %s\n", data.Globals.ReleaseLabel))
	buf.WriteString(fmt.Sprintf("source: %s", self.gitRev))
	if self.vcsHadModifications {
		buf.WriteString(" (in progress)")
	}
	buf.WriteString("\n")
	buf.WriteString(fmt.Sprintf(
		"build: %s-%s w/ %s\n",
		self.buildOsTarget, self.buildArch, self.goVersion,
	))

	return buf.String()
}
