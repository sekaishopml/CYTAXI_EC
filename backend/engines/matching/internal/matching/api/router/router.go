package router

import (
	"net/http"

	"github.com/sekaishopml/cytaxi/backend/engines/matching/internal/matching/api/handler"
)

func New(mux *http.ServeMux, h *handler.Handler) http.Handler {
	// Health
	mux.HandleFunc("GET /health", h.Health)

	// Query endpoints (paths sin prefijo /api/v1/matching)
	mux.HandleFunc("GET /matching/{matching_id}", h.GetMatching)
	mux.HandleFunc("GET /matching/{matching_id}/candidates", h.GetCandidates)

	// Command endpoints
	mux.HandleFunc("POST /matching/start", h.HandleStartMatching)
	mux.HandleFunc("POST /matching/select", h.HandleSelectDriver)
	mux.HandleFunc("POST /start", h.HandleStartMatching)
	mux.HandleFunc("POST /select", h.HandleSelectDriver)

	// Rutas con prefijo /api/v1/matching (acceso directo)
	mux.HandleFunc("POST /api/v1/matching/start", h.HandleStartMatching)
	mux.HandleFunc("POST /api/v1/matching/select", h.HandleSelectDriver)
	mux.HandleFunc("GET /api/v1/matching/{matching_id}", h.GetMatching)
	mux.HandleFunc("GET /api/v1/matching/{matching_id}/candidates", h.GetCandidates)

	return mux
}
