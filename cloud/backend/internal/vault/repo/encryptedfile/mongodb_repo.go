// cloud/backend/internal/vault/repo/encryptedfile/mongodb_repo.go
package encryptedfile

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.uber.org/zap"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
	domain "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/vault/domain/encryptedfile"
)

// encryptedFileRepository implements the domain.Repository interface
type encryptedFileRepository struct {
	logger     *zap.Logger
	collection *mongo.Collection
	database   *mongo.Database
}

// NewRepository creates a new repository for encrypted files
func NewRepository(
	cfg *config.Configuration,
	logger *zap.Logger,
	dbClient *mongo.Client,
) domain.Repository {
	// Initialize the MongoDB database
	database := dbClient.Database(cfg.DB.EncryptionName)

	// Initialize the MongoDB collection for file metadata
	collection := database.Collection("encrypted_files")

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

	return &encryptedFileRepository{
		logger:     logger.With(zap.String("component", "encrypted-file-repository")),
		collection: collection,
		database:   database,
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

	// Generate a unique storage path for the file
	userID := file.UserID.Hex()
	storagePath := fmt.Sprintf("%s/%s", userID, file.FileID)
	file.StoragePath = storagePath

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

	// Use the existing storage path
	file.StoragePath = existingFile.StoragePath

	// If a new encrypted content is provided, update the file in GridFS
	if encryptedContent != nil {
		// Delete the existing file from GridFS
		cursor, err := repo.database.Collection("encryptedFiles.files").Find(
			ctx,
			bson.M{"filename": existingFile.StoragePath},
		)
		if err != nil {
			repo.logger.Error("Failed to find existing GridFS file", zap.Error(err))
		} else {
			var existingFiles []struct {
				ID primitive.ObjectID `bson:"_id"`
			}
			if err := cursor.All(ctx, &existingFiles); err == nil {

			}
			cursor.Close(ctx)
		}

		if err != nil {
			return fmt.Errorf("failed to open GridFS upload stream: %w", err)
		}

	} else {
		// If no new content, keep the existing size
		file.EncryptedSize = existingFile.EncryptedSize
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
