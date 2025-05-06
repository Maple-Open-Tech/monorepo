package e2ee

import (
	"encoding/base64"
	"fmt"
)

// Client is the main struct for interacting with the E2EE system
type Client struct {
	Config ClientConfig
	Keys   *KeySet
}

// NewClient creates a new E2EE client with the provided configuration
func NewClient(config ClientConfig) *Client {
	if config.ServerURL == "" {
		config.ServerURL = DefaultServerURL
	}

	return &Client{
		Config: config,
	}
}

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
	keyEncryptionKey := deriveKeyFromPassword(password, salt)

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
// 4. Returns the recovery key for the user to save
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

	// Store the keys in the client
	c.Keys = keySet

	// Send the registration request
	_, err = sendRegistrationRequest(c.Config, payload) // Assuming sendRegistrationRequest takes *RegistrationPayload
	if err != nil {
		return "", fmt.Errorf("failed to send registration: %v", err)
	}

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
