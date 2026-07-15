package cmd

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

type TripServer struct {
	service port.TripService
	logger  *slog.Logger
}

func NewTripServer(service port.TripService, logger *slog.Logger) *TripServer {
	return &TripServer{service: service, logger: logger}
}

func (s *TripServer) HandleCreateTrip(w http.ResponseWriter, r *http.Request) {
	var req struct {
		CustomerID   string  `json:"customer_id"`
		Phone        string  `json:"phone"`
		PassengerName string `json:"passenger_name"`
		OriginAddr   string  `json:"origin_address"`
		OriginLat    float64 `json:"origin_lat"`
		OriginLng    float64 `json:"origin_lng"`
		DestAddr     string  `json:"dest_address"`
		DestLat      float64 `json:"dest_lat"`
		DestLng      float64 `json:"dest_lng"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
		return
	}

	cmd := command.CreateTrip{
		CustomerID: valueobject.CustomerID(req.CustomerID),
		Passenger: passenger.Passenger{
			ID:    valueobject.CustomerID(req.CustomerID),
			Phone: req.Phone,
			Name:  req.PassengerName,
		},
		Pickup: stop.NewStop(req.OriginAddr, valueobject.Coordinates{Lat: req.OriginLat, Lng: req.OriginLng}),
		Destination: destination.Destination{
			Address:  req.DestAddr,
			Location: valueobject.Coordinates{Lat: req.DestLat, Lng: req.DestLng},
		},
	}

	if err := s.service.Create(r.Context(), cmd); err != nil {
		s.logger.Error("create trip failed", "error", err)
		http.Error(w, `{"error":"trip creation failed"}`, http.StatusInternalServerError)
		return
	}

	s.logger.Info("trip created", "customer_id", req.CustomerID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"status": "created",
		"message": "Trip request received",
	})
}

func (s *TripServer) HandleGetTrip(w http.ResponseWriter, r *http.Request) {
	tripID := r.PathValue("trip_id")
	result, err := s.service.GetTrip(r.Context(), query.GetTrip{TripID: valueobject.TripID(tripID)})
	if err != nil {
		http.Error(w, `{"error":"trip not found"}`, http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, result.Trip)
}

func (s *TripServer) HandleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"status": "ok", "service": "trip-engine"})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
