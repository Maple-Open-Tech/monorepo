package keys

// MasterKey represents the root encryption key for a user
type MasterKey struct {
	Key []byte `json:"key" bson:"key"`
}

// EncryptedMasterKey is the master key encrypted with the key encryption key
type EncryptedMasterKey struct {
	Ciphertext []byte `json:"ciphertext" bson:"ciphertext"`
	Nonce      []byte `json:"nonce" bson:"nonce"`
}

// KeyEncryptionKey derived from user password
type KeyEncryptionKey struct {
	Key  []byte `json:"key" bson:"key"`
	Salt []byte `json:"salt" bson:"salt"`
}

// PublicKey for asymmetric encryption
type PublicKey struct {
	Key            []byte `json:"key" bson:"key"`
	VerificationID string `json:"verification_id" bson:"verification_id"`
}

// PrivateKey for asymmetric decryption
type PrivateKey struct {
	Key []byte `json:"key" bson:"key"`
}

// EncryptedPrivateKey is the private key encrypted with the master key
type EncryptedPrivateKey struct {
	Ciphertext []byte `json:"ciphertext" bson:"ciphertext"`
	Nonce      []byte `json:"nonce" bson:"nonce"`
}

// RecoveryKey for account recovery
type RecoveryKey struct {
	Key []byte `json:"key" bson:"key"`
}

// EncryptedRecoveryKey is the recovery key encrypted with the master key
type EncryptedRecoveryKey struct {
	Ciphertext []byte `json:"ciphertext" bson:"ciphertext"`
	Nonce      []byte `json:"nonce" bson:"nonce"`
}

// CollectionKey encrypts files in a collection
type CollectionKey struct {
	Key          []byte `json:"key" bson:"key"`
	CollectionID string `json:"collection_id" bson:"collection_id"`
}

// EncryptedCollectionKey is the collection key encrypted with master key
type EncryptedCollectionKey struct {
	Ciphertext []byte `json:"ciphertext" bson:"ciphertext"`
	Nonce      []byte `json:"nonce" bson:"nonce"`
}

// FileKey encrypts a specific file
type FileKey struct {
	Key    []byte `json:"key" bson:"key"`
	FileID string `json:"file_id" bson:"file_id"`
}

// EncryptedFileKey is the file key encrypted with collection key
type EncryptedFileKey struct {
	Ciphertext []byte `json:"ciphertext" bson:"ciphertext"`
	Nonce      []byte `json:"nonce" bson:"nonce"`
}

// MasterKeyEncryptedWithRecoveryKey allows account recovery
type MasterKeyEncryptedWithRecoveryKey struct {
	Ciphertext []byte `json:"ciphertext" bson:"ciphertext"`
	Nonce      []byte `json:"nonce" bson:"nonce"`
}
