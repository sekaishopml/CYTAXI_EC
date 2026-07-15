package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/sekaishopml/cytaxi/backend/engines/matching/internal/matching/application/port"
	"github.com/sekaishopml/cytaxi/backend/engines/matching/internal/matching/application/query"
	"github.com/sekaishopml/cytaxi/backend/engines/matching/internal/matching/domain/valueobject"
)

type Handler struct {
	service port.MatchingService
	logger  *slog.Logger
}

func New(service port.MatchingService, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok","service":"matching-engine"}`))
}

func (h *Handler) GetMatching(w http.ResponseWriter, r *http.Request) {
	id := valueobject.MatchingID(r.PathValue("matching_id"))
	m, err := h.service.GetMatching(r.Context(), query.GetMatching{MatchingID: id})
	if err != nil {
		http.Error(w, "matching not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, m)
}

func (h *Handler) GetCandidates(w http.ResponseWriter, r *http.Request) {
	id := valueobject.MatchingID(r.PathValue("matching_id"))
	cs, err := h.service.GetCandidates(r.Context(), query.GetCandidates{MatchingID: id})
	if err != nil {
		http.Error(w, "candidates not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, cs)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
