package router

import (
	"net/http"

	"github.com/sekaishopml/cytaxi/backend/engines/notification/internal/notification/api/handler"
)

func New(mux *http.ServeMux, h *handler.Handler) http.Handler {
	mux.HandleFunc("GET /health", h.Health)
	mux.HandleFunc("GET /notifications/{notification_id}", h.GetNotification)
	mux.HandleFunc("GET /recipients/{recipient_id}/notifications", h.GetHistory)
	mux.HandleFunc("GET /notifications/templates", h.GetTemplates)
	return mux
}
