// github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/interface/http/gateway/register.go
package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	_ "time/tzdata"

	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/interface/http/middleware"
	sv_gateway "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/service/gateway"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/httperror"
)

type GatewayUserRegisterHTTPHandler struct {
	logger     *zap.Logger
	dbClient   *mongo.Client
	service    sv_gateway.GatewayUserRegisterService
	middleware middleware.Middleware
}

func NewGatewayUserRegisterHTTPHandler(
	logger *zap.Logger,
	dbClient *mongo.Client,
	service sv_gateway.GatewayUserRegisterService,
	middleware middleware.Middleware,
) *GatewayUserRegisterHTTPHandler {
	return &GatewayUserRegisterHTTPHandler{
		logger:     logger,
		dbClient:   dbClient,
		service:    service,
		middleware: middleware,
	}
}

func (*GatewayUserRegisterHTTPHandler) Pattern() string {
	return "POST /ipe/api/v1/register"
}

func (r *GatewayUserRegisterHTTPHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Apply MaplesSend middleware before handling the request
	r.middleware.Attach(r.Execute)(w, req)
}

func (h *GatewayUserRegisterHTTPHandler) unmarshalRegisterCustomerRequest(
	ctx context.Context,
	r *http.Request,
) (*sv_gateway.RegisterCustomerRequestIDO, error) {
	// Initialize our array which will store all the results from the remote server.
	var requestData sv_gateway.RegisterCustomerRequestIDO

	defer r.Body.Close()

	var rawJSON bytes.Buffer
	teeReader := io.TeeReader(r.Body, &rawJSON) // TeeReader allows you to read the JSON and capture it

	// Read the JSON string and convert it into our golang stuct else we need
	// to send a `400 Bad Request` errror message back to the client,
	err := json.NewDecoder(teeReader).Decode(&requestData) // [1]
	if err != nil {
		h.logger.Error("decoding error",
			zap.Any("err", err),
			zap.String("json", rawJSON.String()),
		)
		return nil, httperror.NewForSingleField(http.StatusBadRequest, "non_field_error", "payload structure is wrong")
	}

	// Defensive Code: For security purposes we need to remove all whitespaces from the email and lower the characters.
	requestData.Email = strings.ToLower(requestData.Email)
	requestData.Email = strings.ReplaceAll(requestData.Email, " ", "")

	return &requestData, nil
}

func (h *GatewayUserRegisterHTTPHandler) Execute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	data, err := h.unmarshalRegisterCustomerRequest(ctx, r)
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
	transactionFunc := func(sessCtx context.Context) (any, error) {
		resp, err := h.service.Execute(sessCtx, data)
		if err != nil {
			return nil, err
		}
		return resp, nil
	}

	// Start a transaction
	result, err := session.WithTransaction(ctx, transactionFunc)
	if err != nil {
		h.logger.Error("session failed error",
			zap.Any("error", err))
		httperror.ResponseError(w, err)
		return
	}

	resp := result.(*sv_gateway.RegisterCustomerResponseIDO)

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
