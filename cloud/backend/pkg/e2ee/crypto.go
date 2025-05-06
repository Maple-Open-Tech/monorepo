package e2ee

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/nacl/box"
	"golang.org/x/crypto/nacl/secretbox"
)

// generateRandomBytes creates a slice of random bytes of the specified length
func generateRandomBytes(length int) ([]byte, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to generate random bytes: %v", err)
	}
	return bytes, nil
}

// generateMasterKey creates a new random master key
// WHY: The master key is the root of our encryption hierarchy. It's randomly
// generated (not derived from a password) for maximum security.
func generateMasterKey() ([]byte, error) {
	return generateRandomBytes(32) // 256-bit key
}

// generateKeyPair creates a new public/private key pair
// WHY: These keys enable asymmetric encryption - others can use your public key
// to encrypt data that only you can decrypt with your private key.
func generateKeyPair() (publicKey, privateKey []byte, err error) {
	pub, priv, err := box.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate keypair: %v", err)
	}
	return pub[:], priv[:], nil
}

// generateRecoveryKey creates a new random recovery key
// WHY: The recovery key provides a backup method to recover your master key
// if you forget your password.
func generateRecoveryKey() ([]byte, error) {
	return generateRandomBytes(32) // 256-bit key
}

// generateSalt creates a random salt for password hashing
// WHY: A salt ensures that the same password for different users produces
// different keys, protecting against rainbow table attacks.
func generateSalt() ([]byte, error) {
	return generateRandomBytes(16)
}

// deriveKeyFromPassword converts a user password into a cryptographic key
// WHY: This strengthens weak passwords and ensures the key is the right size
// regardless of password length.
func deriveKeyFromPassword(password string, salt []byte) []byte {
	return argon2.IDKey(
		[]byte(password),
		salt,
		3,       // Time cost (iterations)
		64*1024, // Memory cost (64 MB)
		4,       // Parallelism (4 threads)
		32,      // Key length (256 bits)
	)
}

// encryptData securely encrypts data with a key
// WHY: This handles the actual encryption process, automatically generating a
// random nonce for security.
func encryptData(data, key []byte) ([]byte, error) {
	var nonce [24]byte
	if _, err := rand.Read(nonce[:]); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %v", err)
	}

	var keyArray [32]byte
	copy(keyArray[:], key)

	// Prepend nonce to ciphertext for storage
	encrypted := secretbox.Seal(nonce[:], data, &nonce, &keyArray)
	return encrypted, nil
}

// decryptData decrypts data that was encrypted with encryptData
// WHY: This reverses the encryption process, extracting the nonce and using it
// with the provided key to decrypt the data.
func decryptData(encrypted, key []byte) ([]byte, error) {
	if len(encrypted) < 24 {
		return nil, fmt.Errorf("encrypted data too short")
	}

	var nonce [24]byte
	copy(nonce[:], encrypted[:24])

	var keyArray [32]byte
	copy(keyArray[:], key)

	decrypted, ok := secretbox.Open(nil, encrypted[24:], &nonce, &keyArray)
	if !ok {
		return nil, fmt.Errorf("decryption failed")
	}

	return decrypted, nil
}

// createVerificationID generates a human-readable identifier from a public key
// WHY: This creates a fingerprint that users can compare to verify identities
// during sharing, preventing man-in-the-middle attacks.
func createVerificationID(publicKey []byte) (string, error) {
	// Hash the public key
	hash := sha256.Sum256(publicKey)

	// Use a word list for better human readability
	// In production, you would use a larger word list
	words := []string{
		"able", "acid", "also", "apex", "aqua", "arch", "atom", "aunt",
		"back", "base", "bath", "bear", "bell", "best", "bird", "blue",
		"boat", "body", "bone", "book", "born", "both", "bowl", "bulk",
		"burn", "bush", "busy", "calm", "came", "camp", "card", "care",
	}

	var result []string
	// Use 4 bytes from the hash to select words
	for i := 0; i < 4; i++ {
		wordIndex := int(hash[i]) % len(words)
		result = append(result, words[wordIndex])
	}

	return strings.Join(result, "-"), nil
}
