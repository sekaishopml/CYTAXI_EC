package cmd

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/sekaishopml/cytaxi/backend/engines/trust/internal/trust/infrastructure/trust"
)

type TrustServer struct {
	manager *trust.Manager
	logger  *slog.Logger
}

func NewTrustServer(manager *trust.Manager, logger *slog.Logger) *TrustServer {
	return &TrustServer{manager: manager, logger: logger}
}

func (s *TrustServer) HandleRate(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TripID  string `json:"trip_id"`
		FromID  string `json:"from_id"`
		ToID    string `json:"to_id"`
		Score   int    `json:"score"`
		Comment string `json:"comment"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	if req.Score < 1 || req.Score > 5 {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "score must be 1-5"})
		return
	}

	rating := s.manager.SubmitRating(req.TripID, req.FromID, req.ToID, req.Score, req.Comment)
	s.logger.Info("rating submitted", "to", req.ToID, "score", req.Score)
	writeJSON(w, http.StatusCreated, rating)
}

func (s *TrustServer) HandleReport(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ReportedBy  string `json:"reported_by"`
		ReportedID  string `json:"reported_id"`
		TripID      string `json:"trip_id"`
		Type        string `json:"type"`
		Severity    string `json:"severity"`
		Description string `json:"description"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	incident := s.manager.ReportIncident(
		req.ReportedBy, req.ReportedID, req.TripID,
		trust.IncidentType(req.Type), req.Severity, req.Description,
	)
	s.logger.Info("incident reported", "by", req.ReportedBy, "against", req.ReportedID, "severity", req.Severity)
	writeJSON(w, http.StatusCreated, incident)
}

func (s *TrustServer) HandleGetTrustScore(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("user_id")
	ts := s.manager.GetTrustScore(userID)
	writeJSON(w, http.StatusOK, ts)
}

func (s *TrustServer) HandleAppeal(w http.ResponseWriter, r *http.Request) {
	var req struct {
		IncidentID string `json:"incident_id"`
		FromID     string `json:"from_id"`
		Reason     string `json:"reason"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	appeal := s.manager.FileAppeal(req.IncidentID, req.FromID, req.Reason)
	writeJSON(w, http.StatusCreated, appeal)
}

func (s *TrustServer) HandleResolveIncident(w http.ResponseWriter, r *http.Request) {
	var req struct {
		IncidentID  string `json:"incident_id"`
		Resolution  string `json:"resolution"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	incident, err := s.manager.ResolveIncident(req.IncidentID, req.Resolution)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, incident)
}

func (s *TrustServer) HandleGetIncidents(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID != "" {
		incidents := s.manager.GetIncidentsByUser(userID)
		writeJSON(w, http.StatusOK, incidents)
		return
	}

	open := s.manager.GetOpenIncidents()
	writeJSON(w, http.StatusOK, open)
}

func (s *TrustServer) HandleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"status": "ok", "service": "trust-safety"})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
