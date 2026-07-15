package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/sekaishopml/cytaxi/backend/engines/mobility/internal/mobility/application/dispatcher"
	"github.com/sekaishopml/cytaxi/backend/engines/mobility/internal/mobility/application/pipeline"
	"github.com/sekaishopml/cytaxi/backend/engines/mobility/internal/mobility/config"
	"github.com/sekaishopml/cytaxi/backend/engines/mobility/internal/mobility/infrastructure/strategy"
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

	strategyRegistry := strategy.NewStrategyRegistry()
	strategyRegistry.Register(&strategy.NearestDriver{})
	strategyRegistry.Register(&strategy.HighestRated{})
	strategyRegistry.Register(strategy.NewBalancedScore(0.6, 0.4))

	selStrategy, _ := strategyRegistry.Get(cfg.DefaultStrategy)
	decPipeline := pipeline.NewDecisionPipeline(selStrategy)

	var candidateFinder dispatcher.CandidateFinder = nil
	coordinator := dispatcher.NewDispatcherCoordinator(decPipeline, candidateFinder, logger)
	_ = coordinator

	logger.Info("mobility engine starting",
		"port", cfg.Port,
		"default_strategy", cfg.DefaultStrategy,
	)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok","service":"mobility-engine"}`))
	})

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
