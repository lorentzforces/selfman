package cli

import (
	"fmt"

	"github.com/lorentzforces/selfman/internal/config"
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
	configData, err := config.Produce()
	if err != nil {
		return err
	}

	results := listApplications(configData)
	for _, line := range results {
		fmt.Println(line)
	}
	return nil
}

func listApplications(config config.Config) []string {
	results := make([]string, 0, len(config.AppConfigs))
	for _, app := range config.AppConfigs {
		results = append(results, app.Name)
	}
	return results
}
