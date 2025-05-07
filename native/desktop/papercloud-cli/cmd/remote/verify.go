// native/desktop/papercloud-cli/cmd/remote/verify.go
package remote

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
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
			fmt.Println("Verifying email...")

			if verificationCode == "" {
				log.Fatal("Verification code is required")
			}

			client := createE2EEClient()

			// Call the VerifyEmail method
			response, err := client.VerifyEmail(verificationCode)
			if err != nil {
				log.Fatalf("Failed to verify email: %v", err)
			}

			// Display success message
			fmt.Printf("Email verification successful! %s\n", response.Message)
		},
	}

	// Define command flags
	cmd.Flags().StringVarP(&verificationCode, "code", "c", "", "Email verification code (required)")

	// Mark required flags
	cmd.MarkFlagRequired("code")

	return cmd
}
