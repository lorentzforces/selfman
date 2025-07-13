package cli

import (
	"fmt"
	"slices"
	"strings"

	"github.com/lorentzforces/selfman/internal/data"
	"github.com/spf13/cobra"
)

func CreateListCmd() SelfmanCommand {
	return SelfmanCommand{
		cobraCmd: &cobra.Command{
			Use: "list",
			Short: "List all applications managed by selfman",
			Aliases: []string{ "ls" },
		},
		runFunc: runListCmd,
	}
}

func runListCmd(cmd *cobra.Command, args []string) (*SelfmanResult, error) {
	if err := validatePrereqs(); err != nil { return nil, err }
	configData, err := data.Produce()
	if err != nil { return nil, err }

	results := listApplications(configData)
	return &SelfmanResult{
		textOutput: listCmdResult{ results },
		operations: nil,
	}, nil
}

type listCmdResult struct {
	results []listResult
}

// TODO: use tabwriter or something similar to format this better
func (self listCmdResult) String() string {
	var buf strings.Builder
	for _, result := range self.results {
		buf.WriteString(fmt.Sprintf("%s @ %s (%s)\n", result.name, result.version, result.status))
	}

	return buf.String()
}

type listResult struct {
	name string
	version string
	status string
}

func listApplications(selfmanData data.Selfman) []listResult {
	results := make([]listResult, 0, len(selfmanData.AppConfigs))
	for _, app := range selfmanData.AppConfigs {
		_, status := selfmanData.AppStatus(app.Name)
		results = append(
			results,
			listResult{ name: app.Name, version: app.Version, status: status.Label() },
		)
	}

	slices.SortFunc(results, func(a, b listResult) int {
		switch {
		case a.name < b.name: return -1
		case b.name > a.name: return 1
		default: return 0
		}
	})

	return results
}
