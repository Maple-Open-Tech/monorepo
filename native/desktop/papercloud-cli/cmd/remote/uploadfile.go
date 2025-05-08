// cmd/remote/uploadfile.go
package remote

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	pref "github.com/Maple-Open-Tech/monorepo/native/desktop/papercloud-cli/internal/common/preferences"
	"github.com/Maple-Open-Tech/monorepo/native/desktop/papercloud-cli/pkg/e2ee"
)

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

			// Check if user is authenticated
			preferences := pref.PreferencesInstance()
			if preferences.LoginResponse == nil || preferences.LoginResponse.AccessToken == "" {
				fmt.Println("Error: You need to be logged in to upload files. Please login first.")
				return
			}

			// Create E2EE client
			client := createE2EEClient()

			// Prepare file metadata
			metadata := &e2ee.FileMetadata{
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
				var customMap map[string]string
				err := json.Unmarshal([]byte(customMetadata), &customMap)
				if err != nil {
					fmt.Printf("Error: Invalid custom metadata format: %v\n", err)
					return
				}
				metadata.CustomMetadata = customMap
			}

			// Generate a unique file ID
			fileID := generateFileID(filePath, *metadata)
			fmt.Printf("Generated file ID: %s\n", fileID)

			// Upload the file
			fmt.Println("Starting file encryption and upload...")
			response, err := client.UploadEncryptedFile(filePath, fileID, metadata)
			if err != nil {
				fmt.Printf("Error: Failed to encrypt and upload file: %v\n", err)
				return
			}

			fmt.Println("File successfully encrypted and uploaded!")
			fmt.Printf("Server ID: %s\n", response.ID)
			fmt.Printf("File ID: %s\n", response.FileID)
			fmt.Printf("Created At: %s\n", response.CreatedAt.Format(time.RFC3339))
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
func generateFileID(filePath string, metadata e2ee.FileMetadata) string {
	// Create a unique identifier based on file path, size, and current time
	uniqueStr := fmt.Sprintf("%s_%d_%d", filePath, metadata.OriginalSize, time.Now().UnixNano())

	// For demonstration, generate a simple hash-like ID (in production, use proper hashing)
	hashBytes := make([]byte, 16)
	for i := 0; i < len(uniqueStr) && i < len(hashBytes); i++ {
		hashBytes[i%len(hashBytes)] ^= uniqueStr[i]
	}

	// Convert to hex string
	hexID := fmt.Sprintf("%x", hashBytes)
	return hexID[:16] // Return first 16 hex chars (8 bytes)
}
