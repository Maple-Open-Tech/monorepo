package unifiedhttp

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type GetHealthCheckHTTPHandler struct {
	logger *slog.Logger
}

func NewGetHealthCheckHTTPHandler(
	logger *slog.Logger,
) *GetHealthCheckHTTPHandler {
	return &GetHealthCheckHTTPHandler{logger}
}

type HealthCheckResponseIDO struct {
	Status string `json:"status"`
}

func (h *GetHealthCheckHTTPHandler) Execute(w http.ResponseWriter, r *http.Request) {

	response := HealthCheckResponseIDO{Status: "running"}
	json.NewEncoder(w).Encode(response)
}
