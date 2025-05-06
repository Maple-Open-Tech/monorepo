package client

import (
	"log"

	"github.com/spf13/cobra"
)

func VerifyEmailCmd() *cobra.Command {

	var cmd = &cobra.Command{
		Use:   "verifyemail",
		Short: "Send email activation code",
		Long:  `Command will execute submitting your email activation code to the backend to finalize registration`,
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("Echo..")
		},
	}

	return cmd
}
