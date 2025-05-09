// cloud/backend/internal/papercloud/repo/file/impl.go
package file

import (
	"fmt"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.uber.org/zap"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
	dom_file "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/papercloud/domain/file"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/papercloud/repo/file/metadata"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/papercloud/repo/file/storage"
	s3storage "github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/storage/object/s3"
)

// Composite repository implementing the domain FileRepository interface
type fileRepositoryImpl struct {
	logger   *zap.Logger
	metadata metadata.FileMetadataRepository
	storage  storage.FileStorageRepository
}

// Constructor function for FileRepository
func NewFileRepository(
	logger *zap.Logger,
	metadataRepo metadata.FileMetadataRepository,
	storageRepo storage.FileStorageRepository,
) dom_file.FileRepository {
	return &fileRepositoryImpl{
		logger:   logger.With(zap.String("repository", "file")),
		metadata: metadataRepo,
		storage:  storageRepo,
	}
}

// Constructor for metadata repository
func NewFileMetadataRepository(
	cfg *config.Configuration,
	logger *zap.Logger,
	client *mongo.Client,
) metadata.FileMetadataRepository {
	return metadata.NewRepository(cfg, logger, client)
}

// Constructor for storage repository
func NewFileStorageRepository(
	cfg *config.Configuration,
	logger *zap.Logger,
	s3 s3storage.S3ObjectStorage,
) storage.FileStorageRepository {
	return storage.NewRepository(cfg, logger, s3)
}

// Create implements the FileRepository.Create method
func (repo *fileRepositoryImpl) Create(file *dom_file.File) error {
	// Generate ID if not provided
	if file.ID == "" {
		file.ID = uuid.New().String()
	}

	// If FileID is not set, generate one - ideally this comes from client
	if file.FileID == "" {
		file.FileID = uuid.New().String()
	}

	// Store metadata
	return repo.metadata.Create(file)
}

// Get implements the FileRepository.Get method
func (repo *fileRepositoryImpl) Get(id string) (*dom_file.File, error) {
	return repo.metadata.Get(id)
}

// GetByCollection implements the FileRepository.GetByCollection method
func (repo *fileRepositoryImpl) GetByCollection(collectionID string) ([]*dom_file.File, error) {
	return repo.metadata.GetByCollection(collectionID)
}

// Update implements the FileRepository.Update method
func (repo *fileRepositoryImpl) Update(file *dom_file.File) error {
	return repo.metadata.Update(file)
}

// Delete implements the FileRepository.Delete method
func (repo *fileRepositoryImpl) Delete(id string) error {
	// First get the file to get its storage path
	file, err := repo.metadata.Get(id)
	if err != nil {
		return err
	}
	if file == nil {
		repo.logger.Info("file not found for deletion", zap.String("id", id))
		return nil
	}

	// Delete from S3
	if err := repo.storage.DeleteEncryptedData(file.StoragePath); err != nil {
		repo.logger.Error("failed to delete file data",
			zap.String("id", id),
			zap.Error(err))
		return err
	}

	// Delete metadata
	return repo.metadata.Delete(id)
}

// StoreEncryptedData implements the FileRepository.StoreEncryptedData method
func (repo *fileRepositoryImpl) StoreEncryptedData(fileID string, encryptedData []byte) error {
	// Get file metadata to ensure it exists
	file, err := repo.metadata.Get(fileID)
	if err != nil {
		return err
	}
	if file == nil {
		return fmt.Errorf("file not found: %s", fileID)
	}

	// Store the encrypted data
	storagePath, err := repo.storage.StoreEncryptedData(file.OwnerID, file.FileID, encryptedData)
	if err != nil {
		return err
	}

	// Update the file storage path and size
	file.StoragePath = storagePath
	file.EncryptedSize = int64(len(encryptedData))

	// Update metadata
	return repo.metadata.Update(file)
}

// GetEncryptedData implements the FileRepository.GetEncryptedData method
func (repo *fileRepositoryImpl) GetEncryptedData(fileID string) ([]byte, error) {
	// Get file metadata to ensure it exists and get storage path
	file, err := repo.metadata.Get(fileID)
	if err != nil {
		return nil, err
	}
	if file == nil {
		return nil, fmt.Errorf("file not found: %s", fileID)
	}

	// Retrieve the encrypted data
	return repo.storage.GetEncryptedData(file.StoragePath)
}
