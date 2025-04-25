package me

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/maplesend/interface/http/middleware"
)

// MeCombinedHandler routes requests to the appropriate handler based on HTTP method
type MeCombinedHandler struct {
	logger        *zap.Logger
	getHandler    *GetMeHTTPHandler
	updateHandler *PutUpdateMeHTTPHandler
	deleteHandler *DeleteMeHTTPHandler
	middleware    middleware.Middleware
}

func NewMeCombinedHandler(
	config *config.Configuration,
	logger *zap.Logger,
	getHandler *GetMeHTTPHandler,
	updateHandler *PutUpdateMeHTTPHandler,
	deleteHandler *DeleteMeHTTPHandler,
	middleware middleware.Middleware,
) *MeCombinedHandler {
	return &MeCombinedHandler{
		logger:        logger,
		getHandler:    getHandler,
		updateHandler: updateHandler,
		deleteHandler: deleteHandler,
		middleware:    middleware,
	}
}

func (h *MeCombinedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.middleware.Attach(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.getHandler.Execute(w, r)
		case http.MethodPut:
			h.updateHandler.Execute(w, r)
		case http.MethodDelete:
			h.deleteHandler.Execute(w, r)
		default:
			w.Header().Set("Allow", "GET, PUT, DELETE")
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})(w, r)
}

func (*MeCombinedHandler) Pattern() string {
	return "/maplesend/api/v1/me"
}
