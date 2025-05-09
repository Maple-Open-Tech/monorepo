package file

import "time"

// File represents an encrypted file
type File struct {
	ID            string
	CollectionID  string
	Size          int64
	MimeType      string // Encrypted
	Name          string // Encrypted
	Metadata      []byte // Encrypted
	EncryptedKey  []byte // File key encrypted with collection key
	EncryptedData []byte // File data encrypted with file key
	Checksum      string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
