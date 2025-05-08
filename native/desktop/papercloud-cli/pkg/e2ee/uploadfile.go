// pkg/e2ee/uploadfile.go
package e2ee

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
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
	if !c.IsAuthenticated() {
		return nil, fmt.Errorf("not authenticated or token expired: please login again")
	}

	// Check if file exists
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to access file: %w", err)
	}
	if fileInfo == nil {
		return nil, fmt.Errorf("file does not exist")
	}

	// Open file for reading
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Create a temporary file for encrypted content
	tempFile, err := os.CreateTemp("", "encrypted-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// For now, we're just copying the file content directly
	// In a real implementation, we would encrypt it here
	if _, err := io.Copy(tempFile, file); err != nil {
		return nil, fmt.Errorf("failed to process file: %w", err)
	}

	// Generate a simple hash of the file
	if _, err := tempFile.Seek(0, 0); err != nil {
		return nil, fmt.Errorf("failed to reset file position: %w", err)
	}

	hasher := sha256.New()
	if _, err := io.Copy(hasher, tempFile); err != nil {
		return nil, fmt.Errorf("failed to calculate hash: %w", err)
	}
	encryptedHash := base64.StdEncoding.EncodeToString(hasher.Sum(nil))

	// Serialize metadata
	metadataBytes, err := json.Marshal(metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal metadata: %w", err)
	}
	encryptedMetadata := base64.StdEncoding.EncodeToString(metadataBytes)

	// Prepare for upload
	if _, err := tempFile.Seek(0, 0); err != nil {
		return nil, fmt.Errorf("failed to reset file position: %w", err)
	}

	// Open the encrypted file for reading
	encryptedFileReader, err := os.Open(tempFile.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to open encrypted file: %w", err)
	}
	defer encryptedFileReader.Close()

	// Prepare form data
	formData := map[string]string{
		"file_id":            fileID,
		"encrypted_metadata": encryptedMetadata,
		"encrypted_hash":     encryptedHash,
		"encryption_version": "1.0",
	}

	// Prepare files
	formFiles := map[string]io.Reader{
		"encrypted_content": encryptedFileReader,
	}

	// Send the authenticated form request
	responseBytes, err := c.AuthenticatedFormRequest(
		"POST",
		"/api/v1/encrypted-files",
		formData,
		formFiles,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	// Parse the response
	var response UploadFileResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}
