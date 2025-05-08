// cloud/backend/internal/vault/usecase/encryptedfile/create.go
package encryptedfile

import (
	"context"
	"fmt"
	"io"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
	domain "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/vault/domain/encryptedfile"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/httperror"
)

// CreateEncryptedFileUseCase defines operations for creating a new encrypted file
type CreateEncryptedFileUseCase interface {
	Execute(ctx context.Context, userID primitive.ObjectID, fileID string, encryptedMetadata string, encryptedHash string, encryptionVersion string, encryptedContent io.Reader) (*domain.EncryptedFile, error)
}

type createEncryptedFileUseCaseImpl struct {
	config     *config.Configuration
	logger     *zap.Logger
	repository domain.Repository
}

// NewCreateEncryptedFileUseCase creates a new instance of the use case
func NewCreateEncryptedFileUseCase(
	config *config.Configuration,
	logger *zap.Logger,
	repository domain.Repository,
) CreateEncryptedFileUseCase {
	return &createEncryptedFileUseCaseImpl{
		config:     config,
		logger:     logger.With(zap.String("component", "create-encrypted-file-usecase")),
		repository: repository,
	}
}

// Execute handles the creation of a new encrypted file
func (uc *createEncryptedFileUseCaseImpl) Execute(
	ctx context.Context,
	userID primitive.ObjectID,
	fileID string,
	encryptedMetadata string,
	encryptedHash string,
	encryptionVersion string,
	encryptedContent io.Reader,
) (*domain.EncryptedFile, error) {
	// Validate inputs
	if userID.IsZero() {
		return nil, httperror.NewForBadRequestWithSingleField("user_id", "User ID cannot be empty")
	}

	if fileID == "" {
		return nil, httperror.NewForBadRequestWithSingleField("file_id", "File ID cannot be empty")
	}

	if encryptedContent == nil {
		return nil, httperror.NewForBadRequestWithSingleField("content", "File content cannot be empty")
	}

	// Check if a file with the same userID and fileID already exists
	existingFile, err := uc.repository.GetByFileID(ctx, userID, fileID)
	if err != nil {
		uc.logger.Error("Failed to check for existing file",
			zap.String("userID", userID.Hex()),
			zap.String("fileID", fileID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to check for existing file: %w", err)
	}

	if existingFile != nil {
		return nil, httperror.NewForBadRequestWithSingleField("file_id", "A file with this ID already exists")
	}

	// Create a new encrypted file entity
	file := &domain.EncryptedFile{
		ID:                primitive.NewObjectID(),
		UserID:            userID,
		FileID:            fileID,
		EncryptedMetadata: encryptedMetadata,
		EncryptedHash:     encryptedHash,
		EncryptionVersion: encryptionVersion,
		// StoragePath and EncryptedSize will be set by the repository
	}

	// Store the file
	err = uc.repository.Create(ctx, file, encryptedContent)
	if err != nil {
		uc.logger.Error("Failed to create encrypted file",
			zap.String("userID", userID.Hex()),
			zap.String("fileID", fileID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to create encrypted file: %w", err)
	}

	uc.logger.Info("Successfully created encrypted file",
		zap.String("id", file.ID.Hex()),
		zap.String("userID", userID.Hex()),
		zap.String("fileID", fileID),
	)

	return file, nil
}
