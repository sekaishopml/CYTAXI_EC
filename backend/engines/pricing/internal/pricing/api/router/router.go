package router

import (
	"net/http"

	"github.com/sekaishopml/cytaxi/backend/engines/pricing/internal/pricing/api/handler"
)

func New(mux *http.ServeMux, h *handler.Handler) http.Handler {
	// Health
	mux.HandleFunc("GET /health", h.Health)

	// Query endpoints (paths sin prefijo)
	mux.HandleFunc("GET /fares/{fare_id}", h.GetFare)
	mux.HandleFunc("GET /trips/{trip_id}/fares", h.GetHistory)
	mux.HandleFunc("GET /promotions", h.GetPromotions)

	// Command endpoints
	mux.HandleFunc("POST /pricing/estimate", h.HandleEstimate)
	mux.HandleFunc("POST /estimate", h.HandleEstimate)
	mux.HandleFunc("POST /fares", h.HandleEstimate)

	// Rutas con prefijo /api/v1/pricing (acceso directo)
	mux.HandleFunc("POST /api/v1/pricing/estimate", h.HandleEstimate)
	mux.HandleFunc("GET /api/v1/fares/{fare_id}", h.GetFare)

	return mux
}
