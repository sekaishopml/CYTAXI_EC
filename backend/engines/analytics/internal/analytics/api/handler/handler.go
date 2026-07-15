package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/sekaishopml/cytaxi/backend/engines/analytics/internal/analytics/application/port"
	"github.com/sekaishopml/cytaxi/backend/engines/analytics/internal/analytics/application/query"
	"github.com/sekaishopml/cytaxi/backend/engines/analytics/internal/analytics/domain/valueobject"
)

type Handler struct {
	service port.AnalyticsService
	logger  *slog.Logger
}

func New(service port.AnalyticsService, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok","service":"analytics-engine"}`))
}

func (h *Handler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	d, err := h.service.GetDashboard(r.Context(), query.GetDashboard{})
	if err != nil {
		http.Error(w, "dashboard not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, d)
}

func (h *Handler) GetMetrics(w http.ResponseWriter, r *http.Request) {
	m, err := h.service.GetMetrics(r.Context(), query.GetMetrics{DateRange: valueobject.DateRange{}})
	if err != nil {
		http.Error(w, "metrics not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, m)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
