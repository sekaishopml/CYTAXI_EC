package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/sekaishopml/cytaxi/backend/engines/trust/internal/trust/application/port"
	"github.com/sekaishopml/cytaxi/backend/engines/trust/internal/trust/application/query"
	"github.com/sekaishopml/cytaxi/backend/engines/trust/internal/trust/domain/valueobject"
)

type Handler struct {
	service port.TrustService
	logger  *slog.Logger
}

func New(service port.TrustService, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok","service":"trust-engine"}`))
}

func (h *Handler) GetIdentity(w http.ResponseWriter, r *http.Request) {
	id := valueobject.IdentityID(r.PathValue("identity_id"))
	identity, err := h.service.GetIdentity(r.Context(), query.GetIdentity{IdentityID: id})
	if err != nil {
		http.Error(w, "identity not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, identity)
}

func (h *Handler) GetTrustScore(w http.ResponseWriter, r *http.Request) {
	id := valueobject.IdentityID(r.PathValue("identity_id"))
	tp, err := h.service.GetTrustScore(r.Context(), query.GetTrustScore{IdentityID: id})
	if err != nil {
		http.Error(w, "trust score not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, tp)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
