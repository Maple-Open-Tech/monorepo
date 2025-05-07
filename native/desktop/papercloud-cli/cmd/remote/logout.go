// native/desktop/papercloud-cli/cmd/remote/logout.go
package remote

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

func LogoutUserCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "logout",
		Short: "Log user out of account",
		Long: `
Command will log out the current user by removing stored authentication tokens.
This will invalidate the current session and require logging in again for future API calls.

Example:
  # Logout the current user
  logout
`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Logging out...")

			client := createE2EEClient()

			// Call the Logout method
			err := client.Logout()
			if err != nil {
				log.Fatalf("Failed to logout: %v", err)
			}

			fmt.Println("✓ Successfully logged out")
		},
	}

	return cmd
}
