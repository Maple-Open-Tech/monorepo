// github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/maplesend/interface/http/me/verifyprofile.go
package me

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.uber.org/zap"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/maplesend/interface/http/middleware"
	svc_me "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/maplesend/service/me"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/httperror"
)

type PostVerifyProfileHTTPHandler struct {
	config     *config.Configuration
	logger     *zap.Logger
	dbClient   *mongo.Client
	service    svc_me.VerifyProfileService
	middleware middleware.Middleware
}

func NewPostVerifyProfileHTTPHandler(
	config *config.Configuration,
	logger *zap.Logger,
	dbClient *mongo.Client,
	service svc_me.VerifyProfileService,
	middleware middleware.Middleware,
) *PostVerifyProfileHTTPHandler {
	return &PostVerifyProfileHTTPHandler{
		config:     config,
		logger:     logger,
		dbClient:   dbClient,
		service:    service,
		middleware: middleware,
	}
}

func (*PostVerifyProfileHTTPHandler) Pattern() string {
	return "/maplesend/api/v1/verify-profile"
}

func (r *PostVerifyProfileHTTPHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Apply MaplesSend middleware before handling the request
	r.middleware.Attach(r.Execute)(w, req)
}

func (h *PostVerifyProfileHTTPHandler) unmarshalRequest(
	ctx context.Context,
	r *http.Request,
) (*svc_me.VerifyProfileRequestDTO, error) {
	// Initialize our structure which will store the parsed request data
	var requestData svc_me.VerifyProfileRequestDTO

	defer r.Body.Close()

	var rawJSON bytes.Buffer
	teeReader := io.TeeReader(r.Body, &rawJSON) // TeeReader allows you to read the JSON and capture it

	// Read the JSON string and convert it into our golang struct
	err := json.NewDecoder(teeReader).Decode(&requestData)
	if err != nil {
		h.logger.Error("decoding error",
			zap.Any("err", err),
			zap.String("json", rawJSON.String()),
		)
		return nil, httperror.NewForSingleField(http.StatusBadRequest, "non_field_error", "payload structure is wrong")
	}

	return &requestData, nil
}

func (h *PostVerifyProfileHTTPHandler) Execute(w http.ResponseWriter, r *http.Request) {
	// Set response content type
	w.Header().Set("Content-Type", "application/json")

	ctx := r.Context()

	req, err := h.unmarshalRequest(ctx, r)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	////
	//// Start the transaction.
	////

	session, err := h.dbClient.StartSession()
	if err != nil {
		h.logger.Error("start session error",
			zap.Any("error", err))
		httperror.ResponseError(w, err)
		return
	}
	defer session.EndSession(ctx)

	// Define a transaction function with a series of operations
	transactionFunc := func(sessCtx context.Context) (interface{}, error) {

		// Call service
		result, err := h.service.Execute(sessCtx, req)
		if err != nil {
			h.logger.Error("failed to verify profile",
				zap.Any("error", err))
			return nil, err
		}
		return result, nil
	}

	// Start a transaction
	result, txErr := session.WithTransaction(ctx, transactionFunc)
	if txErr != nil {
		h.logger.Error("session failed error",
			zap.Any("error", txErr))
		httperror.ResponseError(w, txErr)
		return
	}

	// Encode response
	resp := result.(*svc_me.VerifyProfileResponseDTO)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.logger.Error("failed to encode response",
			zap.Any("error", err))
		httperror.ResponseError(w, err)
		return
	}
}
