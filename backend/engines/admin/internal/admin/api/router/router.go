package router

import (
	"net/http"

	"github.com/sekaishopml/cytaxi/backend/engines/admin/internal/admin/api/handler"
)

func New(mux *http.ServeMux, h *handler.Handler) http.Handler {
	mux.HandleFunc("GET /health", h.Health)
	mux.HandleFunc("GET /admin/roles", h.GetRoles)
	mux.HandleFunc("GET /admin/feature-flags", h.GetFeatureFlags)
	return mux
}
