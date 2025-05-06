package client

import (
	"log"

	"github.com/spf13/cobra"
)

func HealthCheckCmd() *cobra.Command {

	var cmd = &cobra.Command{
		Use:   "healthcheck",
		Short: "Check server status",
		Long:  `Command will execute call to backend server to check the status of the server.`,
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("Health check..")
		},
	}

	return cmd
}
