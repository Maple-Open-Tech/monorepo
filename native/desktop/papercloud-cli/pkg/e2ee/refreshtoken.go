// pkg/e2ee/refreshtoken.go
package e2ee

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"

	pref "github.com/Maple-Open-Tech/monorepo/native/desktop/papercloud-cli/internal/common/preferences"
)

// RefreshTokenRequest is the payload for token refresh
type RefreshTokenRequest struct {
	Value string `json:"value"` // Server expects "value", not "refresh_token"
}

// RefreshTokens attempts to refresh the access token using the refresh token
func (c *Client) RefreshTokens() (bool, error) {
	logger := zap.S() // Use sugared logger
	logger.Debug("Attempting to refresh tokens")

	// Get current preferences
	preferences := pref.PreferencesInstance()

	// Check if refresh token exists
	if preferences.LoginResponse == nil ||
		preferences.LoginResponse.RefreshToken == "" {
		logger.Warn("No refresh token available in preferences to attempt refresh.")
		return false, fmt.Errorf("no refresh token available")
	}
	logger.Debugw("Refresh token found in preferences",
		"refresh_token_empty", preferences.LoginResponse.RefreshToken == "")

	// Create refresh token request payload
	payload := &RefreshTokenRequest{
		Value: preferences.LoginResponse.RefreshToken, // Key changed to "value"
	}
	logger.Debugw("Refresh token request payload prepared",
		"payload_value_set", payload.Value != "")

	// Get the HTTP client to use
	httpClient := c.Config.HTTPClient // Renamed to avoid conflict with package name
	if httpClient == nil {
		logger.Debug("HTTPClient not configured in Client.Config, using default HTTP client for token refresh")
		httpClient = defaultHTTPClient()
	} else {
		logger.Debug("Using pre-configured HTTPClient from Client.Config for token refresh")
	}

	// Prepare server URL - this should match your API endpoint exactly
	serverURL := c.Config.ServerURL
	if serverURL == "" {
		serverURL = DefaultServerURL
		logger.Debugw("ServerURL not configured in Client.Config, using default server URL",
			"default_server_url", DefaultServerURL)
	}
	endpoint := fmt.Sprintf("%s/iam/api/v1/token/refresh", serverURL)
	logger.Infow("Preparing to refresh token", "endpoint", endpoint)

	// Convert payload to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		logger.Errorw("Failed to marshal refresh token request payload", "error", err)
		return false, fmt.Errorf("failed to marshal refresh token data: %w", err)
	}

	// Create request
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		logger.Errorw("Failed to create HTTP request for token refresh",
			"error", err,
			"endpoint", endpoint,
		)
		return false, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	logger.Infow("Sending token refresh request", "method", req.Method, "url", req.URL.String())

	// Send request
	resp, err := httpClient.Do(req)
	if err != nil {
		logger.Errorw("Failed to send token refresh request",
			"error", err,
			"endpoint", endpoint,
		)
		return false, fmt.Errorf("failed to send refresh token request: %w", err)
	}
	defer resp.Body.Close()
	logger.Debugw("Token refresh request sent", "endpoint", endpoint)

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Errorw("Failed to read response body from token refresh",
			"error", err,
			"status_code", resp.StatusCode,
		)
		return false, fmt.Errorf("failed to read response body: %w", err)
	}
	logger.Debugw("Token refresh response body read",
		"status_code", resp.StatusCode,
		"body_length", len(body))

	logger.Infow("Received token refresh response", "status_code", resp.StatusCode)

	// Check response status
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		logger.Errorw("Token refresh request failed with unexpected status code",
			"status_code", resp.StatusCode,
			"response_body", string(body),
			"endpoint", endpoint,
		)
		return false, fmt.Errorf("refresh token failed with status %d: %s",
			resp.StatusCode, string(body))
	}

	// Parse the response
	var response struct {
		AccessToken            string    `json:"access_token"`
		AccessTokenExpiryTime  time.Time `json:"access_token_expiry_time"`
		RefreshToken           string    `json:"refresh_token"`
		RefreshTokenExpiryTime time.Time `json:"refresh_token_expiry_time"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		logger.Errorw("Failed to parse JSON response from token refresh",
			"error", err,
			"response_body", string(body),
			"status_code", resp.StatusCode,
		)
		return false, fmt.Errorf("failed to parse refresh token response: %w", err)
	}
	logger.Debugw("Token refresh response parsed successfully",
		"access_token_present", response.AccessToken != "",
		"refresh_token_present", response.RefreshToken != "",
		"access_token_expiry", response.AccessTokenExpiryTime,
		"refresh_token_expiry", response.RefreshTokenExpiryTime,
	)

	// Update tokens in preferences
	err = preferences.SetLoginResponse(
		response.AccessToken,
		response.AccessTokenExpiryTime,
		response.RefreshToken,
		response.RefreshTokenExpiryTime,
	)
	if err != nil {
		logger.Errorw("Failed to save refreshed tokens to preferences", "error", err)
		return false, fmt.Errorf("failed to save refreshed tokens: %w", err)
	}

	logger.Info("Tokens refreshed and saved successfully")
	return true, nil
}
