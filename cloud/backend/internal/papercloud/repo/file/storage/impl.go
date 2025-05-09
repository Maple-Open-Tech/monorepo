// cloud/backend/internal/papercloud/repo/file/storage/impl.go
package storage

import (
	"context"
	"fmt"
	"io"
	"time"

	"go.uber.org/zap"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
	s3storage "github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/storage/object/s3"
)

type FileStorageRepository interface {
	StoreEncryptedData(ownerID string, fileID string, encryptedData []byte) (string, error)
	GetEncryptedData(storagePath string) ([]byte, error)
	DeleteEncryptedData(storagePath string) error
	GeneratePresignedURL(storagePath string, duration time.Duration) (string, error)
}

type fileStorageRepositoryImpl struct {
	Logger  *zap.Logger
	Storage s3storage.S3ObjectStorage
}

func NewRepository(cfg *config.Configuration, logger *zap.Logger, s3 s3storage.S3ObjectStorage) FileStorageRepository {
	return &fileStorageRepositoryImpl{
		Logger:  logger.With(zap.String("repository", "file_storage")),
		Storage: s3,
	}
}

// StoreEncryptedData uploads encrypted file data to S3 and returns the storage path
func (impl *fileStorageRepositoryImpl) StoreEncryptedData(ownerID string, fileID string, encryptedData []byte) (string, error) {
	ctx := context.Background()

	// Generate a storage path using a deterministic pattern
	storagePath := fmt.Sprintf("users/%s/files/%s", ownerID, fileID)

	// Always store encrypted data as private
	err := impl.Storage.UploadContentWithVisibility(ctx, storagePath, encryptedData, false)
	if err != nil {
		impl.Logger.Error("Failed to store encrypted data",
			zap.String("fileID", fileID),
			zap.String("ownerID", ownerID),
			zap.Error(err))
		return "", err
	}

	return storagePath, nil
}

// GetEncryptedData retrieves encrypted file data from S3
func (impl *fileStorageRepositoryImpl) GetEncryptedData(storagePath string) ([]byte, error) {
	ctx := context.Background()

	// Get the encrypted data
	reader, err := impl.Storage.GetBinaryData(ctx, storagePath)
	if err != nil {
		impl.Logger.Error("Failed to get encrypted data",
			zap.String("storagePath", storagePath),
			zap.Error(err))
		return nil, err
	}
	defer reader.Close()

	// Read all data into memory
	data, err := io.ReadAll(reader)
	if err != nil {
		impl.Logger.Error("Failed to read encrypted data",
			zap.String("storagePath", storagePath),
			zap.Error(err))
		return nil, err
	}

	return data, nil
}

// DeleteEncryptedData removes encrypted file data from S3
func (impl *fileStorageRepositoryImpl) DeleteEncryptedData(storagePath string) error {
	ctx := context.Background()

	// Delete the encrypted data
	err := impl.Storage.DeleteByKeys(ctx, []string{storagePath})
	if err != nil {
		impl.Logger.Error("Failed to delete encrypted data",
			zap.String("storagePath", storagePath),
			zap.Error(err))
		return err
	}

	return nil
}

// GeneratePresignedURL creates a time-limited URL for downloading the file directly
func (impl *fileStorageRepositoryImpl) GeneratePresignedURL(storagePath string, duration time.Duration) (string, error) {
	ctx := context.Background()

	// Generate presigned URL
	url, err := impl.Storage.GetPresignedURL(ctx, storagePath, duration)
	if err != nil {
		impl.Logger.Error("Failed to generate presigned URL",
			zap.String("storagePath", storagePath),
			zap.Error(err))
		return "", err
	}

	return url, nil
}
