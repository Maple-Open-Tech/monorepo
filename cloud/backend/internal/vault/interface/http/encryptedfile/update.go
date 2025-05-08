// cloud/backend/internal/vault/interface/http/encryptedfile/update.go
package encryptedfile

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
	svc "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/vault/service/encryptedfile"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/httperror"
)

// UpdateEncryptedFileHandler handles HTTP requests to update an encrypted file
type UpdateEncryptedFileHandler struct {
	config        *config.Configuration
	logger        *zap.Logger
	updateService svc.UpdateEncryptedFileService
}

// NewUpdateEncryptedFileHandler creates a new handler for file updates
func NewUpdateEncryptedFileHandler(
	config *config.Configuration,
	logger *zap.Logger,
	updateService svc.UpdateEncryptedFileService,
) *UpdateEncryptedFileHandler {
	return &UpdateEncryptedFileHandler{
		config:        config,
		logger:        logger.With(zap.String("handler", "update-encrypted-file")),
		updateService: updateService,
	}
}

// Pattern returns the URL pattern for this handler
func (h *UpdateEncryptedFileHandler) Pattern() string {
	return "PUT /vault/api/v1/encrypted-files/{id}"
}

// ServeHTTP handles HTTP requests
func (h *UpdateEncryptedFileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Extract file ID from URL path
	path := strings.Split(r.URL.Path, "/")
	if len(path) < 4 {
		httperror.ResponseError(w, httperror.NewForBadRequestWithSingleField("id", "File ID is required"))
		return
	}
	idStr := path[3]

	// Convert string ID to ObjectID
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		httperror.ResponseError(w, httperror.NewForBadRequestWithSingleField("id", "Invalid file ID format"))
		return
	}

	// Parse multipart form to get file and metadata
	if err := r.ParseMultipartForm(32 << 20); err != nil { // 32MB max
		h.logger.Error("Failed to parse multipart form", zap.Error(err))
		httperror.ResponseError(w, httperror.NewForBadRequestWithSingleField("content", "Invalid multipart form"))
		return
	}

	// Extract form fields
	encryptedMetadata := r.FormValue("encrypted_metadata")
	encryptedHash := r.FormValue("encrypted_hash")

	// Get file content (optional for updates)
	var fileContent io.Reader
	file, fileHeader, err := r.FormFile("encrypted_content")
	if err == nil && fileHeader != nil {
		fileContent = file
		defer file.Close()
	}

	// Call service to update the file
	result, err := h.updateService.Execute(
		ctx,
		id,
		encryptedMetadata,
		encryptedHash,
		fileContent,
	)

	if err != nil {
		h.logger.Error("Failed to update encrypted file", zap.Error(err))
		httperror.ResponseError(w, err)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := FileResponse{
		ID:                result.ID,
		UserID:            result.UserID,
		FileID:            result.FileID,
		EncryptedMetadata: result.EncryptedMetadata,
		EncryptionVersion: result.EncryptionVersion,
		EncryptedHash:     result.EncryptedHash,
		CreatedAt:         result.CreatedAt,
		ModifiedAt:        result.ModifiedAt,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error("Failed to encode response", zap.Error(err))
	}
}
