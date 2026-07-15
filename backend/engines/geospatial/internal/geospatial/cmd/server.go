package cmd

import (
	"encoding/json"
	"log/slog"
	"math"
	"net/http"

	"github.com/sekaishopml/cytaxi/backend/engines/geospatial/internal/geospatial/domain/types"
)

type GeospatialServer struct {
	logger *slog.Logger
}

func NewGeospatialServer(logger *slog.Logger) *GeospatialServer {
	return &GeospatialServer{logger: logger}
}

func (s *GeospatialServer) HandleUpdateLocation(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Lat     float64 `json:"lat"`
		Lng     float64 `json:"lng"`
		DestLat float64 `json:"dest_lat"`
		DestLng float64 `json:"dest_lng"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	origin := types.Coordinates{Lat: req.Lat, Lng: req.Lng}
	dest := types.Coordinates{Lat: req.DestLat, Lng: req.DestLng}

	distanceMeters := origin.DistanceTo(dest)
	distanceKM := distanceMeters / 1000
	etaSeconds := int((distanceMeters / 8.33) * 1.2) // avg speed ~30km/h with traffic factor

	s.logger.Info("location updated", "lat", req.Lat, "lng", req.Lng, "distance_km", distanceKM, "eta_sec", etaSeconds)

	writeJSON(w, http.StatusOK, map[string]any{
		"status":         "updated",
		"distance_meters": math.Round(distanceMeters),
		"distance_km":    math.Round(distanceKM*100) / 100,
		"eta_seconds":    etaSeconds,
		"eta_minutes":    etaSeconds / 60,
	})
}

func (s *GeospatialServer) HandleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"status": "ok", "service": "geospatial-engine"})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
