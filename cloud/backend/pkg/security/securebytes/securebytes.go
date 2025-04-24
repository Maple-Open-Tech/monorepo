package securebytes

import (
	"errors"

	"github.com/awnumar/memguard"
)

// SecureBytes is used to store a byte slice securely in memory.
type SecureBytes struct {
	buffer *memguard.LockedBuffer
}

// NewSecureBytes creates a new SecureBytes instance from the given byte slice.
func NewSecureBytes(b []byte) (*SecureBytes, error) {
	if len(b) == 0 {
		return nil, errors.New("byte slice cannot be empty")
	}

	buffer := memguard.NewBuffer(len(b))

	copy(buffer.Bytes(), b)

	buffer.Lock()

	return &SecureBytes{buffer: buffer}, nil
}

// Bytes returns the securely stored byte slice.
func (sb *SecureBytes) Bytes() []byte {
	return sb.buffer.Bytes()
}

// Wipe removes the byte slice from memory and makes it unrecoverable.
func (sb *SecureBytes) Wipe() error {
	sb.buffer.Wipe()
	sb.buffer = nil
	return nil
}
