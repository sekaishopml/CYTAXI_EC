package router

import (
	"net/http"

	"github.com/sekaishopml/cytaxi/backend/engines/matching/internal/matching/api/handler"
)

func New(mux *http.ServeMux, h *handler.Handler) http.Handler {
	mux.HandleFunc("GET /health", h.Health)
	mux.HandleFunc("GET /matching/{matching_id}", h.GetMatching)
	mux.HandleFunc("GET /matching/{matching_id}/candidates", h.GetCandidates)
	return mux
}
