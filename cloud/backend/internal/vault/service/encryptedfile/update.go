// cloud/backend/internal/vault/service/encryptedfile/update.go
package encryptedfile

import (
	"context"
	"io"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config/constants"
	domain "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/vault/domain/encryptedfile"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/vault/usecase/encryptedfile"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/httperror"
)

// UpdateEncryptedFileService defines operations for updating an encrypted file
type UpdateEncryptedFileService interface {
	Execute(ctx context.Context, id primitive.ObjectID, encryptedMetadata string, encryptedHash string, encryptedContent io.Reader) (*domain.EncryptedFile, error)
}

type updateEncryptedFileServiceImpl struct {
	config         *config.Configuration
	logger         *zap.Logger
	getByIDUseCase encryptedfile.GetEncryptedFileByIDUseCase
	updateUseCase  encryptedfile.UpdateEncryptedFileUseCase
}

// NewUpdateEncryptedFileService creates a new instance of the service
func NewUpdateEncryptedFileService(
	config *config.Configuration,
	logger *zap.Logger,
	getByIDUseCase encryptedfile.GetEncryptedFileByIDUseCase,
	updateUseCase encryptedfile.UpdateEncryptedFileUseCase,
) UpdateEncryptedFileService {
	return &updateEncryptedFileServiceImpl{
		config:         config,
		logger:         logger.With(zap.String("component", "update-encrypted-file-service")),
		getByIDUseCase: getByIDUseCase,
		updateUseCase:  updateUseCase,
	}
}

// Execute updates an encrypted file after verifying ownership
func (s *updateEncryptedFileServiceImpl) Execute(
	ctx context.Context,
	id primitive.ObjectID,
	encryptedMetadata string,
	encryptedHash string,
	encryptedContent io.Reader,
) (*domain.EncryptedFile, error) {
	// First get the file to verify ownership
	file, err := s.getByIDUseCase.Execute(ctx, id)
	if err != nil {
		return nil, err
	}

	if file == nil {
		return nil, httperror.NewForBadRequestWithSingleField("id", "File not found")
	}

	// Verify that the authenticated user has access to this file
	userID, ok := ctx.Value(constants.SessionFederatedUserID).(primitive.ObjectID)
	if ok && !userID.IsZero() && file.UserID != userID {
		s.logger.Warn("Unauthorized file update attempt",
			zap.String("file_id", id.Hex()),
			zap.String("file_owner", file.UserID.Hex()),
			zap.String("requester", userID.Hex()),
		)
		return nil, httperror.NewForForbiddenWithSingleField("message", "You do not have permission to update this file")
	}

	// Update the file using the use case
	return s.updateUseCase.Execute(ctx, id, encryptedMetadata, encryptedHash, encryptedContent)
}
