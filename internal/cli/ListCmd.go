package cli

import (
	"fmt"

	"github.com/lorentzforces/selfman/internal/config"
	"github.com/lorentzforces/selfman/internal/run"
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
	fmt.Printf("==DEBUG== Resolved config: %s\n", configData)
	if err != nil {
		return err
	}

	appConfigs, err := config.LoadAppConfigs(*configData.AppConfigDir)
	run.AssertNoErr(err)

	fmt.Printf("==DEBUG== # of app  configs found: %d\n", len(appConfigs))
	for _, appConfig := range appConfigs {
		fmt.Printf("==DEBUG== App Config: %#v\n", appConfig)
	}

	return nil
}
