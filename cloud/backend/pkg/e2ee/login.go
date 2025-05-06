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

// LoginPayload contains all data sent to the server during login
type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse contains the server's response after a successful login
type LoginResponse struct {
	AccessToken            string    `json:"access_token"`
	AccessTokenExpiryTime  time.Time `json:"access_token_expiry_time"`
	RefreshToken           string    `json:"refresh_token"`
	RefreshTokenExpiryTime time.Time `json:"refresh_token_expiry_time"`

	// Optional fields that may be returned based on your API implementation
	Salt                string `json:"salt,omitempty"`
	EncryptedMasterKey  string `json:"encryptedMasterKey,omitempty"`
	EncryptedPrivateKey string `json:"encryptedPrivateKey,omitempty"`
	PublicKey           string `json:"publicKey,omitempty"`
}

// Login authenticates the user and retrieves encryption keys
func (c *Client) Login(email, password string) (*LoginResponse, error) {
	// Create the login payload
	payload := &LoginPayload{
		Email:    email,
		Password: password,
	}

	// Send the login request
	responseData, err := sendLoginRequest(c.Config, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to login: %v", err)
	}

	// Parse the response
	var response LoginResponse
	if err := json.Unmarshal(responseData, &response); err != nil {
		return nil, fmt.Errorf("failed to parse login response: %v", err)
	}

	// If the encryption keys are returned, decrypt and store them
	if response.Salt != "" && response.EncryptedMasterKey != "" {
		// Decode from base64
		salt, err := base64.StdEncoding.DecodeString(response.Salt)
		if err != nil {
			return nil, fmt.Errorf("failed to decode salt: %v", err)
		}

		encryptedMasterKey, err := base64.StdEncoding.DecodeString(response.EncryptedMasterKey)
		if err != nil {
			return nil, fmt.Errorf("failed to decode encrypted master key: %v", err)
		}

		// Derive the key encryption key from the password
		keyEncryptionKey := deriveKeyFromPassword(password, salt)

		// Decrypt the master key
		masterKey, err := decryptData(encryptedMasterKey, keyEncryptionKey)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt master key: %v", err)
		}

		if response.EncryptedPrivateKey != "" {
			encryptedPrivateKey, err := base64.StdEncoding.DecodeString(response.EncryptedPrivateKey)
			if err != nil {
				return nil, fmt.Errorf("failed to decode encrypted private key: %v", err)
			}

			// Decrypt the private key using the master key
			privateKey, err := decryptData(encryptedPrivateKey, masterKey)
			if err != nil {
				return nil, fmt.Errorf("failed to decrypt private key: %v", err)
			}

			// Store the keys for subsequent operations
			c.Keys = &KeySet{
				MasterKey:  masterKey,
				PrivateKey: privateKey,
			}

			// If public key is provided, store it too
			if response.PublicKey != "" {
				publicKey, err := base64.StdEncoding.DecodeString(response.PublicKey)
				if err != nil {
					return nil, fmt.Errorf("failed to decode public key: %v", err)
				}
				c.Keys.PublicKey = publicKey
			}
		}
	}

	return &response, nil
}

// sendLoginRequest sends the login payload to the server
func sendLoginRequest(config ClientConfig, payload *LoginPayload) ([]byte, error) {
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
	endpoint := fmt.Sprintf("%s/iam/api/v1/login", serverURL)

	// Convert payload to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal login data: %v", err)
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
		return nil, fmt.Errorf("failed to send login request: %v", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// Check response status
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("login failed with status %d: %s",
			resp.StatusCode, body)
	}

	return body, nil
}
