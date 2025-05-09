// monorepo/cloud/backend/internal/papercloud/domain/file/model.go
package file

import "time"

// File represents an encrypted file stored in the system
// The actual file content is stored in S3 while this entity
// maintains the minimal metadata needed
type File struct {
	ID           string
	CollectionID string

	// Encrypted file identifier (client-generated)
	// This would be a client-side generated id that is
	// meaningful to the client but opaque to the server
	FileID string `bson:"file_id" json:"file_id"`

	// The path/key in S3 storage where the encrypted file is stored
	StoragePath string `bson:"storage_path" json:"storage_path"`

	// Size of the encrypted file in bytes
	EncryptedSize int64 `bson:"encrypted_size" json:"encrypted_size"`

	// Encrypted metadata (JSON blob encrypted by client)
	// Contains file name, original size, modification date, content type, etc.
	// This is encrypted on the client side and opaque to the server
	EncryptedMetadata string `bson:"encrypted_metadata" json:"encrypted_metadata"`

	// Version identifier for the encryption scheme used
	EncryptionVersion string `bson:"encryption_version" json:"encryption_version"`

	// Hash of the encrypted file for integrity checking
	EncryptedHash string `bson:"encrypted_hash" json:"encrypted_hash"`

	// When was this file uploaded
	CreatedAt time.Time `bson:"created_at" json:"created_at"`

	// When was this file last modified
	ModifiedAt time.Time `bson:"modified_at" json:"modified_at"`
}
