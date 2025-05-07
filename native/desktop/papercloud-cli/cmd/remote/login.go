// native/desktop/papercloud-cli/cmd/remote/login.go
package remote

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

func LoginUserCmd() *cobra.Command {
	var email, password string

	var cmd = &cobra.Command{
		Use:   "login",
		Short: "Log user into account",
		Long: `
Command will execute login command and user will get credentials to make API calls to their account.

After registration and email verification, use this command to log in to your account.
You'll receive authentication tokens that will be stored securely for making API calls.

Examples:
  # Login with email and password
  login --email user@example.com --password mysecret

  # Login with email and be prompted for password (more secure)
  login --email user@example.com

  # Using short flags
  login -e user@example.com -p mysecret
`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Logging in...")

			if email == "" {
				log.Fatal("Email is required")
			}

			// If password not provided, prompt for it securely
			if password == "" {
				log.Fatal("Password is required")
			}

			// Sanitize inputs
			email = strings.ToLower(strings.TrimSpace(email))
			password = strings.TrimSpace(password)

			client := createE2EEClient()

			// Call the Login method
			response, err := client.Login(email, password)
			if err != nil {
				log.Fatalf("Failed to login: %v", err)
			}

			// Calculate token expiry duration for display
			duration := response.AccessTokenExpiryTime.Sub(time.Now()).Round(time.Second)

			// Print success message with expiry information
			fmt.Printf("✓ Successfully logged in as %s\n", email)
			fmt.Printf("- Access token expires in: %s\n", duration)

			// Add the home directory information
			homeDir, err := os.UserHomeDir()
			if err == nil {
				fmt.Printf("- Tokens saved to: %s/.maple/auth.json\n", homeDir)
			}
		},
	}

	// Define command flags
	cmd.Flags().StringVarP(&email, "email", "e", "", "Email address for the user (required)")
	cmd.Flags().StringVarP(&password, "password", "p", "", "Password for the user (will prompt if not provided)")

	// Mark required flags
	cmd.MarkFlagRequired("email")

	return cmd
}
