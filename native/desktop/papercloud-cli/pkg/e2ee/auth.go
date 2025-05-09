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
	"go.uber.org/zap"
)

var logger *zap.Logger

func init() {
	// Initialize a default logger. In a real application, this would be
	// configured and potentially passed in or set up in main.
	var err error
	// Using NewDevelopmentConfig for more verbose, human-readable output.
	// For production, consider zap.NewProductionConfig() or a custom configuration.
	config := zap.NewDevelopmentConfig()
	// AddCallerSkip(1) so that the caller location is correctly reported as the line
	// where logger.Info/Debug/Error etc. is called, not from within zap's internals.
	logger, err = config.Build(zap.AddCallerSkip(1))
	if err != nil {
		// Fallback or panic if logger initialization is critical.
		// Using fmt.Printf for this fallback as logger might not be available.
		fmt.Printf("WARN: failed to initialize zap logger: %v. Falling back to NopLogger.\n", err)
		logger = zap.NewNop()
	}
	// Optionally, replace global logger if other parts of the application use zap.L() or zap.S()
	// zap.ReplaceGlobals(logger)
}

// IsAuthenticated checks if the user is authenticated with valid tokens
// and automatically attempts to refresh the token if expired
func (c *Client) IsAuthenticated() bool {
	logger.Debug("Checking authentication status")
	preferences := pref.PreferencesInstance()

	if preferences.LoginResponse == nil {
		logger.Info("User not authenticated: LoginResponse is nil")
		return false
	}
	if preferences.LoginResponse.AccessToken == "" {
		logger.Info("User not authenticated: AccessToken is empty")
		return false
	}

	// Check token expiry
	if time.Now().After(preferences.LoginResponse.AccessTokenExpiryTime) {
		logger.Info("Access token expired",
			zap.Time("expiry_time", preferences.LoginResponse.AccessTokenExpiryTime),
			zap.Time("current_time", time.Now()))

		// Check if refresh token exists and is not expired
		if preferences.LoginResponse.RefreshToken == "" {
			logger.Info("No refresh token available")
			return false
		}

		if time.Now().After(preferences.LoginResponse.RefreshTokenExpiryTime) {
			logger.Info("Refresh token expired",
				zap.Time("refresh_expiry_time", preferences.LoginResponse.RefreshTokenExpiryTime))
			return false
		}

		// Try to refresh the tokens
		logger.Info("Attempting to refresh tokens")
		success, err := c.RefreshTokens()
		if err != nil {
			logger.Error("Failed to refresh tokens", zap.Error(err))
			return false
		}

		if !success {
			logger.Info("Token refresh was not successful")
			return false
		}

		logger.Info("Successfully refreshed tokens")
		return true // Successfully refreshed tokens
	}

	logger.Debug("User is authenticated with valid tokens")
	return true
}

// AuthenticatedRequest sends an authenticated request to the server
// and automatically attempts to refresh tokens if authentication fails
func (c *Client) AuthenticatedRequest(method, endpoint string, payload interface{}) ([]byte, error) {
	logger.Info("AuthenticatedRequest initiated",
		zap.String("method", method),
		zap.String("endpoint", endpoint))

	if !c.IsAuthenticated() {
		err := fmt.Errorf("not authenticated or token expired: please login again")
		logger.Warn("Authentication check failed in AuthenticatedRequest", zap.Error(err))
		return nil, err
	}

	// Execute the request with current token
	responseBytes, statusCode, err := c.executeRequest(method, endpoint, payload, false)

	// If we get a 401, try to refresh the token and retry once
	if statusCode == http.StatusUnauthorized {
		logger.Info("Received 401 Unauthorized, attempting to refresh token and retry")
		success, err := c.RefreshTokens()
		if err != nil || !success {
			logger.Error("Token refresh failed", zap.Error(err))
			return nil, fmt.Errorf("authentication failed and token refresh failed: %w", err)
		}

		// Retry with new token
		logger.Info("Token refreshed successfully, retrying request")
		responseBytes, statusCode, err = c.executeRequest(method, endpoint, payload, true)
	}

	return responseBytes, err
}

// executeRequest executes the actual request
// isRetry indicates if this is a retry after token refresh
func (c *Client) executeRequest(method, endpoint string, payload interface{}, isRetry bool) ([]byte, int, error) {
	preferences := pref.PreferencesInstance()

	httpClient := c.Config.HTTPClient
	if httpClient == nil {
		logger.Debug("Using default HTTP client for AuthenticatedRequest")
		httpClient = defaultHTTPClient()
	}

	serverURL := c.Config.ServerURL
	if serverURL == "" {
		serverURL = DefaultServerURL
		logger.Debug("Using default server URL", zap.String("url", serverURL))
	}
	fullURL := fmt.Sprintf("%s%s", serverURL, endpoint)
	logger.Debug("Request details",
		zap.String("url", fullURL),
		zap.String("method", method))

	var reqBodyBytes []byte
	if payload != nil {
		var err error
		reqBodyBytes, err = json.Marshal(payload)
		if err != nil {
			logger.Error("Failed to marshal request payload", zap.Error(err), zap.Any("payload_type", fmt.Sprintf("%T", payload)))
			return nil, 0, fmt.Errorf("failed to marshal request data: %w", err)
		}
		logger.Debug("Request payload marshaled", zap.Int("payload_size_bytes", len(reqBodyBytes)))
	}

	req, err := http.NewRequest(method, fullURL, bytes.NewBuffer(reqBodyBytes))
	if err != nil {
		logger.Error("Failed to create HTTP request",
			zap.String("url", fullURL),
			zap.String("method", method),
			zap.Error(err))
		return nil, 0, fmt.Errorf("failed to create request: %w", err)
	}

	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Get the latest token from preferences
	accessToken := preferences.LoginResponse.AccessToken
	authHeader := fmt.Sprintf("JWT %s", accessToken)
	req.Header.Set("Authorization", authHeader)

	if isRetry {
		logger.Debug("Retry request with refreshed token",
			zap.String("token_prefix", "JWT "+accessToken[:10]+"..."))
	}
	logger.Debug("Request headers set", zap.Any("headers", req.Header))

	logger.Info("Sending authenticated HTTP request",
		zap.String("method", req.Method),
		zap.String("url", req.URL.String()))

	resp, err := httpClient.Do(req)
	if err != nil {
		logger.Error("HTTP request failed",
			zap.String("url", fullURL),
			zap.String("method", method),
			zap.Error(err))
		return nil, 0, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	statusCode := resp.StatusCode
	logger.Debug("Received HTTP response",
		zap.Int("status_code", statusCode),
		zap.String("status", resp.Status))

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Failed to read response body",
			zap.Int("status_code", statusCode),
			zap.Error(err))
		return nil, statusCode, fmt.Errorf("failed to read response body: %w", err)
	}

	if statusCode < 200 || statusCode >= 300 {
		err := fmt.Errorf("request failed with status %d: %s", statusCode, string(responseBody))
		logger.Warn("HTTP request returned non-2xx status",
			zap.Int("status_code", statusCode),
			zap.ByteString("response_body", responseBody),
			zap.Error(err))
		return nil, statusCode, err
	}

	logger.Info("AuthenticatedRequest completed successfully",
		zap.Int("status_code", statusCode),
		zap.Int("response_size_bytes", len(responseBody)))
	return responseBody, statusCode, nil
}

// AuthenticatedFormRequest sends a multipart form request with authentication
// and automatically attempts to refresh tokens if authentication fails
func (c *Client) AuthenticatedFormRequest(method, endpoint string, formData map[string]string, formFiles map[string]io.Reader) ([]byte, error) {
	formFieldKeys := make([]string, 0, len(formData))
	for k := range formData {
		formFieldKeys = append(formFieldKeys, k)
	}
	formFileKeys := make([]string, 0, len(formFiles))
	for k := range formFiles {
		formFileKeys = append(formFileKeys, k)
	}

	logger.Info("AuthenticatedFormRequest initiated",
		zap.String("method", method),
		zap.String("endpoint", endpoint),
		zap.Strings("form_data_keys", formFieldKeys),
		zap.Strings("form_file_keys", formFileKeys))

	if !c.IsAuthenticated() {
		err := fmt.Errorf("not authenticated or token expired: please login again")
		logger.Warn("Authentication check failed in AuthenticatedFormRequest", zap.Error(err))
		return nil, err
	}

	// Execute the request with current token
	responseBytes, statusCode, err := c.executeFormRequest(method, endpoint, formData, formFiles, false)

	// If we get a 401, try to refresh the token and retry once
	if statusCode == http.StatusUnauthorized {
		logger.Info("Received 401 Unauthorized, attempting to refresh token and retry")
		success, err := c.RefreshTokens()
		if err != nil || !success {
			logger.Error("Token refresh failed", zap.Error(err))
			return nil, fmt.Errorf("authentication failed and token refresh failed: %w", err)
		}

		// Retry with new token
		logger.Info("Token refreshed successfully, retrying request")
		responseBytes, statusCode, err = c.executeFormRequest(method, endpoint, formData, formFiles, true)
	}

	return responseBytes, err
}

// executeFormRequest executes the actual form request
// isRetry indicates if this is a retry after token refresh
func (c *Client) executeFormRequest(method, endpoint string, formData map[string]string, formFiles map[string]io.Reader, isRetry bool) ([]byte, int, error) {
	preferences := pref.PreferencesInstance()

	pipeReader, pipeWriter := io.Pipe()
	multipartWriter := multipart.NewWriter(pipeWriter)

	// Goroutine to write multipart data to the pipe
	go func() {
		// Ensure pipeWriter is always closed
		defer pipeWriter.Close()

		var opError error // To store error from multipart operations

		// Add form fields
		for key, value := range formData {
			logger.Debug("Writing form field to multipart", zap.String("key", key))
			if opError = multipartWriter.WriteField(key, value); opError != nil {
				errMsg := fmt.Errorf("failed to write form field '%s': %w", key, opError)
				logger.Error("Error writing form field in goroutine", zap.Error(errMsg))
				pipeWriter.CloseWithError(errMsg) // Close pipe with error
				return                            // Exit goroutine
			}
		}

		// Add form files
		for name, reader := range formFiles {
			logger.Debug("Creating form file in multipart", zap.String("name", name))
			fileWriter, errCreate := multipartWriter.CreateFormFile(name, name) // Use 'name' as filename in form
			if errCreate != nil {
				errMsg := fmt.Errorf("failed to create form file '%s': %w", name, errCreate)
				logger.Error("Error creating form file in goroutine", zap.Error(errMsg))
				pipeWriter.CloseWithError(errMsg)
				return
			}
			logger.Debug("Copying file content to multipart", zap.String("name", name))
			if _, errCopy := io.Copy(fileWriter, reader); errCopy != nil {
				errMsg := fmt.Errorf("failed to copy content for form file '%s': %w", name, errCopy)
				logger.Error("Error copying file content in goroutine", zap.Error(errMsg))
				pipeWriter.CloseWithError(errMsg)
				return
			}
		}

		// Close multipart writer to finalize the body
		logger.Debug("Closing multipart writer in goroutine")
		if opError = multipartWriter.Close(); opError != nil {
			errMsg := fmt.Errorf("failed to close multipart writer: %w", opError)
			logger.Error("Error closing multipart writer in goroutine", zap.Error(errMsg))
			pipeWriter.CloseWithError(errMsg)
			return
		}
		logger.Debug("Multipart stream successfully written by goroutine. Pipe will be closed by defer.")
	}()

	httpClient := c.Config.HTTPClient
	if httpClient == nil {
		logger.Debug("Using default HTTP client for AuthenticatedFormRequest")
		httpClient = defaultHTTPClient()
	}

	serverURL := c.Config.ServerURL
	if serverURL == "" {
		serverURL = DefaultServerURL
		logger.Debug("Using default server URL for form request", zap.String("url", serverURL))
	}
	fullURL := fmt.Sprintf("%s%s", serverURL, endpoint)

	req, err := http.NewRequest(method, fullURL, pipeReader)
	if err != nil {
		logger.Error("Failed to create HTTP form request",
			zap.String("url", fullURL),
			zap.String("method", method),
			zap.Error(err))
		_ = pipeWriter.CloseWithError(fmt.Errorf("http request creation failed, aborting pipe: %w", err))
		return nil, 0, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", multipartWriter.FormDataContentType())

	// Get the latest token from preferences
	accessToken := preferences.LoginResponse.AccessToken
	authHeader := fmt.Sprintf("JWT %s", accessToken)
	req.Header.Set("Authorization", authHeader)

	if isRetry {
		logger.Debug("Retry request with refreshed token",
			zap.String("token_prefix", "JWT "+accessToken[:10]+"..."))
	}
	logger.Debug("Form request headers set", zap.Any("headers", req.Header))

	logger.Info("Sending authenticated form HTTP request",
		zap.String("method", req.Method),
		zap.String("url", req.URL.String()))

	resp, err := httpClient.Do(req)
	if err != nil {
		logger.Error("HTTP form request failed",
			zap.String("url", fullURL),
			zap.String("method", method),
			zap.Error(err))
		return nil, 0, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	statusCode := resp.StatusCode
	logger.Debug("Received HTTP response for form request",
		zap.Int("status_code", statusCode),
		zap.String("status", resp.Status))

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Failed to read response body for form request",
			zap.Int("status_code", statusCode),
			zap.Error(err))
		return nil, statusCode, fmt.Errorf("failed to read response body: %w", err)
	}

	if statusCode < 200 || statusCode >= 300 {
		err := fmt.Errorf("form request failed with status %d: %s", statusCode, string(responseBody))
		logger.Warn("HTTP form request returned non-2xx status",
			zap.Int("status_code", statusCode),
			zap.ByteString("response_body", responseBody),
			zap.Error(err))
		return nil, statusCode, err
	}

	logger.Info("AuthenticatedFormRequest completed successfully",
		zap.Int("status_code", statusCode),
		zap.Int("response_size_bytes", len(responseBody)))
	return responseBody, statusCode, nil
}
