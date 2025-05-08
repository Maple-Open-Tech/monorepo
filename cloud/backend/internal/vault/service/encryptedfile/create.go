// cloud/backend/internal/vault/service/encryptedfile/create.go
package encryptedfile

import (
	"context"
	"fmt"
	"io"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config/constants"
	domain "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/vault/domain/encryptedfile"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/vault/usecase/encryptedfile"
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
	repository    domain.Repository
}

// NewCreateEncryptedFileService creates a new instance of the service
func NewCreateEncryptedFileService(
	config *config.Configuration,
	logger *zap.Logger,
	createUseCase encryptedfile.CreateEncryptedFileUseCase,
	repository domain.Repository,
) CreateEncryptedFileService {
	return &createEncryptedFileServiceImpl{
		config:        config,
		logger:        logger.With(zap.String("component", "create-encrypted-file-service")),
		createUseCase: createUseCase,
		repository:    repository,
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

	// Validate inputs - moved from usecase to service
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
	existingFile, err := s.repository.GetByFileID(ctx, userID, fileID)
	if err != nil {
		s.logger.Error("Failed to check for existing file",
			zap.String("userID", userID.Hex()),
			zap.String("fileID", fileID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to check for existing file: %w", err)
	}

	if existingFile != nil {
		return nil, httperror.NewForBadRequestWithSingleField("file_id", "A file with this ID already exists")
	}

	// Set default encryption version if not provided
	if encryptionVersion == "" {
		encryptionVersion = "1.0"
	}

	// Create a new encrypted file entity
	now := time.Now()
	file := &domain.EncryptedFile{
		ID:                primitive.NewObjectID(),
		UserID:            userID,
		FileID:            fileID,
		EncryptedMetadata: encryptedMetadata,
		EncryptedHash:     encryptedHash,
		EncryptionVersion: encryptionVersion,
		CreatedAt:         now,
		ModifiedAt:        now,
	}

	// Delegate to the use case for repository interaction
	err = s.createUseCase.Execute(ctx, file, encryptedContent)
	if err != nil {
		s.logger.Error("Failed to create encrypted file",
			zap.String("userID", userID.Hex()),
			zap.String("fileID", fileID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to create encrypted file: %w", err)
	}

	s.logger.Info("Successfully created encrypted file",
		zap.String("id", file.ID.Hex()),
		zap.String("userID", userID.Hex()),
		zap.String("fileID", fileID),
	)

	return file, nil
}
