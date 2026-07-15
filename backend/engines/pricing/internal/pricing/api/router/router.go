package router

import (
	"net/http"

	"github.com/sekaishopml/cytaxi/backend/engines/pricing/internal/pricing/api/handler"
)

func New(mux *http.ServeMux, h *handler.Handler) http.Handler {
	mux.HandleFunc("GET /health", h.Health)
	mux.HandleFunc("GET /fares/{fare_id}", h.GetFare)
	mux.HandleFunc("GET /trips/{trip_id}/fares", h.GetHistory)
	mux.HandleFunc("GET /promotions", h.GetPromotions)
	return mux
}
