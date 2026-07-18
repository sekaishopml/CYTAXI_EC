package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/sekaishopml/cytaxi/backend/engines/matching/internal/matching/api/handler"
	"github.com/sekaishopml/cytaxi/backend/engines/matching/internal/matching/api/router"
	"github.com/sekaishopml/cytaxi/backend/engines/matching/internal/matching/config"
	"github.com/sekaishopml/cytaxi/backend/engines/matching/internal/matching/application/service"
	"github.com/sekaishopml/cytaxi/backend/engines/matching/internal/matching/infrastructure/repository"
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

	// Repositorios in-memory
	matchRepo := repository.NewInMemoryMatchingRepository()
	candidateRepo := repository.NewInMemoryCandidateRepository()

	// Servicio
	matchingService := service.NewMatchingService(matchRepo, candidateRepo, logger)

	mux := http.NewServeMux()
	h := handler.New(matchingService, logger)
	r := router.New(mux, h)

	logger.Info("matching engine starting", "port", cfg.Port)

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
