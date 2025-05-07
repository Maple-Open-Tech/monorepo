// cloud/backend/internal/encryption/usecase/encryptedfile/getbyid.go
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

// GetEncryptedFileByIDUseCase defines operations for retrieving an encrypted file by ID
type GetEncryptedFileByIDUseCase interface {
	Execute(ctx context.Context, id primitive.ObjectID) (*domain.EncryptedFile, error)
}

type getEncryptedFileByIDUseCaseImpl struct {
	config     *config.Configuration
	logger     *zap.Logger
	repository domain.Repository
}

// NewGetEncryptedFileByIDUseCase creates a new instance of the use case
func NewGetEncryptedFileByIDUseCase(
	config *config.Configuration,
	logger *zap.Logger,
	repository domain.Repository,
) GetEncryptedFileByIDUseCase {
	return &getEncryptedFileByIDUseCaseImpl{
		config:     config,
		logger:     logger.With(zap.String("component", "get-encrypted-file-by-id-usecase")),
		repository: repository,
	}
}

// Execute retrieves an encrypted file by its ID
func (uc *getEncryptedFileByIDUseCaseImpl) Execute(
	ctx context.Context,
	id primitive.ObjectID,
) (*domain.EncryptedFile, error) {
	// Validate inputs
	if id.IsZero() {
		return nil, httperror.NewForBadRequestWithSingleField("id", "File ID cannot be empty")
	}

	// Retrieve the file
	file, err := uc.repository.GetByID(ctx, id)
	if err != nil {
		uc.logger.Error("Failed to get encrypted file",
			zap.String("id", id.Hex()),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get encrypted file: %w", err)
	}

	if file == nil {
		return nil, httperror.NewForBadRequestWithSingleField("id", "File not found")
	}

	return file, nil
}
