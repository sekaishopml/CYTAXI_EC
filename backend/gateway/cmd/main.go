package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/sekaishopml/cytaxi/backend/gateway/internal/config"
	"github.com/sekaishopml/cytaxi/backend/gateway/internal/handler"
	"github.com/sekaishopml/cytaxi/backend/gateway/internal/middleware"
	"github.com/sekaishopml/cytaxi/backend/gateway/internal/router"
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

	mwChain := middleware.NewChain(logger)
	mwChain.Use(middleware.Recovery)
	mwChain.Use(middleware.Correlation)
	mwChain.Use(middleware.CORS)
	mwChain.Use(middleware.RequestLogger)
	mwChain.Use(middleware.RateLimiter(cfg.RateLimitRPS))
	mwChain.Use(middleware.AuthJWT(cfg.AuthSecret))

	gwRouter := router.New(cfg.BackendHosts, mwChain, logger)
	gwRouter.RegisterRoutes()

	handler := gwRouter.Handler()

	logger.Info("api gateway starting",
		"port", cfg.Port,
		"backends", len(cfg.BackendHosts),
	)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: handler,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("gateway error", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	logger.Info("shutting down gateway")
	srv.Shutdown(context.Background())
}
