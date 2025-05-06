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
