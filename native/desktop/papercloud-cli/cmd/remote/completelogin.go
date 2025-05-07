// native/desktop/papercloud-cli/cmd/remote/login.go
package remote

import (
	"fmt"
	"log"
	"strings"

	"github.com/Maple-Open-Tech/monorepo/native/desktop/papercloud-cli/pkg/e2ee"
	"github.com/spf13/cobra"
)

func CompleteLoginCmd() *cobra.Command {
	var email, password string

	var cmd = &cobra.Command{
		Use:   "completelogin",
		Short: "Finish the login",
		Long: `
After registration and email verification, use this command to log in to your account.

Examples:
  # Login with email and password
  go run main.go completelogin --email user@example.com

`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Logging in...")

			if email == "" {
				log.Fatal("Email is required")
			}

			// Sanitize inputs
			email = strings.ToLower(strings.TrimSpace(email))

			client := createE2EEClient()

			req := preferences.VerifyOTTResponse

			ottResponse := &e2ee.VerifyOTTResponse{
				Salt:                req.Salt,
				PublicKey:           req.PublicKey,
				EncryptedMasterKey:  req.EncryptedMasterKey,
				EncryptedPrivateKey: req.EncryptedPrivateKey,
				EncryptedChallenge:  req.EncryptedChallenge,
				ChallengeID:         req.ChallengeID,
			}

			// Step 3: Verify password locally and complete the login
			loginResponse, err := client.VerifyPasswordAndCompleteLogin(email, password, ottResponse)
			if err != nil {
				// Wrap error with context of the overall Login operation
				log.Fatalf("Failed verify and complete login: %v", err)
			}

			if err := preferences.SetLoginResponse(
				loginResponse.AccessToken,
				loginResponse.AccessTokenExpiryTime,
				loginResponse.RefreshToken,
				loginResponse.RefreshTokenExpiryTime,
			); err != nil {
				log.Fatalf("Failed to set login response in preferences: %v", err)
			}

			fmt.Println("Login successful!")
			fmt.Printf("Access Token: %s\n", loginResponse.AccessToken)
			// Optionally print other relevant info, but be careful with sensitive data like refresh tokens or keys in logs.
		},
	}

	// Define command flags
	cmd.Flags().StringVarP(&email, "email", "e", "", "Email address for the user (required)")
	cmd.Flags().StringVarP(&password, "password", "p", "", "Password for the user (will prompt if not provided)")

	// Mark required flags
	cmd.MarkFlagRequired("email")

	return cmd
}
