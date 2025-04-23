package unifiedhttp

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type GetVersionHTTPHandler struct {
	logger *slog.Logger
}

func NewGetVersionHTTPHandler(
	logger *slog.Logger,
) *GetVersionHTTPHandler {
	return &GetVersionHTTPHandler{logger}
}

type VersionResponseIDO struct {
	Version string `json:"version"`
}

func (h *GetVersionHTTPHandler) Execute(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("Version requested")

	response := VersionResponseIDO{Version: "1.0"}
	json.NewEncoder(w).Encode(response)
}
