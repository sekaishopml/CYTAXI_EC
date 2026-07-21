package router

import (
	"net/http"

	"github.com/sekaishopml/cytaxi/backend/engines/trip/internal/trip/api/handler"
)

func New(mux *http.ServeMux, h *handler.Handler) http.Handler {
	mux.HandleFunc("GET /health", h.Health)
	mux.HandleFunc("GET /trips/{trip_id}", h.GetTrip)
	mux.HandleFunc("GET /{trip_id}", h.GetTrip)
	mux.HandleFunc("GET /customers/{customer_id}/trips", h.GetTripHistory)
	mux.HandleFunc("GET /trips/active", h.GetActiveTrips)

	mux.HandleFunc("POST /trip/request", h.HandleCreateTrip)
	mux.HandleFunc("POST /request", h.HandleCreateTrip)
	mux.HandleFunc("POST /trips", h.HandleCreateTrip)

	mux.HandleFunc("POST /trips/{trip_id}/cancel", h.HandleCancelTrip)
	mux.HandleFunc("POST /trip/cancel", h.HandleCancelTrip)
	mux.HandleFunc("POST /cancel", h.HandleCancelTrip)

	mux.HandleFunc("PUT /trips/{trip_id}/destination", h.HandleChangeDestination)
	mux.HandleFunc("POST /trip/change-destination", h.HandleChangeDestination)

	mux.HandleFunc("POST /trip/reject-driver", h.HandleRejectDriver)
	mux.HandleFunc("POST /reject-driver", h.HandleRejectDriver)

	// Driver flow: accept / start / complete / location
	mux.HandleFunc("POST /trips/{trip_id}/accept", h.HandleAcceptDriver)
	mux.HandleFunc("POST /trip/accept", h.HandleAcceptDriver)
	mux.HandleFunc("POST /accept", h.HandleAcceptDriver)

	mux.HandleFunc("POST /trips/{trip_id}/start", h.HandleStartTrip)
	mux.HandleFunc("POST /trip/start", h.HandleStartTrip)
	mux.HandleFunc("POST /start", h.HandleStartTrip)

	mux.HandleFunc("POST /trips/{trip_id}/complete", h.HandleCompleteTrip)
	mux.HandleFunc("POST /trip/complete", h.HandleCompleteTrip)
	mux.HandleFunc("POST /trip/finish", h.HandleCompleteTrip)
	mux.HandleFunc("POST /complete", h.HandleCompleteTrip)
	mux.HandleFunc("POST /finish", h.HandleCompleteTrip)

	mux.HandleFunc("POST /trips/{trip_id}/location", h.HandleLocation)
	mux.HandleFunc("POST /trip/location", h.HandleLocation)
	mux.HandleFunc("POST /location", h.HandleLocation)

	// Full paths for gateway compat
	mux.HandleFunc("POST /api/v1/trip/request", h.HandleCreateTrip)
	mux.HandleFunc("POST /api/v1/trip/cancel", h.HandleCancelTrip)
	mux.HandleFunc("PUT /api/v1/trip/change-destination", h.HandleChangeDestination)
	mux.HandleFunc("POST /api/v1/trip/reject-driver", h.HandleRejectDriver)
	mux.HandleFunc("POST /api/v1/trip/accept", h.HandleAcceptDriver)
	mux.HandleFunc("POST /api/v1/trip/start", h.HandleStartTrip)
	mux.HandleFunc("POST /api/v1/trip/complete", h.HandleCompleteTrip)
	mux.HandleFunc("POST /api/v1/trip/location", h.HandleLocation)

	return mux
}
