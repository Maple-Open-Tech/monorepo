// native/desktop/papercloud-cli/cmd/remote/remote.go
package remote

import (
	"github.com/spf13/cobra"

	pref "github.com/Maple-Open-Tech/monorepo/native/desktop/papercloud-cli/internal/common/preferences"
)

var (
	preferences *pref.Preferences
)

// Initialize function will be called when every command gets called.
func init() {
	preferences = pref.PreferencesInstance()
}

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
	cmd.AddCommand(RequestLoginOneTimeTokenUserCmd())
	cmd.AddCommand(VerifyLoginOneTimeTokenUserCmd())
	cmd.AddCommand(CompleteLoginCmd())
	cmd.AddCommand(MeCmd())
	cmd.AddCommand(UploadFileCmd())
	// cmd.AddCommand(LogoutUserCmd())

	return cmd
}
