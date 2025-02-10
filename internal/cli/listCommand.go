package cli

import (
	"fmt"
	"slices"

	"github.com/lorentzforces/selfman/internal/data"
	"github.com/spf13/cobra"
)

func CreateListCmd() *cobra.Command {
	return &cobra.Command{
		Use: "list",
		Short: "List all applications managed by selfman (currently debug placeholder)",
		RunE: runListCmd,
	}
}

func runListCmd(cmd *cobra.Command, args []string) error {
	configData, err := data.Produce()
	if err != nil {
		return err
	}

	// TODO: use tabwriter or something similar to format this better
	results := listApplications(configData)
	for _, result := range results {
		fmt.Println(result)
	}
	return nil
}

type listResult struct {
	name string
	status data.AppStatus
}

func (self listResult) String() string {
	return fmt.Sprintf("%s (%s)", self.name, self.status)
}

func listApplications(selfmanData data.Selfman) []listResult {
	results := make([]listResult, 0, len(selfmanData.AppConfigs))
	for _, app := range selfmanData.AppConfigs {
		status := selfmanData.AppStatus(app.Name)
		results = append(results, listResult{name: app.Name, status: status})
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
