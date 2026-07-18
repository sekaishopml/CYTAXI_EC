package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/sekaishopml/cytaxi/backend/engines/trip/internal/trip/api/handler"
	"github.com/sekaishopml/cytaxi/backend/engines/trip/internal/trip/api/router"
	"github.com/sekaishopml/cytaxi/backend/engines/trip/internal/trip/config"
	"github.com/sekaishopml/cytaxi/backend/engines/trip/internal/trip/application/service"
	"github.com/sekaishopml/cytaxi/backend/engines/trip/internal/trip/infrastructure/repository"
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

	// Repositorios in-memory (sustituibles por PostgreSQL sin cambiar la interfaz)
	tripRepo := repository.NewInMemoryTripRepository()
	timelineRepo := repository.NewInMemoryTimelineRepository()
	assignmentRepo := repository.NewInMemoryAssignmentRepository()

	// Servicio de aplicación
	tripService := service.NewTripService(tripRepo, timelineRepo, assignmentRepo, logger)

	mux := http.NewServeMux()
	h := handler.New(tripService, logger)
	r := router.New(mux, h)

	logger.Info("trip engine starting", "port", cfg.Port)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: r,
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
