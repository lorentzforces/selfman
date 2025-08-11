package cli

import (
	"fmt"
	"strings"

	"github.com/lorentzforces/selfman/internal/data"
	"github.com/spf13/cobra"
)

func CreateCheckCmd() SelfmanCommand {
	return SelfmanCommand{
		cobraCmd: &cobra.Command{
			Use: "check [flags] app-name",
			Short: "Get detailed information about an application",
		},
		runFunc: runCheckCmd,
	}
}

func runCheckCmd(cmd *cobra.Command, args []string) (*SelfmanResult, error) {
	if err := validatePrereqs(); err != nil { return nil, err }
	selfmanData, err := data.Produce()
	if err != nil { return nil, err }

	if len(args) < 1 {
		return nil,
			fmt.Errorf("Check command expects an application name, but one was not provided")
	}
	result, err := checkApp(args[0], selfmanData)
	if err != nil { return nil, err }

	return &SelfmanResult{
		textOutput: result,
		operations: nil,
	}, nil
}

type checkAppResult struct {
	appName string
	appIsLib bool
	status data.AppStatus
}

func (self checkAppResult) String() string {
	versionsString := "None!"
	if len(self.status.AvailableVersions) > 0 {
		versionsString = strings.Join(self.status.AvailableVersions, ", ")
	}
	resultString := fmt.Sprintf(
		"📋 %s\n" +
		"  version: %s\n\n" +
		"Overall status: %s\n" +
		"  Source present: %t\n" +
		"  Target present: %t\n" +
		"  Bin link present: %t\n",
		self.appName, self.status.DesiredVersion, self.status.Label(),
		self.status.SourcePresent, self.status.TargetPresent, self.status.LinkPresent,
	)

	if self.appIsLib {
		resultString += fmt.Sprintf("  Lib link present: %t\n", self.status.LibLinkPresent)
	}

	resultString += fmt.Sprintf("Available versions (locally): %s\n", versionsString)

	return resultString
}

func checkApp(name string, selfmanData data.Selfman) (checkAppResult, error) {
	app, status := selfmanData.AppStatus(name)
	if !status.IsConfigured {
		return checkAppResult{}, fmt.Errorf(
			"No configuration for an app named %s was found",
			name,
		)
	}

	return checkAppResult{
		appName: name,
		appIsLib: app.LinkSourceAsLib,
		status: status,
	}, nil
}
