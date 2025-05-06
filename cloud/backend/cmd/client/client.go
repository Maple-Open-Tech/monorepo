package client

import (
	"github.com/spf13/cobra"
)

func ClientCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "client",
		Short: "Execute commands related to client making API calls",
		Run: func(cmd *cobra.Command, args []string) {
			// Do nothing...
		},
	}

	// Add rcca-related commands
	cmd.AddCommand(HealthCheckCmd())
	cmd.AddCommand(RegisterUserCmd())
	cmd.AddCommand(VerifyEmailCmd())
	cmd.AddCommand(LoginUserCmd())
	cmd.AddCommand(LogoutUserCmd())
	cmd.AddCommand(EchoCmd())
	return cmd
}
