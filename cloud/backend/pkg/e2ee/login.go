// cloud/backend/pkg/e2ee/login.go
package e2ee

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// LoginOTTRequest is the payload to request a one-time token
type LoginOTTRequest struct {
	Email string `json:"email"`
}

// LoginOTTResponse is the server's response to an OTT request
type LoginOTTResponse struct {
	Message string `json:"message"`
}

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

// CompleteLoginRequest is the payload for completing login
type CompleteLoginRequest struct {
	Email         string `json:"email"`
	ChallengeID   string `json:"challengeId"`
	DecryptedData string `json:"decryptedData"`
}

// LoginResponse contains the server's response after a successful login
type LoginResponse struct {
	AccessToken            string    `json:"access_token"`
	AccessTokenExpiryTime  time.Time `json:"access_token_expiry_time"`
	RefreshToken           string    `json:"refresh_token"`
	RefreshTokenExpiryTime time.Time `json:"refresh_token_expiry_time"`
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
		return fmt.Errorf("failed to marshal OTT request data: %v", err)
	}

	// Create request
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send OTT request: %v", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("OTT request failed with status %d: %s",
			resp.StatusCode, body)
	}

	return nil
}

// VerifyLoginOTT verifies a one-time token and initiates the password verification
func (c *Client) VerifyLoginOTT(email, ott string) (*VerifyOTTResponse, error) {
	// Create request payload
	payload := &VerifyOTTRequest{
		Email: email,
		OTT:   ott,
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
		return nil, fmt.Errorf("failed to marshal OTT verification data: %v", err)
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
		return nil, fmt.Errorf("failed to send OTT verification request: %v", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OTT verification failed with status %d: %s",
			resp.StatusCode, body)
	}

	// Parse the response
	var response VerifyOTTResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse OTT verification response: %v", err)
	}

	return &response, nil
}

// VerifyPasswordAndCompleteLogin verifies password locally and completes the login
func (c *Client) VerifyPasswordAndCompleteLogin(email, password string, ottResponse *VerifyOTTResponse) (*LoginResponse, error) {
	// Step 1: Decode the received data
	salt, err := base64.StdEncoding.DecodeString(ottResponse.Salt)
	if err != nil {
		return nil, fmt.Errorf("failed to decode salt: %v", err)
	}

	encryptedMasterKey, err := base64.StdEncoding.DecodeString(ottResponse.EncryptedMasterKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode encrypted master key: %v", err)
	}

	encryptedPrivateKey, err := base64.StdEncoding.DecodeString(ottResponse.EncryptedPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode encrypted private key: %v", err)
	}

	encryptedChallenge, err := base64.StdEncoding.DecodeString(ottResponse.EncryptedChallenge)
	if err != nil {
		return nil, fmt.Errorf("failed to decode encrypted challenge: %v", err)
	}

	// Step 2: Derive the key encryption key from the password
	keyEncryptionKey := deriveKeyFromPassword(password, salt)

	// Step 3: Attempt to decrypt the master key (this verifies the password)
	masterKey, err := decryptData(encryptedMasterKey, keyEncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("incorrect password: %v", err)
	}

	// Step 4: Decrypt the private key using the master key
	privateKey, err := decryptData(encryptedPrivateKey, masterKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt private key: %v", err)
	}

	// Step 5: Store the keys for subsequent operations
	c.Keys = &KeySet{
		MasterKey:  masterKey,
		PrivateKey: privateKey,
	}

	// Step 6: Decrypt the challenge using the master key
	decryptedChallenge, err := decryptData(encryptedChallenge, masterKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt challenge: %v", err)
	}

	// Step 7: Create the complete login request with the decrypted challenge
	completeLoginPayload := &CompleteLoginRequest{
		Email:         email,
		ChallengeID:   ottResponse.ChallengeID,
		DecryptedData: base64.StdEncoding.EncodeToString(decryptedChallenge),
	}

	return c.completeLogin(completeLoginPayload)
}

// completeLogin sends the decrypted challenge to complete the login process
func (c *Client) completeLogin(payload *CompleteLoginRequest) (*LoginResponse, error) {
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
	endpoint := fmt.Sprintf("%s/iam/api/v1/complete-login", serverURL)

	// Convert payload to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal login completion data: %v", err)
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
		return nil, fmt.Errorf("failed to send login completion request: %v", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("login completion failed with status %d: %s",
			resp.StatusCode, body)
	}

	// Parse the response
	var response LoginResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse login completion response: %v", err)
	}

	// Save tokens to local storage
	if err := SaveTokens(
		payload.Email,
		response.AccessToken,
		response.RefreshToken,
		response.AccessTokenExpiryTime,
	); err != nil {
		return nil, fmt.Errorf("failed to save tokens: %v", err)
	}

	return &response, nil
}

// Login provides a simplified interface for the multi-step login process
func (c *Client) Login(email, password string) (*LoginResponse, error) {
	// Step 1: Request a one-time token
	if err := c.RequestLoginOTT(email); err != nil {
		return nil, fmt.Errorf("failed to request OTT: %v", err)
	}

	// In a CLI application, prompt the user for the OTT
	fmt.Println("Please check your email for a one-time token and enter it when prompted.")
	var ott string
	fmt.Print("Enter OTT: ")
	fmt.Scanln(&ott)

	// Step 2: Verify the OTT and get the encrypted keys and challenge
	ottResponse, err := c.VerifyLoginOTT(email, ott)
	if err != nil {
		return nil, fmt.Errorf("failed to verify OTT: %v", err)
	}

	// Step 3: Verify password locally and complete the login
	return c.VerifyPasswordAndCompleteLogin(email, password, ottResponse)
}
