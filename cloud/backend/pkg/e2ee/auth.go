package e2ee

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// AuthenticatedRequest sends an authenticated request to the server
func (c *Client) AuthenticatedRequest(method, endpoint string, payload interface{}) ([]byte, error) {
	// Load tokens
	tokens, err := LoadTokens()
	if err != nil {
		return nil, fmt.Errorf("failed to load tokens: %v", err)
	}

	// Check if tokens are expired
	if IsTokenExpired(tokens) {
		// Try to refresh the token
		refreshResp, err := c.RefreshToken(tokens.RefreshToken)
		if err != nil {
			return nil, fmt.Errorf("session expired and token refresh failed: %v", err)
		}

		// Update tokens
		if err := SaveTokens(
			tokens.Email,
			refreshResp.AccessToken,
			refreshResp.RefreshToken,
			refreshResp.AccessTokenExpiryTime,
		); err != nil {
			return nil, fmt.Errorf("failed to save refreshed tokens: %v", err)
		}

		// Update tokens object
		tokens.AccessToken = refreshResp.AccessToken
		tokens.RefreshToken = refreshResp.RefreshToken
		tokens.ExpiresAt = refreshResp.AccessTokenExpiryTime
	}

	// Get the HTTP client to use
	client := c.Config.HTTPClient
	if client == nil {
		client = defaultHTTPClient()
	}

	// Prepare URL
	serverURL := c.Config.ServerURL
	if serverURL == "" {
		serverURL = DefaultServerURL
	}
	url := fmt.Sprintf("%s%s", serverURL, endpoint)

	// Convert payload to JSON if provided
	var body []byte
	if payload != nil {
		var err error
		body, err = json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request data: %v", err)
		}
	}

	// Create request
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Set headers
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Authorization", "JWT "+tokens.AccessToken)

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Read response
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// Check response status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("request failed with status %d: %s",
			resp.StatusCode, responseBody)
	}

	return responseBody, nil
}
