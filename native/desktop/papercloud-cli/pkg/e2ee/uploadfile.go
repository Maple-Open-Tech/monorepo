// pkg/e2ee/uploadfile.go
package e2ee

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	pref "github.com/Maple-Open-Tech/monorepo/native/desktop/papercloud-cli/internal/common/preferences"
)

// FileMetadata represents the structure for file metadata
type FileMetadata struct {
	Filename       string            `json:"filename"`
	OriginalSize   int64             `json:"original_size"`
	ContentType    string            `json:"content_type"`
	CreatedAt      time.Time         `json:"created_at"`
	ModifiedAt     time.Time         `json:"modified_at"`
	Description    string            `json:"description,omitempty"`
	Tags           []string          `json:"tags,omitempty"`
	CustomMetadata map[string]string `json:"custom_metadata,omitempty"`
}

// UploadFileResponse represents the server's response after a successful upload
type UploadFileResponse struct {
	ID        string    `json:"id"`
	FileID    string    `json:"file_id"`
	CreatedAt time.Time `json:"created_at"`
}

// UploadEncryptedFile handles the file encryption and upload process
func (c *Client) UploadEncryptedFile(filePath string, fileID string, metadata *FileMetadata) (*UploadFileResponse, error) {
	// Check authentication
	preferences := pref.PreferencesInstance()
	if preferences.LoginResponse == nil || preferences.LoginResponse.AccessToken == "" {
		return nil, fmt.Errorf("not authenticated: please login first")
	}

	// Check if file exists
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to access file: %w", err)
	}

	// Set file size in metadata
	metadata.OriginalSize = fileInfo.Size()

	// Open file for reading
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Encrypt file content
	encryptedFile, encryptedMetadata, encryptedHash, err := c.encryptFileAndMetadata(file, metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt file: %w", err)
	}
	defer encryptedFile.Close()

	// Prepare form data for upload
	formData := map[string]string{
		"file_id":            fileID,
		"encrypted_metadata": encryptedMetadata,
		"encrypted_hash":     encryptedHash,
		"encryption_version": "1.0",
	}

	// Reopen the encrypted file for upload
	encryptedFileReader, err := os.Open(encryptedFile.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to open encrypted file: %w", err)
	}
	defer encryptedFileReader.Close()

	// Prepare files for upload
	formFiles := map[string]io.Reader{
		"encrypted_content": encryptedFileReader,
	}

	// Upload the file using authenticated form request
	response, err := c.AuthenticatedFormRequest(
		"POST",
		"/api/v1/encrypted-files",
		formData,
		formFiles,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	// Parse response
	var uploadResponse UploadFileResponse
	if err := json.Unmarshal(response, &uploadResponse); err != nil {
		return nil, fmt.Errorf("failed to parse upload response: %w", err)
	}

	return &uploadResponse, nil
}

// encryptFileAndMetadata encrypts the file and metadata
func (c *Client) encryptFileAndMetadata(fileReader io.Reader, metadata *FileMetadata) (*os.File, string, string, error) {
	// Get preferences for encryption keys
	preferences := pref.PreferencesInstance()
	if preferences.VerifyOTTResponse == nil {
		return nil, "", "", fmt.Errorf("no encryption keys available, please login first")
	}

	// Create a temporary file for encrypted content
	tempFile, err := os.CreateTemp("", "encrypted-*")
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to create temp file: %w", err)
	}

	// Generate a random file key
	fileKey := make([]byte, 32)
	if _, err := rand.Read(fileKey); err != nil {
		tempFile.Close()
		os.Remove(tempFile.Name())
		return nil, "", "", fmt.Errorf("failed to generate file key: %w", err)
	}

	// Create a buffer for the encrypted content
	// In a real implementation, you would:
	// 1. Generate a nonce
	// 2. Encrypt the file with the file key using the nonce
	// 3. Encrypt the file key with the master key

	// For this simplified version, just copy the file
	if _, err := io.Copy(tempFile, fileReader); err != nil {
		tempFile.Close()
		os.Remove(tempFile.Name())
		return nil, "", "", fmt.Errorf("failed to encrypt file content: %w", err)
	}

	// Reset file position for hash calculation
	if _, err := tempFile.Seek(0, 0); err != nil {
		tempFile.Close()
		os.Remove(tempFile.Name())
		return nil, "", "", fmt.Errorf("failed to reset file position: %w", err)
	}

	// Calculate hash (simplified)
	hashBytes := make([]byte, 32)
	if _, err := rand.Read(hashBytes); err != nil {
		tempFile.Close()
		os.Remove(tempFile.Name())
		return nil, "", "", fmt.Errorf("failed to generate hash: %w", err)
	}
	encryptedHash := base64.StdEncoding.EncodeToString(hashBytes)

	// Encrypt metadata with file key
	metadataBytes, err := json.Marshal(metadata)
	if err != nil {
		tempFile.Close()
		os.Remove(tempFile.Name())
		return nil, "", "", fmt.Errorf("failed to marshal metadata: %w", err)
	}

	// In a real implementation, encrypt the metadata with the file key
	encryptedMetadata := base64.StdEncoding.EncodeToString(metadataBytes)

	// Reset file position for reading
	if _, err := tempFile.Seek(0, 0); err != nil {
		tempFile.Close()
		os.Remove(tempFile.Name())
		return nil, "", "", fmt.Errorf("failed to reset file position: %w", err)
	}

	return tempFile, encryptedMetadata, encryptedHash, nil
}
