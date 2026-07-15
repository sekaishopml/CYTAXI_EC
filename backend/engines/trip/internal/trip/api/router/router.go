package router

import (
	"net/http"

	"github.com/sekaishopml/cytaxi/backend/engines/trip/internal/trip/api/handler"
)

func New(mux *http.ServeMux, h *handler.Handler) http.Handler {
	mux.HandleFunc("GET /health", h.Health)
	mux.HandleFunc("GET /trips/{trip_id}", h.GetTrip)
	mux.HandleFunc("GET /customers/{customer_id}/trips", h.GetTripHistory)
	mux.HandleFunc("GET /trips/active", h.GetActiveTrips)
	return mux
}
