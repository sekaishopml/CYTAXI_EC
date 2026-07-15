package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/sekaishopml/cytaxi/backend/engines/admin/internal/admin/application/port"
	"github.com/sekaishopml/cytaxi/backend/engines/admin/internal/admin/application/query"
)

type Handler struct {
	service port.AdminService
	logger  *slog.Logger
}

func New(service port.AdminService, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok","service":"admin-engine"}`))
}

func (h *Handler) GetRoles(w http.ResponseWriter, r *http.Request) {
	roles, err := h.service.GetRoles(r.Context(), query.GetRoles{})
	if err != nil {
		http.Error(w, "roles not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, roles)
}

func (h *Handler) GetFeatureFlags(w http.ResponseWriter, r *http.Request) {
	flags, err := h.service.GetFeatureFlags(r.Context(), query.GetFeatureFlags{})
	if err != nil {
		http.Error(w, "feature flags not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, flags)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
