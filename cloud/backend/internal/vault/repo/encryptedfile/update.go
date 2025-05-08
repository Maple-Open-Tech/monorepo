// cloud/backend/internal/vault/repo/encryptedfile/mongodb_repo.go
package encryptedfile

import (
	"context"
	"fmt"
	"io"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.uber.org/zap"

	domain "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/vault/domain/encryptedfile"
)

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
