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
	if length <= 0 {
		return nil, fmt.Errorf("error in generateRandomBytes: requested byte length %d must be positive", length)
	}
	bytes := make([]byte, length)
	// rand.Read calls io.ReadFull, so it will return an error if not all bytes are read.
	// The number of bytes read (n) will be equal to length if err is nil.
	_, err := rand.Read(bytes)
	if err != nil {
		return nil, fmt.Errorf("error in generateRandomBytes: crypto/rand.Read failed while attempting to generate %d random bytes: %w", length, err)
	}
	return bytes, nil
}

// generateMasterKey creates a new random master key
// WHY: The master key is the root of our encryption hierarchy. It's randomly
// generated (not derived from a password) for maximum security.
func generateMasterKey() ([]byte, error) {
	key, err := generateRandomBytes(32) // 256-bit key
	if err != nil {
		return nil, fmt.Errorf("error in generateMasterKey: failed to generate 32-byte master key: %w", err)
	}
	return key, nil
}

// generateKeyPair creates a new public/private key pair
// WHY: These keys enable asymmetric encryption - others can use your public key
// to encrypt data that only you can decrypt with your private key.
func generateKeyPair() (publicKey, privateKey []byte, err error) {
	pub, priv, err := box.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, fmt.Errorf("error in generateKeyPair: nacl/box.GenerateKey failed to generate key pair: %w", err)
	}
	if pub == nil || priv == nil {
		// This check is defensive, as box.GenerateKey should ensure non-nil keys if err is nil.
		return nil, nil, fmt.Errorf("error in generateKeyPair: nacl/box.GenerateKey returned nil key(s) without an error, which is unexpected")
	}
	return pub[:], priv[:], nil
}

// generateRecoveryKey creates a new random recovery key
// WHY: The recovery key provides a backup method to recover your master key
// if you forget your password.
func generateRecoveryKey() ([]byte, error) {
	key, err := generateRandomBytes(32) // 256-bit key
	if err != nil {
		return nil, fmt.Errorf("error in generateRecoveryKey: failed to generate 32-byte recovery key: %w", err)
	}
	return key, nil
}

// generateSalt creates a random salt for password hashing
// WHY: A salt ensures that the same password for different users produces
// different keys, protecting against rainbow table attacks.
func generateSalt() ([]byte, error) {
	salt, err := generateRandomBytes(16) // Argon2 recommends at least 8 bytes, 16 is common.
	if err != nil {
		return nil, fmt.Errorf("error in generateSalt: failed to generate 16-byte salt: %w", err)
	}
	return salt, nil
}

// deriveKeyFromPassword converts a user password into a cryptographic key
// WHY: This strengthens weak passwords and ensures the key is the right size
// regardless of password length.
func deriveKeyFromPassword(password string, salt []byte) ([]byte, error) {
	const minSaltLength = 8 // Argon2 minimum salt length requirement.
	if salt == nil {
		return nil, fmt.Errorf("error in deriveKeyFromPassword: salt cannot be nil")
	}
	if len(salt) < minSaltLength {
		return nil, fmt.Errorf("error in deriveKeyFromPassword: salt length %d is less than the required minimum of %d bytes for Argon2", len(salt), minSaltLength)
	}

	// Hardcoded Argon2 parameters. Ensure they are valid if ever changed.
	// E.g., keyLen must be >= 1.
	const keyLength = 32
	const timeCost = 3
	const memoryCost = 64 * 1024
	const parallelism = 4

	key := argon2.IDKey(
		[]byte(password),
		salt,
		timeCost,
		memoryCost,
		parallelism,
		keyLength,
	)
	// argon2.IDKey panics on invalid parameters (e.g., nil salt, short salt, invalid time/memory/parallelism/keyLen).
	// We've checked for nil/short salt. Other parameters are hardcoded to valid values.
	return key, nil
}

// encryptData securely encrypts data with a key
// WHY: This handles the actual encryption process, automatically generating a
// random nonce for security.
func encryptData(data, key []byte) ([]byte, error) {
	const expectedKeyLength = 32
	if key == nil {
		return nil, fmt.Errorf("error in encryptData: provided key is nil")
	}
	if len(key) != expectedKeyLength {
		return nil, fmt.Errorf("error in encryptData: provided key length is %d bytes, but expected %d bytes for nacl/secretbox", len(key), expectedKeyLength)
	}

	var nonce [24]byte // nacl/secretbox uses a 24-byte nonce.
	if _, err := rand.Read(nonce[:]); err != nil {
		return nil, fmt.Errorf("error in encryptData: failed to generate 24-byte nonce using crypto/rand.Read for secretbox.Seal: %w", err)
	}

	var keyArray [expectedKeyLength]byte
	copy(keyArray[:], key) // Assumes key is validated to be expectedKeyLength bytes.

	// Prepend nonce to ciphertext for storage.
	// secretbox.Seal will panic if key or nonce have incorrect length,
	// but we ensure [32]byte for keyArray and [24]byte for nonce.
	encrypted := secretbox.Seal(nonce[:], data, &nonce, &keyArray)
	return encrypted, nil
}

// decryptData decrypts data that was encrypted with encryptData
// WHY: This reverses the encryption process, extracting the nonce and using it
// with the provided key to decrypt the data.
func decryptData(encrypted, key []byte) ([]byte, error) {
	const expectedKeyLength = 32
	const nonceSize = 24

	if key == nil {
		return nil, fmt.Errorf("error in decryptData: provided key is nil")
	}
	if len(key) != expectedKeyLength {
		return nil, fmt.Errorf("error in decryptData: provided key length is %d bytes, but expected %d bytes for nacl/secretbox", len(key), expectedKeyLength)
	}
	if encrypted == nil {
		return nil, fmt.Errorf("error in decryptData: encrypted data is nil")
	}

	if len(encrypted) < nonceSize {
		return nil, fmt.Errorf("error in decryptData: encrypted data too short (length %d bytes), requires at least %d bytes for the nonce", len(encrypted), nonceSize)
	}

	var nonce [nonceSize]byte
	copy(nonce[:], encrypted[:nonceSize])

	var keyArray [expectedKeyLength]byte
	copy(keyArray[:], key) // Assumes key is validated to be expectedKeyLength bytes.

	ciphertext := encrypted[nonceSize:]
	decrypted, ok := secretbox.Open(nil, ciphertext, &nonce, &keyArray)
	if !ok {
		// This is a generic failure; could be bad key, bad nonce, or corrupted data.
		return nil, fmt.Errorf("error in decryptData: nacl/secretbox.Open failed to decrypt data. This may be due to an incorrect key, corrupted ciphertext, or an invalid/reused nonce")
	}

	return decrypted, nil
}

// createVerificationID generates a human-readable identifier from a public key
// WHY: This creates a fingerprint that users can compare to verify identities
// during sharing, preventing man-in-the-middle attacks.
func createVerificationID(publicKey []byte) (string, error) {
	if publicKey == nil {
		return "", fmt.Errorf("error in createVerificationID: public key cannot be nil")
	}
	if len(publicKey) == 0 {
		return "", fmt.Errorf("error in createVerificationID: public key cannot be empty")
	}

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

	if len(words) == 0 {
		// This is an internal check, should not happen with a hardcoded list.
		return "", fmt.Errorf("error in createVerificationID: internal configuration error - word list for verification ID is empty")
	}

	var result []string
	const numWordsToSelect = 4 // We want a 4-word identifier.

	// Ensure the hash is long enough to pick numWordsToSelect bytes. SHA256 (32 bytes) is sufficient.
	if len(hash) < numWordsToSelect {
		// This should not happen with sha256.Sum256 and numWordsToSelect=4.
		return "", fmt.Errorf("error in createVerificationID: hash length (%d bytes) is insufficient to select %d words for the ID", len(hash), numWordsToSelect)
	}

	// Use N bytes from the hash to select words
	for i := 0; i < numWordsToSelect; i++ {
		wordIndex := int(hash[i]) % len(words) // Safe as len(words) > 0 checked above.
		result = append(result, words[wordIndex])
	}

	return strings.Join(result, "-"), nil
}
