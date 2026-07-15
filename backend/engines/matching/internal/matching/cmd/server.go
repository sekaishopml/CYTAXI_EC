package cmd

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"

	"github.com/sekaishopml/cytaxi/backend/engines/matching/internal/matching/application/command"
	"github.com/sekaishopml/cytaxi/backend/engines/matching/internal/matching/application/port"
	"github.com/sekaishopml/cytaxi/backend/engines/matching/internal/matching/domain/valueobject"
)

type MatchingServer struct {
	service port.MatchingService
	logger  *slog.Logger
}

func NewMatchingServer(service port.MatchingService, logger *slog.Logger) *MatchingServer {
	return &MatchingServer{service: service, logger: logger}
}

type MockCandidate struct {
	DriverID   string  `json:"driver_id"`
	Name       string  `json:"name"`
	Distance   float64 `json:"distance_meters"`
	ETA        int     `json:"eta_seconds"`
	Score      float64 `json:"score"`
	Vehicle    string  `json:"vehicle"`
	Plate      string  `json:"plate"`
	Rating     float64 `json:"rating"`
}

func (s *MatchingServer) HandleStartMatching(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TripID    string  `json:"trip_id"`
		PickupLat float64 `json:"pickup_lat"`
		PickupLng float64 `json:"pickup_lng"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	cmd := command.StartMatching{
		TripID:    valueobject.TripID(req.TripID),
		PickupLat: req.PickupLat,
		PickupLng: req.PickupLng,
		Strategy:  "balanced",
	}

	m, err := s.service.StartMatching(r.Context(), cmd)
	if err != nil {
		s.logger.Error("matching start failed", "error", err)
		http.Error(w, `{"error":"matching failed"}`, http.StatusInternalServerError)
		return
	}

	candidates := generateMockCandidates(3)

	s.logger.Info("matching started", "trip_id", req.TripID, "candidates", len(candidates))

	writeJSON(w, http.StatusOK, map[string]any{
		"status":      "searching",
		"matching_id": m.ID,
		"candidates":  candidates,
	})
}

func (s *MatchingServer) HandleGetCandidates(w http.ResponseWriter, r *http.Request) {
	id := valueobject.MatchingID(r.PathValue("matching_id"))
	candidates := generateMockCandidates(3)

	s.logger.Info("candidates list", "matching_id", id, "count", len(candidates))
	writeJSON(w, http.StatusOK, map[string]any{
		"matching_id": id,
		"candidates":  candidates,
	})
}

func (s *MatchingServer) HandleSelectDriver(w http.ResponseWriter, r *http.Request) {
	var req struct {
		MatchingID string `json:"matching_id"`
		DriverID   string `json:"driver_id"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	s.logger.Info("driver selected", "matching_id", req.MatchingID, "driver_id", req.DriverID)

	writeJSON(w, http.StatusOK, map[string]any{
		"status":    "assigned",
		"driver_id": req.DriverID,
		"message":   "Driver assigned successfully",
	})
}

func (s *MatchingServer) HandleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"status": "ok", "service": "matching-engine"})
}

func generateMockCandidates(count int) []MockCandidate {
	names := []string{"Carlos M.", "Ana P.", "Luis R.", "María G.", "José V."}
	vehicles := []string{"Toyota Corolla", "Hyundai Elantra", "Kia Rio", "Chevrolet Aveo", "Nissan Versa"}
	plates := []string{"ABC-1234", "XYZ-5678", "DEF-9012", "GHI-3456", "JKL-7890"}

	var result []MockCandidate
	for i := 0; i < count && i < len(names); i++ {
		result = append(result, MockCandidate{
			DriverID: fmt.Sprintf("drv_%d", 1000+i),
			Name:     names[i],
			Distance: 300 + rand.Float64()*2000,
			ETA:      60 + rand.Intn(300),
			Score:    70 + rand.Float64()*30,
			Vehicle:  vehicles[i],
			Plate:    plates[i],
			Rating:   4.0 + rand.Float64(),
		})
	}
	return result
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
