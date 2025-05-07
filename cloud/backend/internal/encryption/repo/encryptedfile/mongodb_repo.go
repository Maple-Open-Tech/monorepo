// cloud/backend/internal/encryption/repo/encryptedfile/mongodb_repo.go
package encryptedfile

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
	domain "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/encryption/domain/encryptedfile"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/storage/object/s3"
)

// encryptedFileRepository implements the domain.Repository interface
type encryptedFileRepository struct {
	logger     *zap.Logger
	collection *mongo.Collection
	s3Storage  *s3ObjectStorage
}

// NewRepository creates a new repository for encrypted files
func NewRepository(
	cfg *config.Configuration,
	logger *zap.Logger,
	dbClient *mongo.Client,
	s3Provider s3.S3ObjectStorage,
) domain.Repository {
	// Initialize the MongoDB collection
	collection := dbClient.Database(cfg.DB.MapleAuthName).Collection("encrypted_files")

	// Create indexes for efficient queries
	indexModels := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "file_id", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "created_at", Value: -1},
			},
		},
	}

	_, err := collection.Indexes().CreateMany(context.Background(), indexModels)
	if err != nil {
		logger.Error("Failed to create indexes for encrypted files collection", zap.Error(err))
	}

	// Initialize the S3 storage
	s3Storage := NewS3ObjectStorage(cfg, logger, s3Provider)

	return &encryptedFileRepository{
		logger:     logger.With(zap.String("component", "encrypted-file-repository")),
		collection: collection,
		s3Storage:  s3Storage,
	}
}

// Create stores a new encrypted file
func (repo *encryptedFileRepository) Create(
	ctx context.Context,
	file *domain.EncryptedFile,
	encryptedContent io.Reader,
) error {
	// Generate a new ID if not provided
	if file.ID == primitive.NilObjectID {
		file.ID = primitive.NewObjectID()
	}

	// Set creation and modification times
	now := time.Now()
	file.CreatedAt = now
	file.ModifiedAt = now

	// Upload the encrypted content to S3
	userID := file.UserID.Hex()
	storagePath, err := repo.s3Storage.UploadFile(ctx, userID, file.FileID, encryptedContent)
	if err != nil {
		return fmt.Errorf("failed to upload encrypted file: %w", err)
	}

	// Set the storage path in the metadata
	file.StoragePath = storagePath

	// Insert the metadata into MongoDB
	_, err = repo.collection.InsertOne(ctx, file)
	if err != nil {
		// If MongoDB insert fails, try to clean up the uploaded file
		cleanupErr := repo.s3Storage.DeleteFile(ctx, storagePath)
		if cleanupErr != nil {
			repo.logger.Error("Failed to clean up S3 file after MongoDB insertion error",
				zap.String("storagePath", storagePath),
				zap.Error(cleanupErr),
			)
		}

		return fmt.Errorf("failed to insert encrypted file metadata: %w", err)
	}

	repo.logger.Debug("Successfully created encrypted file",
		zap.String("id", file.ID.Hex()),
		zap.String("userID", userID),
		zap.String("fileID", file.FileID),
	)

	return nil
}

// GetByID retrieves an encrypted file by its ID
func (repo *encryptedFileRepository) GetByID(
	ctx context.Context,
	id primitive.ObjectID,
) (*domain.EncryptedFile, error) {
	var file domain.EncryptedFile

	err := repo.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&file)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil // Not found
		}
		return nil, fmt.Errorf("failed to get encrypted file: %w", err)
	}

	return &file, nil
}

// GetByFileID retrieves an encrypted file by user ID and file ID
func (repo *encryptedFileRepository) GetByFileID(
	ctx context.Context,
	userID primitive.ObjectID,
	fileID string,
) (*domain.EncryptedFile, error) {
	var file domain.EncryptedFile

	err := repo.collection.FindOne(
		ctx,
		bson.M{
			"user_id": userID,
			"file_id": fileID,
		},
	).Decode(&file)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil // Not found
		}
		return nil, fmt.Errorf("failed to get encrypted file: %w", err)
	}

	return &file, nil
}

// UpdateByID updates an encrypted file
func (repo *encryptedFileRepository) UpdateByID(
	ctx context.Context,
	file *domain.EncryptedFile,
	encryptedContent io.Reader,
) error {
	// Get the existing file to retrieve the storage path
	existingFile, err := repo.GetByID(ctx, file.ID)
	if err != nil {
		return err
	}

	if existingFile == nil {
		return fmt.Errorf("file not found")
	}

	// Update modification time
	file.ModifiedAt = time.Now()
	file.CreatedAt = existingFile.CreatedAt // Preserve creation time

	// If a new encrypted content is provided, update the file in S3
	if encryptedContent != nil {
		userID := file.UserID.Hex()

		// Upload the new content (this will create a new object in S3)
		newStoragePath, err := repo.s3Storage.UploadFile(ctx, userID, file.FileID, encryptedContent)
		if err != nil {
			return fmt.Errorf("failed to upload updated encrypted file: %w", err)
		}

		// Update the storage path
		file.StoragePath = newStoragePath

		// Delete the old file from S3 (after successful upload)
		err = repo.s3Storage.DeleteFile(ctx, existingFile.StoragePath)
		if err != nil {
			// Log but don't fail the update if cleanup fails
			repo.logger.Warn("Failed to delete old encrypted file version",
				zap.String("oldStoragePath", existingFile.StoragePath),
				zap.Error(err),
			)
		}
	} else {
		// If no new content, keep the existing storage path
		file.StoragePath = existingFile.StoragePath
	}

	// Update the metadata in MongoDB
	_, err = repo.collection.ReplaceOne(
		ctx,
		bson.M{"_id": file.ID},
		file,
	)

	if err != nil {
		return fmt.Errorf("failed to update encrypted file metadata: %w", err)
	}

	repo.logger.Debug("Successfully updated encrypted file",
		zap.String("id", file.ID.Hex()),
		zap.String("userID", file.UserID.Hex()),
		zap.String("fileID", file.FileID),
	)

	return nil
}

// DeleteByID deletes an encrypted file
func (repo *encryptedFileRepository) DeleteByID(
	ctx context.Context,
	id primitive.ObjectID,
) error {
	// First get the file to retrieve the storage path
	file, err := repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if file == nil {
		return fmt.Errorf("file not found")
	}

	// Delete from MongoDB first
	_, err = repo.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("failed to delete encrypted file metadata: %w", err)
	}

	// Then delete from S3
	err = repo.s3Storage.DeleteFile(ctx, file.StoragePath)
	if err != nil {
		// Log the error but don't fail the operation
		// The metadata is already deleted, and we don't want to block the client
		repo.logger.Error("Failed to delete encrypted file content",
			zap.String("id", id.Hex()),
			zap.String("storagePath", file.StoragePath),
			zap.Error(err),
		)
	}

	repo.logger.Debug("Successfully deleted encrypted file",
		zap.String("id", id.Hex()),
		zap.String("userID", file.UserID.Hex()),
		zap.String("fileID", file.FileID),
	)

	return nil
}

// ListByUserID lists all encrypted files for a user
func (repo *encryptedFileRepository) ListByUserID(
	ctx context.Context,
	userID primitive.ObjectID,
) ([]*domain.EncryptedFile, error) {
	// Define query options for sorting by creation time
	findOptions := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})

	// Execute the query
	cursor, err := repo.collection.Find(
		ctx,
		bson.M{"user_id": userID},
		findOptions,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to list encrypted files: %w", err)
	}
	defer cursor.Close(ctx)

	// Decode the results
	var files []*domain.EncryptedFile
	if err := cursor.All(ctx, &files); err != nil {
		return nil, fmt.Errorf("failed to decode encrypted files: %w", err)
	}

	return files, nil
}

// DownloadContent downloads the encrypted content of a file
func (repo *encryptedFileRepository) DownloadContent(
	ctx context.Context,
	file *domain.EncryptedFile,
) (io.ReadCloser, error) {
	// Use the S3 storage to download the file
	content, err := repo.s3Storage.DownloadFile(ctx, file.StoragePath)
	if err != nil {
		return nil, fmt.Errorf("failed to download encrypted file: %w", err)
	}

	return content, nil
}

// GetDownloadURL generates a presigned URL for direct download
func (repo *encryptedFileRepository) GetDownloadURL(
	ctx context.Context,
	file *domain.EncryptedFile,
	expiryDuration time.Duration,
) (string, error) {
	// Use the S3 storage to generate a download URL
	url, err := repo.s3Storage.GetDownloadURL(ctx, file.StoragePath, expiryDuration)
	if err != nil {
		return "", fmt.Errorf("failed to generate download URL: %w", err)
	}

	return url, nil
}
