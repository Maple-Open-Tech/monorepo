// cloud/backend/internal/encryption/usecase/encryptedfile/download.go
package encryptedfile

import (
	"context"
	"fmt"
	"io"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
	domain "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/encryption/domain/encryptedfile"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/httperror"
)

// DownloadEncryptedFileUseCase defines operations for downloading encrypted file content
type DownloadEncryptedFileUseCase interface {
	Execute(ctx context.Context, id primitive.ObjectID) (io.ReadCloser, error)
}

type downloadEncryptedFileUseCaseImpl struct {
	config     *config.Configuration
	logger     *zap.Logger
	repository domain.Repository
}

// NewDownloadEncryptedFileUseCase creates a new instance of the use case
func NewDownloadEncryptedFileUseCase(
	config *config.Configuration,
	logger *zap.Logger,
	repository domain.Repository,
) DownloadEncryptedFileUseCase {
	return &downloadEncryptedFileUseCaseImpl{
		config:     config,
		logger:     logger.With(zap.String("component", "download-encrypted-file-usecase")),
		repository: repository,
	}
}

// Execute downloads the encrypted content of a file
func (uc *downloadEncryptedFileUseCaseImpl) Execute(
	ctx context.Context,
	id primitive.ObjectID,
) (io.ReadCloser, error) {
	// Validate inputs
	if id.IsZero() {
		return nil, httperror.NewForBadRequestWithSingleField("id", "File ID cannot be empty")
	}

	// Get the file metadata
	file, err := uc.repository.GetByID(ctx, id)
	if err != nil {
		uc.logger.Error("Failed to get file metadata for download",
			zap.String("id", id.Hex()),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get file metadata: %w", err)
	}

	if file == nil {
		return nil, httperror.NewForBadRequestWithSingleField("id", "File not found")
	}

	// Download the content
	content, err := uc.repository.DownloadContent(ctx, file)
	if err != nil {
		uc.logger.Error("Failed to download encrypted file content",
			zap.String("id", id.Hex()),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to download file content: %w", err)
	}

	uc.logger.Debug("Successfully downloaded encrypted file content",
		zap.String("id", id.Hex()),
		zap.String("userID", file.UserID.Hex()),
		zap.String("fileID", file.FileID),
	)

	return content, nil
}
