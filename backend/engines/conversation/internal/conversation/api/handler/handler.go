package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/sekaishopml/cytaxi/backend/engines/conversation/internal/conversation/application/dto"
	"github.com/sekaishopml/cytaxi/backend/engines/conversation/internal/conversation/application/port"
)

type Handler struct {
	inputPort port.MessageInputPort
	logger    *slog.Logger
}

func New(inputPort port.MessageInputPort, logger *slog.Logger) *Handler {
	return &Handler{
		inputPort: inputPort,
		logger:    logger,
	}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"status":    "ok",
		"service":   "conversation-engine",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *Handler) IncomingMessage(w http.ResponseWriter, r *http.Request) {
	var req dto.IncomingMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Phone == "" || req.Content == "" {
		http.Error(w, "phone and content are required", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	if err := h.inputPort.HandleIncomingMessage(ctx, req.Phone, req.Content); err != nil {
		h.logger.Error("failed to handle message", "error", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(dto.IncomingMessageResponse{
		Status:  "accepted",
		Message: "message received",
	})
}

func (h *Handler) NotFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{
		"error": fmt.Sprintf("route %s %s not found", r.Method, r.URL.Path),
	})
}
