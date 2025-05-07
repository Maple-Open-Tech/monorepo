// native/desktop/papercloud-cli/cmd/remote/remote.go
package remote

import (
	"github.com/spf13/cobra"

	pref "github.com/Maple-Open-Tech/monorepo/native/desktop/papercloud-cli/internal/common/preferences"
)

var (
	preferences *pref.Preferences
)

func RemoteCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "remote",
		Short: "Execute commands related to making remote API calls",
		Run: func(cmd *cobra.Command, args []string) {
			// Do nothing...
		},
	}

	// Add Remote-related commands
	cmd.AddCommand(HealthCheckCmd())
	cmd.AddCommand(EchoCmd())
	cmd.AddCommand(RegisterUserCmd())
	cmd.AddCommand(VerifyEmailCmd())
	cmd.AddCommand(LoginUserCmd())
	cmd.AddCommand(LogoutUserCmd())

	return cmd
}
