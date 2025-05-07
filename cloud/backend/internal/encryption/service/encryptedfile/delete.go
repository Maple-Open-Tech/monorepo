// cloud/backend/internal/encryption/service/encryptedfile/delete.go
package encryptedfile

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config/constants"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/encryption/usecase/encryptedfile"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/httperror"
)

// DeleteEncryptedFileService defines operations for deleting an encrypted file
type DeleteEncryptedFileService interface {
	Execute(ctx context.Context, id primitive.ObjectID) error
}

type deleteEncryptedFileServiceImpl struct {
	config         *config.Configuration
	logger         *zap.Logger
	getByIDUseCase encryptedfile.GetEncryptedFileByIDUseCase
	deleteUseCase  encryptedfile.DeleteEncryptedFileUseCase
}

// NewDeleteEncryptedFileService creates a new instance of the service
func NewDeleteEncryptedFileService(
	config *config.Configuration,
	logger *zap.Logger,
	getByIDUseCase encryptedfile.GetEncryptedFileByIDUseCase,
	deleteUseCase encryptedfile.DeleteEncryptedFileUseCase,
) DeleteEncryptedFileService {
	return &deleteEncryptedFileServiceImpl{
		config:         config,
		logger:         logger.With(zap.String("component", "delete-encrypted-file-service")),
		getByIDUseCase: getByIDUseCase,
		deleteUseCase:  deleteUseCase,
	}
}

// Execute deletes an encrypted file after verifying ownership
func (s *deleteEncryptedFileServiceImpl) Execute(
	ctx context.Context,
	id primitive.ObjectID,
) error {
	// First get the file to verify ownership
	file, err := s.getByIDUseCase.Execute(ctx, id)
	if err != nil {
		return err
	}

	if file == nil {
		return httperror.NewForBadRequestWithSingleField("id", "File not found")
	}

	// Verify that the authenticated user has access to this file
	userID, ok := ctx.Value(constants.SessionFederatedUserID).(primitive.ObjectID)
	if ok && !userID.IsZero() && file.UserID != userID {
		s.logger.Warn("Unauthorized file deletion attempt",
			zap.String("file_id", id.Hex()),
			zap.String("file_owner", file.UserID.Hex()),
			zap.String("requester", userID.Hex()),
		)
		return httperror.NewForForbiddenWithSingleField("message", "You do not have permission to delete this file")
	}

	// Delete the file using the use case
	return s.deleteUseCase.Execute(ctx, id)
}
