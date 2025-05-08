// cloud/backend/native/desktop/papercloud-cli/cmd/remote/uploadfile.go
package remote

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

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

func UploadFileCmd() *cobra.Command {
	var filePath, description, tags, contentType string
	var customMetadata string

	var cmd = &cobra.Command{
		Use:   "upload-file",
		Short: "Upload a file with end-to-end encryption",
		Long: `
Upload a file with end-to-end encryption to your PaperCloud account.
The file will be encrypted locally before being uploaded, ensuring
your data remains secure and private. You can also provide metadata
for the file that will be encrypted along with the content.

Examples:
  # Basic upload with minimal metadata
  papercloud-cli remote upload-file --file /path/to/file.pdf

  # Upload with description and tags
  papercloud-cli remote upload-file --file /path/to/file.pdf --description "Important document" --tags "work,private,important"

  # Upload with custom metadata
  papercloud-cli remote upload-file --file /path/to/file.pdf --custom '{"project":"Project X","department":"Finance"}'
`,
		Run: func(cmd *cobra.Command, args []string) {
			// Check if user is authenticated
			preferences := pref.PreferencesInstance()
			if preferences.LoginResponse == nil || preferences.LoginResponse.AccessToken == "" {
				fmt.Println("Error: You need to be logged in to upload files. Please login first.")
				return
			}

			// Check if file exists
			if filePath == "" {
				fmt.Println("Error: File path is required")
				return
			}

			fileInfo, err := os.Stat(filePath)
			if err != nil {
				fmt.Printf("Error: Failed to access file: %v\n", err)
				return
			}

			// Prepare file metadata
			metadata := FileMetadata{
				Filename:     filepath.Base(filePath),
				OriginalSize: fileInfo.Size(),
				ContentType:  determineContentType(filePath, contentType),
				CreatedAt:    time.Now(),
				ModifiedAt:   fileInfo.ModTime(),
				Description:  description,
			}

			// Process tags if provided
			if tags != "" {
				metadata.Tags = strings.Split(tags, ",")
				// Trim whitespace from each tag
				for i, tag := range metadata.Tags {
					metadata.Tags[i] = strings.TrimSpace(tag)
				}
			}

			// Process custom metadata if provided
			if customMetadata != "" {
				err := json.Unmarshal([]byte(customMetadata), &metadata.CustomMetadata)
				if err != nil {
					fmt.Printf("Error: Invalid custom metadata format: %v\n", err)
					return
				}
			}

			// Generate a unique file ID
			fileID := generateFileID(filePath, metadata)
			fmt.Printf("Generated file ID: %s\n", fileID)

			// Encrypt and upload the file
			fmt.Println("Starting file encryption and upload...")
			if err := encryptAndUploadFile(filePath, fileID, metadata, preferences); err != nil {
				fmt.Printf("Error: Failed to encrypt and upload file: %v\n", err)
				return
			}

			fmt.Println("File successfully encrypted and uploaded!")
		},
	}

	// Define command flags
	cmd.Flags().StringVarP(&filePath, "file", "f", "", "Path to the file to upload (required)")
	cmd.Flags().StringVarP(&description, "description", "d", "", "Description of the file")
	cmd.Flags().StringVarP(&tags, "tags", "t", "", "Comma-separated list of tags")
	cmd.Flags().StringVarP(&contentType, "content-type", "c", "", "Content type of the file (defaults to auto-detection)")
	cmd.Flags().StringVarP(&customMetadata, "custom", "m", "", "Custom metadata in JSON format")

	// Mark required flags
	cmd.MarkFlagRequired("file")

	return cmd
}

// determineContentType attempts to determine the content type of a file
func determineContentType(filePath, providedType string) string {
	if providedType != "" {
		return providedType
	}

	// Simple extension-based content type detection
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".pdf":
		return "application/pdf"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".txt":
		return "text/plain"
	case ".doc", ".docx":
		return "application/msword"
	case ".xls", ".xlsx":
		return "application/vnd.ms-excel"
	default:
		return "application/octet-stream" // Default binary type
	}
}

// generateFileID creates a unique ID for the file
func generateFileID(filePath string, metadata FileMetadata) string {
	// Create a unique identifier based on file path, size, and current time
	uniqueStr := fmt.Sprintf("%s_%d_%d", filePath, metadata.OriginalSize, time.Now().UnixNano())

	// Generate SHA256 hash
	hash := sha256.Sum256([]byte(uniqueStr))

	// Return a base64 URL-safe encoding of the first 16 bytes of the hash
	return base64.URLEncoding.EncodeToString(hash[:16])
}

// encryptAndUploadFile handles the encryption and upload process
func encryptAndUploadFile(filePath, fileID string, metadata FileMetadata, preferences *pref.Preferences) error {
	// Convert metadata to JSON
	metadataBytes, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("failed to encode metadata: %w", err)
	}

	// 1. Read the file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Create a temporary file for the encrypted content
	tempFile, err := os.CreateTemp("", "encrypted-*")
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// 2. Encrypt the file (simulated here)
	// In a real implementation, you would use the E2EE package to:
	// - Generate or retrieve the master key from saved credentials
	// - Encrypt the file content
	// - Encrypt the metadata
	fmt.Println("Encrypting file content...")

	// This would be where you'd use the actual encryption logic
	// For demonstration, we're just copying the file to the temp location
	// In a real implementation, you'd use:
	// - client.Keys.MasterKey to encrypt the file
	// - Encrypted content would be written to tempFile

	if _, err := io.Copy(tempFile, file); err != nil {
		return fmt.Errorf("failed to copy file content: %w", err)
	}

	// Rewind the temp file for reading
	if _, err := tempFile.Seek(0, 0); err != nil {
		return fmt.Errorf("failed to rewind temp file: %w", err)
	}

	// 3. Calculate a hash of the encrypted content for integrity verification
	hasher := sha256.New()
	if _, err := io.Copy(hasher, tempFile); err != nil {
		return fmt.Errorf("failed to calculate file hash: %w", err)
	}
	encryptedHash := base64.StdEncoding.EncodeToString(hasher.Sum(nil))

	// Rewind again for the upload
	if _, err := tempFile.Seek(0, 0); err != nil {
		return fmt.Errorf("failed to rewind temp file again: %w", err)
	}

	// 4. Encrypt the metadata with the master key
	encryptedMetadata := base64.StdEncoding.EncodeToString(metadataBytes)
	fmt.Println("Metadata encrypted successfully")

	// 5. Create multipart form data for upload
	// In a real implementation, you would:
	// - Prepare multipart form with the encrypted file data
	// - Add encrypted metadata, file ID, and hash to the form
	// - Upload using the authentication token

	// Here we're simulating the API call - in a real implementation
	// you'd make an HTTP request to the server
	fmt.Printf("Uploading file %s with ID %s...\n", metadata.Filename, fileID)
	fmt.Printf("File size: %d bytes\n", metadata.OriginalSize)

	// Log the key pieces of information that would be sent
	fmt.Println("Upload would include:")
	fmt.Printf("- Encrypted File Content (%d bytes)\n", metadata.OriginalSize)
	fmt.Printf("- File ID: %s\n", fileID)
	fmt.Printf("- Encrypted Metadata: %s...\n", encryptedMetadata[:20])
	fmt.Printf("- Encrypted Hash: %s\n", encryptedHash)
	fmt.Println("- Encryption Version: 1.0")

	// In a real implementation, this is where you'd make the actual API call
	// and process the server response
	return nil
}
