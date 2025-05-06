package client

import (
	"log"

	"github.com/spf13/cobra"
)

func EchoCmd() *cobra.Command {

	var cmd = &cobra.Command{
		Use:   "echo",
		Short: "Echo text to backend",
		Long:  `Command will execute submitting any text to the server and the server will respond back with the text you submitted`,
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("Echo..")
		},
	}

	return cmd
}
