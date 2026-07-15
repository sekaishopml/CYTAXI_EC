package cmd

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/sekaishopml/cytaxi/backend/engines/pricing/internal/pricing/application/command"
	"github.com/sekaishopml/cytaxi/backend/engines/pricing/internal/pricing/application/port"
	"github.com/sekaishopml/cytaxi/backend/engines/pricing/internal/pricing/application/query"
	"github.com/sekaishopml/cytaxi/backend/engines/pricing/internal/pricing/domain/valueobject"
)

type PricingServer struct {
	service port.PricingService
	logger  *slog.Logger
}

func NewPricingServer(service port.PricingService, logger *slog.Logger) *PricingServer {
	return &PricingServer{service: service, logger: logger}
}

func (s *PricingServer) HandleEstimate(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TripID     string  `json:"trip_id"`
		DistanceKM float64 `json:"distance_km"`
		DurationSec int   `json:"duration_sec"`
		Region     string  `json:"region"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
		return
	}

	cmd := command.CalculateFare{
		TripID:      req.TripID,
		DistanceKM:  req.DistanceKM,
		DurationSec: req.DurationSec,
		Region:      req.Region,
	}

	fare, err := s.service.CalculateFare(r.Context(), cmd)
	if err != nil {
		s.logger.Error("fare estimate failed", "error", err)
		http.Error(w, `{"error":"estimate failed"}`, http.StatusInternalServerError)
		return
	}

	s.logger.Info("fare estimated", "trip_id", req.TripID, "total", fare.Components.Total.Amount)

	writeJSON(w, http.StatusOK, map[string]any{
		"status": "estimated",
		"fare": map[string]any{
			"fare_id":    fare.ID,
			"base":       fare.Components.BaseFare.Amount,
			"distance":   fare.Components.DistanceFare.Amount,
			"time":       fare.Components.TimeFare.Amount,
			"subtotal":   fare.Components.Subtotal.Amount,
			"total":      fare.Components.Total.Amount,
			"currency":   fare.Components.Total.Currency,
		},
	})
}

func (s *PricingServer) HandleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"status": "ok", "service": "pricing-engine"})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
