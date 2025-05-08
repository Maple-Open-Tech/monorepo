package unifiedhttp

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

// curl http://localhost:8000/papercloud/api/v1/version
type GetIncomePropertyEvaluatorVersionHTTPHandler struct {
	log *zap.Logger
}

func NewGetIncomePropertyEvaluatorVersionHTTPHandler(
	log *zap.Logger,
) *GetIncomePropertyEvaluatorVersionHTTPHandler {
	return &GetIncomePropertyEvaluatorVersionHTTPHandler{log}
}

type IncomePropertyEvaluatorVersionResponseIDO struct {
	Version string `json:"version"`
}

func (h *GetIncomePropertyEvaluatorVersionHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	response := IncomePropertyEvaluatorVersionResponseIDO{Version: "v1.0.0"}
	json.NewEncoder(w).Encode(response)
}

func (*GetIncomePropertyEvaluatorVersionHTTPHandler) Pattern() string {
	return "/papercloud/api/v1/version"
}
