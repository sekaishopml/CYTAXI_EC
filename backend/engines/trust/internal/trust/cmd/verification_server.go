package cmd

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/sekaishopml/cytaxi/backend/engines/trust/internal/trust/infrastructure/kyc"
)

type VerificationServer struct {
	manager *kyc.Manager
	logger  *slog.Logger
}

func NewVerificationServer(manager *kyc.Manager, logger *slog.Logger) *VerificationServer {
	return &VerificationServer{manager: manager, logger: logger}
}

func (s *VerificationServer) HandleStart(w http.ResponseWriter, r *http.Request) {
	var req struct {
		DriverID   string `json:"driver_id"`
		DriverName string `json:"driver_name"`
		Phone      string `json:"phone"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	v := s.manager.StartVerification(req.DriverID, req.DriverName, req.Phone)
	s.logger.Info("verification started", "ver_id", v.ID, "driver_id", req.DriverID)
	writeJSON(w, http.StatusCreated, v)
}

func (s *VerificationServer) HandleUploadDocument(w http.ResponseWriter, r *http.Request) {
	var req struct {
		VerificationID string            `json:"verification_id"`
		Type           kyc.DocumentType  `json:"type"`
		URL            string            `json:"url"`
		Filename       string            `json:"filename"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	doc := kyc.Document{
		ID: fmt.Sprintf("doc_%d", time.Now().UnixNano()),
		DriverID: "", Type: req.Type, URL: req.URL, Filename: req.Filename,
	}

	if err := s.manager.AddDocument(req.VerificationID, doc); err != nil {
		writeJSON(w, http.StatusNotFound, map[string]any{"error": err.Error()})
		return
	}

	s.logger.Info("document uploaded", "ver_id", req.VerificationID, "type", req.Type)
	writeJSON(w, http.StatusOK, map[string]any{"status": "uploaded", "document": doc})
}

func (s *VerificationServer) HandleApprove(w http.ResponseWriter, r *http.Request) {
	var req struct {
		VerificationID string `json:"verification_id"`
		Reviewer       string `json:"reviewer"`
		Notes          string `json:"notes"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	if err := s.manager.Approve(req.VerificationID, req.Reviewer, req.Notes); err != nil {
		writeJSON(w, http.StatusNotFound, map[string]any{"error": err.Error()})
		return
	}

	s.logger.Info("verification approved", "ver_id", req.VerificationID, "reviewer", req.Reviewer)
	writeJSON(w, http.StatusOK, map[string]any{"status": "approved"})
}

func (s *VerificationServer) HandleReject(w http.ResponseWriter, r *http.Request) {
	var req struct {
		VerificationID string `json:"verification_id"`
		Reviewer       string `json:"reviewer"`
		Reason         string `json:"reason"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	if err := s.manager.Reject(req.VerificationID, req.Reviewer, req.Reason); err != nil {
		writeJSON(w, http.StatusNotFound, map[string]any{"error": err.Error()})
		return
	}

	s.logger.Info("verification rejected", "ver_id", req.VerificationID, "reason", req.Reason)
	writeJSON(w, http.StatusOK, map[string]any{"status": "rejected"})
}

func (s *VerificationServer) HandleGetStatus(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("verification_id")
	v, err := s.manager.GetVerification(id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, v)
}

func (s *VerificationServer) HandleListPending(w http.ResponseWriter, r *http.Request) {
	pending := s.manager.ListByStatus(kyc.StatusPending)
	inReview := s.manager.ListByStatus(kyc.StatusInReview)
	writeJSON(w, http.StatusOK, map[string]any{"pending": len(pending), "in_review": len(inReview), "verifications": append(pending, inReview...)})
}

func (s *VerificationServer) HandleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"status": "ok", "service": "driver-verification"})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
