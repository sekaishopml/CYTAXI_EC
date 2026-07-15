package cmd

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/sekaishopml/cytaxi/backend/engines/geospatial/internal/geospatial/domain/types"
	"github.com/sekaishopml/cytaxi/backend/engines/geospatial/internal/geospatial/infrastructure/cache"
	"github.com/sekaishopml/cytaxi/backend/engines/geospatial/internal/geospatial/infrastructure/real"
)

type GeoServer struct {
	provider types.GeospatialProvider
	cache    *cache.GeospatialCache
	logger   *slog.Logger
}

func NewGeoServer(provider types.GeospatialProvider, logger *slog.Logger) *GeoServer {
	return &GeoServer{
		provider: provider,
		cache:    cache.NewGeospatialCache(5 * time.Minute),
		logger:   logger,
	}
}

func (s *GeoServer) HandleSearch(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if q == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "q parameter required"})
		return
	}

	start := time.Now()
	key := "search:" + q

	result, err := s.cache.GetOrFetch(r.Context(), key, func() (any, error) {
		return s.provider.SearchPlaces(types.PlaceSearchRequest{Query: q, MaxResult: 5})
	})
	if err != nil {
		s.logger.Error("search failed", "q", q, "error", err)
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}

	s.logger.Info("geo search", "q", q, "provider", s.provider.Name(), "latency", time.Since(start).String())
	writeJSON(w, http.StatusOK, result)
}

func (s *GeoServer) HandleGeocode(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	if address == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "address required"})
		return
	}

	decoded, _ := url.QueryUnescape(address)
	key := "geocode:" + decoded

	result, err := s.cache.GetOrFetch(r.Context(), key, func() (any, error) {
		return s.provider.Geocode(decoded)
	})
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (s *GeoServer) HandleReverseGeocode(w http.ResponseWriter, r *http.Request) {
	lat, lng := r.URL.Query().Get("lat"), r.URL.Query().Get("lng")
	key := fmt.Sprintf("reverse:%s,%s", lat, lng)

	result, err := s.cache.GetOrFetch(r.Context(), key, func() (any, error) {
		latF := parseFloat(lat)
		lngF := parseFloat(lng)
		return s.provider.ReverseGeocode(types.Coordinates{Lat: latF, Lng: lngF})
	})
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (s *GeoServer) HandleRoute(w http.ResponseWriter, r *http.Request) {
	var req struct {
		OriginLat, OriginLng     float64 `json:"origin_lat,origin_lng"`
		DestLat, DestLng         float64 `json:"dest_lat,dest_lng"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	key := fmt.Sprintf("route:%f,%f-%f,%f", req.OriginLat, req.OriginLng, req.DestLat, req.DestLng)

	result, err := s.cache.GetOrFetch(r.Context(), key, func() (any, error) {
		return s.provider.FindRoute(types.RouteRequest{
			Origin:      types.Coordinates{Lat: req.OriginLat, Lng: req.OriginLng},
			Destination: types.Coordinates{Lat: req.DestLat, Lng: req.DestLng},
			Mode:        types.TravelModeDriving,
		})
	})
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (s *GeoServer) HandleDistance(w http.ResponseWriter, r *http.Request) {
	var req struct {
		OriginLat, OriginLng float64
		DestLat, DestLng     float64
	}
	json.NewDecoder(r.Body).Decode(&req)

	origin := types.Coordinates{Lat: req.OriginLat, Lng: req.OriginLng}
	dest := types.Coordinates{Lat: req.DestLat, Lng: req.DestLng}
	dist := origin.DistanceTo(dest)

	writeJSON(w, http.StatusOK, map[string]any{
		"distance_meters": int(dist),
		"distance_km":     fmt.Sprintf("%.2f", dist/1000),
	})
}

func (s *GeoServer) HandleETA(w http.ResponseWriter, r *http.Request) {
	var req struct {
		OriginLat, OriginLng float64
		DestLat, DestLng     float64
	}
	json.NewDecoder(r.Body).Decode(&req)

	origin := types.Coordinates{Lat: req.OriginLat, Lng: req.OriginLng}
	dest := types.Coordinates{Lat: req.DestLat, Lng: req.DestLng}
	dist := origin.DistanceTo(dest)
	etaSec := int(dist / 8.33)

	writeJSON(w, http.StatusOK, map[string]any{
		"eta_seconds": etaSec,
		"eta_minutes": etaSec / 60,
		"distance_km": fmt.Sprintf("%.2f", dist/1000),
	})
}

func (s *GeoServer) HandleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"status": "ok", "service": "geospatial-engine", "provider": s.provider.Name()})
}

func parseFloat(s string) float64 {
	var f float64
	fmt.Sscanf(s, "%f", &f)
	return f
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
