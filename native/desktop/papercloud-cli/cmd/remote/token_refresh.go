// native/desktop/papercloud-cli/cmd/remote/token_refresh.go
package remote

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	pref "github.com/Maple-Open-Tech/monorepo/native/desktop/papercloud-cli/internal/common/preferences"
)

// RefreshTokenRequest represents the request format for token refresh
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// RefreshTokenResponse represents the response format from token refresh
type RefreshTokenResponse struct {
	AccessToken            string    `json:"access_token"`
	RefreshToken           string    `json:"refresh_token"`
	AccessTokenExpiryTime  time.Time `json:"access_token_expiry_time"`
	RefreshTokenExpiryTime time.Time `json:"refresh_token_expiry_time"`
}

// RefreshTokens attempts to refresh the access token using the refresh token
// Returns true if successful, false otherwise
func RefreshTokens() (bool, error) {
	preferences := pref.PreferencesInstance()

	// Check if refresh token exists and is not expired
	if preferences.LoginResponse == nil ||
		preferences.LoginResponse.RefreshToken == "" ||
		time.Now().After(preferences.LoginResponse.RefreshTokenExpiryTime) {
		return false, fmt.Errorf("refresh token is missing or expired")
	}

	// Get server URL
	serverURL := preferences.CloudProviderAddress
	if serverURL == "" {
		serverURL = "http://localhost:8000" // Default if not configured
	}

	// Create refresh token request payload
	reqBody := RefreshTokenRequest{
		RefreshToken: preferences.LoginResponse.RefreshToken,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return false, fmt.Errorf("failed to marshal refresh token request: %v", err)
	}

	// Create request to refresh token endpoint
	refreshURL := fmt.Sprintf("%s/iam/api/v1/token/refresh", serverURL)
	req, err := http.NewRequest("POST", refreshURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return false, fmt.Errorf("failed to create refresh token request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Execute request
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to execute refresh token request: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return false, fmt.Errorf("refresh token request failed with status %d: %s",
			resp.StatusCode, string(body))
	}

	// Parse response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("failed to read refresh token response: %v", err)
	}

	var tokenResp RefreshTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return false, fmt.Errorf("failed to parse refresh token response: %v", err)
	}

	// Update stored tokens
	err = preferences.SetLoginResponse(
		tokenResp.AccessToken,
		tokenResp.AccessTokenExpiryTime,
		tokenResp.RefreshToken,
		tokenResp.RefreshTokenExpiryTime,
	)
	if err != nil {
		return false, fmt.Errorf("failed to save new tokens: %v", err)
	}

	fmt.Println("Access token refreshed successfully")
	return true, nil
}
