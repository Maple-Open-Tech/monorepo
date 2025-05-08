// pkg/e2ee/auth.go
package e2ee

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	pref "github.com/Maple-Open-Tech/monorepo/native/desktop/papercloud-cli/internal/common/preferences"
)

// IsAuthenticated checks if the user is authenticated with valid tokens
func (c *Client) IsAuthenticated() bool {
	preferences := pref.PreferencesInstance()

	// Check if tokens exist and are not expired
	if preferences.LoginResponse == nil ||
		preferences.LoginResponse.AccessToken == "" {
		return false
	}

	// Check token expiry
	if time.Now().After(preferences.LoginResponse.AccessTokenExpiryTime) {
		return false
	}

	return true
}

// AuthenticatedRequest sends an authenticated request to the server
func (c *Client) AuthenticatedRequest(method, endpoint string, payload interface{}) ([]byte, error) {
	// Check authentication
	if !c.IsAuthenticated() {
		return nil, fmt.Errorf("not authenticated or token expired: please login again")
	}

	preferences := pref.PreferencesInstance()

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
			return nil, fmt.Errorf("failed to marshal request data: %w", err)
		}
	}

	// Create request
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Authorization", "JWT "+preferences.LoginResponse.AccessToken)

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check response status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("request failed with status %d: %s",
			resp.StatusCode, string(responseBody))
	}

	return responseBody, nil
}

// AuthenticatedFormRequest sends a multipart form request with authentication
func (c *Client) AuthenticatedFormRequest(method, endpoint string, formData map[string]string, formFiles map[string]io.Reader) ([]byte, error) {
	// Check authentication
	if !c.IsAuthenticated() {
		return nil, fmt.Errorf("not authenticated or token expired: please login again")
	}

	preferences := pref.PreferencesInstance()

	// Create pipe for streaming the form data
	pipeReader, pipeWriter := io.Pipe()
	multipartWriter := multipart.NewWriter(pipeWriter)

	// Write form data in a goroutine
	go func() {
		defer pipeWriter.Close()

		// Add form fields
		for key, value := range formData {
			if err := multipartWriter.WriteField(key, value); err != nil {
				pipeWriter.CloseWithError(fmt.Errorf("failed to write form field %s: %w", key, err))
				return
			}
		}

		// Add form files
		for name, reader := range formFiles {
			fileWriter, err := multipartWriter.CreateFormFile(name, name)
			if err != nil {
				pipeWriter.CloseWithError(fmt.Errorf("failed to create form file %s: %w", name, err))
				return
			}

			if _, err := io.Copy(fileWriter, reader); err != nil {
				pipeWriter.CloseWithError(fmt.Errorf("failed to write form file %s: %w", name, err))
				return
			}
		}

		if err := multipartWriter.Close(); err != nil {
			pipeWriter.CloseWithError(fmt.Errorf("failed to close multipart writer: %w", err))
			return
		}
	}()

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

	// Create request
	req, err := http.NewRequest(method, url, pipeReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", multipartWriter.FormDataContentType())
	req.Header.Set("Authorization", "JWT "+preferences.LoginResponse.AccessToken)

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check response status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("request failed with status %d: %s",
			resp.StatusCode, string(responseBody))
	}

	return responseBody, nil
}
