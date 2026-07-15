package router

import (
	"net/http"

	"github.com/sekaishopml/cytaxi/backend/engines/analytics/internal/analytics/api/handler"
)

func New(mux *http.ServeMux, h *handler.Handler) http.Handler {
	mux.HandleFunc("GET /health", h.Health)
	mux.HandleFunc("GET /analytics/dashboard", h.GetDashboard)
	mux.HandleFunc("GET /analytics/metrics", h.GetMetrics)
	return mux
}
