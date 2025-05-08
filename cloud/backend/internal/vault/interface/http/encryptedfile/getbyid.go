// cloud/backend/internal/vault/interface/http/encryptedfile/getbyid.go
package encryptedfile

import (
	"encoding/json"
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
	svc "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/vault/service/encryptedfile"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/httperror"
)

// GetEncryptedFileByIDHandler handles HTTP requests to get an encrypted file by ID
type GetEncryptedFileByIDHandler struct {
	config         *config.Configuration
	logger         *zap.Logger
	getByIDService svc.GetEncryptedFileByIDService
}

// NewGetEncryptedFileByIDHandler creates a new handler for getting a file by ID
func NewGetEncryptedFileByIDHandler(
	config *config.Configuration,
	logger *zap.Logger,
	getByIDService svc.GetEncryptedFileByIDService,
) *GetEncryptedFileByIDHandler {
	return &GetEncryptedFileByIDHandler{
		config:         config,
		logger:         logger.With(zap.String("handler", "get-encrypted-file-by-id")),
		getByIDService: getByIDService,
	}
}

// Pattern returns the URL pattern for this handler
func (h *GetEncryptedFileByIDHandler) Pattern() string {
	return "GET /vault/api/v1/encrypted-files/{id}"
}

// ServeHTTP handles HTTP requests
func (h *GetEncryptedFileByIDHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	// Call service to get the file
	file, err := h.getByIDService.Execute(ctx, id)
	if err != nil {
		h.logger.Error("Failed to get encrypted file by ID", zap.Error(err))
		httperror.ResponseError(w, err)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := FileResponse{
		ID:                file.ID,
		UserID:            file.UserID,
		FileID:            file.FileID,
		EncryptedMetadata: file.EncryptedMetadata,
		EncryptionVersion: file.EncryptionVersion,
		EncryptedHash:     file.EncryptedHash,
		CreatedAt:         file.CreatedAt,
		ModifiedAt:        file.ModifiedAt,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error("Failed to encode response", zap.Error(err))
	}
}
