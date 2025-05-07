// cloud/backend/internal/encryption/service/encryptedfile/dto.go
package encryptedfile

import (
	"io"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	domain "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/encryption/domain/encryptedfile"
)

// CreateFileRequest represents the data needed to create a new encrypted file
type CreateFileRequest struct {
	UserID            primitive.ObjectID
	FileID            string
	EncryptedMetadata string
	EncryptedHash     string
	EncryptionVersion string
	EncryptedContent  io.Reader
}

// UpdateFileRequest represents the data needed to update an encrypted file
type UpdateFileRequest struct {
	ID                primitive.ObjectID
	EncryptedMetadata string
	EncryptedHash     string
	EncryptedContent  io.Reader
}

// FileResponse represents file metadata returned by the service
type FileResponse struct {
	ID                primitive.ObjectID `json:"id"`
	UserID            primitive.ObjectID `json:"user_id"`
	FileID            string             `json:"file_id"`
	EncryptedMetadata string             `json:"encrypted_metadata"`
	EncryptionVersion string             `json:"encryption_version"`
	EncryptedHash     string             `json:"encrypted_hash"`
	CreatedAt         time.Time          `json:"created_at"`
	ModifiedAt        time.Time          `json:"modified_at"`
}

// FilesListResponse represents a list of file metadata
type FilesListResponse struct {
	Files []*FileResponse `json:"files"`
}

// FileURLResponse represents a presigned download URL for a file
type FileURLResponse struct {
	URL       string    `json:"url"`
	ExpiresAt time.Time `json:"expires_at"`
}

// DomainToFileResponse converts a domain file model to a FileResponse DTO
func DomainToFileResponse(file *domain.EncryptedFile) *FileResponse {
	return &FileResponse{
		ID:                file.ID,
		UserID:            file.UserID,
		FileID:            file.FileID,
		EncryptedMetadata: file.EncryptedMetadata,
		EncryptionVersion: file.EncryptionVersion,
		EncryptedHash:     file.EncryptedHash,
		CreatedAt:         file.CreatedAt,
		ModifiedAt:        file.ModifiedAt,
	}
}
