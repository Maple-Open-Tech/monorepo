package e2ee

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io" // Use io instead of deprecated io/ioutil
	"net/http"
)

// PrepareRegistration generates all keys and prepares the registration payload
// without actually sending it to the server yet
// This allows the application to inspect or modify the payload before submission
func (c *Client) PrepareRegistration(
	// --- Authentication ---
	email, password string,
	// --- Application and PII ---
	betaAccessCode, firstName, lastName, phone, country, timezone string,
	agreeTermsOfService, agreePromotions, agreeToTracking bool,
	module int,
) (*KeySet, *RegistrationPayload, error) {
	// Generate all necessary keys
	masterKey, err := generateMasterKey()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate master key: %v", err)
	}

	publicKey, privateKey, err := generateKeyPair()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate keypair: %v", err)
	}

	recoveryKey, err := generateRecoveryKey()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate recovery key: %v", err)
	}

	// Create a salt for password hashing
	salt, err := generateSalt()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate salt: %v", err)
	}

	// Derive the key encryption key from the password
	keyEncryptionKey, err := deriveKeyFromPassword(password, salt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to derive key encryption key: %v", err)
	}

	// Encrypt the master key with the key encryption key
	encryptedMasterKey, err := encryptData(masterKey, keyEncryptionKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to encrypt master key: %v", err)
	}

	// Encrypt the private key with the master key
	encryptedPrivateKey, err := encryptData(privateKey, masterKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to encrypt private key: %v", err)
	}

	// Encrypt the recovery key with the master key
	encryptedRecoveryKey, err := encryptData(recoveryKey, masterKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to encrypt recovery key: %v", err)
	}

	// Encrypt the master key with the recovery key
	masterKeyEncryptedWithRecoveryKey, err := encryptData(masterKey, recoveryKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to encrypt master key with recovery key: %v", err)
	}

	// Create verification ID
	verificationID, err := createVerificationID(publicKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create verification ID: %v", err)
	}

	// Store the keys for client-side use (temporarily, Register will store them permanently on the client)
	keySet := &KeySet{
		MasterKey:   masterKey,
		PublicKey:   publicKey,
		PrivateKey:  privateKey,
		RecoveryKey: recoveryKey,
	}

	// Create the registration payload
	payload := &RegistrationPayload{
		// --- E2EE Fields ---
		Salt:                              base64.StdEncoding.EncodeToString(salt),
		PublicKey:                         base64.StdEncoding.EncodeToString(publicKey),
		EncryptedMasterKey:                base64.StdEncoding.EncodeToString(encryptedMasterKey),
		EncryptedPrivateKey:               base64.StdEncoding.EncodeToString(encryptedPrivateKey),
		EncryptedRecoveryKey:              base64.StdEncoding.EncodeToString(encryptedRecoveryKey),
		MasterKeyEncryptedWithRecoveryKey: base64.StdEncoding.EncodeToString(masterKeyEncryptedWithRecoveryKey),
		VerificationID:                    verificationID,

		// --- Application and PII Fields ---
		Email:               email, // Also used for auth identification
		BetaAccessCode:      betaAccessCode,
		FirstName:           firstName,
		LastName:            lastName,
		Phone:               phone,
		Country:             country,
		Timezone:            timezone,
		AgreeTermsOfService: agreeTermsOfService,
		AgreePromotions:     agreePromotions,
		AgreeToTrackingAcrossThirdPartyAppsAndServices: agreeToTracking,
		Module: module,
	}

	return keySet, payload, nil
}

// Register performs the full registration process:
// 1. Generates all necessary keys (via PrepareRegistration)
// 2. Creates the registration payload (via PrepareRegistration)
// 3. Sends the registration request to the server
// 4. Stores relevant E2EE state on the client upon success
// 5. Returns the recovery key for the user to save
func (c *Client) Register(
	// --- Authentication ---
	email, password string,
	// --- Application and PII ---
	betaAccessCode, firstName, lastName, phone, country, timezone string,
	agreeTermsOfService, agreePromotions, agreeToTracking bool,
	module int,
) (string, error) {
	// Prepare registration data
	keySet, payload, err := c.PrepareRegistration(
		email, password,
		betaAccessCode, firstName, lastName, phone, country, timezone,
		agreeTermsOfService, agreePromotions, agreeToTracking,
		module,
	)
	if err != nil {
		return "", fmt.Errorf("failed to prepare registration: %v", err)
	}

	// Store the raw keys temporarily (will be persisted if registration succeeds)
	// These keys are needed for immediate client-side operations.
	c.Keys = keySet

	// Send the registration request
	_, err = sendRegistrationRequest(c.Config, payload)
	if err != nil {
		// If registration fails, clear the temporarily stored keys
		c.Keys = nil
		return "", fmt.Errorf("failed to send registration: %v", err)
	}

	// --- Store E2EE fields from the payload on successful registration ---
	// These fields represent the state sent to the server and may be needed
	// for future operations or client state persistence.
	// Assuming Client struct has fields like:
	// Salt string // Base64 encoded salt used for KDF
	// StoredPublicKey string // Base64 encoded public key
	// StoredEncryptedMasterKey string // Base64 encoded master key encrypted with KEK
	// StoredEncryptedPrivateKey string // Base64 encoded private key encrypted with master key
	// StoredEncryptedRecoveryKey string // Base64 encoded recovery key encrypted with master key
	// StoredMasterKeyEncryptedWithRecoveryKey string // Base64 encoded master key encrypted with recovery key
	// StoredVerificationID string // Verification ID derived from public key

	c.Salt = payload.Salt
	c.StoredPublicKey = payload.PublicKey
	c.StoredEncryptedMasterKey = payload.EncryptedMasterKey
	c.StoredEncryptedPrivateKey = payload.EncryptedPrivateKey
	c.StoredEncryptedRecoveryKey = payload.EncryptedRecoveryKey
	c.StoredMasterKeyEncryptedWithRecoveryKey = payload.MasterKeyEncryptedWithRecoveryKey
	c.StoredVerificationID = payload.VerificationID

	// Optionally, persist the entire client state (including c.Keys and the stored fields above)
	// if err := c.SaveState(); err != nil {
	//     // Log or handle state saving error. Registration succeeded server-side.
	//     // Depending on requirements, this might warrant specific handling.
	//     log.Printf("Warning: failed to save client state after successful registration: %v", err)
	// }
	// --- End of storing E2EE fields ---

	// Return the recovery key as a base64 string for the user to save
	recoveryKeyString := base64.StdEncoding.EncodeToString(keySet.RecoveryKey)
	return recoveryKeyString, nil
}

// GetRecoveryKeyInfo provides explanatory text about the recovery key
// This is separate so applications can customize how they present this info
func (c *Client) GetRecoveryKeyInfo(recoveryKeyString string) string {
	return `
===================================================================
								⚠️ IMPORTANT ⚠️
===================================================================
Save your recovery key in a secure location. This is shown ONLY ONCE:

` + recoveryKeyString + `

===================================================================
If you forget your password, this key is your ONLY way to recover
your account. Without it, your data will be permanently inaccessible.
Consider writing it down and storing it in a safe, or using a
password manager.
===================================================================
`
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
	body, err := io.ReadAll(resp.Body) // Use io.ReadAll instead of deprecated ioutil.ReadAll
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// Check response status
	if resp.StatusCode != http.StatusCreated {
		// Consider providing more context from the body if possible
		// e.g., attempt to unmarshal body into an error struct
		return nil, fmt.Errorf("registration failed with status %d: %s",
			resp.StatusCode, string(body)) // Convert body to string for readability
	}

	return body, nil
}
