package e2ee

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"       // Changed from io/ioutil
	"net/http" // Added for potential future use if needed, matching ioutil deprecation recommendation
	"time"
)

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

	MasterKeyEncrypted   string // e.g., the encrypted master key
	MasterKeySalt        string // salt used for key derivation
	RecoveryKeyEncrypted string // encrypted recovery key
	PublicKey            string // user's public key
	PrivateKeyEncrypted  string // encrypted private key
}

// VerifyPasswordAndCompleteLogin verifies password locally and completes the login
func (c *Client) VerifyPasswordAndCompleteLogin(email, password string, ottResponse *VerifyOTTResponse) (*LoginResponse, error) {

	fmt.Print("VerifyPasswordAndCompleteLogin is starting...")

	// Create a censored version of the email for logging
	censoredEmail := censorEmail(email)

	// Step 1: Decode the received data
	salt, err := base64.StdEncoding.DecodeString(ottResponse.Salt)
	if err != nil {
		// Added censored email for privacy and context about the source data (first few chars might be helpful but potentially noisy)
		return nil, fmt.Errorf("VerifyPasswordAndCompleteLogin: failed to decode base64 salt (len %d) for email %s: %w", len(ottResponse.Salt), censoredEmail, err)
	}

	encryptedMasterKey, err := base64.StdEncoding.DecodeString(ottResponse.EncryptedMasterKey)
	if err != nil {
		return nil, fmt.Errorf("VerifyPasswordAndCompleteLogin: failed to decode base64 encrypted master key (len %d) for email %s: %w", len(ottResponse.EncryptedMasterKey), censoredEmail, err)
	}

	encryptedPrivateKey, err := base64.StdEncoding.DecodeString(ottResponse.EncryptedPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("VerifyPasswordAndCompleteLogin: failed to decode base64 encrypted private key (len %d) for email %s: %w", len(ottResponse.EncryptedPrivateKey), censoredEmail, err)
	}

	// Step 2: Derive the key encryption key from the password
	// Be careful about logging anything related to the password itself.
	keyEncryptionKey, err := deriveKeyFromPassword(password, salt)
	if err != nil {
		// Added censored email
		return nil, fmt.Errorf("VerifyPasswordAndCompleteLogin: failed to derive key encryption key from password (salt len %d) for email %s: %w", len(salt), censoredEmail, err)
	}

	// Step 3: Attempt to decrypt the master key (this verifies the password)
	masterKey, err := decryptData(encryptedMasterKey, keyEncryptionKey)
	if err != nil {
		// This error specifically indicates a likely password mismatch.
		// Added censored email and length of key tried to decrypt
		return nil, fmt.Errorf("VerifyPasswordAndCompleteLogin: failed to decrypt master key (len %d), likely incorrect password for email %s: %w", len(encryptedMasterKey), censoredEmail, err)
	}

	// Step 4: Decrypt the private key using the master key
	privateKey, err := decryptData(encryptedPrivateKey, masterKey)
	if err != nil {
		// Added censored email and length of key tried to decrypt
		return nil, fmt.Errorf("VerifyPasswordAndCompleteLogin: failed to decrypt private key (len %d) using master key for email %s: %w", len(encryptedPrivateKey), censoredEmail, err)
	}

	// Step 5: Store the keys for subsequent operations
	c.Keys = &KeySet{
		MasterKey:  masterKey,
		PrivateKey: privateKey,
	}
	fmt.Printf("VerifyPasswordAndCompleteLogin: Successfully decrypted and stored keys for email %s\n", censoredEmail) // Optional success log with censored email
	fmt.Println("VerifyPasswordAndCompleteLogin is starting to decrypt the encrypted challenge...")

	// Step 6: Decrypt the challenge
	fmt.Println("Starting challenge decryption with these values:")
	fmt.Printf("- Challenge ID: %s\n", ottResponse.ChallengeID)
	fmt.Printf("- EncryptedChallenge (base64) length: %d\n", len(ottResponse.EncryptedChallenge))
	fmt.Printf("- Private key length: %d\n", len(privateKey))

	// Use the original base64 string, not the decoded bytes
	decryptedChallenge, err := decryptChallengeWithPrivateKey(ottResponse.EncryptedChallenge, privateKey)
	if err != nil {
		fmt.Printf("DECRYPTION ERROR: %v\n", err)
		// More detailed debugging info
		fmt.Printf("- Original challenge string (first 20 chars): %s\n",
			ottResponse.EncryptedChallenge[:min(20, len(ottResponse.EncryptedChallenge))])
		return nil, fmt.Errorf("VerifyPasswordAndCompleteLogin: failed to decrypt server challenge (base64 len %d, challengeId %s) using private key for email %s: %w",
			len(ottResponse.EncryptedChallenge), ottResponse.ChallengeID, censoredEmail, err)
	}

	fmt.Printf("Successfully decrypted challenge (length: %d bytes)\n", len(decryptedChallenge))

	// Step 7: Create the complete login request with the decrypted challenge
	completeLoginPayload := &CompleteLoginRequest{
		Email:         email, // Use original email for the request payload
		ChallengeID:   ottResponse.ChallengeID,
		DecryptedData: base64.StdEncoding.EncodeToString(decryptedChallenge),
	}

	// Call completeLogin and wrap potential errors with context
	loginResponse, err := c.completeLogin(completeLoginPayload)
	if err != nil {
		// Added censored email and challenge ID for context
		return nil, fmt.Errorf("VerifyPasswordAndCompleteLogin: failed during final login completion step (challengeId %s) for email %s: %w", ottResponse.ChallengeID, censoredEmail, err)
	}

	return loginResponse, nil
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
		return nil, fmt.Errorf("completeLogin: failed to marshal CompleteLoginRequest payload for email %s (challengeId %s): %w", censorEmail(payload.Email), payload.ChallengeID, err)
	}

	// Create request
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("completeLogin: failed to create POST request for %s: %w", endpoint, err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("completeLogin: failed to send request to %s: %w", endpoint, err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body) // Changed from ioutil.ReadAll
	if err != nil {
		// Log the status code even if reading the body fails
		return nil, fmt.Errorf("completeLogin: failed to read response body from %s (status %d): %w", endpoint, resp.StatusCode, err)
	}

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("completeLogin: request to %s failed with status %d: %s",
			endpoint, resp.StatusCode, string(body))
	}

	// Parse the response
	var response LoginResponse
	if err := json.Unmarshal(body, &response); err != nil {
		// Log the raw body for debugging if unmarshalling fails
		return nil, fmt.Errorf("completeLogin: failed to parse LoginResponse JSON from %s: %w. Raw body: %s", endpoint, err, string(body))
	}

	// Save tokens to local storage
	if err := SaveTokens(
		payload.Email,
		response.AccessToken,
		response.RefreshToken,
		response.AccessTokenExpiryTime,
	); err != nil {
		// Log the error but potentially return the successful response anyway,
		// as the login itself succeeded, only saving tokens failed.
		// Or return the error if saving tokens is critical.
		// Current implementation returns the error.
		return nil, fmt.Errorf("completeLogin: login successful for %s, but failed to save tokens: %w", payload.Email, err)
	}

	// fmt.Printf("completeLogin: Login successful and tokens saved for email %s\n", payload.Email) // Optional success log
	return &response, nil
}

// Login provides a simplified interface for the multi-step login process
func (c *Client) Login(email, password string) (*LoginResponse, error) {
	// Step 1: Request a one-time token
	if err := c.RequestLoginOTT(email); err != nil {
		// Wrap error with context of the overall Login operation
		return nil, fmt.Errorf("Login: step 1 (RequestLoginOTT) failed for email %s: %w", email, err)
	}

	// In a CLI application, prompt the user for the OTT
	fmt.Println("Please check your email for a one-time token and enter it when prompted.")
	var ott string
	fmt.Print("Enter OTT: ")
	_, err := fmt.Scanln(&ott) // Check Scanln error
	if err != nil {
		return nil, fmt.Errorf("Login: failed to read OTT input: %w", err)
	}
	if ott == "" {
		return nil, fmt.Errorf("Login: OTT input cannot be empty")
	}

	// Step 2: Verify the OTT and get the encrypted keys and challenge
	ottResponse, err := c.VerifyLoginOTT(email, ott)
	if err != nil {
		// Wrap error with context of the overall Login operation
		return nil, fmt.Errorf("Login: step 2 (VerifyLoginOTT) failed for email %s: %w", email, err)
	}

	// Step 3: Verify password locally and complete the login
	loginResponse, err := c.VerifyPasswordAndCompleteLogin(email, password, ottResponse)
	if err != nil {
		// Wrap error with context of the overall Login operation
		return nil, fmt.Errorf("Login: step 3 (VerifyPasswordAndCompleteLogin) failed for email %s: %w", censorEmail(email), err)
	}

	return loginResponse, nil
}
