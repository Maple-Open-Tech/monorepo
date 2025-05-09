// monorepo/cloud/backend/internal/papercloud/domain/collection/model.go
package collection

import "time"

const (
	CollectionTypeFolder = "folder"
	CollectionTypeAlbum  = "album"
)

// Collection represents a folder or album
type Collection struct {
	ID           string
	Name         string // Encrypted
	Path         string // Encrypted
	Type         string // "folder" or "album"
	CreatedAt    time.Time
	UpdatedAt    time.Time
	EncryptedKey []byte // Collection key encrypted with master key
	SharedWith   []Share
}

// Share represents a shared collection
type Share struct {
	UserID          string
	PublicKey       []byte
	EncryptedKey    []byte // Collection key encrypted with recipient's public key
	PermissionLevel string
	CreatedAt       time.Time
}
