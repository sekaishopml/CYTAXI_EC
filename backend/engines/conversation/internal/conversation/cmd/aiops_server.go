package cmd

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/sekaishopml/cytaxi/backend/engines/conversation/internal/conversation/infrastructure/aiops"
)

type AIOpsServer struct {
	manager *aiops.Manager
	logger  *slog.Logger
}

func NewAIOpsServer(manager *aiops.Manager, logger *slog.Logger) *AIOpsServer {
	return &AIOpsServer{manager: manager, logger: logger}
}

func (s *AIOpsServer) HandleHealth(w http.ResponseWriter, r *http.Request) {
	status := s.manager.GetStatus()
	status["status"] = "ok"
	status["service"] = "ai-operations"
	writeJSON(w, http.StatusOK, status)
}

func (s *AIOpsServer) HandleRecommendations(w http.ResponseWriter, r *http.Request) {
	recs := s.manager.GetRecommendations()
	writeJSON(w, http.StatusOK, map[string]any{"recommendations": recs, "count": len(recs)})
}

func (s *AIOpsServer) HandleIncidents(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var req struct {
			Type        string `json:"type"`
			Severity    string `json:"severity"`
			Source      string `json:"source"`
			Description string `json:"description"`
		}
		json.NewDecoder(r.Body).Decode(&req)
		incident := s.manager.DetectIncident(req.Type, req.Severity, req.Source, req.Description)
		s.logger.Info("incident detected", "id", incident.ID, "type", req.Type)

		// Return runbook if available
		runbook := s.manager.GetRunbook(req.Type)
		writeJSON(w, http.StatusCreated, map[string]any{
			"incident": incident,
			"runbook":  runbook,
		})
		return
	}

	incidents := s.manager.GetIncidents()
	writeJSON(w, http.StatusOK, map[string]any{"incidents": incidents, "count": len(incidents)})
}

func (s *AIOpsServer) HandleRunbooks(w http.ResponseWriter, r *http.Request) {
	incidentType := r.URL.Query().Get("type")
	if incidentType != "" {
		runbook := s.manager.GetRunbook(incidentType)
		writeJSON(w, http.StatusOK, runbook)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"message": "Use ?type=service_down|high_latency|payment_failure|matching_failure"})
}

func (s *AIOpsServer) HandleStatus(w http.ResponseWriter, r *http.Request) {
	status := s.manager.GetStatus()
	writeJSON(w, http.StatusOK, status)
}

func (s *AIOpsServer) HandleAccept(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RecommendationID string `json:"recommendation_id"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	if err := s.manager.AcceptRecommendation(req.RecommendationID); err != nil {
		writeJSON(w, http.StatusNotFound, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"status": "accepted"})
}

func (s *AIOpsServer) HandleKnowledge(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	entries := s.manager.GetKnowledge(query)
	writeJSON(w, http.StatusOK, map[string]any{"articles": entries, "count": len(entries)})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
