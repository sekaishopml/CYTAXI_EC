package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/sekaishopml/cytaxi/backend/engines/trip/internal/trip/application/command"
	"github.com/sekaishopml/cytaxi/backend/engines/trip/internal/trip/application/port"
	"github.com/sekaishopml/cytaxi/backend/engines/trip/internal/trip/application/query"
	"github.com/sekaishopml/cytaxi/backend/engines/trip/internal/trip/domain/destination"
	"github.com/sekaishopml/cytaxi/backend/engines/trip/internal/trip/domain/passenger"
	"github.com/sekaishopml/cytaxi/backend/engines/trip/internal/trip/domain/stop"
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

// CreateTripRequest es el body esperado desde el frontend
type CreateTripRequest struct {
	CustomerID     string  `json:"customer_id"`
	Phone          string  `json:"phone"`
	PassengerName  string  `json:"passenger_name"`
	OriginAddress  string  `json:"origin_address"`
	OriginLat      float64 `json:"origin_lat"`
	OriginLng      float64 `json:"origin_lng"`
	DestAddress    string  `json:"dest_address"`
	DestLat        float64 `json:"dest_lat"`
	DestLng        float64 `json:"dest_lng"`
}

// HandleCreateTrip POST /trip/request — Endpoint principal para crear un viaje
func (h *Handler) HandleCreateTrip(w http.ResponseWriter, r *http.Request) {
	var req CreateTripRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid json", "details": err.Error()})
		return
	}

	customerID := req.CustomerID
	if customerID == "" && req.Phone != "" {
		customerID = "cust_" + req.Phone
	}
	if customerID == "" {
		customerID = "cust_anonymous"
	}

	cmd := command.CreateTrip{
		CustomerID: valueobject.CustomerID(customerID),
		Passenger: passenger.Passenger{
			ID:    valueobject.CustomerID(customerID),
			Phone: req.Phone,
			Name:  req.PassengerName,
		},
		Pickup: stop.NewStop(req.OriginAddress, valueobject.Coordinates{Lat: req.OriginLat, Lng: req.OriginLng}),
		Destination: destination.Destination{
			Address:  req.DestAddress,
			Location: valueobject.Coordinates{Lat: req.DestLat, Lng: req.DestLng},
		},
	}

	t, err := h.service.Create(r.Context(), cmd)
	if err != nil {
		h.logger.Error("create trip failed", "error", err)
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "could not create trip", "details": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"trip_id":     string(t.ID),
		"status":      string(t.Status),
		"customer_id": string(t.CustomerID),
		"message":     "Trip created",
	})
}

// CancelTripRequest body para cancelar un viaje
type CancelTripRequest struct {
	TripID string `json:"trip_id"`
	Reason string `json:"reason"`
	By     string `json:"by"`
}

// HandleCancelTrip POST /trips/{trip_id}/cancel o POST /trip/cancel
func (h *Handler) HandleCancelTrip(w http.ResponseWriter, r *http.Request) {
	var req CancelTripRequest
	tripID := r.PathValue("trip_id")
	json.NewDecoder(r.Body).Decode(&req)
	if tripID == "" { tripID = req.TripID }
	if tripID == "" { writeJSON(w, http.StatusBadRequest, map[string]any{"error": "trip_id required"}); return }
	if req.Reason == "" { req.Reason = "user_cancelled" }
	if req.By == "" { req.By = "customer" }

	if err := h.service.Cancel(r.Context(), command.CancelTrip{
		TripID: valueobject.TripID(tripID),
		Reason: req.Reason,
		By:     req.By,
	}); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"status": "cancelled", "trip_id": tripID})
}

// ChangeDestRequest body
type ChangeDestRequest struct {
	TripID      string  `json:"trip_id"`
	DestAddress string  `json:"dest_address"`
	DestLat     float64 `json:"dest_lat"`
	DestLng     float64 `json:"dest_lng"`
}

// HandleChangeDestination PUT /trips/{trip_id}/destination
func (h *Handler) HandleChangeDestination(w http.ResponseWriter, r *http.Request) {
	var req ChangeDestRequest
	tripID := r.PathValue("trip_id")
	json.NewDecoder(r.Body).Decode(&req)
	if tripID == "" { tripID = req.TripID }
	if tripID == "" { writeJSON(w, http.StatusBadRequest, map[string]any{"error": "trip_id required"}); return }

	if err := h.service.ChangeDestination(r.Context(), command.ChangeDestination{
		TripID: valueobject.TripID(tripID),
		Destination: destination.Destination{
			Address:  req.DestAddress,
			Location: valueobject.Coordinates{Lat: req.DestLat, Lng: req.DestLng},
		},
	}); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"status": "destination_changed", "trip_id": tripID})
}

// RejectDriverRequest body
type RejectDriverRequest struct {
	TripID string `json:"trip_id"`
	Reason string `json:"reason"`
}

// HandleRejectDriver POST /trip/reject-driver
func (h *Handler) HandleRejectDriver(w http.ResponseWriter, r *http.Request) {
	var req RejectDriverRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil { writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid json"}); return }
	if req.TripID == "" { writeJSON(w, http.StatusBadRequest, map[string]any{"error": "trip_id required"}); return }
	if req.Reason == "" { req.Reason = "rejected_by_customer" }

	if err := h.service.Reject(r.Context(), command.RejectTrip{
		TripID: valueobject.TripID(req.TripID),
		Reason: req.Reason,
	}); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"status": "driver_rejected", "trip_id": req.TripID})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
