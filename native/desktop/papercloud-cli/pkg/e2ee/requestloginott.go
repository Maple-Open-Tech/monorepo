package e2ee

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"       // Changed from io/ioutil
	"net/http" // Added for potential future use if needed, matching ioutil deprecation recommendation
)

// LoginOTTRequest is the payload to request a one-time token
type LoginOTTRequest struct {
	Email string `json:"email"`
}

// LoginOTTResponse is the server's response to an OTT request
type LoginOTTResponse struct {
	Message string `json:"message"`
}

// RequestLoginOTT requests a one-time token sent to the user's email
func (c *Client) RequestLoginOTT(email string) error {
	// Create request payload
	payload := &LoginOTTRequest{
		Email: email,
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
	endpoint := fmt.Sprintf("%s/iam/api/v1/request-login-ott", serverURL)

	// Convert payload to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("RequestLoginOTT: failed to marshal LoginOTTRequest payload for email %s: %w", censorEmail(email), err)
	}

	// Create request
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("RequestLoginOTT: failed to create POST request for %s: %w", endpoint, err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("RequestLoginOTT: failed to send request to %s: %w", endpoint, err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body) // Changed from ioutil.ReadAll
	if err != nil {
		// Log the status code even if reading the body fails
		return fmt.Errorf("RequestLoginOTT: failed to read response body from %s (status %d): %w", endpoint, resp.StatusCode, err)
	}

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("RequestLoginOTT: request to %s failed with status %d: %s",
			endpoint, resp.StatusCode, string(body))
	}

	// Optionally log success or parse the LoginOTTResponse if needed
	var response LoginOTTResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("RequestLoginOTT: failed to parse success response body: %w", err)
	}
	fmt.Printf("RequestLoginOTT: success: %s\n", response.Message)

	return nil
}
