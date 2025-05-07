// cloud/backend/internal/encryption/service/encryptedfile/auth.go
package encryptedfile

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config/constants"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/encryption/usecase/encryptedfile"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/httperror"
)

// FileAuthorizationMiddleware provides authorization checks for file operations
type FileAuthorizationMiddleware struct {
	logger             *zap.Logger
	getFileByIDUseCase encryptedfile.GetEncryptedFileByIDUseCase
}

// NewFileAuthorizationMiddleware creates a new file authorization middleware
func NewFileAuthorizationMiddleware(
	logger *zap.Logger,
	getFileByIDUseCase encryptedfile.GetEncryptedFileByIDUseCase,
) *FileAuthorizationMiddleware {
	return &FileAuthorizationMiddleware{
		logger:             logger.With(zap.String("component", "file-authorization-middleware")),
		getFileByIDUseCase: getFileByIDUseCase,
	}
}

// CheckFileAccess verifies that the authenticated user has permission to access the file
func (m *FileAuthorizationMiddleware) CheckFileAccess(ctx context.Context, fileID primitive.ObjectID) error {
	// Extract authenticated user ID from context
	userID, ok := ctx.Value(constants.SessionFederatedUserID).(primitive.ObjectID)
	if !ok || userID.IsZero() {
		m.logger.Error("User ID not found in context")
		return httperror.NewForUnauthorizedWithSingleField("message", "Authentication required")
	}

	// Get the file
	file, err := m.getFileByIDUseCase.Execute(ctx, fileID)
	if err != nil {
		m.logger.Error("Failed to get file for authorization check",
			zap.Error(err),
			zap.String("file_id", fileID.Hex()),
		)
		return err
	}

	if file == nil {
		return httperror.NewForBadRequestWithSingleField("id", "File not found")
	}

	// Check if the authenticated user is the owner of the file
	if file.UserID != userID {
		m.logger.Warn("Unauthorized file access attempt",
			zap.String("file_id", fileID.Hex()),
			zap.String("file_owner", file.UserID.Hex()),
			zap.String("requester", userID.Hex()),
		)
		return httperror.NewForForbiddenWithSingleField("message", "You do not have permission to access this file")
	}

	return nil
}
