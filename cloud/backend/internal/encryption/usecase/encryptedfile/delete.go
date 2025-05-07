// cloud/backend/internal/encryption/usecase/encryptedfile/delete.go
package encryptedfile

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
	domain "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/encryption/domain/encryptedfile"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/httperror"
)

// DeleteEncryptedFileUseCase defines operations for deleting an encrypted file
type DeleteEncryptedFileUseCase interface {
	Execute(ctx context.Context, id primitive.ObjectID) error
}

type deleteEncryptedFileUseCaseImpl struct {
	config     *config.Configuration
	logger     *zap.Logger
	repository domain.Repository
}

// NewDeleteEncryptedFileUseCase creates a new instance of the use case
func NewDeleteEncryptedFileUseCase(
	config *config.Configuration,
	logger *zap.Logger,
	repository domain.Repository,
) DeleteEncryptedFileUseCase {
	return &deleteEncryptedFileUseCaseImpl{
		config:     config,
		logger:     logger.With(zap.String("component", "delete-encrypted-file-usecase")),
		repository: repository,
	}
}

// Execute deletes an encrypted file
func (uc *deleteEncryptedFileUseCaseImpl) Execute(
	ctx context.Context,
	id primitive.ObjectID,
) error {
	// Validate inputs
	if id.IsZero() {
		return httperror.NewForBadRequestWithSingleField("id", "File ID cannot be empty")
	}

	// Check if the file exists
	file, err := uc.repository.GetByID(ctx, id)
	if err != nil {
		uc.logger.Error("Failed to check for file existence",
			zap.String("id", id.Hex()),
			zap.Error(err),
		)
		return fmt.Errorf("failed to check for file existence: %w", err)
	}

	if file == nil {
		return httperror.NewForBadRequestWithSingleField("id", "File not found")
	}

	// Delete the file
	err = uc.repository.DeleteByID(ctx, id)
	if err != nil {
		uc.logger.Error("Failed to delete encrypted file",
			zap.String("id", id.Hex()),
			zap.Error(err),
		)
		return fmt.Errorf("failed to delete encrypted file: %w", err)
	}

	uc.logger.Info("Successfully deleted encrypted file",
		zap.String("id", id.Hex()),
		zap.String("userID", file.UserID.Hex()),
		zap.String("fileID", file.FileID),
	)

	return nil
}
