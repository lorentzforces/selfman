package cli

import (
	"fmt"

	"github.com/lorentzforces/selfman/internal/data"
	"github.com/spf13/cobra"
)

func CreateCheckCmd() *cobra.Command {
	return &cobra.Command{
		Use: "check [flags] app-name",
		Short: "Get detailed information about an application",
		RunE: runCheckCmd,
	}
}

func runCheckCmd(cmd *cobra.Command, args []string) error {
	if err := validatePrereqs(); err != nil { return err }
	selfmanData, err := data.Produce()
	if err != nil { return err }

	if len(args) < 1 {
		return fmt.Errorf("Check command expects an application name, but one was not provided")
	}
	result, err := checkApp(args[0], selfmanData)

	if err == nil { fmt.Println(result) }

	return err
}

type checkAppResult struct {
	appName string
	status data.AppStatus
}

func (self checkAppResult) String() string {
	return fmt.Sprintf(
		"📋 %s\n\n" +
		"Overall status: %s\n" +
		"  Source present: %t\n" +
		"  Target present: %t\n" +
		"  Bin link present: %t\n",
		self.appName, self.status.Label(), self.status.SourcePresent,
		self.status.TargetPresent, self.status.LinkPresent,
	)
}

func checkApp(name string, selfmanData data.Selfman) (checkAppResult, error) {
	_, status := selfmanData.AppStatus(name)
	if !status.IsConfigured {
		return checkAppResult{}, fmt.Errorf(
			"No configuration for an app named %s was found",
			name,
		)
	}

	return checkAppResult{
		appName: name,
		status: status,
	}, nil
}
