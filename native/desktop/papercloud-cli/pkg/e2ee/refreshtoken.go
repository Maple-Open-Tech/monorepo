// pkg/e2ee/refreshtoken.go
package e2ee

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	pref "github.com/Maple-Open-Tech/monorepo/native/desktop/papercloud-cli/internal/common/preferences"
)

// RefreshTokenRequest is the payload for token refresh
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// RefreshTokenResponse is the server's response to token refresh
type RefreshTokenResponse struct {
	AccessToken            string    `json:"access_token"`
	AccessTokenExpiryTime  time.Time `json:"access_token_expiry_time"`
	RefreshToken           string    `json:"refresh_token"`
	RefreshTokenExpiryTime time.Time `json:"refresh_token_expiry_time"`
}

// RefreshTokens attempts to refresh the access token using the refresh token
// Returns true if successful, false otherwise
func (c *Client) RefreshTokens() (bool, error) {
	// Get current preferences
	preferences := pref.PreferencesInstance()

	// Check if refresh token exists
	if preferences.LoginResponse == nil ||
		preferences.LoginResponse.RefreshToken == "" {
		return false, fmt.Errorf("no refresh token available")
	}

	// Create refresh token request payload
	payload := &RefreshTokenRequest{
		RefreshToken: preferences.LoginResponse.RefreshToken,
	}

	// Get the HTTP client to use
	client := c.Config.HTTPClient
	if client == nil {
		client = defaultHTTPClient()
	}

	// Prepare server URL
	serverURL := c.Config.ServerURL
	if serverURL == "" {
		serverURL = DefaultServerURL
	}
	endpoint := fmt.Sprintf("%s/iam/api/v1/token/refresh", serverURL)

	// Convert payload to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return false, fmt.Errorf("failed to marshal refresh token data: %w", err)
	}

	// Create request
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return false, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to send refresh token request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check response status
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("refresh token failed with status %d: %s",
			resp.StatusCode, string(body))
	}

	// Parse the response
	var response RefreshTokenResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return false, fmt.Errorf("failed to parse refresh token response: %w", err)
	}

	// Update tokens in preferences
	err = preferences.SetLoginResponse(
		response.AccessToken,
		response.AccessTokenExpiryTime,
		response.RefreshToken,
		response.RefreshTokenExpiryTime,
	)
	if err != nil {
		return false, fmt.Errorf("failed to save refreshed tokens: %w", err)
	}

	return true, nil
}
