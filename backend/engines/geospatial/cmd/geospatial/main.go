package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	geocmd "github.com/sekaishopml/cytaxi/backend/engines/geospatial/internal/geospatial/cmd"
	"github.com/sekaishopml/cytaxi/backend/engines/geospatial/internal/geospatial/config"
	"github.com/sekaishopml/cytaxi/backend/engines/geospatial/internal/geospatial/domain/service"
	"github.com/sekaishopml/cytaxi/backend/engines/geospatial/internal/geospatial/infrastructure/real"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cfg, err := config.Load()
	if err != nil {
		slog.Error("config load failed", "error", err)
		os.Exit(1)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// Seleccionar provider
	var provider service.GeospatialProvider
	switch cfg.Provider {
	case "openstreetmap", "osm", "real", "nominatim":
		provider = real.NewProvider(nil)
		logger.Info("geospatial provider selected", "provider", "openstreetmap", "geocoder", "nominatim.openstreetmap.org", "router", "router.project-osrm.org")
	case "googlemaps", "google", "google_maps":
		// Fallback a OpenStreetMap mientras Google Maps está en stub
		provider = real.NewProvider(nil)
		logger.Warn("google_maps requested but using openstreetmap fallback (adapter not implemented yet)")
	default:
		provider = real.NewProvider(nil)
		logger.Warn("unknown provider, defaulting to openstreetmap", "provider", cfg.Provider)
	}

	logger.Info("geospatial engine starting", "port", cfg.Port, "provider", provider.Name(), "env", cfg.Env)

	// Servidor HTTP con todos los handlers montados
	mux := http.NewServeMux()
	server := geocmd.NewGeoServer(provider, logger)
	// Los handlers se montan SIN prefijo /api/v1/geo porque el gateway lo remueve
	mux.HandleFunc("GET /health", server.HandleHealth)
	mux.HandleFunc("GET /search", server.HandleSearch)
	mux.HandleFunc("GET /geocode", server.HandleGeocode)
	mux.HandleFunc("GET /reverse", server.HandleReverseGeocode)
	mux.HandleFunc("POST /route", server.HandleRoute)
	mux.HandleFunc("POST /distance", server.HandleDistance)
	mux.HandleFunc("POST /eta", server.HandleETA)

	// Rutas legacy sin prefijo (compatibilidad con versiones anteriores)
	mux.HandleFunc("GET /api/search", server.HandleSearch)
	mux.HandleFunc("GET /api/geocode", server.HandleGeocode)
	mux.HandleFunc("GET /api/reverse", server.HandleReverseGeocode)
	mux.HandleFunc("POST /api/route", server.HandleRoute)
	mux.HandleFunc("POST /api/distance", server.HandleDistance)
	mux.HandleFunc("POST /api/eta", server.HandleETA)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: mux,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	logger.Info("shutting down")
	srv.Shutdown(context.Background())
}
