package e2ee

import (
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

// Client is the main struct for interacting with the E2EE system
type Client struct {
	Config                                  ClientConfig
	Keys                                    *KeySet
	Salt                                    string
	StoredPublicKey                         string
	StoredEncryptedMasterKey                string
	StoredEncryptedPrivateKey               string
	StoredEncryptedRecoveryKey              string
	StoredMasterKeyEncryptedWithRecoveryKey string
	StoredVerificationID                    string
}

// NewClient creates a new E2EE client with the provided configuration
func NewClient(config ClientConfig) *Client {
	if config.ServerURL == "" {
		config.ServerURL = DefaultServerURL
	}

	// Note: The new fields (Salt, Stored*, etc.) will have their zero values initially.
	// They are expected to be populated later, e.g., during registration or login.
	return &Client{
		Config: config,
	}
}
