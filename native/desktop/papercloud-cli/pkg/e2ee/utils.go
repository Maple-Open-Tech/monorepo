package e2ee

import (
	"encoding/base64"
	"fmt" // Changed from io/ioutil
	// Added for potential future use if needed, matching ioutil deprecation recommendation
	"strings"

	"golang.org/x/crypto/nacl/box"
)

func censorEmail(email string) string {
	atIndex := strings.Index(email, "@")
	if atIndex <= 0 { // Also handles cases like "@domain.com" or no "@"
		return "***" // Or return email if it's definitely not an email format? "***" is safer.
	}

	localPart := email[:atIndex]
	domainPart := email[atIndex+1:] // Skip the '@'

	prefixLen := 3 // Number of characters to keep at the start
	if len(localPart) <= prefixLen {
		// If the local part is short or equal to the prefix length, show it all
		return localPart + "***@" + domainPart
	}

	// Otherwise, show the prefix and hide the rest
	return localPart[:prefixLen] + "***@" + domainPart
}

// Add this helper function to the client code:
func decryptChallengeWithPrivateKey(encryptedChallengeBase64 string, privateKey []byte) ([]byte, error) {
	// Decode the base64 challenge
	encryptedData, err := base64.StdEncoding.DecodeString(encryptedChallengeBase64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 challenge: %w", err)
	}

	// We need at least 56 bytes (32 for ephemeral pubkey + 24 for nonce)
	if len(encryptedData) < 56 {
		return nil, fmt.Errorf("encrypted data too short, expected at least 56 bytes, got %d", len(encryptedData))
	}

	// Extract the ephemeral public key
	var ephemeralPub [32]byte
	copy(ephemeralPub[:], encryptedData[:32])

	// Extract the nonce
	var nonce [24]byte
	copy(nonce[:], encryptedData[32:56])

	// Prepare private key in correct format
	var privKey [32]byte
	copy(privKey[:], privateKey)

	// Decrypt the message
	decrypted, ok := box.Open(nil, encryptedData[56:], &nonce, &ephemeralPub, &privKey)
	if !ok {
		return nil, fmt.Errorf("failed to decrypt challenge")
	}

	return decrypted, nil
}

// Helper function to find minimum of two ints
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
