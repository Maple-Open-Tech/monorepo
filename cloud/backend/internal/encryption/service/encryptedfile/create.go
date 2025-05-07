// cloud/backend/internal/encryption/service/encryptedfile/create.go
package encryptedfile

import (
	"context"
	"io"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config/constants"
	domain "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/encryption/domain/encryptedfile"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/encryption/usecase/encryptedfile"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/httperror"
)

// CreateEncryptedFileService defines operations for creating an encrypted file
type CreateEncryptedFileService interface {
	Execute(ctx context.Context, userID primitive.ObjectID, fileID string, encryptedMetadata string, encryptedHash string, encryptionVersion string, encryptedContent io.Reader) (*domain.EncryptedFile, error)
}

type createEncryptedFileServiceImpl struct {
	config        *config.Configuration
	logger        *zap.Logger
	createUseCase encryptedfile.CreateEncryptedFileUseCase
}

// NewCreateEncryptedFileService creates a new instance of the service
func NewCreateEncryptedFileService(
	config *config.Configuration,
	logger *zap.Logger,
	createUseCase encryptedfile.CreateEncryptedFileUseCase,
) CreateEncryptedFileService {
	return &createEncryptedFileServiceImpl{
		config:        config,
		logger:        logger.With(zap.String("component", "create-encrypted-file-service")),
		createUseCase: createUseCase,
	}
}

// Execute handles the creation of a new encrypted file
func (s *createEncryptedFileServiceImpl) Execute(
	ctx context.Context,
	userID primitive.ObjectID,
	fileID string,
	encryptedMetadata string,
	encryptedHash string,
	encryptionVersion string,
	encryptedContent io.Reader,
) (*domain.EncryptedFile, error) {
	// Extract authenticated user ID from context if not provided
	if userID.IsZero() {
		contextUserID, ok := ctx.Value(constants.SessionFederatedUserID).(primitive.ObjectID)
		if !ok || contextUserID.IsZero() {
			s.logger.Error("User ID not provided and not found in context")
			return nil, httperror.NewForBadRequestWithSingleField("user_id", "User ID is required")
		}
		userID = contextUserID
	}

	// Delegate to the use case
	return s.createUseCase.Execute(
		ctx,
		userID,
		fileID,
		encryptedMetadata,
		encryptedHash,
		encryptionVersion,
		encryptedContent,
	)
}
