package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/sekaishopml/cytaxi/backend/engines/notification/internal/notification/application/port"
	"github.com/sekaishopml/cytaxi/backend/engines/notification/internal/notification/application/query"
	"github.com/sekaishopml/cytaxi/backend/engines/notification/internal/notification/domain/valueobject"
)

type Handler struct {
	service port.NotificationService
	logger  *slog.Logger
}

func New(service port.NotificationService, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok","service":"notification-engine"}`))
}

func (h *Handler) GetNotification(w http.ResponseWriter, r *http.Request) {
	id := valueobject.NotificationID(r.PathValue("notification_id"))
	n, err := h.service.Get(r.Context(), query.GetNotification{NotificationID: id})
	if err != nil {
		http.Error(w, "notification not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, n)
}

func (h *Handler) GetHistory(w http.ResponseWriter, r *http.Request) {
	recipientID := valueobject.RecipientID(r.PathValue("recipient_id"))
	history, err := h.service.GetHistory(r.Context(), query.GetNotificationHistory{RecipientID: recipientID})
	if err != nil {
		http.Error(w, "history not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, history)
}

func (h *Handler) GetTemplates(w http.ResponseWriter, r *http.Request) {
	templates, err := h.service.GetTemplates(r.Context(), query.GetTemplates{})
	if err != nil {
		http.Error(w, "templates not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, templates)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
