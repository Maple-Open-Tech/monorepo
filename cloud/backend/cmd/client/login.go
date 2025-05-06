package client

import (
	"log"

	"github.com/spf13/cobra"
)

func LoginUserCmd() *cobra.Command {
	var email, password string

	var cmd = &cobra.Command{
		Use:   "register",
		Short: "Log user into account",
		Long:  `Command will execute login command and user will get credentials to make API calls to their account`,
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("Logging in...")

			if email == "" || password == "" {
				log.Fatal("Email and password required")
			}

			if err := loginUser(email, password); err != nil {
				log.Fatalf("Failed to login user: %v", err)
			}
		},
	}

	// Define command flags
	cmd.Flags().StringVarP(&email, "email", "e", "", "Email address for the user (required)")
	cmd.Flags().StringVarP(&password, "password", "p", "", "Password for the user (required)")

	// Mark required flags
	cmd.MarkFlagRequired("email")
	cmd.MarkFlagRequired("password")

	return cmd
}

func loginUser(email, plainPassword string) error {

	return nil
}
