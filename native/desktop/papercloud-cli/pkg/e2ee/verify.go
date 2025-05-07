package e2ee

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// VerifyEmailPayload contains data sent to the server for email verification
type VerifyEmailPayload struct {
	Code string `json:"code"`
}

// VerifyEmailResponse contains the server's response after verification
type VerifyEmailResponse struct {
	Message           string `json:"message"`
	FederatedUserRole int8   `json:"user_role"`
}

// VerifyEmail sends the verification code to the server to verify the user's email
func (c *Client) VerifyEmail(code string) (*VerifyEmailResponse, error) {
	// Create the verification payload
	payload := &VerifyEmailPayload{
		Code: code,
	}

	// Send the verification request
	responseData, err := sendVerificationRequest(c.Config, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to verify email: %v", err)
	}

	// Parse the response
	var response VerifyEmailResponse
	if err := json.Unmarshal(responseData, &response); err != nil {
		return nil, fmt.Errorf("failed to parse verification response: %v", err)
	}

	return &response, nil
}

// sendVerificationRequest sends the verification payload to the server
func sendVerificationRequest(config ClientConfig, payload *VerifyEmailPayload) ([]byte, error) {
	// Get the HTTP client to use
	client := config.HTTPClient
	if client == nil {
		client = defaultHTTPClient()
	}

	// Prepare server URL
	serverURL := config.ServerURL
	if serverURL == "" {
		serverURL = DefaultServerURL
	}
	endpoint := fmt.Sprintf("%s/iam/api/v1/verify-email-code", serverURL)

	// Convert payload to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal verification data: %v", err)
	}

	// Create request
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send verification request: %v", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// Check response status
	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("verification failed with status %d: %s",
			resp.StatusCode, body)
	}

	return body, nil
}
