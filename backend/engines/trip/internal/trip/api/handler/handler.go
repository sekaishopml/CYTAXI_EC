package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/sekaishopml/cytaxi/backend/engines/trip/internal/trip/application/port"
	"github.com/sekaishopml/cytaxi/backend/engines/trip/internal/trip/application/query"
	"github.com/sekaishopml/cytaxi/backend/engines/trip/internal/trip/domain/valueobject"
)

type Handler struct {
	service port.TripService
	logger  *slog.Logger
}

func New(service port.TripService, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok","service":"trip-engine"}`))
}

func (h *Handler) GetTrip(w http.ResponseWriter, r *http.Request) {
	tripID := valueobject.TripID(r.PathValue("trip_id"))
	result, err := h.service.GetTrip(r.Context(), query.GetTrip{TripID: tripID})
	if err != nil {
		http.Error(w, "trip not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, result.Trip)
}

func (h *Handler) GetTripHistory(w http.ResponseWriter, r *http.Request) {
	customerID := valueobject.CustomerID(r.PathValue("customer_id"))
	result, err := h.service.GetTripHistory(r.Context(), query.GetTripHistory{CustomerID: customerID})
	if err != nil {
		http.Error(w, "history not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (h *Handler) GetActiveTrips(w http.ResponseWriter, r *http.Request) {
	result, err := h.service.GetActiveTrips(r.Context(), query.GetActiveTrips{})
	if err != nil {
		http.Error(w, "active trips not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
