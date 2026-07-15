package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/sekaishopml/cytaxi/backend/engines/pricing/internal/pricing/application/port"
	"github.com/sekaishopml/cytaxi/backend/engines/pricing/internal/pricing/application/query"
	"github.com/sekaishopml/cytaxi/backend/engines/pricing/internal/pricing/domain/valueobject"
)

type Handler struct {
	service port.PricingService
	logger  *slog.Logger
}

func New(service port.PricingService, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok","service":"pricing-engine"}`))
}

func (h *Handler) GetFare(w http.ResponseWriter, r *http.Request) {
	fareID := valueobject.FareID(r.PathValue("fare_id"))
	f, err := h.service.GetFare(r.Context(), query.GetFare{FareID: fareID})
	if err != nil {
		http.Error(w, "fare not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, f)
}

func (h *Handler) GetHistory(w http.ResponseWriter, r *http.Request) {
	tripID := r.PathValue("trip_id")
	fares, err := h.service.GetFareHistory(r.Context(), query.GetFareHistory{TripID: tripID})
	if err != nil {
		http.Error(w, "history not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, fares)
}

func (h *Handler) GetPromotions(w http.ResponseWriter, r *http.Request) {
	promos, err := h.service.GetPromotions(r.Context(), query.GetPromotions{Active: true})
	if err != nil {
		http.Error(w, "promotions not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, promos)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
