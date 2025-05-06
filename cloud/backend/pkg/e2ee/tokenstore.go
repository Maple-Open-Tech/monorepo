package e2ee

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

// TokenStore represents stored authentication tokens
type TokenStore struct {
	Email        string    `json:"email"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// SaveTokens saves authentication tokens to a local file
func SaveTokens(email, accessToken, refreshToken string, expiresAt time.Time) error {
	// Create token store
	store := &TokenStore{
		Email:        email,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}

	// Convert to JSON
	jsonData, err := json.Marshal(store)
	if err != nil {
		return fmt.Errorf("failed to marshal token data: %v", err)
	}

	// Get user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %v", err)
	}

	// Create .maple directory if it doesn't exist
	configDir := filepath.Join(homeDir, ".maple")
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}

	// Write tokens to file
	tokenPath := filepath.Join(configDir, "auth.json")
	if err := ioutil.WriteFile(tokenPath, jsonData, 0600); err != nil {
		return fmt.Errorf("failed to write token file: %v", err)
	}

	return nil
}

// LoadTokens loads authentication tokens from a local file
func LoadTokens() (*TokenStore, error) {
	// Get user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %v", err)
	}

	// Check if token file exists
	tokenPath := filepath.Join(homeDir, ".maple", "auth.json")
	if _, err := os.Stat(tokenPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("not logged in: no token file found")
	}

	// Read token file
	jsonData, err := ioutil.ReadFile(tokenPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read token file: %v", err)
	}

	// Parse JSON
	var store TokenStore
	if err := json.Unmarshal(jsonData, &store); err != nil {
		return nil, fmt.Errorf("failed to parse token data: %v", err)
	}

	return &store, nil
}

// ClearTokens removes the stored tokens
func ClearTokens() error {
	// Get user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %v", err)
	}

	// Remove token file
	tokenPath := filepath.Join(homeDir, ".maple", "auth.json")
	if err := os.Remove(tokenPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove token file: %v", err)
	}

	return nil
}

// IsTokenExpired checks if the access token is expired
func IsTokenExpired(store *TokenStore) bool {
	// Add a small buffer to ensure we don't use tokens that are about to expire
	return store.ExpiresAt.Add(-30 * time.Second).Before(time.Now())
}
