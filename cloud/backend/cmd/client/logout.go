package client

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/spf13/cobra"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/e2ee"
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
			log.Println("Logging out...")

			if err := logoutUser(); err != nil {
				log.Fatalf("Failed to logout: %v", err)
			}
		},
	}

	return cmd
}

func logoutUser() error {
	// Load current tokens
	tokens, err := e2ee.LoadTokens()
	if err != nil {
		// If no tokens exist, we're already logged out
		if strings.Contains(err.Error(), "not logged in") {
			fmt.Println("You are not currently logged in.")
			return nil
		}
		return fmt.Errorf("failed to load current session: %w", err)
	}

	cfg := config.NewProvider()

	// Create a new E2EE client
	client := e2ee.NewClient(e2ee.ClientConfig{
		ServerURL: fmt.Sprintf("%s://%s:%s", "http", cfg.App.IP, cfg.App.Port),
	})

	// Call the server's logout endpoint
	// Note: This is optional, as many REST APIs are stateless and don't require server-side logout
	if tokens.AccessToken != "" {
		// Create request with Authorization header
		req, err := http.NewRequest("POST", fmt.Sprintf("%s/iam/api/v1/logout", client.Config.ServerURL), nil)
		if err != nil {
			return fmt.Errorf("failed to create request: %v", err)
		}
		req.Header.Set("Authorization", "JWT "+tokens.AccessToken)

		// Send request
		httpClient := defaultHTTPClient()
		resp, err := httpClient.Do(req)
		if err != nil {
			// Log the error but continue with local logout
			log.Printf("Warning: Failed to notify server about logout: %v", err)
		} else {
			defer resp.Body.Close()
			// We don't need to process the response
		}
	}

	// Clear local tokens
	if err := e2ee.ClearTokens(); err != nil {
		return fmt.Errorf("failed to clear tokens: %w", err)
	}

	fmt.Println("✓ Successfully logged out")
	return nil
}
