package client

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/e2ee"
)

func VerifyEmailCmd() *cobra.Command {
	var verificationCode string

	var cmd = &cobra.Command{
		Use:   "verifyemail",
		Short: "Verify email with activation code",
		Long: `
Command will submit your email activation code to the backend to finalize registration.

After registration, you will receive an email with a verification code.
Use this command to submit that code and activate your account.

Examples:
  # Verify email with a code
  verifyemail --code 123456

  # Alternatively, use the short flag
  verifyemail -c 123456
`,
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("Verifying email...")

			if verificationCode == "" {
				log.Fatal("Verification code is required")
			}

			if err := verifyEmail(verificationCode); err != nil {
				log.Fatalf("Failed to verify email: %v", err)
			}
		},
	}

	// Define command flags
	cmd.Flags().StringVarP(&verificationCode, "code", "c", "", "Email verification code (required)")

	// Mark required flags
	cmd.MarkFlagRequired("code")

	return cmd
}

// verifyEmail handles the actual verification logic
func verifyEmail(verificationCode string) error {
	cfg := config.NewProvider()

	// Create a new E2EE client
	client := e2ee.NewClient(e2ee.ClientConfig{
		ServerURL: fmt.Sprintf("%s://%s:%s", "http", cfg.App.IP, cfg.App.Port),
	})

	// Call the VerifyEmail method
	response, err := client.VerifyEmail(verificationCode)
	if err != nil {
		return fmt.Errorf("e2ee.Client.VerifyEmail failed: %w", err)
	}

	// Display success message
	fmt.Printf("Email verification successful! %s\n", response.Message)

	return nil
}
