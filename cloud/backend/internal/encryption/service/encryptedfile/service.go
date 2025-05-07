// cloud/backend/internal/encryption/service/encryptedfile/service.go
package encryptedfile

import (
	"context"
	"errors"
	"io"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config/constants"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/encryption/usecase/encryptedfile"
)

// ErrUserIDRequired is returned when the user ID is missing.
var ErrUserIDRequired = errors.New("user ID is required")

// EncryptedFileService defines operations for working with encrypted files
type EncryptedFileService interface {
	Create(ctx context.Context, req *CreateFileRequest) (*FileResponse, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*FileResponse, error)
	GetByFileID(ctx context.Context, userID primitive.ObjectID, fileID string) (*FileResponse, error)
	Update(ctx context.Context, req *UpdateFileRequest) (*FileResponse, error)
	Delete(ctx context.Context, id primitive.ObjectID) error
	List(ctx context.Context, userID primitive.ObjectID) (*FilesListResponse, error)
	Download(ctx context.Context, id primitive.ObjectID) (io.ReadCloser, error)
	GetDownloadURL(ctx context.Context, id primitive.ObjectID, expiryDuration time.Duration) (*FileURLResponse, error)
}

type encryptedFileServiceImpl struct {
	config                    *config.Configuration
	logger                    *zap.Logger
	createFileUseCase         encryptedfile.CreateEncryptedFileUseCase
	getFileByIDUseCase        encryptedfile.GetEncryptedFileByIDUseCase
	getFileByFileIDUseCase    encryptedfile.GetEncryptedFileByFileIDUseCase
	updateFileUseCase         encryptedfile.UpdateEncryptedFileUseCase
	deleteFileUseCase         encryptedfile.DeleteEncryptedFileUseCase
	listFilesUseCase          encryptedfile.ListEncryptedFilesUseCase
	downloadFileUseCase       encryptedfile.DownloadEncryptedFileUseCase
	getFileDownloadURLUseCase encryptedfile.GetEncryptedFileDownloadURLUseCase
}

// NewEncryptedFileService creates a new instance of the encrypted file service
func NewEncryptedFileService(
	config *config.Configuration,
	logger *zap.Logger,
	createFileUseCase encryptedfile.CreateEncryptedFileUseCase,
	getFileByIDUseCase encryptedfile.GetEncryptedFileByIDUseCase,
	getFileByFileIDUseCase encryptedfile.GetEncryptedFileByFileIDUseCase,
	updateFileUseCase encryptedfile.UpdateEncryptedFileUseCase,
	deleteFileUseCase encryptedfile.DeleteEncryptedFileUseCase,
	listFilesUseCase encryptedfile.ListEncryptedFilesUseCase,
	downloadFileUseCase encryptedfile.DownloadEncryptedFileUseCase,
	getFileDownloadURLUseCase encryptedfile.GetEncryptedFileDownloadURLUseCase,
) EncryptedFileService {
	return &encryptedFileServiceImpl{
		config:                    config,
		logger:                    logger.With(zap.String("component", "encrypted-file-service")),
		createFileUseCase:         createFileUseCase,
		getFileByIDUseCase:        getFileByIDUseCase,
		getFileByFileIDUseCase:    getFileByFileIDUseCase,
		updateFileUseCase:         updateFileUseCase,
		deleteFileUseCase:         deleteFileUseCase,
		listFilesUseCase:          listFilesUseCase,
		downloadFileUseCase:       downloadFileUseCase,
		getFileDownloadURLUseCase: getFileDownloadURLUseCase,
	}
}

// Create handles creation of a new encrypted file
func (s *encryptedFileServiceImpl) Create(ctx context.Context, req *CreateFileRequest) (*FileResponse, error) {
	// Extract authenticated user ID from context if not provided
	if req.UserID.IsZero() {
		userID, ok := ctx.Value(constants.SessionFederatedUserID).(primitive.ObjectID)
		if !ok || userID.IsZero() {
			s.logger.Error("User ID not provided and not found in context")
			return nil, ErrUserIDRequired
		}
		req.UserID = userID
	}

	// Call use case to create file
	file, err := s.createFileUseCase.Execute(
		ctx,
		req.UserID,
		req.FileID,
		req.EncryptedMetadata,
		req.EncryptedHash,
		req.EncryptionVersion,
		req.EncryptedContent,
	)
	if err != nil {
		s.logger.Error("Failed to create encrypted file",
			zap.Error(err),
			zap.String("user_id", req.UserID.Hex()),
			zap.String("file_id", req.FileID),
		)
		return nil, err
	}

	// Convert domain entity to response DTO
	return DomainToFileResponse(file), nil
}

// GetByID retrieves an encrypted file by its ID
func (s *encryptedFileServiceImpl) GetByID(ctx context.Context, id primitive.ObjectID) (*FileResponse, error) {
	// Call use case to get file
	file, err := s.getFileByIDUseCase.Execute(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get encrypted file by ID",
			zap.Error(err),
			zap.String("id", id.Hex()),
		)
		return nil, err
	}

	// Convert domain entity to response DTO
	return DomainToFileResponse(file), nil
}

// GetByFileID retrieves an encrypted file by user ID and file ID
func (s *encryptedFileServiceImpl) GetByFileID(ctx context.Context, userID primitive.ObjectID, fileID string) (*FileResponse, error) {
	// Extract authenticated user ID from context if not provided
	if userID.IsZero() {
		contextUserID, ok := ctx.Value(constants.SessionFederatedUserID).(primitive.ObjectID)
		if !ok || contextUserID.IsZero() {
			s.logger.Error("User ID not provided and not found in context")
			return nil, ErrUserIDRequired
		}
		userID = contextUserID
	}

	// Call use case to get file
	file, err := s.getFileByFileIDUseCase.Execute(ctx, userID, fileID)
	if err != nil {
		s.logger.Error("Failed to get encrypted file by file ID",
			zap.Error(err),
			zap.String("user_id", userID.Hex()),
			zap.String("file_id", fileID),
		)
		return nil, err
	}

	// Convert domain entity to response DTO
	return DomainToFileResponse(file), nil
}

// Update handles updating an existing encrypted file
func (s *encryptedFileServiceImpl) Update(ctx context.Context, req *UpdateFileRequest) (*FileResponse, error) {
	// Call use case to update file
	file, err := s.updateFileUseCase.Execute(
		ctx,
		req.ID,
		req.EncryptedMetadata,
		req.EncryptedHash,
		req.EncryptedContent,
	)
	if err != nil {
		s.logger.Error("Failed to update encrypted file",
			zap.Error(err),
			zap.String("id", req.ID.Hex()),
		)
		return nil, err
	}

	// Convert domain entity to response DTO
	return DomainToFileResponse(file), nil
}

// Delete handles deletion of an encrypted file
func (s *encryptedFileServiceImpl) Delete(ctx context.Context, id primitive.ObjectID) error {
	// Call use case to delete file
	err := s.deleteFileUseCase.Execute(ctx, id)
	if err != nil {
		s.logger.Error("Failed to delete encrypted file",
			zap.Error(err),
			zap.String("id", id.Hex()),
		)
		return err
	}

	return nil
}

// List retrieves all encrypted files for a user
func (s *encryptedFileServiceImpl) List(ctx context.Context, userID primitive.ObjectID) (*FilesListResponse, error) {
	// Extract authenticated user ID from context if not provided
	if userID.IsZero() {
		contextUserID, ok := ctx.Value(constants.SessionFederatedUserID).(primitive.ObjectID)
		if !ok || contextUserID.IsZero() {
			s.logger.Error("User ID not provided and not found in context")
			return nil, ErrUserIDRequired
		}
		userID = contextUserID
	}

	// Call use case to list files
	files, err := s.listFilesUseCase.Execute(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to list encrypted files",
			zap.Error(err),
			zap.String("user_id", userID.Hex()),
		)
		return nil, err
	}

	// Convert domain entities to response DTOs
	response := &FilesListResponse{
		Files: make([]*FileResponse, len(files)),
	}
	for i, file := range files {
		response.Files[i] = DomainToFileResponse(file)
	}

	return response, nil
}

// Download retrieves the encrypted content of a file
func (s *encryptedFileServiceImpl) Download(ctx context.Context, id primitive.ObjectID) (io.ReadCloser, error) {
	// Call use case to download file
	content, err := s.downloadFileUseCase.Execute(ctx, id)
	if err != nil {
		s.logger.Error("Failed to download encrypted file",
			zap.Error(err),
			zap.String("id", id.Hex()),
		)
		return nil, err
	}

	return content, nil
}

// GetDownloadURL generates a presigned URL for direct download
func (s *encryptedFileServiceImpl) GetDownloadURL(ctx context.Context, id primitive.ObjectID, expiryDuration time.Duration) (*FileURLResponse, error) {
	// Call use case to get download URL
	url, err := s.getFileDownloadURLUseCase.Execute(ctx, id, expiryDuration)
	if err != nil {
		s.logger.Error("Failed to get encrypted file download URL",
			zap.Error(err),
			zap.String("id", id.Hex()),
		)
		return nil, err
	}

	// Create response with URL and expiry time
	response := &FileURLResponse{
		URL:       url,
		ExpiresAt: time.Now().Add(expiryDuration),
	}

	return response, nil
}
