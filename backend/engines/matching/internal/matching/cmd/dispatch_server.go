package cmd

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/sekaishopml/cytaxi/backend/engines/matching/internal/matching/infrastructure/dispatch"
)

type DispatchServer struct {
	manager *dispatch.Manager
	logger  *slog.Logger
}

func NewDispatchServer(manager *dispatch.Manager, logger *slog.Logger) *DispatchServer {
	// Add default zones
	manager.AddZone(dispatch.DispatchZone{ID: "zone_a", Name: "Downtown", Lat: -0.18, Lng: -78.47, RadiusKM: 3, Active: true, Priority: 10})
	manager.AddZone(dispatch.DispatchZone{ID: "zone_b", Name: "Airport Area", Lat: -0.13, Lng: -78.37, RadiusKM: 5, Active: true, Priority: 7})
	return &DispatchServer{manager: manager, logger: logger}
}

func (s *DispatchServer) HandleStart(w http.ResponseWriter, r *http.Request) {
	var req dispatch.DispatchRequest
	json.NewDecoder(r.Body).Decode(&req)

	result := s.manager.StartDispatch(req)

	candidates := s.manager.GetCandidatesForZone(result.Zone, result.MaxRadius)
	ranked := s.manager.ScoreCandidates(result, candidates)
	best := s.manager.GetBestCandidate(ranked)

	s.logger.Info("dispatch started", "id", result.ID, "zone", result.Zone, "candidates", len(ranked))

	writeJSON(w, http.StatusOK, map[string]any{
		"dispatch":   result,
		"candidates": ranked,
		"best":       best,
		"zone":       "zone_a",
	})
}

func (s *DispatchServer) HandleRetry(w http.ResponseWriter, r *http.Request) {
	var req struct {
		DispatchID string `json:"dispatch_id"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	result, err := s.manager.RetryDispatch(req.DispatchID)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error(), "status": "max_retries"})
		return
	}

	s.logger.Info("dispatch retry", "id", req.DispatchID, "status", result.Status)
	writeJSON(w, http.StatusOK, result)
}

func (s *DispatchServer) HandleGetStatus(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("dispatch_id")
	req, err := s.manager.GetRequest(id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, req)
}

func (s *DispatchServer) HandleGetCandidates(w http.ResponseWriter, r *http.Request) {
	zone := r.URL.Query().Get("zone")
	if zone == "" { zone = "zone_a" }
	candidates := s.manager.GetCandidatesForZone(zone, 5000)
	writeJSON(w, http.StatusOK, map[string]any{"zone": zone, "candidates": candidates})
}

func (s *DispatchServer) HandleGetMetrics(w http.ResponseWriter, r *http.Request) {
	metrics := s.manager.GetMetrics()
	writeJSON(w, http.StatusOK, metrics)
}

func (s *DispatchServer) HandleHealth(w http.ResponseWriter, r *http.Request) {
	m := s.manager.GetMetrics()
	writeJSON(w, http.StatusOK, map[string]any{
		"status": "ok", "service": "intelligent-dispatch",
		"total_requests": m.TotalRequests,
		"acceptance_rate": fmt.Sprintf("%.1f%%", m.AcceptanceRate),
	})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
