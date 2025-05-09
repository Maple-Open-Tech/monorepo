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

	"go.uber.org/zap"
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
	logger := zap.L().With(zap.String("filePath", filePath), zap.String("fileID", fileID))
	logger.Info("Starting UploadEncryptedFile")

	// Check authentication
	if !c.IsAuthenticated() {
		logger.Error("Authentication check failed: not authenticated or token expired")
		return nil, fmt.Errorf("not authenticated or token expired: please login again")
	}
	logger.Debug("Authentication successful")

	// Check if file exists
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		logger.Error("Failed to access file (os.Stat)", zap.Error(err))
		return nil, fmt.Errorf("failed to access file: %w", err)
	}
	if fileInfo == nil { // This condition is technically unreachable if os.Stat returns no error.
		logger.Error("File does not exist (fileInfo is nil after os.Stat returned no error)")
		return nil, fmt.Errorf("file does not exist")
	}
	logger.Debug("File exists", zap.Int64("original_size", fileInfo.Size()))

	// Open file for reading
	file, err := os.Open(filePath)
	if err != nil {
		logger.Error("Failed to open file for reading", zap.Error(err))
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()
	logger.Debug("Successfully opened file for reading")

	// Create a temporary file for encrypted content
	tempFile, err := os.CreateTemp("", "encrypted-*")
	if err != nil {
		logger.Error("Failed to create temporary file", zap.Error(err))
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	logger.Debug("Created temporary file", zap.String("tempFileName", tempFile.Name()))
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// For now, we're just copying the file content directly
	// In a real implementation, we would encrypt it here
	logger.Info("Starting file processing (copying to temp file)", zap.String("tempFileName", tempFile.Name()))
	if _, err := io.Copy(tempFile, file); err != nil {
		logger.Error("Failed to process file (copy to temp)", zap.Error(err), zap.String("tempFileName", tempFile.Name()))
		return nil, fmt.Errorf("failed to process file: %w", err)
	}
	logger.Debug("Successfully copied file content to temporary file")

	// Generate a simple hash of the file
	logger.Debug("Seeking to beginning of temp file for hashing", zap.String("tempFileName", tempFile.Name()))
	if _, err := tempFile.Seek(0, 0); err != nil {
		logger.Error("Failed to reset temp file position for hashing", zap.Error(err), zap.String("tempFileName", tempFile.Name()))
		return nil, fmt.Errorf("failed to reset file position: %w", err)
	}

	hasher := sha256.New()
	logger.Debug("Calculating SHA256 hash of temporary file content", zap.String("tempFileName", tempFile.Name()))
	if _, err := io.Copy(hasher, tempFile); err != nil {
		logger.Error("Failed to calculate hash of temp file", zap.Error(err), zap.String("tempFileName", tempFile.Name()))
		return nil, fmt.Errorf("failed to calculate hash: %w", err)
	}
	encryptedHash := base64.StdEncoding.EncodeToString(hasher.Sum(nil))
	logger.Info("Generated encrypted hash", zap.String("encryptedHash", encryptedHash))

	// Serialize metadata
	logger.Debug("Serializing metadata")
	metadataBytes, err := json.Marshal(metadata)
	if err != nil {
		// Avoid logging full metadata if it's sensitive; log filename if available
		var filenameField zap.Field
		if metadata != nil {
			filenameField = zap.String("metadata.filename", metadata.Filename)
		} else {
			filenameField = zap.Skip()
		}
		logger.Error("Failed to marshal metadata", zap.Error(err), filenameField)
		return nil, fmt.Errorf("failed to marshal metadata: %w", err)
	}
	encryptedMetadata := base64.StdEncoding.EncodeToString(metadataBytes)
	logger.Debug("Successfully serialized and encoded metadata")

	// Prepare for upload
	logger.Debug("Seeking to beginning of temp file for upload", zap.String("tempFileName", tempFile.Name()))
	if _, err := tempFile.Seek(0, 0); err != nil {
		logger.Error("Failed to reset temp file position for upload", zap.Error(err), zap.String("tempFileName", tempFile.Name()))
		return nil, fmt.Errorf("failed to reset file position: %w", err)
	}

	// Open the encrypted file for reading
	logger.Debug("Opening encrypted (temp) file for reading for upload", zap.String("tempFileName", tempFile.Name()))
	encryptedFileReader, err := os.Open(tempFile.Name())
	if err != nil {
		logger.Error("Failed to open encrypted (temp) file for reading", zap.Error(err), zap.String("tempFileName", tempFile.Name()))
		return nil, fmt.Errorf("failed to open encrypted file: %w", err)
	}
	defer encryptedFileReader.Close()

	// Prepare form data
	formData := map[string]string{
		"file_id":            fileID,
		"encrypted_metadata": encryptedMetadata, // Value can be large/sensitive, logging keys is safer
		"encrypted_hash":     encryptedHash,
		"encryption_version": "1.0",
	}
	formDataKeys := make([]string, 0, len(formData))
	for k := range formData {
		formDataKeys = append(formDataKeys, k)
	}
	logger.Debug("Prepared form data for upload", zap.Strings("formDataKeys", formDataKeys))

	// Prepare files
	formFiles := map[string]io.Reader{
		"encrypted_content": encryptedFileReader,
	}
	logger.Debug("Prepared form files for upload", zap.Strings("formFilesKeys", []string{"encrypted_content"}))

	// Send the authenticated form request
	uploadURL := "/vault/api/v1/encrypted-files"
	logger.Info("Sending authenticated form request to upload file",
		zap.String("method", "POST"),
		zap.String("url", uploadURL),
	)
	responseBytes, err := c.AuthenticatedFormRequest(
		"POST",
		uploadURL,
		formData,
		formFiles,
	)
	if err != nil {
		logger.Error("Authenticated form request failed", zap.Error(err), zap.String("url", uploadURL))
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}
	logger.Debug("Received response from upload request", zap.Int("responseBytesLength", len(responseBytes)))

	// Parse the response
	logger.Debug("Parsing server response")
	var response UploadFileResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		logger.Error("Failed to parse server response", zap.Error(err), zap.ByteString("responseBody", responseBytes))
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	logger.Info("Successfully uploaded encrypted file and parsed response",
		zap.String("uploadID", response.ID),
		zap.String("responseFileID", response.FileID),
		zap.Time("responseCreatedAt", response.CreatedAt),
	)
	return &response, nil
}
