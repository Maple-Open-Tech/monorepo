// cloud/backend/internal/encryption/service/encryptedfile/getbyid.go
package encryptedfile

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config/constants"
	domain "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/encryption/domain/encryptedfile"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/encryption/usecase/encryptedfile"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/httperror"
)

// GetEncryptedFileByIDService defines operations for retrieving an encrypted file by ID
type GetEncryptedFileByIDService interface {
	Execute(ctx context.Context, id primitive.ObjectID) (*domain.EncryptedFile, error)
}

type getEncryptedFileByIDServiceImpl struct {
	config         *config.Configuration
	logger         *zap.Logger
	getByIDUseCase encryptedfile.GetEncryptedFileByIDUseCase
}

// NewGetEncryptedFileByIDService creates a new instance of the service
func NewGetEncryptedFileByIDService(
	config *config.Configuration,
	logger *zap.Logger,
	getByIDUseCase encryptedfile.GetEncryptedFileByIDUseCase,
) GetEncryptedFileByIDService {
	return &getEncryptedFileByIDServiceImpl{
		config:         config,
		logger:         logger.With(zap.String("component", "get-encrypted-file-by-id-service")),
		getByIDUseCase: getByIDUseCase,
	}
}

// Execute retrieves an encrypted file by its ID and verifies ownership
func (s *getEncryptedFileByIDServiceImpl) Execute(
	ctx context.Context,
	id primitive.ObjectID,
) (*domain.EncryptedFile, error) {
	// Get the file using the use case
	file, err := s.getByIDUseCase.Execute(ctx, id)
	if err != nil {
		return nil, err
	}

	// Verify that the authenticated user has access to this file
	userID, ok := ctx.Value(constants.SessionFederatedUserID).(primitive.ObjectID)
	if ok && !userID.IsZero() && file.UserID != userID {
		s.logger.Warn("Unauthorized file access attempt",
			zap.String("file_id", id.Hex()),
			zap.String("file_owner", file.UserID.Hex()),
			zap.String("requester", userID.Hex()),
		)
		return nil, httperror.NewForForbiddenWithSingleField("message", "You do not have permission to access this file")
	}

	return file, nil
}
