// monorepo/cloud/backend/internal/papercloud/domain/collection/model.go
package collection

import (
	"time"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/iam/domain/keys"
)

const (
	CollectionTypeFolder = "folder"
	CollectionTypeAlbum  = "album"
)

// Collection represents a folder or album
type Collection struct {
	ID        string    `bson:"id" json:"id"`
	OwnerID   string    `bson:"owner_id" json:"owner_id"`
	Name      string    `bson:"name" json:"name"` // Encrypted
	Path      string    `bson:"path" json:"path"` // Encrypted
	Type      string    `bson:"type" json:"type"` // "folder" or "album"
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`

	// Collection key encrypted with owner's master key
	EncryptedKey keys.EncryptedCollectionKey `bson:"encrypted_key" json:"encrypted_key"`

	// Collection shares (users with access)
	SharedWith []Share `bson:"shared_with" json:"shared_with"`
}

// Share represents a shared collection
type Share struct {
	UserID          string    `bson:"user_id" json:"user_id"`
	PublicKey       []byte    `bson:"public_key" json:"public_key"`
	EncryptedKey    []byte    `bson:"encrypted_key" json:"encrypted_key"` // Collection key encrypted with recipient's public key
	PermissionLevel string    `bson:"permission_level" json:"permission_level"`
	CreatedAt       time.Time `bson:"created_at" json:"created_at"`
}
