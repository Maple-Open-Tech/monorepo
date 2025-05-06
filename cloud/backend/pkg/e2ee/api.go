package e2ee

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// HTTPClient interface allows for custom HTTP clients and easier testing
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// defaultHTTPClient returns the default HTTP client
func defaultHTTPClient() HTTPClient {
	return &http.Client{}
}

// sendRegistrationRequest sends the registration payload to the server
// This is separated from the registration logic to allow for easier
// customization and testing
func sendRegistrationRequest(config ClientConfig, payload *RegistrationPayload) ([]byte, error) {
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
	endpoint := fmt.Sprintf("%s/iam/api/v1/register", serverURL)

	// Convert payload to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal registration data: %v", err)
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
		return nil, fmt.Errorf("failed to send registration request: %v", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// Check response status
	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("registration failed with status %d: %s",
			resp.StatusCode, body)
	}

	return body, nil
}
