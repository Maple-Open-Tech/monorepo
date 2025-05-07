package e2ee

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"       // Changed from io/ioutil
	"net/http" // Added for potential future use if needed, matching ioutil deprecation recommendation
)

// VerifyOTTRequest is the payload to verify a one-time token
type VerifyOTTRequest struct {
	Email string `json:"email"`
	OTT   string `json:"ott"`
}

// VerifyOTTResponse contains encrypted keys and challenge
type VerifyOTTResponse struct {
	Salt                string `json:"salt"`
	PublicKey           string `json:"publicKey"`
	EncryptedMasterKey  string `json:"encryptedMasterKey"`
	EncryptedPrivateKey string `json:"encryptedPrivateKey"`
	EncryptedChallenge  string `json:"encryptedChallenge"`
	ChallengeID         string `json:"challengeId"`
}

// VerifyLoginOTT verifies a one-time token and initiates the password verification
func (c *Client) VerifyLoginOTT(email, ott string) (*VerifyOTTResponse, error) {
	// Create request payload
	payload := &VerifyOTTRequest{
		Email: email,
		OTT:   ott, // Be careful not to log the actual OTT value in production logs
	}

	// Get HTTP client
	client := c.Config.HTTPClient
	if client == nil {
		client = defaultHTTPClient()
	}

	// Prepare server URL
	serverURL := c.Config.ServerURL
	if serverURL == "" {
		serverURL = DefaultServerURL
	}
	endpoint := fmt.Sprintf("%s/iam/api/v1/verify-login-ott", serverURL)

	// Convert payload to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("VerifyLoginOTT: failed to marshal VerifyOTTRequest payload for email %s: %w", censorEmail(email), err)
	}

	// Create request
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("VerifyLoginOTT: failed to create POST request for %s: %w", endpoint, err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("VerifyLoginOTT: failed to send request to %s: %w", endpoint, err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body) // Changed from ioutil.ReadAll
	if err != nil {
		// Log the status code even if reading the body fails
		return nil, fmt.Errorf("VerifyLoginOTT: failed to read response body from %s (status %d): %w", endpoint, resp.StatusCode, err)
	}

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("VerifyLoginOTT: request to %s failed with status %d: %s",
			endpoint, resp.StatusCode, string(body))
	}

	// Parse the response
	var response VerifyOTTResponse
	if err := json.Unmarshal(body, &response); err != nil {
		// Log the raw body for debugging if unmarshalling fails
		return nil, fmt.Errorf("VerifyLoginOTT: failed to parse VerifyOTTResponse JSON from %s: %w. Raw body: %s", endpoint, err, string(body))
	}
	fmt.Printf("VerifyLoginOTT: success: %s\n", response)

	return &response, nil
}
