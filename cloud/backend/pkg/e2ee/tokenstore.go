package e2ee

import (
	"encoding/json"
	"fmt" // Keep for now, will replace specific calls
	"log" // Import log package for potential future structured logging
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
	jsonData, err := json.MarshalIndent(store, "", "  ") // Use MarshalIndent for readability
	if err != nil {
		// Log the error with more context if needed, but return a user-friendly error
		log.Printf("Error marshalling token data for email %s: %v", email, err) // Example verbose log
		return fmt.Errorf("internal error: failed to prepare token data: %w", err)
	}

	// Get user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Printf("Error getting user home directory: %v", err)
		return fmt.Errorf("failed to determine home directory: %w", err)
	}

	// Define config directory and token path
	configDir := filepath.Join(homeDir, ".maple")
	tokenPath := filepath.Join(configDir, "auth.json")

	// Create .maple directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0700); err != nil {
		log.Printf("Error creating config directory '%s': %v", configDir, err)
		return fmt.Errorf("failed to create config directory '%s': %w", configDir, err)
	}

	// Write tokens to file using os.WriteFile (replaces deprecated ioutil.WriteFile)
	if err := os.WriteFile(tokenPath, jsonData, 0600); err != nil {
		log.Printf("Error writing token file to '%s': %v", tokenPath, err)
		return fmt.Errorf("failed to write token file '%s': %w", tokenPath, err)
	}

	log.Printf("Tokens successfully saved for email %s to %s", email, tokenPath)
	return nil
}

// LoadTokens loads authentication tokens from a local file
func LoadTokens() (*TokenStore, error) {
	// Get user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Printf("Error getting user home directory: %v", err)
		return nil, fmt.Errorf("failed to determine home directory: %w", err)
	}

	// Define token path
	tokenPath := filepath.Join(homeDir, ".maple", "auth.json")

	// Check if token file exists
	if _, err := os.Stat(tokenPath); err != nil {
		if os.IsNotExist(err) {
			log.Printf("Token file not found at '%s'. User likely not logged in.", tokenPath)
			return nil, fmt.Errorf("not logged in: token file not found at '%s'", tokenPath) // More specific message
		}
		// Handle other potential errors from os.Stat (e.g., permission denied)
		log.Printf("Error checking token file status '%s': %v", tokenPath, err)
		return nil, fmt.Errorf("failed to access token file '%s': %w", tokenPath, err)
	}

	// Read token file using os.ReadFile (replaces deprecated ioutil.ReadFile)
	jsonData, err := os.ReadFile(tokenPath)
	if err != nil {
		log.Printf("Error reading token file '%s': %v", tokenPath, err)
		return nil, fmt.Errorf("failed to read token file '%s': %w", tokenPath, err)
	}

	// Parse JSON
	var store TokenStore
	if err := json.Unmarshal(jsonData, &store); err != nil {
		log.Printf("Error parsing token data from '%s': %v", tokenPath, err)
		// Consider logging snippet of invalid JSON if safe, but be careful with tokens
		return nil, fmt.Errorf("failed to parse token data from '%s': %w", tokenPath, err)
	}

	log.Printf("Tokens successfully loaded from %s for email %s", tokenPath, store.Email)
	return &store, nil
}

// ClearTokens removes the stored tokens
func ClearTokens() error {
	// Get user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Printf("Error getting user home directory: %v", err)
		return fmt.Errorf("failed to determine home directory: %w", err)
	}

	// Define token path
	tokenPath := filepath.Join(homeDir, ".maple", "auth.json")

	// Attempt to remove token file
	err = os.Remove(tokenPath)
	if err != nil && !os.IsNotExist(err) {
		// Log error only if it's not "file not found"
		log.Printf("Error removing token file '%s': %v", tokenPath, err)
		return fmt.Errorf("failed to remove token file '%s': %w", tokenPath, err)
	}

	if err == nil {
		log.Printf("Successfully removed token file '%s'", tokenPath)
	} else {
		// Log if the file didn't exist anyway
		log.Printf("Token file '%s' did not exist or was already removed.", tokenPath)
	}

	return nil
}

// IsTokenExpired checks if the access token is expired
func IsTokenExpired(store *TokenStore) bool {
	if store == nil {
		log.Println("Warning: IsTokenExpired called with nil TokenStore")
		return true // Treat nil store as expired/invalid
	}
	// Add a small buffer (e.g., 30 seconds) to account for clock skew and request time
	// Consider making the buffer duration configurable
	buffer := 30 * time.Second
	expiryThreshold := store.ExpiresAt.Add(-buffer)
	isExpired := expiryThreshold.Before(time.Now())

	if isExpired {
		log.Printf("Token for %s is expired or expiring soon (expires: %s, threshold: %s, now: %s)",
			store.Email, store.ExpiresAt.Format(time.RFC3339), expiryThreshold.Format(time.RFC3339), time.Now().Format(time.RFC3339))
	}

	return isExpired
}
