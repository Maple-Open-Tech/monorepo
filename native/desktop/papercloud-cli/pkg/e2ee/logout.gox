package e2ee

import (
	"fmt"
	"net/http"
	"strings"
)

// Logout ends the user's session both on the server and locally
func (c *Client) Logout() error {
	// Load current tokens
	tokens, err := LoadTokens()
	if err != nil {
		// If no tokens exist, we're already logged out
		if strings.Contains(err.Error(), "not logged in") {
			return nil // Already logged out
		}
		return fmt.Errorf("failed to load current session: %w", err)
	}

	// Call the server's logout endpoint if we have an access token
	if tokens.AccessToken != "" {
		// Create request with Authorization header
		req, err := http.NewRequest("POST", fmt.Sprintf("%s/iam/api/v1/logout", c.Config.ServerURL), nil)
		if err != nil {
			return fmt.Errorf("failed to create request: %v", err)
		}
		req.Header.Set("Authorization", "JWT "+tokens.AccessToken)

		// Send request
		httpClient := c.Config.HTTPClient
		if httpClient == nil {
			httpClient = defaultHTTPClient()
		}

		resp, err := httpClient.Do(req)
		if err != nil {
			// Log the error but continue with local logout
			// We don't return the error here because we still want to clear local tokens
			fmt.Printf("Warning: Failed to notify server about logout: %v\n", err)
		} else {
			defer resp.Body.Close()
			// Check for server error
			if resp.StatusCode >= 400 {
				fmt.Printf("Warning: Server returned status %d during logout\n", resp.StatusCode)
				// Continue with local logout regardless
			}
		}
	}

	// Clear local tokens
	if err := ClearTokens(); err != nil {
		return fmt.Errorf("failed to clear tokens: %w", err)
	}

	return nil
}
