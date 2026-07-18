package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/sekaishopml/cytaxi/backend/engines/pricing/internal/pricing/api/handler"
	"github.com/sekaishopml/cytaxi/backend/engines/pricing/internal/pricing/api/router"
	"github.com/sekaishopml/cytaxi/backend/engines/pricing/internal/pricing/config"
	"github.com/sekaishopml/cytaxi/backend/engines/pricing/internal/pricing/application/service"
	"github.com/sekaishopml/cytaxi/backend/engines/pricing/internal/pricing/infrastructure/repository"
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
	fareRepo := repository.NewInMemoryFareRepository()
	promoRepo := repository.NewInMemoryPromotionRepository()
	couponRepo := repository.NewInMemoryCouponRepository()

	// Servicio
	pricingService := service.NewPricingService(fareRepo, promoRepo, couponRepo, logger)

	mux := http.NewServeMux()
	h := handler.New(pricingService, logger)
	r := router.New(mux, h)

	logger.Info("pricing engine starting", "port", cfg.Port)

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
