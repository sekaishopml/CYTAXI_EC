package cmd

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/sekaishopml/cytaxi/backend/engines/conversation/internal/conversation/application/port"
	"github.com/sekaishopml/cytaxi/backend/engines/conversation/internal/conversation/application/usecase"
	"github.com/sekaishopml/cytaxi/backend/engines/conversation/internal/conversation/domain/entity"
)

type ConversationServer struct {
	pipeline *usecase.MessagePipeline
	sessionUC *usecase.SessionUseCase
	logger   *slog.Logger
}

func NewConversationServer(pipeline *usecase.MessagePipeline, sessionUC *usecase.SessionUseCase, logger *slog.Logger) *ConversationServer {
	return &ConversationServer{pipeline: pipeline, sessionUC: sessionUC, logger: logger}
}

func (s *ConversationServer) HandleStartConversation(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Phone   string `json:"phone"`
		Message string `json:"message"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
		return
	}

	session, err := s.sessionUC.CreateSession(r.Context(), req.Phone)
	if err != nil {
		http.Error(w, `{"error":"failed to create session"}`, http.StatusInternalServerError)
		return
	}

	s.logger.Info("conversation started", "phone", req.Phone, "session_id", session.ID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"status":      "ok",
		"session_id":  session.ID,
		"conversation_id": session.ConversationID,
	})
}

func (s *ConversationServer) HandleIncomingMessage(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Phone   string `json:"phone"`
		Content string `json:"content"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	err := s.pipeline.ProcessIncoming(r.Context(), req.Phone, req.Content)
	if err != nil {
		s.logger.Error("message processing failed", "error", err)
		http.Error(w, `{"error":"processing failed"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"status":  "accepted",
		"message": "Message processed",
	})
}

func (s *ConversationServer) HandleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"status": "ok", "service": "conversation-engine"})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
