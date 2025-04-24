package gateway

import (
	"net/http"
	_ "time/tzdata"

	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/v2/mongo"

	sv_gateway "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/maplesend/service/gateway"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/httperror"
)

type GatewayLogoutHTTPHandler struct {
	logger   *zap.Logger
	dbClient *mongo.Client
	service  sv_gateway.GatewayLogoutService
}

func NewGatewayLogoutHTTPHandler(
	logger *zap.Logger,
	dbClient *mongo.Client,
	service sv_gateway.GatewayLogoutService,
) *GatewayLogoutHTTPHandler {
	return &GatewayLogoutHTTPHandler{
		logger:   logger,
		dbClient: dbClient,
		service:  service,
	}
}

func (h *GatewayLogoutHTTPHandler) Execute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if err := h.service.Execute(ctx); err != nil {
		httperror.ResponseError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)

}
